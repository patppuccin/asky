package asky

import (
	"bufio"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/mattn/go-runewidth"
	"golang.org/x/term"
)

// echoMode controls how typed characters are displayed during input.
type echoMode uint8

const (
	echoNormal echoMode = iota // characters echoed as-is
	echoMask                   // characters echoed as *
	echoSilent                 // nothing echoed
)

// input renders an interactive single-line text prompt.
// Construct one with [Input], [InputSecret], or [InputSilent].
type input struct {
	cfg          Config
	prefix       string
	label        string
	placeholder  string
	defaultValue string
	echo         echoMode
	validator    func(string) (string, bool)
}

// Input returns a builder for an interactive single-line text prompt.
//
//	name, err := asky.Input().WithLabel("Project name").Render()
//	if errors.Is(err, asky.ErrInterrupted) { ... }
func Input() *input {
	return &input{
		cfg:   pkgConfig,
		label: "Enter value",
		echo:  echoNormal,
	}
}

// InputSecret returns a builder for a masked prompt.
// Characters are echoed as * so the user receives visual feedback.
//
//	pass, err := asky.InputSecret().WithLabel("Password").Render()
func InputSecret() *input {
	i := Input()
	i.echo = echoMask
	return i
}

// InputSilent returns a builder for a completely silent prompt.
// Nothing is echoed as the user types and cursor position is not exposed.
//
//	pass, err := asky.InputSilent().WithLabel("Password").Render()
func InputSilent() *input {
	i := Input()
	i.echo = echoSilent
	return i
}

// WithStyles overrides the [StyleMap] for this prompt.
func (i *input) WithStyles(s *StyleMap) *input {
	i.cfg.Styles = s
	return i
}

// WithPrefix overrides the default prompt prefix symbol.
func (i *input) WithPrefix(p string) *input {
	i.prefix = p
	return i
}

// WithLabel sets the prompt label shown to the user.
func (i *input) WithLabel(l string) *input {
	i.label = l
	return i
}

// WithPlaceholder sets placeholder text shown when the input is empty.
func (i *input) WithPlaceholder(p string) *input {
	i.placeholder = p
	return i
}

// WithDefaultValue sets a default value used when the user submits empty input.
func (i *input) WithDefaultValue(v string) *input {
	i.defaultValue = v
	return i
}

// WithValidator sets a validation function called on every keystroke and on submit.
// Returns a message and false to block submission, or a message and true to allow.
func (i *input) WithValidator(fn func(string) (string, bool)) *input {
	i.validator = fn
	return i
}

// Render displays the interactive prompt and blocks until the user submits or
// cancels. Returns the entered string, or [ErrInterrupted] if Ctrl+C is pressed.
//
// In accessible mode, input is collected line-by-line.
// Validation is checked on Enter and the prompt reprints until satisfied.
func (i *input) Render() (string, error) {
	if i.cfg.Accessible {
		return i.renderAccessible()
	}
	return i.renderInteractive()
}

// renderAccessible collects input without cursor magic.
// Plain input echoes characters as typed using bufio.
// Secret echoes * per character; silent echoes nothing.
// Validation is checked on Enter and the prompt reprints on failure.
func (i *input) renderAccessible() (string, error) {
	prefix := pick(i.prefix, "(?)")
	promptLine := safeStyle(i.cfg.Styles.InputPrefix).Sprint(prefix) + " " +
		safeStyle(i.cfg.Styles.InputLabel).Sprint(i.label)

	placeholderLine := i.buildPlaceholderLine()

	for {
		stdOutput.Write([]byte(promptLine + "\n"))
		if placeholderLine != "" {
			stdOutput.Write([]byte(placeholderLine + "\n"))
		}

		var result string

		if i.echo != echoNormal {
			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

			type readResult struct {
				b   []byte
				err error
			}
			ch := make(chan readResult, 1)
			go func() {
				b, err := term.ReadPassword(int(os.Stdin.Fd()))
				ch <- readResult{b, err}
			}()

			select {
			case <-sigCh:
				signal.Stop(sigCh)
				return "", ErrInterrupted
			case r := <-ch:
				signal.Stop(sigCh)
				if r.err != nil {
					if isInterrupt(r.err) {
						return "", ErrInterrupted
					}
					return "", r.err
				}
				if i.echo == echoMask {
					stdOutput.Write([]byte(strings.Repeat("*", len(r.b)) + "\n"))
				} else {
					stdOutput.Write([]byte("\n"))
				}
				result = string(r.b)
			}
		} else {
			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

			type readResult struct {
				line string
				err  error
			}
			ch := make(chan readResult, 1)
			go func() {
				line, err := bufio.NewReader(os.Stdin).ReadString('\n')
				ch <- readResult{line, err}
			}()

			select {
			case <-sigCh:
				signal.Stop(sigCh)
				return "", ErrInterrupted
			case r := <-ch:
				signal.Stop(sigCh)
				if r.err != nil {
					if isInterrupt(r.err) {
						return "", ErrInterrupted
					}
					return "", r.err
				}
				result = strings.TrimRight(r.line, "\r\n")
			}
		}

		if i.validator != nil {
			msg, ok := i.validator(result)
			if !ok {
				stdOutput.Write([]byte(safeStyle(i.cfg.Styles.InputValidationFail).Sprint(msg) + "\n\n"))
				continue
			}
		}

		if result == "" && i.defaultValue != "" {
			result = i.defaultValue
		}

		return result, nil
	}
}

// renderInteractive renders the animated single-line prompt with live redraws.
func (i *input) renderInteractive() (string, error) {
	if err := reserveLines(5); err != nil {
		return "", ErrTerminalTooSmall
	}

	prefix := pick(i.prefix, "(?)")
	promptLine := safeStyle(i.cfg.Styles.InputPrefix).Sprint(prefix) + " " +
		safeStyle(i.cfg.Styles.InputLabel).Sprint(i.label)
	placeholderLine := i.buildPlaceholderLine()
	helpLine := safeStyle(i.cfg.Styles.InputHelp).Sprint("enter to confirm  •  ctrl+c to cancel")

	var inBuf []rune
	cursorPos := 0
	interrupted := false
	receivedInput := false

	// displayBuf returns the string to render based on echo mode.
	displayBuf := func() string {
		switch i.echo {
		case echoMask:
			return strings.Repeat("*", len(inBuf))
		case echoSilent:
			return ""
		default:
			return string(inBuf)
		}
	}

	// redraws prompt and input line, then repositions the cursor on the input.
	redraw := func(validationMsg string, validOK *bool) {
		// Prompt + placeholder + input line
		stdOutput.Write([]byte(ansiHideCursor + ansiRestoreCursor + "\r" + ansiClearLine + promptLine))
		if len(inBuf) == 0 {
			stdOutput.Write([]byte(" " + placeholderLine))
		} else {
			stdOutput.Write([]byte(" " + safeStyle(i.cfg.Styles.InputText).Sprint(displayBuf())))
		}

		// Spacer Line
		stdOutput.Write([]byte("\n\r" + ansiClearLine))

		// Validation Line
		stdOutput.Write([]byte("\n\r" + ansiClearLine))
		if receivedInput && i.validator != nil && validOK != nil && !*validOK && validationMsg != "" {
			stdOutput.Write([]byte(safeStyle(i.cfg.Styles.InputValidationFail).Sprint(validationMsg)))
		}

		// Help Line
		stdOutput.Write([]byte("\n\r" + ansiClearLine + helpLine))

		// Reposition cursor at input line
		stdOutput.Write([]byte(ansiRestoreCursor + "\r" + promptLine + " "))
		if len(inBuf) > 0 {
			stdOutput.Write([]byte(safeStyle(i.cfg.Styles.InputText).Sprint(displayBuf())))
			// Reposition cursor within the input for plain mode only.
			// Use visual column width (runewidth) so wide chars (CJK, emoji) land correctly.
			if i.echo == echoNormal && cursorPos < len(inBuf) {
				afterCols := runewidth.StringWidth(string(inBuf[cursorPos:]))
				ansiCursorLeft(afterCols)
			}
		}

		stdOutput.Write([]byte(ansiShowCursor))
	}

	stdOutput.Write([]byte(ansiHideCursor + ansiSaveCursor))
	defer stdOutput.Write([]byte(ansiRestoreCursor + ansiClearScreen + ansiReset + ansiShowCursor))

	redraw("", nil)

	err := listenKeys(func(ev keyEvent) (stop bool) {
		switch ev.code {
		case keyCtrlC:
			interrupted = true
			return true

		case keyEnter:
			if i.validator != nil {
				msg, ok := i.validator(string(inBuf))
				if !ok {
					redraw(msg, &ok)
					return false
				}
			}
			if len(inBuf) == 0 && i.defaultValue != "" {
				inBuf = []rune(i.defaultValue)
			}
			return true

		case keyLeft:
			if i.echo == echoNormal && cursorPos > 0 {
				cursorPos--
			}

		case keyRight:
			if i.echo == echoNormal && cursorPos < len(inBuf) {
				cursorPos++
			}

		case keyHome, keyCtrlHome:
			if i.echo == echoNormal {
				cursorPos = 0
			}

		case keyEnd, keyCtrlEnd:
			if i.echo == echoNormal {
				cursorPos = len(inBuf)
			}

		case keyCtrlLeft:
			if i.echo == echoNormal && cursorPos > 0 {
				cursorPos--
				for cursorPos > 0 && inBuf[cursorPos-1] == ' ' {
					cursorPos--
				}
				for cursorPos > 0 && inBuf[cursorPos-1] != ' ' {
					cursorPos--
				}
			}

		case keyCtrlRight:
			if i.echo == echoNormal && cursorPos < len(inBuf) {
				for cursorPos < len(inBuf) && inBuf[cursorPos] == ' ' {
					cursorPos++
				}
				for cursorPos < len(inBuf) && inBuf[cursorPos] != ' ' {
					cursorPos++
				}
			}

		case keyBackspace:
			if i.echo != echoNormal {
				if len(inBuf) > 0 {
					inBuf = inBuf[:len(inBuf)-1]
					cursorPos = len(inBuf)
				}
			} else if cursorPos > 0 {
				inBuf = append(inBuf[:cursorPos-1], inBuf[cursorPos:]...)
				cursorPos--
			}

		case keyDelete:
			if i.echo == echoNormal && cursorPos < len(inBuf) {
				inBuf = append(inBuf[:cursorPos], inBuf[cursorPos+1:]...)
			}

		case keyRune:
			inBuf = append(inBuf[:cursorPos], append([]rune{ev.r}, inBuf[cursorPos:]...)...)
			cursorPos++
		}

		receivedInput = true

		if i.validator != nil {
			msg, ok := i.validator(string(inBuf))
			redraw(msg, &ok)
		} else {
			redraw("", nil)
		}
		return false
	})

	if err != nil {
		return "", err
	}
	if interrupted {
		return "", ErrInterrupted
	}

	return strings.TrimRight(string(inBuf), "\r\n"), nil
}

// buildPlaceholderLine composes the styled placeholder string from the prompt's
// placeholder and defaultValue fields. Returns empty string if neither is set.
func (i *input) buildPlaceholderLine() string {
	switch {
	case i.placeholder != "" && i.defaultValue != "":
		return safeStyle(i.cfg.Styles.InputPlaceholder).Sprint(i.placeholder + " (default: " + i.defaultValue + ")")
	case i.placeholder != "":
		return safeStyle(i.cfg.Styles.InputPlaceholder).Sprint(i.placeholder)
	case i.defaultValue != "":
		return safeStyle(i.cfg.Styles.InputPlaceholder).Sprint("default: " + i.defaultValue)
	default:
		return ""
	}
}

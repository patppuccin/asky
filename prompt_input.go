package asky

import (
	"bufio"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"golang.org/x/term"
)

// input renders an interactive single-line text prompt.
// Construct one with [Input], [InputSecret], or [InputSilent].
type input struct {
	cfg          Config
	prefix       string
	label        string
	placeholder  string
	defaultValue string
	secret       bool
	silent       bool
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
	}
}

// InputSecret returns a builder for a masked prompt.
// Characters are echoed as * so the user receives visual feedback.
//
//	pass, err := asky.InputSecret().WithLabel("Password").Render()
func InputSecret() *input {
	i := Input()
	i.secret = true
	return i
}

// InputSilent returns a builder for a completely silent prompt.
// Nothing is echoed as the user types
// The cursor position is not exposed as well.
//
//	pass, err := asky.InputSilent().WithLabel("Password").Render()
func InputSilent() *input {
	i := Input()
	i.silent = true
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
// Returns a message and a boolean to block submission (false) or allow (true).
func (i *input) WithValidator(fn func(string) (string, bool)) *input {
	i.validator = fn
	return i
}

// Render displays the interactive prompt and blocks until the user submits or
// cancels. Returns the entered string, or [ErrInterrupted] if Ctrl+C is pressed.
//
// In the accessible mode, input is collected line-by-line
// Validation is checked on Enter and the prompt reprints until satisfied.
func (i *input) Render() (string, error) {
	if i.cfg.Accessible {
		return i.renderAccessible()
	}
	return i.renderInteractive()
}

// renderAccessible collects input without cursor magic.
// Plain input echoes characters as typed using bufio.
// Secret echoes * per character, silent echoes nothing — both use term.ReadPassword.
// On Enter, validation is checked and the prompt reprints on failure.
func (i *input) renderAccessible() (string, error) {
	prefix := pick(i.prefix, "(?)")
	promptLine := safeStyle(i.cfg.Styles.InputPrefix).Sprint(prefix) + " " +
		safeStyle(i.cfg.Styles.InputLabel).Sprint(i.label)

	for {
		stdOutput.Write([]byte(promptLine + "\n"))

		var placeholderLine string
		switch {
		case i.placeholder != "" && i.defaultValue != "":
			placeholderLine = safeStyle(i.cfg.Styles.InputPlaceholder).Sprint(i.placeholder + " (default: " + i.defaultValue + ")")
		case i.placeholder != "":
			placeholderLine = safeStyle(i.cfg.Styles.InputPlaceholder).Sprint(i.placeholder)
		case i.defaultValue != "":
			placeholderLine = safeStyle(i.cfg.Styles.InputPlaceholder).Sprint("default: " + i.defaultValue)
		}
		if placeholderLine != "" {
			stdOutput.Write([]byte(placeholderLine + "\n"))
		}

		var result string

		if i.secret || i.silent {
			type readResult struct {
				b   []byte
				err error
			}

			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

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
					return "", r.err
				}
				if i.secret {
					stdOutput.Write([]byte(strings.Repeat("*", len(r.b)) + "\n"))
				} else {
					stdOutput.Write([]byte("\n"))
				}
				result = string(r.b)
			}
		} else {
			type readResult struct {
				line string
				err  error
			}

			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

			ch := make(chan readResult, 1)
			go func() {
				reader := bufio.NewReader(os.Stdin)
				line, err := reader.ReadString('\n')
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

		// Validate input
		if i.validator != nil {
			msg, ok := i.validator(result)
			if !ok {
				stdOutput.Write([]byte(safeStyle(i.cfg.Styles.InputValidationFail).Sprint(msg) + "\n"))
				continue
			}
		}

		// Apply default after validation passes
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

	var placeholderLine string
	switch {
	case i.placeholder != "" && i.defaultValue != "":
		placeholderLine = safeStyle(i.cfg.Styles.InputPlaceholder).Sprint(i.placeholder + " (default: " + i.defaultValue + ")")
	case i.placeholder != "":
		placeholderLine = safeStyle(i.cfg.Styles.InputPlaceholder).Sprint(i.placeholder)
	case i.defaultValue != "":
		placeholderLine = safeStyle(i.cfg.Styles.InputPlaceholder).Sprint("default: " + i.defaultValue)
	}

	helpLine := safeStyle(i.cfg.Styles.InputHelp).Sprint("enter (confirm) | ctrl+c (cancel)")

	var inBuf []rune
	cursorPos := 0
	interrupted := false
	receivedInput := false

	// displayBuf returns the string to render based on mode.
	displayBuf := func() string {
		switch {
		case i.silent:
			return ""
		case i.secret:
			return strings.Repeat("*", len(inBuf))
		default:
			return string(inBuf)
		}
	}

	redraw := func(validationMsg string, validOK *bool) {
		stdOutput.Write([]byte(ansiHideCursor + ansiRestoreCursor + ansiClearLine))

		// Prompt + input line
		stdOutput.Write([]byte("\r" + promptLine))
		if len(inBuf) == 0 {
			stdOutput.Write([]byte(" " + placeholderLine + ansiClearLine))
		} else {
			stdOutput.Write([]byte(" " + safeStyle(i.cfg.Styles.InputText).Sprint(displayBuf()) + ansiClearLine))
		}

		// Empty line before validation
		stdOutput.Write([]byte("\n\r" + ansiClearLine))

		// Validation line
		stdOutput.Write([]byte("\n\r" + ansiClearLine))
		if i.validator != nil && validationMsg != "" && receivedInput {
			if validOK != nil && !*validOK {
				stdOutput.Write([]byte(safeStyle(i.cfg.Styles.InputValidationFail).Sprint(validationMsg)))
			} else {
				stdOutput.Write([]byte(safeStyle(i.cfg.Styles.InputValidationPass).Sprint(validationMsg)))
			}
		}

		// Help line
		stdOutput.Write([]byte("\n\r" + helpLine + ansiClearLine))

		// Reposition cursor at input line
		stdOutput.Write([]byte(ansiRestoreCursor + "\r" + promptLine + " "))
		if len(inBuf) > 0 {
			stdOutput.Write([]byte(safeStyle(i.cfg.Styles.InputText).Sprint(displayBuf())))
			// Only reposition cursor for plain input & secret; for don't expose position
			if !i.secret && cursorPos < len(inBuf) {
				ansiCursorLeft(len(inBuf) - cursorPos)
			}
		}
		stdOutput.Write([]byte(ansiShowCursor))
	}

	stdOutput.Write([]byte(ansiHideCursor + ansiSaveCursor))
	defer stdOutput.Write([]byte(ansiRestoreCursor + ansiClearScreen + ansiReset + ansiShowCursor))

	redraw("", nil)

	err := keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		receivedInput = true

		switch key.Code {
		case keys.CtrlC:
			interrupted = true
			return true, nil

		case keys.Enter:
			if i.validator != nil {
				msg, ok := i.validator(string(inBuf))
				if !ok {
					redraw(msg, &ok)
					return false, nil
				}
			}
			if len(inBuf) == 0 && i.defaultValue != "" {
				inBuf = []rune(i.defaultValue)
			}
			return true, nil

		case keys.Left:
			if !i.secret && !i.silent && cursorPos > 0 {
				cursorPos--
			}

		case keys.Right:
			if !i.secret && !i.silent && cursorPos < len(inBuf) {
				cursorPos++
			}

		case keys.Backspace:
			if i.secret || i.silent {
				// Always delete from end — position not exposed
				if len(inBuf) > 0 {
					inBuf = inBuf[:len(inBuf)-1]
					cursorPos = len(inBuf)
				}
			} else {
				if cursorPos > 0 {
					inBuf = append(inBuf[:cursorPos-1], inBuf[cursorPos:]...)
					cursorPos--
				}
			}

		case keys.Delete:
			if !i.secret && !i.silent && cursorPos < len(inBuf) {
				inBuf = append(inBuf[:cursorPos], inBuf[cursorPos+1:]...)
			}

		case keys.Space:
			inBuf = append(inBuf[:cursorPos], append([]rune{' '}, inBuf[cursorPos:]...)...)
			cursorPos++

		case keys.RuneKey:
			if len(key.Runes) > 0 {
				inBuf = append(inBuf[:cursorPos], append([]rune{key.Runes[0]}, inBuf[cursorPos:]...)...)
				cursorPos++
			}
		}

		if i.validator != nil && receivedInput {
			msg, ok := i.validator(string(inBuf))
			redraw(msg, &ok)
		} else {
			redraw("", nil)
		}
		return false, nil
	})

	if err != nil {
		return "", err
	}
	if interrupted {
		return "", ErrInterrupted
	}

	return strings.TrimRight(string(inBuf), "\r\n"), nil
}

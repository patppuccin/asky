package asky

import (
	"strings"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
)

// --- Definition ------------------------------------------
type SecureInput struct {
	theme       *Theme
	style       *Style
	prefix      string
	label       string
	description string
	placeholder string
	noEcho      bool
	validator   func(string) (string, bool)
}

// --- Initiation ------------------------------------------
func NewSecureInput() *SecureInput {
	return &SecureInput{
		prefix:    "[?] ",
		label:     "Enter secure input",
		noEcho:    false,
		validator: nil,
	}
}

// --- Configuration ---------------------------------------
func (si *SecureInput) WithTheme(theme Theme) *SecureInput      { si.theme = &theme; return si }
func (si *SecureInput) WithStyle(style Style) *SecureInput      { si.style = &style; return si }
func (si *SecureInput) WithPrefix(p string) *SecureInput        { si.prefix = p; return si }
func (si *SecureInput) WithLabel(p string) *SecureInput         { si.label = p; return si }
func (si *SecureInput) WithDescription(txt string) *SecureInput { si.description = txt; return si }
func (si *SecureInput) WithPlaceholder(txt string) *SecureInput { si.placeholder = txt; return si }
func (si *SecureInput) WithNoEcho() *SecureInput                { si.noEcho = true; return si }
func (si *SecureInput) WithValidator(fn func(string) (string, bool)) *SecureInput {
	si.validator = fn
	return si
}

// --- Presentation ----------------------------------------
func (si *SecureInput) Render() (string, error) {
	// Setup theme and style (apply defaults if not set)
	if si.theme == nil {
		si.theme = &ThemeDefault
	}
	if si.style == nil {
		si.style = StyleDefault(si.theme)
	}

	// Ensure terminal is large enough for the prompt
	if err := makeSpace(8); err != nil {
		return "", ErrTerminalTooSmall
	}

	// State variables for this render cycle
	interrupted := false   // true if user aborted (Ctrl+C)
	receivedInput := false // turns true after user provides input event
	var inBuf []rune       // Input buffer to store user input
	cursorPos := 0         // Cursor position

	// Line constructors
	descriptionLine := si.style.InputDesc.Sprint(si.description)
	promptLine := si.style.InputPrefix.Sprint(si.prefix) + si.style.InputLabel.Sprint(si.label)
	placeholderLine := si.style.InputPlaceholder.Sprint(si.placeholder)
	helpLine := si.style.InputHelp.Sprint("Type to input . Enter to confirm")

	// Prompt Redraw Renderer
	redraw := func(input []rune, cursor int, validationMsg string, ok *bool) {
		stdOutput.Write([]byte(ansiHideCursor + ansiRestoreCursor + ansiClearLine + "\n\r"))
		if si.description != "" {
			stdOutput.Write([]byte(descriptionLine + "\n\r"))
		}
		stdOutput.Write([]byte(promptLine + ansiClearLine))
		if len(input) == 0 {
			stdOutput.Write([]byte(placeholderLine))
		}
		stdOutput.Write([]byte("\n\n\r" + ansiClearLine))
		if si.validator != nil && validationMsg != "" && receivedInput {
			if ok != nil && !*ok {
				stdOutput.Write([]byte(si.style.InputValidationFail.Sprint(validationMsg)))
			} else {
				stdOutput.Write([]byte(si.style.InputValidationPass.Sprint(validationMsg)))
			}
			stdOutput.Write([]byte(ansiClearLine))
		}
		stdOutput.Write([]byte("\n\n\r" + helpLine + ansiClearLine))
		stdOutput.Write([]byte(ansiRestoreCursor + "\n\r"))
		if si.description != "" {
			stdOutput.Write([]byte(descriptionLine + "\n\r"))
		}
		stdOutput.Write([]byte(promptLine))
		if len(input) != 0 && !si.noEcho {
			stdOutput.Write([]byte(si.style.InputText.Sprint(strings.Repeat("*", len(input))) + ansiClearLine))
			if cursor < len(input) {
				ansiCursorLeft(len(input) - cursor)
			}
		}
		stdOutput.Write([]byte(ansiShowCursor))
	}

	// Helper: Reset cursor after prompt render
	resetState := func() {
		stdOutput.Write([]byte(ansiRestoreCursor + ansiClearScreen + ansiReset + ansiShowCursor))
	}

	// Save state before prompt & defer reset
	stdOutput.Write([]byte(ansiHideCursor + ansiSaveCursor))
	defer resetState()

	// Prompt Initial Renderer
	redraw([]rune{}, 0, "", nil)

	// Intercept keyboard events & handle them
	err := keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		receivedInput = true
		switch key.Code {
		case keys.CtrlC:
			interrupted = true
			return true, nil
		case keys.Enter:
			if si.validator != nil {
				msg, ok := si.validator(string(inBuf))
				if !ok {
					redraw(inBuf, cursorPos, msg, &ok)
					return false, nil // block submit
				}
			}
			return true, nil
		case keys.Left:
			if cursorPos > 0 {
				cursorPos--
			}
		case keys.Right:
			if cursorPos < len(inBuf) {
				cursorPos++
			}
		case keys.Backspace:
			if cursorPos > 0 {
				inBuf = append(inBuf[:cursorPos-1], inBuf[cursorPos:]...)
				cursorPos--
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

		// live redraw with validator feedback
		if si.validator != nil {
			msg, ok := si.validator(string(inBuf))
			redraw(inBuf, cursorPos, msg, &ok)
		} else {
			redraw(inBuf, cursorPos, "", nil)
		}
		return false, nil
	})

	// Handle errors
	if err != nil {
		return "", err
	}
	if interrupted {
		return "", ErrInterrupted
	}

	// Return the input
	return strings.TrimRight(string(inBuf), "\r\n"), nil
}

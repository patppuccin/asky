package asky

import (
	"os"
	"strings"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
)

// Definition ----------------------------------------------
type SecureInput struct {
	theme           Theme
	promptSymbol    string
	promptText      string
	promptSeparator string
	helpText        string
	validator       func(string) (string, bool)
}

// Initialization ------------------------------------------
func NewSecureInput() *SecureInput {
	return &SecureInput{
		theme:           ThemeDefault,
		promptSymbol:    "[?] ",
		promptText:      "Enter secure input",
		promptSeparator: ": ",
		helpText:        "",
		validator:       nil,
	}
}

// Configuration -------------------------------------------
func (si *SecureInput) WithTheme(th Theme) *SecureInput        { si.theme = th; return si }
func (si *SecureInput) WithPromptSymbol(p string) *SecureInput { si.promptSymbol = p; return si }
func (si *SecureInput) WithPromptText(p string) *SecureInput   { si.promptText = p; return si }
func (si *SecureInput) WithPromptSeparator(sep string) *SecureInput {
	si.promptSeparator = sep
	return si
}
func (si *SecureInput) WithHelpText(txt string) *SecureInput { si.helpText = txt; return si }
func (si *SecureInput) WithValidator(fn func(string) (string, bool)) *SecureInput {
	si.validator = fn
	return si
}

// Presentation --------------------------------------------
func (si *SecureInput) Render() (string, error) {

	var inBuf []rune // Input buffer to store user input
	cursorPos := 0   // Cursor position
	saveCursor()     // Save cursor state before prompt

	// Help line construction
	helpLine := ""
	if si.helpText != "" {
		helpLine += si.theme.MutedStyle(si.helpText)
	}

	// Prompt line construction
	promptLine := si.theme.SecondaryStyle(si.promptSymbol)
	promptLine += si.theme.PrimaryStyle(si.promptText + si.promptSeparator)

	// Prompt Initial Renderer
	os.Stdout.WriteString("\n")
	if si.validator != nil {
		os.Stdout.WriteString("\n")
	}
	os.Stdout.WriteString(helpLine + "\n")
	os.Stdout.WriteString(promptLine)

	// Prompt Redraw Renderer
	redraw := func(input []rune, cursor int, validationMsg string, ok *bool) {
		restoreCursor()
		os.Stdout.WriteString("\n")
		if si.validator != nil && validationMsg != "" {
			if ok != nil && !*ok {
				os.Stdout.WriteString(si.theme.ErrorStyle(validationMsg))
			} else {
				os.Stdout.WriteString(si.theme.SuccessStyle(validationMsg))
			}
			clearLineTillEnd()
			os.Stdout.WriteString("\n\r")
		}

		os.Stdout.WriteString(helpLine + "\n\r")
		os.Stdout.WriteString(promptLine)
		os.Stdout.WriteString(strings.Repeat("*", len(input)))
		clearLineTillEnd()
		if cursor < len(input) {
			cursorMoveLeft(len(input) - cursor)
		}
	}

	err := keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		switch key.Code {
		case keys.CtrlC, keys.Escape:
			return true, ErrInterrupted
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

	restoreCursor()
	clearTillEnd()
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(inBuf), "\r\n"), nil
}

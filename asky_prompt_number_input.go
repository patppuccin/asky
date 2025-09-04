package asky

import (
	"os"
	"strings"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
)

// Definition ----------------------------------------------
type NumberInput struct {
	theme           Theme
	promptSymbol    string
	promptText      string
	promptSeparator string
	helpText        string
	defaultValue    string
	validator       func(string) (string, bool)
}

// Initialization ------------------------------------------
func NewNumberInput() *NumberInput {
	return &NumberInput{
		theme:           ThemeDefault,
		promptSymbol:    "[?] ",
		promptText:      "Enter text input",
		promptSeparator: ": ",
		helpText:        "",
		defaultValue:    "",
		validator:       nil,
	}
}

// Configuration -------------------------------------------
func (ni *NumberInput) WithTheme(th Theme) *NumberInput        { ni.theme = th; return ni }
func (ni *NumberInput) WithPromptSymbol(p string) *NumberInput { ni.promptSymbol = p; return ni }
func (ni *NumberInput) WithPromptText(p string) *NumberInput   { ni.promptText = p; return ni }
func (ni *NumberInput) WithPromptSeparator(sep string) *NumberInput {
	ni.promptSeparator = sep
	return ni
}
func (ni *NumberInput) WithHelpText(txt string) *NumberInput     { ni.helpText = txt; return ni }
func (ni *NumberInput) WithDefaultValue(val string) *NumberInput { ni.defaultValue = val; return ni }
func (ni *NumberInput) WithValidator(fn func(string) (string, bool)) *NumberInput {
	ni.validator = fn
	return ni
}

// Presentation --------------------------------------------
func (ni *NumberInput) Render() (string, error) {

	var inBuf []rune                      // Input buffer to store user input
	cursorPos := 0                        // Cursor position
	os.Stdout.WriteString(ansiSaveCursor) // Save cursor state before prompt

	// Help line construction
	helpLine := ""
	if ni.helpText != "" {
		helpLine += ni.helpText
	}

	// Prompt line construction
	promptLine := ni.theme.SecondaryStyle(ni.promptSymbol)
	promptLine += ni.theme.PrimaryStyle(ni.promptText)
	if ni.defaultValue != "" {
		promptLine += ni.theme.PrimaryStyle(" (" + ni.defaultValue + ")")
	}
	promptLine += ni.theme.PrimaryStyle(ni.promptSeparator)

	// Prompt Initial Renderer
	os.Stdout.WriteString("\n")
	if ni.validator != nil {
		os.Stdout.WriteString("\n")
	}
	os.Stdout.WriteString(helpLine + "\n")
	os.Stdout.WriteString(promptLine)

	// Prompt Redraw Renderer
	redraw := func(input []rune, cursor int, validationMsg string, ok *bool) {
		os.Stdout.WriteString(ansiRestoreCursor)
		os.Stdout.WriteString("\n")
		if ni.validator != nil && validationMsg != "" {
			if ok != nil && !*ok {
				os.Stdout.WriteString(ni.theme.ErrorStyle(validationMsg))
			} else {
				os.Stdout.WriteString(ni.theme.SuccessStyle(validationMsg))
			}
			os.Stdout.WriteString(ansiClearLineEnd)
			os.Stdout.WriteString("\n\r")
		}

		os.Stdout.WriteString(helpLine + "\n\r")
		os.Stdout.WriteString(promptLine)
		os.Stdout.WriteString(string(input))
		os.Stdout.WriteString(ansiClearLineEnd)
		if cursor < len(input) {
			cursorMoveLeft(len(input) - cursor)
		}
	}

	err := keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		switch key.Code {
		case keys.CtrlC, keys.Escape:
			return true, ErrInterrupted
		case keys.Enter:
			if ni.validator != nil {
				msg, ok := ni.validator(string(inBuf))
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
				r := key.Runes[0]
				if r >= '0' && r <= '9' { // allow only digits
					inBuf = append(inBuf[:cursorPos], append([]rune{r}, inBuf[cursorPos:]...)...)
					cursorPos++
				}
			}
		}

		// live redraw with validator feedback
		if ni.validator != nil {
			msg, ok := ni.validator(string(inBuf))
			redraw(inBuf, cursorPos, msg, &ok)
		} else {
			redraw(inBuf, cursorPos, "", nil)
		}
		return false, nil
	})

	os.Stdout.WriteString(ansiRestoreCursor + ansiClearScreenEnd + ansiReset + ansiShowCursor)
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(inBuf), "\r\n"), nil
}

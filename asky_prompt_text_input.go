package asky

import (
	"os"
	"strings"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
)

// Definition ----------------------------------------------
type TextInput struct {
	theme           Theme
	promptSymbol    string
	promptText      string
	promptSeparator string
	helpText        string
	defaultValue    string
	validator       func(string) (string, bool)
}

// Initialization ------------------------------------------
func NewTextInput() *TextInput {
	return &TextInput{
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
func (ti *TextInput) WithTheme(th Theme) *TextInput        { ti.theme = th; return ti }
func (ti *TextInput) WithPromptSymbol(p string) *TextInput { ti.promptSymbol = p; return ti }
func (ti *TextInput) WithPromptText(p string) *TextInput   { ti.promptText = p; return ti }
func (ti *TextInput) WithPromptSeparator(sep string) *TextInput {
	ti.promptSeparator = sep
	return ti
}
func (ti *TextInput) WithHelpText(txt string) *TextInput     { ti.helpText = txt; return ti }
func (ti *TextInput) WithDefaultValue(val string) *TextInput { ti.defaultValue = val; return ti }
func (ti *TextInput) WithValidator(fn func(string) (string, bool)) *TextInput {
	ti.validator = fn
	return ti
}

// Presentation --------------------------------------------
func (ti *TextInput) Render() (string, error) {

	var inBuf []rune // Input buffer to store user input
	cursorPos := 0   // Cursor position
	interrupted := false
	os.Stdout.WriteString(ansiSaveCursor) // Save cursor state before prompt

	// Help line construction
	helpLine := ""
	if ti.helpText != "" {
		helpLine += ti.theme.MutedStyle(ti.helpText)
	}

	// Prompt line construction
	promptLine := ti.theme.SecondaryStyle(ti.promptSymbol)
	promptLine += ti.theme.PrimaryStyle(ti.promptText)
	if ti.defaultValue != "" {
		promptLine += ti.theme.PrimaryStyle(" (" + ti.defaultValue + ")")
	}
	promptLine += ti.theme.PrimaryStyle(ti.promptSeparator)

	// Prompt Initial Renderer
	os.Stdout.WriteString("\n")
	if ti.validator != nil {
		os.Stdout.WriteString("\n")
	}
	os.Stdout.WriteString(helpLine + "\n")
	os.Stdout.WriteString(promptLine)

	// Prompt Redraw Renderer
	redraw := func(input []rune, cursor int, validationMsg string, ok *bool) {
		os.Stdout.WriteString(ansiRestoreCursor)
		os.Stdout.WriteString("\n")
		if ti.validator != nil && validationMsg != "" {
			if ok != nil && !*ok {
				os.Stdout.WriteString(ti.theme.ErrorStyle(validationMsg))
			} else {
				os.Stdout.WriteString(ti.theme.SuccessStyle(validationMsg))
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
		case keys.CtrlC:
			interrupted = true
			return true, nil
		case keys.Enter:
			if ti.validator != nil {
				msg, ok := ti.validator(string(inBuf))
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
		if ti.validator != nil {
			msg, ok := ti.validator(string(inBuf))
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
	if interrupted {
		return "", ErrInterrupted
	}
	return strings.TrimRight(string(inBuf), "\r\n"), nil
}

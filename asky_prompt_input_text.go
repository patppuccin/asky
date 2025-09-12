package asky

import (
	"strings"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
)

// --- Definition ------------------------------------------
type TextInput struct {
	theme        *Theme
	style        *Style
	prefix       string
	label        string
	description  string
	placeholder  string
	defaultValue string
	validator    func(string) (string, bool)
}

// --- Initiation ------------------------------------------
func NewTextInput() *TextInput {
	return &TextInput{
		prefix:    "[?] ",
		label:     "Enter text input",
		validator: nil,
	}
}

// --- Configuration ---------------------------------------
func (ti *TextInput) WithTheme(theme Theme) *TextInput       { ti.theme = &theme; return ti }
func (ti *TextInput) WithStyle(style Style) *TextInput       { ti.style = &style; return ti }
func (ti *TextInput) WithPrefix(p string) *TextInput         { ti.prefix = p; return ti }
func (ti *TextInput) WithLabel(p string) *TextInput          { ti.label = p; return ti }
func (ti *TextInput) WithDescription(txt string) *TextInput  { ti.description = txt; return ti }
func (ti *TextInput) WithPlaceholder(txt string) *TextInput  { ti.placeholder = txt; return ti }
func (ti *TextInput) WithDefaultValue(val string) *TextInput { ti.defaultValue = val; return ti }
func (ti *TextInput) WithValidator(fn func(string) (string, bool)) *TextInput {
	ti.validator = fn
	return ti
}

// --- Presentation ----------------------------------------
func (ti *TextInput) Render() (string, error) {
	// Setup theme and style (apply defaults if not set)
	if ti.theme == nil {
		ti.theme = &ThemeDefault
	}
	if ti.style == nil {
		ti.style = StyleDefault(ti.theme)
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
	descriptionLine := ti.style.InputDesc.Sprint(ti.description)
	promptLine := ti.style.InputPrefix.Sprint(ti.prefix) + ti.style.InputLabel.Sprint(ti.label)
	var placeholderLine string
	switch {
	case ti.placeholder != "" && ti.defaultValue != "":
		placeholderLine = ti.style.InputPlaceholder.Sprint(ti.placeholder + " (default: " + ti.defaultValue + ")")
	case ti.placeholder != "":
		placeholderLine = ti.style.InputPlaceholder.Sprint(ti.placeholder)
	case ti.defaultValue != "":
		placeholderLine = ti.style.InputPlaceholder.Sprint("default: " + ti.defaultValue)
	}
	helpLine := ti.style.InputHelp.Sprint("Type to input . Enter to confirm")

	// Prompt Redraw Renderer
	redraw := func(input []rune, cursor int, validationMsg string, ok *bool) {
		stdOutput.Write([]byte(ansiHideCursor + ansiRestoreCursor + ansiClearLine + "\n\r"))
		if ti.description != "" {
			stdOutput.Write([]byte(descriptionLine + "\n\r"))
		}
		stdOutput.Write([]byte(promptLine + ansiClearLine))
		if len(input) == 0 {
			stdOutput.Write([]byte(placeholderLine))
		}
		stdOutput.Write([]byte("\n\n\r" + ansiClearLine))
		if ti.validator != nil && validationMsg != "" && receivedInput {
			if ok != nil && !*ok {
				stdOutput.Write([]byte(ti.style.InputValidationFail.Sprint(validationMsg)))
			} else {
				stdOutput.Write([]byte(ti.style.InputValidationPass.Sprint(validationMsg)))
			}
			stdOutput.Write([]byte(ansiClearLine))
		}
		stdOutput.Write([]byte("\n\n\r" + helpLine + ansiClearLine))
		stdOutput.Write([]byte(ansiRestoreCursor + "\n\r"))
		if ti.description != "" {
			stdOutput.Write([]byte(descriptionLine + "\n\r"))
		}
		stdOutput.Write([]byte(promptLine))
		if len(input) != 0 {
			stdOutput.Write([]byte(ti.style.InputText.Sprint(string(input)) + ansiClearLine))
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

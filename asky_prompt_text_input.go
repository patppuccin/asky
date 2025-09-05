package asky

import (
	"os"
	"strings"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
)

// Definition ----------------------------------------------
type TextInput struct {
	theme        Theme
	prefix       string
	label        string
	description  string
	placeholder  string
	defaultValue string
	validator    func(string) (string, bool)
}

// Initialization ------------------------------------------
func NewTextInput() *TextInput {
	return &TextInput{
		theme:        ThemeDefault,
		prefix:       "[?] ",
		label:        "Enter text input",
		description:  "",
		placeholder:  "",
		defaultValue: "",
		validator:    nil,
	}
}

// Configuration -------------------------------------------
func (ti *TextInput) WithTheme(th Theme) *TextInput          { ti.theme = th; return ti }
func (ti *TextInput) WithPrefix(p string) *TextInput         { ti.prefix = p; return ti }
func (ti *TextInput) WithLabel(p string) *TextInput          { ti.label = p; return ti }
func (ti *TextInput) WithDescription(txt string) *TextInput  { ti.description = txt; return ti }
func (ti *TextInput) WithPlaceholder(txt string) *TextInput  { ti.placeholder = txt; return ti }
func (ti *TextInput) WithDefaultValue(val string) *TextInput { ti.defaultValue = val; return ti }
func (ti *TextInput) WithValidator(fn func(string) (string, bool)) *TextInput {
	ti.validator = fn
	return ti
}

// Presentation --------------------------------------------
func (ti *TextInput) Render() (string, error) {
	// Get the style preset
	preset := newPreset(ti.theme)

	// State variables for this render cycle
	interrupted := false   // true if user aborted (Ctrl+C)
	receivedInput := false // turns true after user provides input event
	var inBuf []rune       // Input buffer to store user input
	cursorPos := 0         // Cursor position

	// Line constructors
	descriptionLine := preset.accent.Sprint(ti.description)
	promptLine := preset.primary.Sprint(ti.prefix) + preset.secondary.Sprint(ti.label)
	var placeholderLine string
	switch {
	case ti.placeholder != "" && ti.defaultValue != "":
		placeholderLine = preset.muted.Sprint(ti.placeholder + " (default: " + ti.defaultValue + ")")
	case ti.placeholder != "":
		placeholderLine = preset.muted.Sprint(ti.placeholder)
	case ti.defaultValue != "":
		placeholderLine = preset.muted.Sprint("default: " + ti.defaultValue)
	}
	helpLine := preset.muted.Sprint("Type to input. Enter to confirm")

	// Prompt Redraw Renderer
	redraw := func(input []rune, cursor int, validationMsg string, ok *bool) {
		os.Stdout.WriteString(ansiRestoreCursor + ansiClearLineEnd + "\n\r")
		if ti.description != "" {
			os.Stdout.WriteString(descriptionLine + "\n\r")
		}
		os.Stdout.WriteString(promptLine + ansiClearLineEnd + " ")
		if len(input) == 0 {
			os.Stdout.WriteString(placeholderLine)
		} else {
			os.Stdout.WriteString(preset.neutral.Sprint(string(input)) + ansiClearLineEnd)
			if cursor < len(input) {
				cursorMoveLeft(len(input) - cursor)
			}
		}
		os.Stdout.WriteString("\n\n\r" + ansiClearLineEnd)
		if ti.validator != nil && validationMsg != "" && receivedInput {
			if ok != nil && !*ok {
				os.Stdout.WriteString(preset.err.Sprint(validationMsg))
			} else {
				os.Stdout.WriteString(preset.success.Sprint(validationMsg))
			}
			os.Stdout.WriteString(ansiClearLineEnd)
		}
		os.Stdout.WriteString("\n\n\r" + helpLine + ansiClearLineEnd)
	}

	// Helper: Reset cursor after prompt render
	resetState := func() {
		os.Stdout.WriteString(ansiRestoreCursor + ansiClearScreenEnd + ansiReset + ansiShowCursor)
	}

	// Save state before prompt & defer reset
	os.Stdout.WriteString(ansiHideCursor + ansiSaveCursor)
	defer resetState()

	// Prompt Initial Renderer
	redraw([]rune{}, 0, "", nil)

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

	os.Stdout.WriteString(ansiRestoreCursor + ansiClearScreenEnd + ansiReset + ansiShowCursor)
	if err != nil {
		return "", err
	}
	if interrupted {
		return "", ErrInterrupted
	}
	return strings.TrimRight(string(inBuf), "\r\n"), nil
}

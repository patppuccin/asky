package asky

import (
	"os"
	"strings"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
)

type SecureInput struct {
	promptSymbol string
	promptText   string
	helperText   string
	separator    string
	theme        Theme
	validator    func(string) (string, bool)
}

func NewSecureInput() *SecureInput {
	return &SecureInput{
		promptSymbol: "[?] ",
		separator:    ": ",
		theme:        ThemeDefault,
		validator:    nil,
	}
}

func (ti *SecureInput) WithPromptSymbol(p string) *SecureInput { ti.promptSymbol = p; return ti }
func (ti *SecureInput) WithPromptText(p string) *SecureInput   { ti.promptText = p; return ti }
func (ti *SecureInput) WithHelper(txt string) *SecureInput     { ti.helperText = txt; return ti }
func (ti *SecureInput) WithSeparator(sep string) *SecureInput  { ti.separator = sep; return ti }
func (ti *SecureInput) WithTheme(th Theme) *SecureInput        { ti.theme = th; return ti }
func (ti *SecureInput) WithValidator(fn func(string) (string, bool)) *SecureInput {
	ti.validator = fn
	return ti
}

func saveCursor()       { os.Stdout.Write([]byte("\033[s")) }
func restoreCursor()    { os.Stdout.Write([]byte("\033[u")) }
func clearLineTillEnd() { os.Stdout.Write([]byte("\033[K")) }
func clearTillEnd()     { os.Stdout.Write([]byte("\033[J")) }

func (ti *SecureInput) Render() (string, error) {
	saveCursor()
	helperLine := "\n"
	if ti.helperText != "" {
		helperLine = helperLine + ti.theme.MutedStyle(ti.helperText)
	}
	promptLine := ti.theme.SecondaryStyle(ti.promptSymbol) + ti.theme.PrimaryStyle(ti.promptText+ti.separator)

	// initial paint
	os.Stdout.WriteString(helperLine + "\n")
	os.Stdout.WriteString(promptLine)

	var inBuf []rune

	// factor out repaint
	redraw := func(input []rune, validationMsg string, ok *bool) {
		restoreCursor()
		os.Stdout.WriteString(helperLine)

		if ti.validator != nil && validationMsg != "" {
			spacer := " ("
			if helperLine == "\n" {
				spacer = "("
			}

			if ok != nil && !*ok {
				os.Stdout.WriteString(ti.theme.ErrorStyle(spacer + validationMsg + ")"))
			} else {
				os.Stdout.WriteString(ti.theme.SuccessStyle(spacer + validationMsg + ")"))
			}
		}

		clearLineTillEnd()
		os.Stdout.WriteString("\n\r" + promptLine)
		os.Stdout.WriteString(strings.Repeat("*", len(input)))
		clearLineTillEnd()
	}

	err := keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		switch key.Code {

		case keys.CtrlC, keys.Escape:
			return true, ErrInterrupted

		case keys.Enter:
			if ti.validator != nil {
				msg, ok := ti.validator(string(inBuf))
				if !ok {
					redraw(inBuf, msg, &ok)
					return false, nil // block submit
				}
			}
			return true, nil

		case keys.Backspace:
			if len(inBuf) > 0 {
				inBuf = inBuf[:len(inBuf)-1]
			}

		case keys.Space:
			inBuf = append(inBuf, ' ')

		default:
			if key.Code == keys.RuneKey && len(key.Runes) > 0 {
				inBuf = append(inBuf, key.Runes[0])
			}
		}

		// live redraw with validator feedback
		if ti.validator != nil {
			msg, ok := ti.validator(string(inBuf))
			redraw(inBuf, msg, &ok)
		} else {
			redraw(inBuf, "", nil)
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

// func (ti *SecureInput) Render() (string, error) {
// 	saveCursor()
// 	helperLine := "\n"
// 	if ti.helperText != "" {
// 		helperLine = helperLine + ti.theme.MutedStyle(ti.helperText)
// 	}
// 	promptLine := ti.theme.SecondaryStyle(ti.promptSymbol) + ti.theme.PrimaryStyle(ti.promptText+ti.separator)

// 	os.Stdout.WriteString(helperLine + "\n")
// 	os.Stdout.WriteString(promptLine)

// 	var inBuf []rune

// 	err := keyboard.Listen(func(key keys.Key) (stop bool, err error) {

// 		switch key.Code {
// 		case keys.Enter:
// 			if ti.validator != nil {
// 				msg, ok := ti.validator(string(inBuf))
// 				if !ok {
// 					// redraw error before blocking submit
// 					restoreCursor()
// 					os.Stdout.WriteString(helperLine)
// 					os.Stdout.WriteString(" " + ti.theme.ErrorStyle("("+msg+")"))
// 					clearLineTillEnd()
// 					os.Stdout.WriteString("\n\r" + promptLine)
// 					os.Stdout.WriteString(strings.Repeat("*", len(inBuf)))
// 					clearLineTillEnd()
// 					return false, nil // block submit
// 				}
// 			}
// 			return true, nil

// 		case keys.Backspace:
// 			if len(inBuf) > 0 {
// 				inBuf = inBuf[:len(inBuf)-1]
// 			}

// 		case keys.Space:
// 			inBuf = append(inBuf, ' ')

// 		default:
// 			if key.Code == keys.RuneKey && len(key.Runes) > 0 {
// 				inBuf = append(inBuf, key.Runes[0])
// 			}
// 		}

// 		// redraw prompt + stars, clear rest of line
// 		restoreCursor()
// 		os.Stdout.WriteString(helperLine)
// 		if ti.validator != nil {
// 			spacer := "("
// 			if len(helperLine) > 0 && ti.validator != nil {
// 				spacer = " ("
// 			}
// 			msg, ok := ti.validator(string(inBuf))
// 			if !ok {
// 				os.Stdout.WriteString(ti.theme.ErrorStyle(spacer + msg + ")"))
// 			} else if msg != "" {
// 				os.Stdout.WriteString(ti.theme.SuccessStyle(spacer + msg + ")"))
// 			}
// 		}
// 		clearLineTillEnd()
// 		os.Stdout.WriteString("\n\r" + promptLine)
// 		os.Stdout.WriteString(strings.Repeat("*", len(inBuf)))
// 		clearLineTillEnd()
// 		return false, nil
// 	})

// 	restoreCursor()
// 	clearTillEnd()
// 	if err != nil {
// 		return "", err
// 	}

// 	return strings.TrimRight(string(inBuf), "\r\n"), nil

// }

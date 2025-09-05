package asky

import (
	"os"
	"strings"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
)

type Confirm struct {
	theme         Theme
	prefix        string
	label         string
	description   string
	defaultAnswer bool
}

func NewConfirm() *Confirm {
	return &Confirm{
		theme:         ThemeDefault,
		prefix:        "[?] ",
		label:         "Are you sure?",
		description:   "",
		defaultAnswer: false,
	}
}

func (cf *Confirm) WithTheme(th Theme) *Confirm         { cf.theme = th; return cf }
func (cf *Confirm) WithPrefix(p string) *Confirm        { cf.prefix = p; return cf }
func (cf *Confirm) WithLabel(p string) *Confirm         { cf.label = p; return cf }
func (cf *Confirm) WithDescription(txt string) *Confirm { cf.description = txt; return cf }
func (cf *Confirm) WithDefaultAnswer(val bool) *Confirm { cf.defaultAnswer = val; return cf }

// --- Presentation --------------------------------------------
func (cf *Confirm) Render() (bool, error) {
	// Get the style preset
	preset := newPreset(cf.theme)

	// State variables for this render cycle
	interrupted := false // true if user aborted (Ctrl+C)
	confirm := true

	// Set default answer for state tracking
	if !cf.defaultAnswer {
		confirm = false
	}

	// Line constructors
	descriptionLine := preset.accent.Sprint(cf.description)
	promptLine := preset.primary.Sprint(cf.prefix) + preset.secondary.Sprint(cf.label)
	helpLine := preset.muted.Sprint("← or → to move. Enter to confirm")

	// Helper: Reset cursor state after prompt render
	resetState := func() {
		os.Stdout.WriteString(ansiRestoreCursor + ansiClearScreenEnd + ansiReset + ansiShowCursor)
	}

	// Helper: Redraw the prompt with the current state
	redraw := func() {
		os.Stdout.WriteString(ansiRestoreCursor)
		yesStyle := preset.primary.Sprint(" YES ")
		if confirm {
			yesStyle = preset.highlight.Sprint(" YES ")
		}
		noStyle := preset.primary.Sprint(" NO ")
		if !confirm {
			noStyle = preset.highlight.Sprint(" NO ")
		}

		os.Stdout.WriteString("\n")
		if cf.description != "" {
			os.Stdout.WriteString(descriptionLine + "\n")
		}
		os.Stdout.WriteString("\r" + promptLine + "\n")
		os.Stdout.WriteString("\n\r" + strings.Repeat(" ", len(cf.prefix)))
		os.Stdout.WriteString(yesStyle + "  " + noStyle + "\n\r\n")
		os.Stdout.WriteString(helpLine + "\n")
	}

	// Save cursor state before prompt & defer reset
	os.Stdout.WriteString(ansiSaveCursor + ansiHideCursor)
	defer resetState()

	// Initial render
	redraw()

	// Intercept keyboard events & handle them
	err := keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		switch key.Code {
		case keys.Enter:
			return true, nil
		case keys.CtrlC:
			return true, nil
		case keys.Left, keys.Right:
			confirm = !confirm
		}
		redraw()
		return false, nil
	})

	// Handle errors
	if err != nil {
		return false, err
	}

	// Handle interrupts
	if interrupted {
		return false, ErrInterrupted
	}

	return confirm, nil
}

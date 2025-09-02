package asky

import (
	"os"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
)

type SelectionOption struct {
	Value string
	Label string
}

// Definition ----------------------------------------------
type Select struct {
	theme                 Theme
	promptSymbol          string
	promptText            string
	promptSeparator       string
	helpText              string
	options               []SelectionOption
	defaultOption         []int
	checkedOptionSymbol   string
	uncheckedOptionSymbol string
	allowMultiple         bool
	maxOptionsVisible     int
	selectedOptions       map[int]bool // Track selected options
}

// Initialization ------------------------------------------
func NewSelect() *Select {
	return &Select{
		theme:                 ThemeDefault,
		promptSymbol:          "[?] ",
		promptText:            "Select an option",
		promptSeparator:       ": ",
		helpText:              "",
		defaultOption:         []int{},
		options:               []SelectionOption{},
		checkedOptionSymbol:   " + ",
		uncheckedOptionSymbol: "   ",
		maxOptionsVisible:     10,
		selectedOptions:       make(map[int]bool),
	}
}

func (ss *Select) WithTheme(th Theme) *Select        { ss.theme = th; return ss }
func (ss *Select) WithPromptSymbol(p string) *Select { ss.promptSymbol = p; return ss }
func (ss *Select) WithPromptText(p string) *Select   { ss.promptText = p; return ss }
func (ss *Select) WithPromptSeparator(sep string) *Select {
	ss.promptSeparator = sep
	return ss
}
func (ss *Select) WithHelpText(txt string) *Select { ss.helpText = txt; return ss }
func (ss *Select) WithDefaultOption(val []int) *Select {
	ss.defaultOption = val
	// Initialize selected options based on defaults
	for _, idx := range val {
		if idx >= 0 && idx < len(ss.options) {
			ss.selectedOptions[idx] = true
		}
	}
	return ss
}
func (ss *Select) WithOptions(opts []SelectionOption) *Select { ss.options = opts; return ss }
func (ss *Select) WithMultiSelection() *Select                { ss.allowMultiple = true; return ss }
func (ss *Select) WithCheckedOptionSymbol(sym string) *Select {
	ss.checkedOptionSymbol = sym
	return ss
}
func (ss *Select) WithUncheckedOptionSymbol(sym string) *Select {
	ss.uncheckedOptionSymbol = sym
	return ss
}
func (ss *Select) WithMaxOptionsVisible(n int) *Select {
	ss.maxOptionsVisible = n
	return ss
}

// Presentation --------------------------------------------
func (ss *Select) Render() ([]SelectionOption, error) {
	hideCursor()
	if len(ss.options) == 0 {
		restoreCursor()
		clearTillEnd()
		showCursor()
		return []SelectionOption{}, ErrNoOptions
	}

	// Initialize selected options with defaults
	if len(ss.selectedOptions) == 0 && len(ss.defaultOption) > 0 {
		for _, idx := range ss.defaultOption {
			if idx >= 0 && idx < len(ss.options) {
				ss.selectedOptions[idx] = true
			}
		}
	}

	selIdx := 0                                          // Selection index (defaults to the first option)
	startIdx := 0                                        // Start index of the visible options
	endIdx := min(len(ss.options), ss.maxOptionsVisible) // End index of the visible options
	saveCursor()                                         // Save cursor state before prompt

	// Help line construction
	helpLine := ""
	if ss.helpText != "" {
		helpLine += ss.theme.MutedStyle(ss.helpText)
	}

	// Prompt line construction
	promptLine := ss.theme.SecondaryStyle(ss.promptSymbol)
	promptLine += ss.theme.PrimaryStyle(ss.promptText + ss.promptSeparator + "\n")

	// Info line construction
	infoText := "↑/↓ or j/k to move"
	if ss.allowMultiple {
		infoText += " . space to toggle"
	} else {
		infoText += " . space to select"
	}
	infoText += " . enter to confirm"
	infoLine := ss.theme.MutedStyle(infoText)

	// Initial render - show first option as highlighted
	initialOptions := make([]string, endIdx-startIdx)
	for i := startIdx; i < endIdx; i++ {
		idx := i - startIdx
		if i == selIdx {
			if ss.selectedOptions[i] {
				initialOptions[idx] = ss.theme.SuccessStyle(ss.checkedOptionSymbol + "> " + ss.options[i].Label)
			} else {
				initialOptions[idx] = ss.theme.SecondaryStyle(ss.uncheckedOptionSymbol) + ss.theme.PrimaryStyle("> "+ss.options[i].Label)
			}
		} else {
			if ss.selectedOptions[i] {
				initialOptions[idx] = ss.theme.SuccessStyle(ss.checkedOptionSymbol + ss.options[i].Label)
			} else {
				initialOptions[idx] = ss.theme.SecondaryStyle(ss.uncheckedOptionSymbol) + ss.theme.AccentStyle(ss.options[i].Label)
			}
		}
	}

	// Prompt Initial Renderer
	os.Stdout.WriteString("\n")
	if ss.helpText != "" {
		os.Stdout.WriteString(helpLine + "\n")
	}
	os.Stdout.WriteString(promptLine)

	// Initial options render
	for i := startIdx; i < endIdx; i++ {
		if i == selIdx {
			if ss.selectedOptions[i] {
				os.Stdout.WriteString(ss.theme.SuccessStyle(ss.checkedOptionSymbol + "> " + ss.options[i].Label))
			} else {
				os.Stdout.WriteString(ss.theme.SecondaryStyle(ss.uncheckedOptionSymbol) + ss.theme.PrimaryStyle("> "+ss.options[i].Label))
			}
		} else {
			if ss.selectedOptions[i] {
				os.Stdout.WriteString(ss.theme.SuccessStyle(ss.checkedOptionSymbol + ss.options[i].Label))
			} else {
				os.Stdout.WriteString(ss.theme.SecondaryStyle(ss.uncheckedOptionSymbol) + ss.theme.AccentStyle(ss.options[i].Label))
			}
		}
		os.Stdout.WriteString("\n")
	}
	os.Stdout.WriteString("\n\r" + infoLine)

	// Prompt Redraw Renderer
	redraw := func(sel, start, end int) {
		restoreCursor()

		// Redraw help line
		if ss.helpText != "" {
			os.Stdout.WriteString("\n" + helpLine)
		} else {
			os.Stdout.WriteString("\n")
		}
		clearLineTillEnd()

		// Redraw prompt line
		os.Stdout.WriteString("\r\n" + promptLine)
		clearLineTillEnd()

		// Redraw options
		for i := start; i < end; i++ {
			os.Stdout.WriteString("\r")
			if i == sel {
				if ss.selectedOptions[i] {
					os.Stdout.WriteString(ss.theme.SuccessStyle(ss.checkedOptionSymbol + "> " + ss.options[i].Label))
				} else {
					os.Stdout.WriteString(ss.theme.SecondaryStyle(ss.uncheckedOptionSymbol) + ss.theme.PrimaryStyle("> "+ss.options[i].Label))
				}
			} else {
				if ss.selectedOptions[i] {
					os.Stdout.WriteString(ss.theme.SuccessStyle(ss.checkedOptionSymbol + ss.options[i].Label))
				} else {
					os.Stdout.WriteString(ss.theme.SecondaryStyle(ss.uncheckedOptionSymbol) + ss.theme.AccentStyle(ss.options[i].Label))
				}
			}
			clearLineTillEnd()
			os.Stdout.WriteString("\n")
		}

		// Clear any remaining lines
		for i := end - start; i < ss.maxOptionsVisible; i++ {
			clearLineTillEnd()
			os.Stdout.WriteString("\n")
		}

		os.Stdout.WriteString("\n\r" + infoLine)
		clearLineTillEnd()
	}

	err := keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		switch key.Code {
		case keys.CtrlC, keys.Escape:
			return true, ErrInterrupted
		case keys.Enter:
			return true, nil
		case keys.Up, keys.Left, keys.KeyCode('k'), keys.KeyCode('h'):
			// Navigate up
			if selIdx > 0 {
				selIdx--
				// Scroll up if needed
				if selIdx < startIdx {
					startIdx = selIdx
					endIdx = min(startIdx+ss.maxOptionsVisible, len(ss.options))
				}
			}
		case keys.Down, keys.Right, keys.KeyCode('j'), keys.KeyCode('l'):
			// Navigate down
			if selIdx < len(ss.options)-1 {
				selIdx++
				// Scroll down if needed
				if selIdx >= endIdx {
					endIdx = selIdx + 1
					startIdx = max(0, endIdx-ss.maxOptionsVisible)
				}
			}
		case keys.Space:
			// Toggle selection
			if ss.allowMultiple {
				ss.selectedOptions[selIdx] = !ss.selectedOptions[selIdx]
			} else {
				// Single select - clear all others and select current
				ss.selectedOptions = make(map[int]bool)
				ss.selectedOptions[selIdx] = true
			}
		}

		redraw(selIdx, startIdx, endIdx)
		return false, nil
	})

	if err != nil {
		restoreCursor()
		clearTillEnd()
		showCursor()
		return []SelectionOption{}, err
	}

	if ss.allowMultiple {
		// For multi-select, return all selected options
		var selections []SelectionOption
		for idx, selected := range ss.selectedOptions {
			if selected {
				selections = append(selections, ss.options[idx])
			}
		}
		restoreCursor()
		clearTillEnd()
		showCursor()
		return selections, nil
	}

	restoreCursor()
	clearTillEnd()
	showCursor()
	// Return currently highlighted option if nothing selected
	return []SelectionOption{ss.options[selIdx]}, nil
}

// package asky

// import (
// 	"os"
// 	"strings"

// 	"atomicgo.dev/keyboard"
// 	"atomicgo.dev/keyboard/keys"
// )

// type SelectionOption struct {
// 	Value string
// 	Label string
// }

// // Definition ----------------------------------------------
// type Select struct {
// 	theme                 Theme
// 	promptSymbol          string
// 	promptText            string
// 	promptSeparator       string
// 	helpText              string
// 	options               []SelectionOption
// 	defaultOption         []int
// 	checkedOptionSymbol   string
// 	uncheckedOptionSymbol string
// 	allowMultiple         bool
// 	maxOptionsVisible     int
// }

// // Initialization ------------------------------------------
// func NewSelect() *Select {
// 	return &Select{
// 		theme:                 ThemeDefault,
// 		promptSymbol:          "[?] ",
// 		promptText:            "Select an option",
// 		promptSeparator:       ": ",
// 		helpText:              "",
// 		defaultOption:         []int{},
// 		options:               []SelectionOption{},
// 		checkedOptionSymbol:   " + ",
// 		uncheckedOptionSymbol: "   ",
// 		maxOptionsVisible:     10,
// 	}
// }

// func (ss *Select) WithTheme(th Theme) *Select        { ss.theme = th; return ss }
// func (ss *Select) WithPromptSymbol(p string) *Select { ss.promptSymbol = p; return ss }
// func (ss *Select) WithPromptText(p string) *Select   { ss.promptText = p; return ss }
// func (ss *Select) WithPromptSeparator(sep string) *Select {
// 	ss.promptSeparator = sep
// 	return ss
// }
// func (ss *Select) WithHelpText(txt string) *Select            { ss.helpText = txt; return ss }
// func (ss *Select) WithDefaultOption(val []int) *Select        { ss.defaultOption = val; return ss }
// func (ss *Select) WithOptions(opts []SelectionOption) *Select { ss.options = opts; return ss }
// func (ss *Select) WithMultiSelection() *Select                { ss.allowMultiple = true; return ss }
// func (ss *Select) WithCheckedOptionSymbol(sym string) *Select {
// 	ss.checkedOptionSymbol = sym
// 	return ss
// }
// func (ss *Select) WithUncheckedOptionSymbol(sym string) *Select {
// 	ss.uncheckedOptionSymbol = sym
// 	return ss
// }
// func (ss *Select) WithMaxOptionsVisible(n int) *Select {
// 	ss.maxOptionsVisible = n
// 	return ss
// }

// // Presentation --------------------------------------------
// func (ss *Select) Render() (SelectionOption, error) {

// 	// var inBuf []rune // Input buffer to store user input
// 	selIdx := 0                                          // Selection index (defaults to the first option)
// 	startIdx := 0                                        // Start index of the visible options
// 	endIdx := min(len(ss.options), ss.maxOptionsVisible) // End index of the visible options
// 	saveCursor()                                         // Save cursor state before prompt

// 	// Help line construction
// 	helpLine := ""
// 	if ss.helpText != "" {
// 		helpLine += ss.theme.MutedStyle(ss.helpText)
// 	}

// 	// Prompt line construction
// 	promptLine := ss.theme.SecondaryStyle(ss.promptSymbol)
// 	promptLine += ss.theme.PrimaryStyle(ss.promptText + ss.promptSeparator + "\n")

// 	// info line construction
// 	infoLine := ss.theme.MutedStyle("↑/↓ or j/k to move . space to toggle . enter to confirm")

// 	visibleCount := min(len(ss.options), ss.maxOptionsVisible)
// 	options := make([]string, visibleCount)
// 	for i := range visibleCount {
// 		options[i] = ss.theme.SecondaryStyle(ss.uncheckedOptionSymbol) + ss.theme.AccentStyle(ss.options[i].Label)
// 	}

// 	// Prompt Initial Renderer
// 	os.Stdout.WriteString("\n")
// 	os.Stdout.WriteString(helpLine + "\n")
// 	os.Stdout.WriteString(promptLine)
// 	os.Stdout.WriteString(strings.Join(options, "\n"))
// 	os.Stdout.WriteString("\n\n" + infoLine)

// 	// Prompt Redraw Renderer
// 	redraw := func(sel, start, end int) {
// 		restoreCursor()
// 		os.Stdout.WriteString("\n" + helpLine)
// 		os.Stdout.WriteString("\n\r" + promptLine)
// 		clearLineTillEnd()
// 		for i := start; i < end; i++ {
// 			if i == sel {
// 				os.Stdout.WriteString(ss.theme.SuccessStyle(ss.options[i].Label))
// 			} else {
// 				os.Stdout.WriteString(ss.theme.SecondaryStyle(ss.uncheckedOptionSymbol) + ss.theme.AccentStyle(ss.options[i].Label))
// 			}
// 			os.Stdout.WriteString("\n")
// 		}
// 		clearLineTillEnd()

// 	}

// 	err := keyboard.Listen(func(key keys.Key) (stop bool, err error) {
// 		switch key.Code {
// 		case keys.CtrlC, keys.Escape:
// 			return true, ErrInterrupted
// 		case keys.Enter:
// 			return true, nil
// 		case keys.Up, keys.Left, keys.KeyCode('k'), keys.KeyCode('h'):
// 			// navigate up
// 		case keys.Down, keys.Right, keys.KeyCode('j'), keys.KeyCode('l'):
// 			// navigate down
// 		case keys.Space:
// 			// toggle selection
// 		}

// 		redraw(sel, start, end)
// 	})

// }

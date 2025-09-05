package asky

import (
	"os"
	"strconv"
	"strings"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
)

// Initialization ------------------------------------------
func NewSingleSelect() *SingleSelect {
	return &SingleSelect{
		theme:           ThemeDefault,
		prefix:          "[?] ",
		label:           "Select an option",
		description:     "",
		defaultChoice:   -1,
		choices:         []Choice{},
		cursorIndicator: " * ",
		selectionMarker: " + ",
		disabledMarker:  " x ",
		pageSize:        10,
	}
}

func (ss *SingleSelect) WithTheme(th Theme) *SingleSelect         { ss.theme = th; return ss }
func (ss *SingleSelect) WithPrefix(p string) *SingleSelect        { ss.prefix = p; return ss }
func (ss *SingleSelect) WithLabel(p string) *SingleSelect         { ss.label = p; return ss }
func (ss *SingleSelect) WithDescription(txt string) *SingleSelect { ss.description = txt; return ss }
func (ss *SingleSelect) WithDefaultChoice(idx int) *SingleSelect  { ss.defaultChoice = idx; return ss }
func (ss *SingleSelect) WithChoiceOptional() *SingleSelect        { ss.choiceOptional = true; return ss }
func (ss *SingleSelect) WithChoices(ch []Choice) *SingleSelect    { ss.choices = ch; return ss }
func (ss *SingleSelect) WithPageSize(n int) *SingleSelect         { ss.pageSize = n; return ss }
func (ss *SingleSelect) WithCursorIndicator(ind string) *SingleSelect {
	ss.cursorIndicator = ind
	return ss
}
func (ss *SingleSelect) WithSelectionMarker(mrk string) *SingleSelect {
	ss.selectionMarker = mrk
	return ss
}
func (ss *SingleSelect) WithDisabledMarker(mrk string) *SingleSelect {
	ss.disabledMarker = mrk
	return ss
}

// Presentation --------------------------------------------
func (ss *SingleSelect) Render() (Choice, error) {
	// Get the style preset
	preset := newPreset(ss.theme)

	// Sanity check for no choices
	if len(ss.choices) == 0 {
		return Choice{}, ErrNoOptions
	}

	// State variables for this render cycle
	interrupted := false                               // true if user aborted (Ctrl+C)
	searchQuery := ""                                  // current search text
	searchMode := false                                // whether search mode is active
	filteredChoices := ss.choices                      // visible choices after filtering
	pageSize := min(ss.pageSize, len(filteredChoices)) // items per page
	cursorIdx := 0                                     // index of the highlighted choice
	startIdx := 0                                      // index of first visible choice
	endIdx := min(len(filteredChoices), pageSize)      // index after last visible choice

	// Line constructors
	descriptionLine := preset.accent.Sprint(ss.description)
	promptLine := preset.primary.Sprint(ss.prefix) + preset.secondary.Sprint(ss.label)
	searchLine := preset.primary.Sprint("Search: ")
	helpLineNormalMode := preset.muted.Sprint("↑/↓ or j/k move . space select . enter confirm" + ansiClearLineEnd + "\n\rTAB search")
	helpLineSearchMode := preset.muted.Sprint("↑/↓ move . space select . enter confirm" + ansiClearLineEnd + "\n\rType to search . ESC/TAB nav")

	// Helper: Choice Renderer based on the state, selection & cursor
	renderChoice := func(c Choice, cur, sel bool) string {
		switch {
		case c.Disabled:
			return preset.muted.Sprint(ss.disabledMarker) + preset.disabled.Sprint(c.Label)
		case sel:
			return preset.success.Sprint(ss.selectionMarker + c.Label)
		case cur:
			return preset.primary.Sprint(ss.cursorIndicator + c.Label)
		default:
			return preset.neutral.Sprint(strings.Repeat(" ", len(ss.selectionMarker)) + c.Label)
		}
	}

	// Helper: Filter choices based on search query (for search mode)
	filterChoices := func(query string) []Choice {
		if query == "" {
			return ss.choices
		}

		var filtered []Choice
		query = strings.ToLower(query)

		for _, choice := range ss.choices {
			if strings.Contains(strings.ToLower(choice.Label), query) {
				filtered = append(filtered, choice)
			}
		}
		return filtered
	}

	// Helper: Reset cursor position after filtering
	resetCursorAfterFilter := func() {
		if len(filteredChoices) == 0 {
			cursorIdx = 0
			startIdx = 0
			endIdx = 0
			return
		}

		// Keep cursor in bounds
		if cursorIdx >= len(filteredChoices) {
			cursorIdx = len(filteredChoices) - 1
		}

		// Recalculate pagination
		if cursorIdx < startIdx {
			startIdx = cursorIdx
		}
		if cursorIdx >= startIdx+pageSize {
			startIdx = max(0, cursorIdx-pageSize+1)
		}
		endIdx = min(startIdx+pageSize, len(filteredChoices))
	}

	// Helper: Navigate choices up based on cursor position
	navigateUp := func() {
		if cursorIdx > 0 {
			cursorIdx--
			// Scroll up if needed
			if cursorIdx < startIdx {
				startIdx = cursorIdx
				endIdx = min(startIdx+pageSize, len(filteredChoices))
			}
		}
	}

	// Helper: Navigate choices down based on cursor position
	navigateDown := func() {
		if cursorIdx < len(filteredChoices)-1 {
			cursorIdx++
			// Scroll down if needed
			if cursorIdx >= endIdx {
				endIdx = cursorIdx + 1
				startIdx = max(0, endIdx-pageSize)
			}
		}
	}

	// Helper: Redraw the prompt with the current state
	redraw := func(cursor, start, end int) {
		os.Stdout.WriteString(ansiRestoreCursor)
		os.Stdout.WriteString("\n")
		if ss.description != "" {
			os.Stdout.WriteString(descriptionLine + "\n")
		}
		os.Stdout.WriteString("\r" + promptLine + "\n")

		// Search line with mode indicator
		sl := searchLine
		sl += preset.neutral.Sprint(searchQuery)
		if searchMode {
			sl += preset.muted.Sprint(" ◂ " + strconv.Itoa(len(filteredChoices)) + " hits") // Visual indicator for search mode
		}
		os.Stdout.WriteString("\r" + sl)
		os.Stdout.WriteString(ansiClearLineEnd)
		os.Stdout.WriteString("\n")

		// Redraw options
		for i := start; i < end; i++ {
			c := filteredChoices[i]
			cur := i == cursor
			sel := c.Value == ss.selectedChoice.Value
			os.Stdout.WriteString("\r" + renderChoice(c, cur, sel) + ansiClearLineEnd + "\n")
		}

		// Clear any remaining lines (move to start, clear contents, next line)
		for i := end - start; i < pageSize; i++ {
			os.Stdout.WriteString("\r" + ansiClearLineEnd + "\n")
		}

		// Show appropriate info line
		if searchMode {
			os.Stdout.WriteString(ansiClearLineEnd + "\n\r" + helpLineSearchMode + ansiClearLineEnd)
		} else {
			os.Stdout.WriteString(ansiClearLineEnd + "\n\r" + helpLineNormalMode + ansiClearLineEnd)
		}
	}

	// Helper: Reset cursor after prompt render
	resetState := func() {
		os.Stdout.WriteString(ansiRestoreCursor + ansiClearScreenEnd + ansiReset + ansiShowCursor)
	}

	// Save state before prompt & defer reset
	os.Stdout.WriteString(ansiHideCursor + ansiSaveCursor)
	defer resetState()

	// Initialize the selected choice with the default choice
	if ss.defaultChoice >= 0 && ss.defaultChoice < len(ss.choices) {
		ss.selectedChoice = ss.choices[ss.defaultChoice]
	}

	// Prompt Initial Renderer
	os.Stdout.WriteString("\n")
	if ss.description != "" {
		os.Stdout.WriteString(descriptionLine + "\n")
	}
	os.Stdout.WriteString("\r" + promptLine + "\n")
	os.Stdout.WriteString("\r" + searchLine + "\n")
	for i := startIdx; i < endIdx; i++ {
		c := filteredChoices[i]
		cur := i == cursorIdx
		sel := c.Value == ss.selectedChoice.Value
		os.Stdout.WriteString("\r" + renderChoice(c, cur, sel) + "\n")
	}
	if searchMode {
		os.Stdout.WriteString("\n\r" + helpLineSearchMode)
	} else {
		os.Stdout.WriteString("\n\r" + helpLineNormalMode)
	}

	// Intercept keyboard events & handle them
	err := keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		switch key.Code {
		case keys.Tab:
			searchMode = !searchMode
		case keys.CtrlC:
			interrupted = true
			return true, nil
		case keys.Escape:
			if searchMode {
				searchMode = false // In search mode, ESC exits search mode
			}
		case keys.Enter:
			if len(filteredChoices) == 0 {
				return ss.choiceOptional, nil
			}
			if ss.selectedChoice == (Choice{}) {
				return ss.choiceOptional, nil
			}
			return !ss.selectedChoice.Disabled, nil
		case keys.Up, keys.Left:
			navigateUp()
		case keys.Down, keys.Right:
			navigateDown()
		case keys.Space:
			if len(filteredChoices) > 0 && cursorIdx < len(filteredChoices) && !filteredChoices[cursorIdx].Disabled {
				ss.selectedChoice = filteredChoices[cursorIdx]
			}
		case keys.Backspace:
			if searchMode && len(searchQuery) > 0 {
				searchQuery = searchQuery[:len(searchQuery)-1]
				filteredChoices = filterChoices(searchQuery)
				resetCursorAfterFilter()
			}
		case keys.RuneKey:
			if len(key.Runes) == 0 {
				break
			}

			if searchMode {
				// In search mode, add characters to query
				searchQuery += string(key.Runes[0])
				filteredChoices = filterChoices(searchQuery)
				resetCursorAfterFilter()
			} else {
				// In nav mode, handle vi-style navigation
				switch key.Runes[0] {
				case 'j', 'l':
					navigateDown()
				case 'k', 'h':
					navigateUp()
				}
			}
		}

		redraw(cursorIdx, startIdx, endIdx)
		return false, nil
	})

	// Handle errors
	if err != nil {
		return Choice{}, err
	}

	// Handle interrupts
	if interrupted {
		return Choice{}, ErrInterrupted
	}

	// Restore state & return the selected choice
	return ss.selectedChoice, nil
}

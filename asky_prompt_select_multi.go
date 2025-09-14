package asky

import (
	"os"
	"strconv"
	"strings"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"github.com/mattn/go-runewidth"
)

// --- Definition ------------------------------------------
type multiSelect struct {
	theme            *Theme
	style            *Style
	prefix           string
	label            string
	description      string
	choices          []Choice
	defaultChoices   []int
	optional         bool
	minSelectedCount int
	maxSelectedCount int
	cursorIndicator  string
	selectionMarker  string
	disabledMarker   string
	pageSize         int
	selectedChoices  []Choice
}

// --- Initiation ------------------------------------------
func NewMultiSelect() *multiSelect {
	return &multiSelect{
		prefix:          "[?] ",
		label:           "Select an option",
		defaultChoices:  []int{},
		choices:         []Choice{},
		cursorIndicator: " >",
		selectionMarker: "+ ",
		disabledMarker:  "x ",
		pageSize:        7,
	}
}

// --- Configuration ---------------------------------------
func (ms *multiSelect) WithTheme(theme Theme) *multiSelect      { ms.theme = &theme; return ms }
func (ms *multiSelect) WithStyle(style Style) *multiSelect      { ms.style = &style; return ms }
func (ms *multiSelect) WithPrefix(p string) *multiSelect        { ms.prefix = p; return ms }
func (ms *multiSelect) WithLabel(p string) *multiSelect         { ms.label = p; return ms }
func (ms *multiSelect) WithDescription(txt string) *multiSelect { ms.description = txt; return ms }
func (ms *multiSelect) WithDefaultChoice(idxs []int) *multiSelect {
	if len(idxs) > 0 {
		ms.defaultChoices = idxs
	}
	return ms
}
func (ms *multiSelect) Optional() *multiSelect { ms.optional = true; return ms }
func (ms *multiSelect) WithMinSelectedCount(n int) *multiSelect {
	ms.minSelectedCount = max(n, 0)
	return ms
}
func (ms *multiSelect) WithMaxSelectedCount(n int) *multiSelect {
	ms.maxSelectedCount = max(n, 0)
	return ms
}
func (ms *multiSelect) WithChoices(ch []Choice) *multiSelect { ms.choices = ch; return ms }
func (ms *multiSelect) WithPageSize(n int) *multiSelect      { ms.pageSize = n; return ms }
func (ms *multiSelect) WithCursorIndicator(ind string) *multiSelect {
	ms.cursorIndicator = ind
	return ms
}
func (ms *multiSelect) WithSelectionMarker(mrk string) *multiSelect {
	ms.selectionMarker = mrk
	return ms
}
func (ms *multiSelect) WithDisabledMarker(mrk string) *multiSelect {
	ms.disabledMarker = mrk
	return ms
}

// --- Presentation ----------------------------------------
func (ms *multiSelect) Render() ([]Choice, error) {
	// State variables
	interrupted := false                               // true if user aborted (Ctrl+C)
	searchQuery := ""                                  // current search text
	searchMode := false                                // whether search mode is active
	filteredChoices := ms.choices                      // visible choices after filtering
	pageSize := min(ms.pageSize, len(filteredChoices)) // items per page
	cursorIdx := 0                                     // index of the highlighted choice
	startIdx := 0                                      // index of first visible choice
	endIdx := min(len(filteredChoices), pageSize)      // index after last visible choice
	valMessage := ""                                   // validation message to display
	minRequired := max(ms.minSelectedCount, 0)
	if !ms.optional && minRequired == 0 {
		minRequired = 1
	}
	maxAllowed := max(ms.maxSelectedCount, 0)
	if maxAllowed == 0 {
		maxAllowed = len(ms.choices)
	}

	// Ensure terminal is large enough for the prompt
	if err := makeSpace(9 + pageSize); err != nil {
		return []Choice{}, ErrTerminalTooSmall
	}

	// Sanity check for no choices supplied
	if len(ms.choices) == 0 {
		return []Choice{}, ErrNoSelectionChoices
	}

	// Sanity check to check if min is not greater than max
	if minRequired > maxAllowed {
		return []Choice{}, ErrInvalidSelectionCount
	}

	// Setup theme and style (apply defaults if not set)
	if ms.theme == nil {
		ms.theme = &ThemeDefault
	}
	if ms.style == nil {
		ms.style = StyleDefault(ms.theme)
	}

	// Line constructors
	descriptionLine := ms.style.SelectionDesc.Sprint(ms.description)
	promptLine := ms.style.SelectionPrefix.Sprint(ms.prefix) + ms.style.SelectionLabel.Sprint(ms.label)
	searchLine := ms.style.SelectionSearchLabel.Sprint("Search: ")
	helpLineNormalMode := ms.style.SelectionHelp.Sprint("↑/↓ move . space select . enter confirm" + ansiClearLine + "\n\rtab to search" + ansiClearLine)
	helpLineSearchMode := ms.style.SelectionHelp.Sprint("↑/↓ move . space select . enter confirm" + ansiClearLine + "\n\rtype to search (ESC/TAB nav)" + ansiClearLine)

	// Check if a choice is selected
	isSelected := func(choice Choice) bool {
		for _, selected := range ms.selectedChoices {
			if selected.Value == choice.Value {
				return true
			}
		}
		return false
	}

	// Render choice based on the state, selection & cursor
	renderChoice := func(c Choice, cur, sel bool) string {
		cursorSpacer := strings.Repeat(" ", runewidth.StringWidth(ms.cursorIndicator))
		selectionSpacer := strings.Repeat(" ", runewidth.StringWidth(ms.selectionMarker))
		switch {
		case c.Disabled && cur:
			return ms.style.SelectionDisabledItemMarker.Sprint(ms.cursorIndicator+ms.disabledMarker) +
				ms.style.SelectionDisabledItemLabel.Sprint(c.Label)
		case sel && cur:
			return ms.style.SelectionSelectedItemMarker.Sprint(ms.cursorIndicator+ms.selectionMarker) +
				ms.style.SelectionSelectedItemLabel.Sprint(c.Label)
		case c.Disabled:
			return cursorSpacer +
				ms.style.SelectionDisabledItemMarker.Sprint(ms.disabledMarker) +
				ms.style.SelectionDisabledItemLabel.Sprint(c.Label)
		case sel:
			return cursorSpacer +
				ms.style.SelectionSelectedItemMarker.Sprint(ms.selectionMarker) +
				ms.style.SelectionSelectedItemLabel.Sprint(c.Label)
		case cur:
			return ms.style.SelectionCurrentItemMarker.Sprint(ms.cursorIndicator) + selectionSpacer +
				ms.style.SelectionCurrentItemLabel.Sprint(c.Label)
		default:
			return cursorSpacer + selectionSpacer +
				ms.style.SelectionListItemLabel.Sprint(c.Label)
		}
	}

	// Filter choices based on the search query (for search mode)
	filterChoices := func(query string) []Choice {
		if query == "" {
			return ms.choices
		}

		var filtered []Choice
		query = strings.ToLower(query)

		for _, choice := range ms.choices {
			if strings.Contains(strings.ToLower(choice.Label), query) {
				filtered = append(filtered, choice)
			}
		}
		return filtered
	}

	// Reset cursor position after filtering
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

	// Navigate choices up based on the cursor position
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

	// Navigate choices down based on the cursor position
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

	// Toggle selection state of a choice
	toggleSelection := func(choice Choice) {
		// Remove if already selected
		for i, selected := range ms.selectedChoices {
			if selected.Value == choice.Value {
				ms.selectedChoices = append(ms.selectedChoices[:i], ms.selectedChoices[i+1:]...)
				return
			}
		}
		// Add if not selected and within limits
		if ms.maxSelectedCount == 0 || len(ms.selectedChoices) < ms.maxSelectedCount {
			ms.selectedChoices = append(ms.selectedChoices, choice)
		}
	}

	// Prompt Redraw Renderer
	redraw := func(cursor, start, end int) {
		stdOutput.Write([]byte(ansiRestoreCursor + "\n"))
		if ms.description != "" {
			stdOutput.Write([]byte(descriptionLine + "\n"))
		}
		stdOutput.Write([]byte("\r" + promptLine + "\n"))

		// Search line with mode indicator
		sl := searchLine
		sl += ms.style.SelectionSearchHint.Sprint(searchQuery)
		if searchMode {
			sl += ms.style.SelectionSearchHint.Sprint(" ◂ " + strconv.Itoa(len(filteredChoices)) + " hits")
		}
		// Show selection count
		selectedCount := len(ms.selectedChoices)
		sl += ms.style.SelectionSearchHint.Sprint(" [" + strconv.Itoa(selectedCount) + " selected]")

		os.Stdout.WriteString("\r" + sl)
		os.Stdout.WriteString(ansiClearLine)
		os.Stdout.WriteString("\n")

		// Redraw options
		for i := start; i < end; i++ {
			c := filteredChoices[i]
			cur := i == cursor
			sel := isSelected(c)
			stdOutput.Write([]byte("\r" + renderChoice(c, cur, sel) + ansiClearLine + "\n"))
		}

		// Clear any remaining lines (move to start, clear contents, next line)
		for i := end - start; i < pageSize; i++ {
			stdOutput.Write([]byte("\r" + ansiClearLine + "\n"))
		}

		// Show validation message
		stdOutput.Write([]byte("\n\r" + ms.style.SelectionValidationFail.Sprint(valMessage) + ansiClearLine + "\n\r"))

		// Show appropriate info line
		helpLine := helpLineNormalMode
		if searchMode {
			helpLine = helpLineSearchMode
		}
		stdOutput.Write([]byte(helpLine))
	}

	// Reset cursor after prompt render
	resetState := func() { stdOutput.Write([]byte(ansiRestoreCursor + ansiClearScreen + ansiReset + ansiShowCursor)) }

	// Save state before prompt & defer reset
	stdOutput.Write([]byte(ansiHideCursor + ansiSaveCursor))
	defer resetState()

	// Initialize the selected choices with the default choices
	if len(ms.defaultChoices) > 0 {
		ms.selectedChoices = []Choice{}
		for _, defaultIdx := range ms.defaultChoices {
			if defaultIdx >= 0 && defaultIdx < len(ms.choices) {
				ms.selectedChoices = append(ms.selectedChoices, ms.choices[defaultIdx])
			}
		}
	}

	// Prompt Initial Render
	redraw(cursorIdx, startIdx, endIdx)

	// Intercept keyboard events & handle them
	err := keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		switch key.Code {
		case keys.CtrlC:
			interrupted = true
			return true, nil
		case keys.Up, keys.Left:
			navigateUp()
		case keys.Down, keys.Right:
			navigateDown()
		case keys.Tab:
			searchMode = !searchMode
		case keys.Escape:
			if searchMode {
				searchMode = false // In search mode, ESC exits search mode
			}
		case keys.Enter:
			if len(ms.selectedChoices) < minRequired {
				valMessage = "At least " + strconv.Itoa(minRequired) + " choices must be selected"
			} else {
				return true, nil
			}
		case keys.Space:
			if len(filteredChoices) == 0 {
				valMessage = "No choices available"
				break
			}
			currentChoice := filteredChoices[cursorIdx]
			if currentChoice.Disabled {
				valMessage = "Cannot select a disabled choice"
				break
			}
			alreadySelected := isSelected(currentChoice)
			if !alreadySelected && len(ms.selectedChoices) >= maxAllowed {
				valMessage = "Cannot select more than " + strconv.Itoa(maxAllowed) + " choices"
				break
			}
			toggleSelection(currentChoice)
			if len(ms.selectedChoices) < minRequired {
				valMessage = "At least " + strconv.Itoa(minRequired) + " choices must be selected"
			} else {
				valMessage = ""
			}
		case keys.Backspace:
			if searchMode && len(searchQuery) > 0 {
				searchQuery = searchQuery[:len(searchQuery)-1]
				filteredChoices = filterChoices(searchQuery)
				resetCursorAfterFilter()
			}
		case keys.RuneKey:
			if len(key.Runes) == 0 { // No rune key pressed
				break
			}
			keyPressed := string(key.Runes[0])
			if searchMode { // In search mode, add characters to query
				searchQuery += keyPressed
				filteredChoices = filterChoices(searchQuery)
				resetCursorAfterFilter()
			} else { // In nav mode, handle vi-style navigation
				switch keyPressed {
				case "j", "l":
					navigateDown()
				case "k", "h":
					navigateUp()
				}
			}
		}

		redraw(cursorIdx, startIdx, endIdx)
		return false, nil
	})

	// Handle errors
	if err != nil {
		return []Choice{}, err
	}

	// Handle interrupts
	if interrupted {
		return []Choice{}, ErrInterrupted
	}

	// Restore state & return the selected choices
	return ms.selectedChoices, nil
}

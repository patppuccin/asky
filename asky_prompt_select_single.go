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
type Choice struct {
	Value    string
	Label    string
	Disabled bool
}

type singleSelect struct {
	theme           *Theme
	style           *Style
	prefix          string
	label           string
	description     string
	choices         []Choice
	defaultChoice   int
	optional        bool
	cursorIndicator string
	selectionMarker string
	disabledMarker  string
	pageSize        int
	selectedChoice  Choice
}

// --- Initiation ------------------------------------------
func NewSingleSelect() *singleSelect {
	return &singleSelect{
		prefix:          "[?] ",
		label:           "Select an option",
		defaultChoice:   -1,
		choices:         []Choice{},
		cursorIndicator: " >",
		selectionMarker: "+ ",
		disabledMarker:  "x ",
		pageSize:        10,
	}
}

// --- Configuration ---------------------------------------
func (ss *singleSelect) WithTheme(theme Theme) *singleSelect      { ss.theme = &theme; return ss }
func (ss *singleSelect) WithStyle(style Style) *singleSelect      { ss.style = &style; return ss }
func (ss *singleSelect) WithPrefix(p string) *singleSelect        { ss.prefix = p; return ss }
func (ss *singleSelect) WithLabel(p string) *singleSelect         { ss.label = p; return ss }
func (ss *singleSelect) WithDescription(txt string) *singleSelect { ss.description = txt; return ss }
func (ss *singleSelect) WithDefaultChoice(idx int) *singleSelect {
	ss.defaultChoice = max(0, idx)
	return ss
}
func (ss *singleSelect) Optional() *singleSelect               { ss.optional = true; return ss }
func (ss *singleSelect) WithChoices(ch []Choice) *singleSelect { ss.choices = ch; return ss }
func (ss *singleSelect) WithPageSize(n int) *singleSelect      { ss.pageSize = n; return ss }
func (ss *singleSelect) WithCursorIndicator(ind string) *singleSelect {
	ss.cursorIndicator = ind
	return ss
}
func (ss *singleSelect) WithSelectionMarker(mrk string) *singleSelect {
	ss.selectionMarker = mrk
	return ss
}
func (ss *singleSelect) WithDisabledMarker(mrk string) *singleSelect {
	ss.disabledMarker = mrk
	return ss
}

// --- Presentation ----------------------------------------
func (ss *singleSelect) Render() (Choice, error) {
	// State variables
	interrupted := false                               // true if user aborted (Ctrl+C)
	searchQuery := ""                                  // current search text
	searchMode := false                                // whether search mode is active
	filteredChoices := ss.choices                      // visible choices after filtering
	pageSize := min(ss.pageSize, len(filteredChoices)) // items per page
	cursorIdx := 0                                     // index of the highlighted choice
	startIdx := 0                                      // index of first visible choice
	endIdx := min(len(filteredChoices), pageSize)      // index after last visible choice
	valMessage := ""                                   // validation message to display

	// Ensure terminal is large enough for the prompt
	if err := makeSpace(9 + pageSize); err != nil {
		return Choice{}, ErrTerminalTooSmall
	}

	// Sanity check for no choices
	if len(ss.choices) == 0 {
		return Choice{}, ErrNoSelectionChoices
	}

	// Setup theme and style (apply defaults if not set)
	if ss.theme == nil {
		ss.theme = &ThemeDefault
	}
	if ss.style == nil {
		ss.style = StyleDefault(ss.theme)
	}

	// Line constructors
	descriptionLine := ss.style.SelectionDesc.Sprint(ss.description)
	promptLine := ss.style.SelectionPrefix.Sprint(ss.prefix) + ss.style.SelectionLabel.Sprint(ss.label)
	searchLine := ss.style.SelectionSearchLabel.Sprint("Search: ")
	helpLineNormalMode := ss.style.SelectionHelp.Sprint("↑/↓ move . space select . enter confirm" + ansiClearLine + "\n\rtab to search" + ansiClearLine)
	helpLineSearchMode := ss.style.SelectionHelp.Sprint("↑/↓ move . space select . enter confirm" + ansiClearLine + "\n\rtype to search (ESC/TAB nav)" + ansiClearLine)

	// Render choice based on the state, selection & cursor
	renderChoice := func(c Choice, cur, sel bool) string {
		cursorSpacer := strings.Repeat(" ", runewidth.StringWidth(ss.cursorIndicator))
		selectionSpacer := strings.Repeat(" ", runewidth.StringWidth(ss.selectionMarker))
		switch {
		case c.Disabled && cur:
			return ss.style.SelectionDisabledItemMarker.Sprint(ss.cursorIndicator+ss.disabledMarker) +
				ss.style.SelectionDisabledItemLabel.Sprint(c.Label)
		case sel && cur:
			return ss.style.SelectionSelectedItemMarker.Sprint(ss.cursorIndicator+ss.selectionMarker) +
				ss.style.SelectionSelectedItemLabel.Sprint(c.Label)
		case c.Disabled:
			return cursorSpacer +
				ss.style.SelectionDisabledItemMarker.Sprint(ss.disabledMarker) +
				ss.style.SelectionDisabledItemLabel.Sprint(c.Label)
		case sel:
			return cursorSpacer +
				ss.style.SelectionSelectedItemMarker.Sprint(ss.selectionMarker) +
				ss.style.SelectionSelectedItemLabel.Sprint(c.Label)
		case cur:
			return ss.style.SelectionCurrentItemMarker.Sprint(ss.cursorIndicator) + selectionSpacer +
				ss.style.SelectionCurrentItemLabel.Sprint(c.Label)
		default:
			return cursorSpacer + selectionSpacer +
				ss.style.SelectionListItemLabel.Sprint(c.Label)
		}
	}

	// Filter choices based on the search query (for search mode)
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

	// Prompt Redraw Renderer
	redraw := func(cursor, start, end int) {
		stdOutput.Write([]byte(ansiRestoreCursor + "\n"))
		if ss.description != "" {
			stdOutput.Write([]byte(descriptionLine + "\n"))
		}
		stdOutput.Write([]byte("\r" + promptLine + "\n"))

		// Search line with mode indicator
		sl := searchLine
		sl += ss.style.SelectionSearchHint.Sprint(searchQuery)
		if searchMode {
			sl += ss.style.SelectionSearchHint.Sprint(" ◂ " + strconv.Itoa(len(filteredChoices)) + " hits")
		}
		// Show selection count
		if ss.selectedChoice != (Choice{}) {
			sl += ss.style.SelectionSearchHint.Sprint(" [1 selected]")
		} else {
			sl += ss.style.SelectionSearchHint.Sprint(" [0 selected]")
		}

		os.Stdout.WriteString("\r" + sl)
		os.Stdout.WriteString(ansiClearLine)
		os.Stdout.WriteString("\n")

		// Redraw options
		for i := start; i < end; i++ {
			c := filteredChoices[i]
			cur := i == cursor
			sel := c.Value == ss.selectedChoice.Value
			stdOutput.Write([]byte("\r" + renderChoice(c, cur, sel) + ansiClearLine + "\n"))
		}

		// Clear any remaining lines (move to start, clear contents, next line)
		for i := end - start; i < pageSize; i++ {
			stdOutput.Write([]byte("\r" + ansiClearLine + "\n"))
		}

		// Show validation message
		stdOutput.Write([]byte("\n\r" + ss.style.SelectionValidationFail.Sprint(valMessage) + ansiClearLine + "\n\r"))

		// Show appropriate info line
		helpLine := helpLineNormalMode
		if searchMode {
			helpLine = helpLineSearchMode
		}
		stdOutput.Write([]byte(helpLine))
	}

	// Reset cursor after prompt render
	resetState := func() {
		stdOutput.Write([]byte(ansiRestoreCursor + ansiClearScreen + ansiReset + ansiShowCursor))
	}

	// Save state before prompt & defer reset
	stdOutput.Write([]byte(ansiHideCursor + ansiSaveCursor))
	defer resetState()

	// Initialize the selected choice with the default choice
	if ss.defaultChoice >= 0 && ss.defaultChoice < len(ss.choices) {
		ss.selectedChoice = ss.choices[ss.defaultChoice]
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
			if len(filteredChoices) == 0 || ss.selectedChoice == (Choice{}) {
				if ss.optional {
					return true, nil
				}
				valMessage = "No selection made (required)"
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
			if ss.selectedChoice.Value == currentChoice.Value {
				ss.selectedChoice = Choice{}
			} else {
				ss.selectedChoice = currentChoice
			}
			valMessage = ""
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
		return Choice{}, err
	}

	// Handle interrupts
	if interrupted {
		return Choice{}, ErrInterrupted
	}

	// Restore state & return the selected choice
	return ss.selectedChoice, nil
}

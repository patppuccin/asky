package asky

import (
	"os"
	"strconv"
	"strings"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
)

// Initialization ------------------------------------------
func NewMultiSelect() *MultiSelect {
	return &MultiSelect{
		theme:              ThemeDefault,
		prefix:             "[?] ",
		label:              "Select one or more options",
		separator:          ": ",
		help:               "",
		choices:            []Choice{},
		defaultChoices:     nil,
		minChoicesRequired: 0,
		maxChoicesAllowed:  0,     // 0 = no upper cap
		selectionMarker:    " + ", // shown when selected
		pageSize:           10,
		selectedChoices:    []int{}, // indices into ss.choices
	}
}

func (ms *MultiSelect) WithTheme(th Theme) *MultiSelect       { ms.theme = th; return ms }
func (ms *MultiSelect) WithPrefix(p string) *MultiSelect      { ms.prefix = p; return ms }
func (ms *MultiSelect) WithLabel(p string) *MultiSelect       { ms.label = p; return ms }
func (ms *MultiSelect) WithSeparator(sep string) *MultiSelect { ms.separator = sep; return ms }
func (ms *MultiSelect) WithHelp(txt string) *MultiSelect      { ms.help = txt; return ms }
func (ms *MultiSelect) WithChoices(ch []Choice) *MultiSelect  { ms.choices = ch; return ms }
func (ms *MultiSelect) WithDefaultChoices(indices []int) *MultiSelect {
	ms.defaultChoices = indices
	return ms
}
func (ms *MultiSelect) WithMinChoicesRequired(n int) *MultiSelect {
	ms.minChoicesRequired = n
	return ms
}
func (ms *MultiSelect) WithMaxChoicesAllowed(n int) *MultiSelect { ms.maxChoicesAllowed = n; return ms }
func (ms *MultiSelect) WithSelectionMarker(mrk string) *MultiSelect {
	ms.selectionMarker = mrk
	return ms
}
func (ms *MultiSelect) WithPageSize(n int) *MultiSelect { ms.pageSize = n; return ms }

// Presentation --------------------------------------------
// Render shows a searchable, paginated multi-select.
// Controls: ↑/↓ or j/k move • space toggle • a select all (visible) • n clear all • TAB search • enter confirm • ESC exits search.
func (ms *MultiSelect) Render() ([]Choice, error) {
	// Guard
	if len(ms.choices) == 0 {
		return nil, ErrNoOptions
	}

	// Selection state (use a map for O(1) checks)
	selected := make(map[int]struct{}, len(ms.choices))
	if len(ms.defaultChoices) > 0 {
		for _, idx := range ms.defaultChoices {
			if idx >= 0 && idx < len(ms.choices) && !ms.choices[idx].Disabled {
				selected[idx] = struct{}{}
			}
		}
	}

	searchQuery := ""
	searchMode := false
	filteredChoices := ms.choices

	pageSize := min(ms.pageSize, len(filteredChoices))
	cursorIdx := 0
	startIdx := 0
	endIdx := min(len(filteredChoices), pageSize)

	// Hide & save cursor before prompt
	hideCursor()
	saveCursor()

	// Lines
	helpLine := ms.theme.MutedStyle(ms.help)
	promptLine := ms.theme.SecondaryStyle(ms.prefix) + ms.theme.PrimaryStyle(ms.label+ms.separator)
	searchPrefix := ms.theme.MutedStyle("Search: ")
	infoLineNormal := ms.theme.MutedStyle("↑/↓ or j/k move . space toggle . a all . n none . enter confirm . TAB search")
	infoLineSearch := ms.theme.MutedStyle("↑/↓ move . space toggle . a all . n none . type search . ESC/TAB nav")

	// Helpers ----------------------------------------------

	// isSelected checks if absolute choice index is selected
	isSelected := func(absIdx int) bool {
		_, ok := selected[absIdx]
		return ok
	}

	// countSelected counts current selections
	countSelected := func() int { return len(selected) }

	// canAddMore enforces max cap (0 = no cap)
	canAddMore := func() bool { return ms.maxChoicesAllowed == 0 || countSelected() < ms.maxChoicesAllowed }

	// renderChoice builds a line with selection marker and color by state.
	// We don't show a cursor icon; we rely on color emphasis for cursor.
	renderChoice := func(absIdx int, cur bool) string {
		c := ms.choices[absIdx]
		var line string
		if isSelected(absIdx) {
			line = ms.selectionMarker + c.Label
		} else {
			line = strings.Repeat(" ", len(ms.selectionMarker)) + c.Label
		}

		switch {
		case c.Disabled:
			return ms.theme.MutedStyle(line)
		case cur && isSelected(absIdx):
			return ms.theme.SuccessStyle(line) // selected gets success
		case cur:
			return ms.theme.SecondaryStyle(line) // cursor item emphasized
		case isSelected(absIdx):
			return ms.theme.SuccessStyle(line)
		default:
			return ms.theme.AccentStyle(line)
		}
	}

	// filterChoices returns filtered slice and a mapping from filtered index -> absolute index.
	filterChoices := func(query string) ([]Choice, []int) {
		if query == "" {
			abs := make([]int, len(ms.choices))
			for i := range abs {
				abs[i] = i
			}
			return ms.choices, abs
		}
		var filtered []Choice
		var indexMap []int
		q := strings.ToLower(query)
		for i, c := range ms.choices {
			if strings.Contains(strings.ToLower(c.Label), q) {
				filtered = append(filtered, c)
				indexMap = append(indexMap, i)
			}
		}
		return filtered, indexMap
	}

	// filteredIndexMap maps filtered index -> absolute index in ms.choices
	_, filteredIndexMap := filterChoices(searchQuery)

	resetCursorAfterFilter := func() {
		if len(filteredChoices) == 0 {
			cursorIdx, startIdx, endIdx = 0, 0, 0
			return
		}
		if cursorIdx >= len(filteredChoices) {
			cursorIdx = len(filteredChoices) - 1
		}
		if cursorIdx < startIdx {
			startIdx = cursorIdx
		}
		if cursorIdx >= startIdx+pageSize {
			startIdx = max(0, cursorIdx-pageSize+1)
		}
		endIdx = min(startIdx+pageSize, len(filteredChoices))
	}

	navigateUp := func() {
		if cursorIdx > 0 {
			cursorIdx--
			if cursorIdx < startIdx {
				startIdx = cursorIdx
				endIdx = min(startIdx+pageSize, len(filteredChoices))
			}
		}
	}

	navigateDown := func() {
		if cursorIdx < len(filteredChoices)-1 {
			cursorIdx++
			if cursorIdx >= endIdx {
				endIdx = cursorIdx + 1
				startIdx = max(0, endIdx-pageSize)
			}
		}
	}

	// toggleSelection toggles selection of the cursor item (if not disabled and within max)
	toggleSelection := func() {
		if len(filteredChoices) == 0 {
			return
		}
		absIdx := filteredIndexMap[cursorIdx]
		c := ms.choices[absIdx]
		if c.Disabled {
			return
		}
		if isSelected(absIdx) {
			delete(selected, absIdx)
			return
		}
		if canAddMore() {
			selected[absIdx] = struct{}{}
		}
	}

	// selectAllVisible selects all non-disabled visible choices, respecting max cap
	selectAllVisible := func() {
		for i := startIdx; i < endIdx; i++ {
			absIdx := filteredIndexMap[i]
			c := ms.choices[absIdx]
			if c.Disabled {
				continue
			}
			if !isSelected(absIdx) && canAddMore() {
				selected[absIdx] = struct{}{}
			}
		}
	}

	// clearAll unselects everything
	clearAll := func() {
		for k := range selected {
			delete(selected, k)
		}
	}

	// redraw paints the whole prompt area again
	validationMsg := ""
	redraw := func() {
		restoreCursor()
		os.Stdout.WriteString("\n")
		if ms.help != "" {
			os.Stdout.WriteString(helpLine + "\n")
		}
		os.Stdout.WriteString("\r" + promptLine + "\n")

		line := searchPrefix + searchQuery
		if searchMode {
			line += ms.theme.AccentStyle(" ◂ " + strconv.Itoa(len(filteredChoices)) + " hits")
		}
		os.Stdout.WriteString("\r" + line)
		clearLineTillEnd()
		os.Stdout.WriteString("\n")

		for i := startIdx; i < endIdx; i++ {
			absIdx := filteredIndexMap[i]
			cur := i == cursorIdx
			os.Stdout.WriteString("\r" + renderChoice(absIdx, cur))
			clearLineTillEnd()
			os.Stdout.WriteString("\n")
		}
		// Blank out any leftover lines in the window
		for i := endIdx - startIdx; i < pageSize; i++ {
			os.Stdout.WriteString("\r")
			clearLineTillEnd()
			os.Stdout.WriteString("\n")
		}

		if searchMode {
			os.Stdout.WriteString("\n\r" + infoLineSearch)
		} else {
			os.Stdout.WriteString("\n\r" + infoLineNormal)
		}
		clearLineTillEnd()

		// show selection count and validation hint
		selCount := countSelected()
		status := "selected: " + strconv.Itoa(selCount)
		if ms.minChoicesRequired > 0 {
			status += " | min " + strconv.Itoa(ms.minChoicesRequired)
		}
		if ms.maxChoicesAllowed > 0 {
			status += " | max " + strconv.Itoa(ms.maxChoicesAllowed)
		}
		os.Stdout.WriteString("\n\r" + ms.theme.MutedStyle(status))
		clearLineTillEnd()

		if validationMsg != "" {
			os.Stdout.WriteString("\n\r" + ms.theme.WarningStyle(validationMsg))
			clearLineTillEnd()
		}
	}

	// Initial paint ----------------------------------------
	os.Stdout.WriteString("\n")
	if ms.help != "" {
		os.Stdout.WriteString(helpLine + "\n")
	}
	os.Stdout.WriteString("\r" + promptLine + "\n")
	os.Stdout.WriteString("\r" + searchPrefix + "\n")
	for i := startIdx; i < endIdx; i++ {
		absIdx := filteredIndexMap[i]
		cur := i == cursorIdx
		os.Stdout.WriteString("\r" + renderChoice(absIdx, cur) + "\n")
	}
	os.Stdout.WriteString("\n\r" + infoLineNormal)
	os.Stdout.WriteString("\n\r" + ms.theme.MutedStyle("selected: "+strconv.Itoa(countSelected())))
	// ------------------------------------------------------

	// Input loop -------------------------------------------
	err := keyboard.Listen(func(key keys.Key) (bool, error) {
		switch key.Code {
		case keys.Tab:
			searchMode = !searchMode
			validationMsg = ""
		case keys.CtrlC:
			return true, ErrInterrupted
		case keys.Escape:
			if searchMode {
				searchMode = false
				validationMsg = ""
			}
		case keys.Enter:
			validationMsg = ""
			// enforce min/max
			sz := countSelected()
			if ms.minChoicesRequired > 0 && sz < ms.minChoicesRequired {
				validationMsg = "pick at least " + strconv.Itoa(ms.minChoicesRequired)
				redraw()
				return false, nil
			}
			if ms.maxChoicesAllowed > 0 && sz > ms.maxChoicesAllowed {
				validationMsg = "pick at most " + strconv.Itoa(ms.maxChoicesAllowed)
				redraw()
				return false, nil
			}
			return true, nil

		case keys.Up, keys.Left:
			navigateUp()
		case keys.Down, keys.Right:
			navigateDown()
		case keys.Space:
			toggleSelection()
		case keys.RuneKey:
			if len(key.Runes) == 0 {
				break
			}
			r := key.Runes[0]
			if searchMode {
				// typing modifies search
				searchQuery += string(r)
				filteredChoices, filteredIndexMap = filterChoices(searchQuery)
				pageSize = min(ms.pageSize, len(filteredChoices))
				resetCursorAfterFilter()
			} else {
				switch r {
				case 'j':
					navigateDown()
				case 'k':
					navigateUp()
				case 'h':
					navigateUp()
				case 'l':
					navigateDown()
				case 'a': // select all visible
					selectAllVisible()
				case 'n': // clear all
					clearAll()
				}
			}
		case keys.Backspace:
			if searchMode && len(searchQuery) > 0 {
				searchQuery = searchQuery[:len(searchQuery)-1]
				filteredChoices, filteredIndexMap = filterChoices(searchQuery)
				pageSize = min(ms.pageSize, len(filteredChoices))
				resetCursorAfterFilter()
			}
		}

		redraw()
		return false, nil
	})

	// Cleanup ----------------------------------------------
	if err != nil {
		restoreCursor()
		clearTillEnd()
		showCursor()
		return nil, err
	}

	restoreCursor()
	clearTillEnd()
	showCursor()

	// Build result slice in the original order
	var out []Choice
	for i, c := range ms.choices {
		if isSelected(i) {
			out = append(out, c)
		}
	}
	return out, nil
}

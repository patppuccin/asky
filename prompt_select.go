package asky

import (
	"bufio"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/mattn/go-runewidth"
)

// singleSelect renders an interactive single-selection prompt.
// Construct one with [Select].
type singleSelect struct {
	cfg              Config
	prefix           string
	label            string
	choices          []Choice
	defaultChoiceIdx *int
	cursorIndicator  string
	selectionMarker  string
	disabledMarker   string
	pageSize         int
	selectedChoice   Choice
	validator        func(Choice) (string, bool)
}

// Select returns a builder for an interactive single-selection prompt.
//
//	choice, err := asky.Select().WithLabel("Pick one").WithChoices(choices).Render()
//	if errors.Is(err, asky.ErrInterrupted) { ... }
func Select() *singleSelect {
	return &singleSelect{
		cfg:             pkgConfig,
		label:           "Select an option",
		choices:         []Choice{},
		cursorIndicator: ">",
		selectionMarker: "*",
		disabledMarker:  "-",
		pageSize:        10,
	}
}

// WithStyles overrides the [StyleMap] for this prompt.
func (s *singleSelect) WithStyles(sm *StyleMap) *singleSelect {
	s.cfg.Styles = sm
	return s
}

// WithPrefix overrides the default prompt prefix symbol.
func (s *singleSelect) WithPrefix(p string) *singleSelect {
	s.prefix = p
	return s
}

// WithLabel sets the prompt label shown to the user.
func (s *singleSelect) WithLabel(l string) *singleSelect {
	s.label = l
	return s
}

// WithChoices sets the list of choices available for selection.
func (s *singleSelect) WithChoices(ch []Choice) *singleSelect {
	s.choices = ch
	return s
}

// WithDefaultChoice pre-selects a choice by index.
func (s *singleSelect) WithDefaultChoice(idx int) *singleSelect {
	s.defaultChoiceIdx = &idx
	return s
}

// WithPageSize sets the number of choices visible at once.
func (s *singleSelect) WithPageSize(n int) *singleSelect {
	s.pageSize = n
	return s
}

// WithCursorIndicator overrides the cursor indicator symbol.
func (s *singleSelect) WithCursorIndicator(ind string) *singleSelect {
	s.cursorIndicator = ind
	return s
}

// WithSelectionMarker overrides the selection marker symbol.
func (s *singleSelect) WithSelectionMarker(mrk string) *singleSelect {
	s.selectionMarker = mrk
	return s
}

// WithDisabledMarker overrides the disabled item marker symbol.
func (s *singleSelect) WithDisabledMarker(mrk string) *singleSelect {
	s.disabledMarker = mrk
	return s
}

// WithValidator sets a validator called on enter. Use [ValidateSelectRequired]
// or a custom func(Choice) (string, bool).
func (s *singleSelect) WithValidator(v func(Choice) (string, bool)) *singleSelect {
	s.validator = v
	return s
}

// Render displays the prompt and blocks until the user confirms or cancels.
// Returns the selected [Choice], or [ErrInterrupted] if Ctrl+C is pressed.
//
// In accessible mode, choices are printed as a numbered list and the user
// types the index. In interactive mode, choices are navigated with arrow keys.
func (s *singleSelect) Render() (Choice, error) {
	if len(s.choices) == 0 {
		return Choice{}, ErrNoSelectionChoices
	}
	if s.cfg.Accessible {
		return s.renderAccessible()
	}
	return s.renderInteractive()
}

// renderAccessible prints a numbered list and collects the user's choice by index.
func (s *singleSelect) renderAccessible() (Choice, error) {
	prefix := pick(s.prefix, "(?)")

	// Print header
	stdOutput.Write([]byte(
		safeStyle(s.cfg.Styles.SelectionPrefix).Sprint(prefix+" ") +
			safeStyle(s.cfg.Styles.SelectionLabel).Sprint(s.label) + "\n",
	))

	// Print numbered choices
	width := len(strconv.Itoa(len(s.choices)))
	for i, c := range s.choices {
		num := safeStyle(s.cfg.Styles.SelectionSearchHint).Sprintf("%*d. ", width, i+1)
		var label string
		if c.Disabled {
			label = safeStyle(s.cfg.Styles.SelectionItemDisabledLabel).Sprint(c.Label) +
				safeStyle(s.cfg.Styles.SelectionItemDisabledMarker).Sprint(" (disabled)")
		} else {
			label = safeStyle(s.cfg.Styles.SelectionItemNormalLabel).Sprint(c.Label)
		}
		stdOutput.Write([]byte("  " + num + label + "\n"))
	}

	// Build prompt hint
	hint := ""
	if s.defaultChoiceIdx != nil {
		idx := *s.defaultChoiceIdx
		if idx >= 0 && idx < len(s.choices) {
			hint = safeStyle(s.cfg.Styles.SelectionSearchHint).
				Sprintf("(default: %d) ", idx+1)
		}
	}
	promptStr := safeStyle(s.cfg.Styles.SelectionPrefix).Sprint("> ") + hint

	for {
		stdOutput.Write([]byte(promptStr))

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

		type readResult struct {
			line string
			err  error
		}
		ch := make(chan readResult, 1)
		go func() {
			line, err := bufio.NewReader(os.Stdin).ReadString('\n')
			ch <- readResult{line, err}
		}()

		var line string
		select {
		case <-sigCh:
			signal.Stop(sigCh)
			stdOutput.Write([]byte("\n"))
			return Choice{}, ErrInterrupted
		case r := <-ch:
			signal.Stop(sigCh)
			if r.err != nil {
				if isInterrupt(r.err) {
					stdOutput.Write([]byte("\n"))
					return Choice{}, ErrInterrupted
				}
				return Choice{}, r.err
			}
			line = strings.TrimSpace(strings.TrimRight(r.line, "\r\n"))
		}

		// Empty input — use default if set
		if line == "" {
			if s.defaultChoiceIdx != nil {
				idx := *s.defaultChoiceIdx
				if idx >= 0 && idx < len(s.choices) {
					chosen := s.choices[idx]
					if s.validator != nil {
						if msg, ok := s.validator(chosen); !ok {
							stdOutput.Write([]byte(safeStyle(s.cfg.Styles.SelectionValidationFail).Sprint(msg) + "\n"))
							continue
						}
					}
					return chosen, nil
				}
			}
			stdOutput.Write([]byte(safeStyle(s.cfg.Styles.SelectionValidationFail).Sprint("please enter a number") + "\n"))
			continue
		}

		// Parse number
		n, err := strconv.Atoi(line)
		if err != nil || n < 1 || n > len(s.choices) {
			stdOutput.Write([]byte(
				safeStyle(s.cfg.Styles.SelectionValidationFail).
					Sprintf("enter a number between 1 and %d\n", len(s.choices)),
			))
			continue
		}

		chosen := s.choices[n-1]

		if chosen.Disabled {
			stdOutput.Write([]byte(safeStyle(s.cfg.Styles.SelectionValidationFail).Sprint("that choice is disabled\n")))
			continue
		}

		if s.validator != nil {
			if msg, ok := s.validator(chosen); !ok {
				stdOutput.Write([]byte(safeStyle(s.cfg.Styles.SelectionValidationFail).Sprint(msg) + "\n"))
				continue
			}
		}

		return chosen, nil
	}
}

// renderInteractive renders a navigable list with search. Arrow keys and
// vi-keys move the cursor, space selects, enter confirms.
func (s *singleSelect) renderInteractive() (Choice, error) {
	if err := reserveLines(6 + s.pageSize); err != nil {
		return Choice{}, ErrTerminalTooSmall
	}

	// State
	interrupted := false
	searchQuery := ""
	searchMode := false
	filteredChoices := s.choices
	pageSize := min(s.pageSize, len(filteredChoices))
	cursorIdx := 0
	startIdx := 0
	endIdx := min(len(filteredChoices), pageSize)
	valMessage := ""

	// Line constructors
	prefix := pick(s.prefix, "(?)")
	promptLine := safeStyle(s.cfg.Styles.SelectionPrefix).Sprint(prefix+" ") +
		safeStyle(s.cfg.Styles.SelectionLabel).Sprint(s.label)
	searchLabel := safeStyle(s.cfg.Styles.SelectionSearchLabel).Sprint("Search: ")
	helpNormal := safeStyle(s.cfg.Styles.SelectionHelp).Sprint("↑/↓ move • space select • enter confirm" + ansiClearLine + "\n\rtab to search" + ansiClearLine)
	helpSearch := safeStyle(s.cfg.Styles.SelectionHelp).Sprint("↑/↓ move • space select • enter confirm" + ansiClearLine + "\n\rtype to search (esc/tab nav)" + ansiClearLine)

	renderChoice := func(c Choice, cur, sel bool) string {
		cursorSpacer := strings.Repeat(" ", runewidth.StringWidth(s.cursorIndicator))
		selSpacer := strings.Repeat(" ", runewidth.StringWidth(s.selectionMarker))
		switch {
		case c.Disabled && cur:
			return safeStyle(s.cfg.Styles.SelectionItemDisabledMarker).Sprint(s.cursorIndicator+s.disabledMarker) +
				safeStyle(s.cfg.Styles.SelectionItemDisabledLabel).Sprint(c.Label)
		case sel && cur:
			return safeStyle(s.cfg.Styles.SelectionItemSelectedMarker).Sprint(s.cursorIndicator+s.selectionMarker) +
				safeStyle(s.cfg.Styles.SelectionItemSelectedLabel).Sprint(c.Label)
		case c.Disabled:
			return cursorSpacer +
				safeStyle(s.cfg.Styles.SelectionItemDisabledMarker).Sprint(s.disabledMarker) +
				safeStyle(s.cfg.Styles.SelectionItemDisabledLabel).Sprint(c.Label)
		case sel:
			return cursorSpacer +
				safeStyle(s.cfg.Styles.SelectionItemSelectedMarker).Sprint(s.selectionMarker) +
				safeStyle(s.cfg.Styles.SelectionItemSelectedLabel).Sprint(c.Label)
		case cur:
			return safeStyle(s.cfg.Styles.SelectionItemCurrentMarker).Sprint(s.cursorIndicator) + selSpacer +
				safeStyle(s.cfg.Styles.SelectionItemCurrentLabel).Sprint(c.Label)
		default:
			return cursorSpacer + selSpacer +
				safeStyle(s.cfg.Styles.SelectionItemNormalLabel).Sprint(c.Label)
		}
	}

	filterChoices := func(query string) []Choice {
		if query == "" {
			return s.choices
		}
		var filtered []Choice
		q := strings.ToLower(query)
		for _, c := range s.choices {
			if strings.Contains(strings.ToLower(c.Label), q) {
				filtered = append(filtered, c)
			}
		}
		return filtered
	}

	resetCursor := func() {
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

	redraw := func() {
		stdOutput.Write([]byte(ansiRestoreCursor + "\r" + promptLine + "\n"))

		// Search line
		sl := searchLabel + safeStyle(s.cfg.Styles.SelectionSearchText).Sprint(searchQuery)
		if searchMode {
			sl += safeStyle(s.cfg.Styles.SelectionSearchHint).Sprint(" ◂ " + strconv.Itoa(len(filteredChoices)) + " hits")
		}
		if s.selectedChoice != (Choice{}) {
			sl += safeStyle(s.cfg.Styles.SelectionSearchHint).Sprint(" [1 selected]")
		} else {
			sl += safeStyle(s.cfg.Styles.SelectionSearchHint).Sprint(" [0 selected]")
		}
		stdOutput.Write([]byte("\r" + sl + ansiClearLine + "\n"))

		// Choices
		for i := startIdx; i < endIdx; i++ {
			c := filteredChoices[i]
			stdOutput.Write([]byte("\r" + renderChoice(c, i == cursorIdx, c.Value == s.selectedChoice.Value) + ansiClearLine + "\n"))
		}
		// Pad remaining page lines
		for i := endIdx - startIdx; i < pageSize; i++ {
			stdOutput.Write([]byte("\r" + ansiClearLine + "\n"))
		}

		// Validation message
		stdOutput.Write([]byte("\n\r" + safeStyle(s.cfg.Styles.SelectionValidationFail).Sprint(valMessage) + ansiClearLine + "\n\r"))

		// Help line
		if searchMode {
			stdOutput.Write([]byte(helpSearch))
		} else {
			stdOutput.Write([]byte(helpNormal))
		}
	}

	// Save cursor, defer cleanup
	stdOutput.Write([]byte(ansiHideCursor + ansiSaveCursor))
	defer stdOutput.Write([]byte(ansiRestoreCursor + ansiClearScreen + ansiReset + ansiShowCursor))

	// Apply default selection
	if s.defaultChoiceIdx != nil {
		idx := *s.defaultChoiceIdx
		if idx >= 0 && idx < len(s.choices) {
			s.selectedChoice = s.choices[idx]
		}
	}

	redraw()

	err := listenKeys(func(ev keyEvent) (stop bool) {
		switch ev.code {
		case keyCtrlC:
			interrupted = true
			return true

		case keyUp:
			navigateUp()

		case keyDown:
			navigateDown()

		case keyTab:
			searchMode = !searchMode

		case keyEscape:
			searchMode = false

		case keyEnter:
			if s.validator != nil {
				if msg, ok := s.validator(s.selectedChoice); !ok {
					valMessage = msg
					break
				}
			}
			return true

		case keySpace:
			if len(filteredChoices) == 0 {
				valMessage = "no choices available"
				break
			}
			cur := filteredChoices[cursorIdx]
			if cur.Disabled {
				valMessage = "cannot select a disabled choice"
				break
			}
			if s.selectedChoice.Value == cur.Value {
				s.selectedChoice = Choice{}
			} else {
				s.selectedChoice = cur
			}
			valMessage = ""

		case keyBackspace:
			if searchMode && len(searchQuery) > 0 {
				searchQuery = searchQuery[:len(searchQuery)-1]
				filteredChoices = filterChoices(searchQuery)
				resetCursor()
			}

		case keyRune:
			if searchMode {
				searchQuery += string(ev.r)
				filteredChoices = filterChoices(searchQuery)
				resetCursor()
			} else {
				switch ev.r {
				case 'j', 'l':
					navigateDown()
				case 'k', 'h':
					navigateUp()
				}
			}
		}

		redraw()
		return false
	})

	if err != nil {
		return Choice{}, err
	}
	if interrupted {
		return Choice{}, ErrInterrupted
	}
	return s.selectedChoice, nil
}

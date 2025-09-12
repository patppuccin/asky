package asky

import (
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/mattn/go-runewidth"
)

// --- Presets ---------------------------------------------
var ProgressPatternDefault = ProgressPattern{"+", "-", " [", "] "}
var ProgressPatternHashes = ProgressPattern{"#", "=", " [", "] "}
var ProgressPatternPipes = ProgressPattern{"|", " ", " ", " "}
var ProgressPatternSolid = ProgressPattern{"█", "░", " ", " "}

// --- Definition -------------------------------------------
type ProgressPattern struct {
	DoneChar, PendingChar   string
	BarPadLeft, BarPadRight string
}

type Progress struct {
	theme                      *Theme
	style                      *Style
	prefix, label, description string
	steps, current, width      int
	pattern                    ProgressPattern

	stop bool
	wg   sync.WaitGroup
	mu   sync.Mutex
}

// --- Initiation ------------------------------------------
func NewProgress() *Progress {
	return &Progress{
		prefix:  "[~] ",
		label:   "Activity in progress",
		width:   40,
		steps:   0,
		current: 0,
		pattern: ProgressPatternDefault,
	}
}

// --- Configuration ---------------------------------------
func (pr *Progress) WithTheme(theme Theme) *Progress           { pr.theme = &theme; return pr }
func (pr *Progress) WithStyle(style Style) *Progress           { pr.style = &style; return pr }
func (pr *Progress) WithPrefix(px string) *Progress            { pr.prefix = px; return pr }
func (pr *Progress) WithLabel(lbl string) *Progress            { pr.label = lbl; return pr }
func (pr *Progress) WithDescription(desc string) *Progress     { pr.description = desc; return pr }
func (pr *Progress) WithWidth(width int) *Progress             { pr.width = max(0, width); return pr }
func (pr *Progress) WithSteps(steps int) *Progress             { pr.steps = max(0, steps); return pr }
func (pr *Progress) WithPattern(ptn ProgressPattern) *Progress { pr.pattern = ptn; return pr }

// --- Presentation ----------------------------------------
func (pr *Progress) Start() {
	// Sanity check for no steps or no label
	if pr.steps <= 0 || pr.label == "" {
		return
	}

	// Setup theme and style (apply defaults if not set)
	if pr.theme == nil {
		pr.theme = &ThemeDefault
	}
	if pr.style == nil {
		pr.style = StyleDefault(pr.theme)
	}

	// Set up default bar width if not set
	if pr.width <= 0 {
		pr.width = 30
	}

	// Ensure terminal is large enough for the prompt to render
	_ = makeSpace(4)

	// Prep and save cursor state
	stdOutput.Write([]byte(ansiSaveCursor + ansiHideCursor + ansiClearLine + "\n\r"))
	pr.stop = false

	// Redraw the progress bar with current state.
	redraw := func() {
		// Acquire lock on the progress bar state (defer release)
		pr.mu.Lock()
		defer pr.mu.Unlock()

		// Clamp ratio of current to steps between 0 and 1.
		ratio := float64(pr.current) / float64(pr.steps)
		ratio = min(max(ratio, 0), 1)

		// Format percentage segment (padded to 3 chars).
		percent := strconv.Itoa(int(ratio * 100))
		for runewidth.StringWidth(percent) < 3 {
			percent = " " + percent
		}
		percent += "% "

		// Determine terminal width (if unknown, fallback to 80)
		termWidth, _, _ := getTermDimensions()
		if termWidth <= 0 {
			termWidth = 80
		}

		// Compute available width for the bar from available terminal width
		fixedWidth := runewidth.StringWidth(pr.prefix + pr.label + percent + pr.pattern.BarPadLeft + pr.pattern.BarPadRight)
		availWidth := max(termWidth-fixedWidth, 0)
		barWidth := min(availWidth, pr.width)

		// Calculate filled & pending segments of the bar
		filled := int(ratio * float64(barWidth))
		filled = min(filled, barWidth)
		pending := barWidth - filled

		// Build progress bar segments (with styling)
		doneChars := strings.Repeat(pr.pattern.DoneChar, filled)
		pendingChars := strings.Repeat(pr.pattern.PendingChar, pending)
		bar := pr.style.ProgressBarPad.Sprint(pr.pattern.BarPadLeft) +
			pr.style.ProgressBarDone.Sprint(doneChars) +
			pr.style.ProgressBarPending.Sprint(pendingChars) +
			pr.style.ProgressBarPad.Sprint(pr.pattern.BarPadRight)

		// Redraw the screen: restore cursor, print optional description, then the bar.
		stdOutput.Write([]byte(ansiRestoreCursor + "\n\r"))
		if pr.description != "" {
			stdOutput.Write([]byte(pr.style.ProgressDesc.Sprint(pr.description) + "\n\r"))
		}
		stdOutput.Write([]byte(pr.style.ProgressPrefix.Sprint(pr.prefix)))
		stdOutput.Write([]byte(pr.style.ProgressLabel.Sprint(pr.label)))
		stdOutput.Write([]byte(bar))
		stdOutput.Write([]byte(pr.style.ProgressBarStatus.Sprint(percent) + ansiClearLine))
	}

	// Watch for interrupts and stop the progress
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		pr.Stop()
		os.Exit(1)
	}()

	// Run the progress bar render loop until stop (completion or interrupt)
	pr.wg.Go(func() {
		defer stdOutput.Write([]byte(ansiRestoreCursor + ansiClearScreen + ansiReset + ansiShowCursor))
		for !pr.stop {
			redraw()
			time.Sleep(100 * time.Millisecond)
		}
	})
}

// Trigger stop of the progress bar
func (pr *Progress) Stop() {
	pr.stop = true
	pr.wg.Wait()
}

// Increment the progress bar by one step
func (pr *Progress) Increment() {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	if pr.steps <= 0 {
		return
	}

	pr.current++
	if pr.current > pr.steps {
		pr.current = pr.steps
	}
}

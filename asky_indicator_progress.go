package asky

import (
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

// ProgressPattern defines how the bar looks
type ProgressPattern struct {
	DoneChar    string
	PendingChar string
	BarPadLeft  string
	BarPadRight string
}

// Some built-in patterns
var ProgressPatternDefault = ProgressPattern{
	DoneChar:    "+",
	PendingChar: "-",
	BarPadLeft:  " [",
	BarPadRight: "] ",
}

var ProgressPatternSharp = ProgressPattern{
	DoneChar:    "#",
	PendingChar: "=",
	BarPadLeft:  " [",
	BarPadRight: "] ",
}

var ProgressPatternPipes = ProgressPattern{
	DoneChar:    "|",
	PendingChar: " ",
	BarPadLeft:  " ",
	BarPadRight: " ",
}

var ProgressPatternSolid = ProgressPattern{
	DoneChar:    "█",
	PendingChar: " ",
	BarPadLeft:  " ",
	BarPadRight: " ",
}

// Progress holds state for an animated progress bar
type Progress struct {
	prefix      string
	label       string
	description string
	steps       int
	current     int
	barWidth    int
	theme       Theme
	pattern     ProgressPattern

	stop bool
	wg   sync.WaitGroup
	mu   sync.Mutex
}

// Constructor with defaults
func NewProgress() *Progress {
	return &Progress{
		theme:       ThemeDefault,
		prefix:      "[~] ",
		label:       "Working...",
		description: "",
		steps:       100,
		current:     0,
		barWidth:    30,
		pattern:     ProgressPatternDefault,
	}
}

// Fluent config
func (pr *Progress) WithTheme(h Theme) *Progress             { pr.theme = h; return pr }
func (pr *Progress) WithPrefix(s string) *Progress           { pr.prefix = s; return pr }
func (pr *Progress) WithWidth(w int) *Progress               { pr.barWidth = max(0, w); return pr }
func (pr *Progress) WithPattern(p ProgressPattern) *Progress { pr.pattern = p; return pr }
func (pr *Progress) WithLabel(l string) *Progress            { pr.label = l; return pr }
func (pr *Progress) WithDescription(d string) *Progress      { pr.description = d; return pr }
func (pr *Progress) WithSteps(t int) *Progress               { pr.steps = max(0, t); return pr }

// Start begins the render loop
func (pr *Progress) Start() {
	// Get the style preset
	preset := newPreset(pr.theme)

	// Save cursor state before prompt & defer reset
	pr.stop = false
	os.Stdout.WriteString(ansiSaveCursor + ansiHideCursor + ansiClearLineEnd + "\n\r")

	// Print the description line (no need to redraw on updates)
	if pr.description != "" {
		os.Stdout.WriteString(preset.accent.Sprint(pr.description) + "\n")
	}

	// Helper: Redraw the progress bar with current state
	redraw := func() {
		pr.mu.Lock()
		defer pr.mu.Unlock()

		ratio := float64(pr.current) / float64(pr.steps)
		if ratio < 0 {
			ratio = 0
		}
		if ratio > 1 {
			ratio = 1
		}

		filled := int(ratio * float64(pr.barWidth))
		filled = min(filled, pr.barWidth)
		pending := pr.barWidth - filled

		bar := strings.Repeat(pr.pattern.DoneChar, filled) +
			strings.Repeat(pr.pattern.PendingChar, pending)

		percent := strconv.Itoa(int(ratio * 100))
		for len(percent) < 3 {
			percent = " " + percent
		}
		percent += "% "

		os.Stdout.WriteString(preset.primary.Sprint(pr.prefix))
		os.Stdout.WriteString(preset.secondary.Sprint(pr.label + pr.pattern.BarPadLeft + bar + pr.pattern.BarPadRight))
		os.Stdout.WriteString(preset.primary.Sprint(" " + percent + " " + ansiClearLineEnd + "\r"))
	}

	// Watch for Ctrl+C and set stop flag
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		pr.Stop()
		os.Exit(1)
	}()

	// Run the progress bar render loop until stop (completion or interrupt)
	pr.wg.Go(func() {
		defer os.Stdout.WriteString(ansiRestoreCursor + ansiClearScreenEnd + ansiReset + ansiShowCursor)
		for !pr.stop {
			redraw()
			time.Sleep(100 * time.Millisecond)
		}
	})
}

// Stop finishes the bar and restores cursor
func (pr *Progress) Stop() {
	pr.stop = true
	pr.wg.Wait()
}

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

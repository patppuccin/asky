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
	prefix   string
	label    string
	help     string
	steps    int
	current  int
	barWidth int
	theme    Theme
	pattern  ProgressPattern

	stop bool
	wg   sync.WaitGroup
	mu   sync.Mutex
}

// Constructor with defaults
func NewProgress() *Progress {
	return &Progress{
		theme:    ThemeDefault,
		prefix:   "[~] ",
		label:    "Working...",
		help:     "",
		steps:    100,
		current:  0,
		barWidth: 30,
		pattern:  ProgressPatternDefault,
	}
}

// Fluent config
func (p *Progress) WithTheme(h Theme) *Progress             { p.theme = h; return p }
func (p *Progress) WithPrefix(s string) *Progress           { p.prefix = s; return p }
func (p *Progress) WithWidth(w int) *Progress               { p.barWidth = max(0, w); return p }
func (p *Progress) WithPattern(x ProgressPattern) *Progress { p.pattern = x; return p }
func (p *Progress) WithLabel(t string) *Progress            { p.label = t; return p }
func (p *Progress) WithHelp(t string) *Progress             { p.help = t; return p }
func (p *Progress) WithSteps(t int) *Progress               { p.steps = max(0, t); return p }

// Start begins the render loop
func (p *Progress) Start() {

	// Save cursor state before prompt & defer reset
	p.stop = false
	os.Stdout.WriteString(ansiSaveCursor + ansiHideCursor + ansiClearLineEnd + "\n\r")

	// Print the helper line (no need to redraw on updates)
	if p.help != "" {
		os.Stdout.WriteString(p.theme.MutedStyle(p.help) + "\n")
	}

	// Helper: Redraw the progress bar with current state
	redraw := func() {
		p.mu.Lock()
		defer p.mu.Unlock()

		ratio := float64(p.current) / float64(p.steps)
		if ratio < 0 {
			ratio = 0
		}
		if ratio > 1 {
			ratio = 1
		}

		filled := int(ratio * float64(p.barWidth))
		filled = min(filled, p.barWidth)
		pending := p.barWidth - filled

		bar := strings.Repeat(p.pattern.DoneChar, filled) +
			strings.Repeat(p.pattern.PendingChar, pending)

		percent := strconv.Itoa(int(ratio * 100))
		for len(percent) < 3 {
			percent = " " + percent
		}
		percent += "% "

		os.Stdout.WriteString(p.theme.PrimaryStyle(p.prefix))
		os.Stdout.WriteString(p.theme.SecondaryStyle(p.label))
		os.Stdout.WriteString(p.theme.AccentStyle(p.pattern.BarPadLeft + bar + p.pattern.BarPadRight))
		os.Stdout.WriteString(p.theme.SecondaryStyle(" " + percent + " " + ansiClearLineEnd + "\r"))
	}

	// Watch for Ctrl+C and set stop flag
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		p.Stop(true)
		os.Exit(1)
	}()

	// Run the progress bar render loop until stop (completion or intterupt)
	p.wg.Go(func() {
		defer os.Stdout.WriteString(ansiRestoreCursor + ansiClearScreenEnd + ansiReset + ansiShowCursor)
		for !p.stop {
			redraw()
			time.Sleep(100 * time.Millisecond)
		}
	})
}

// Stop finishes the bar and restores cursor
func (p *Progress) Stop(clear bool) {
	p.stop = true
	p.wg.Wait()
}

func (p *Progress) Increment() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.steps <= 0 {
		return
	}

	p.current++
	if p.current > p.steps {
		p.current = p.steps
	}
}

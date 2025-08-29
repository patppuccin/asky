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
	DoneChar:    "=",
	PendingChar: " ",
	BarPadLeft:  " [",
	BarPadRight: "] ",
}

var ProgressPatternMathSymbols = ProgressPattern{
	DoneChar:    "+",
	PendingChar: "-",
	BarPadLeft:  " [",
	BarPadRight: "] ",
}

var ProgressPatternBars = ProgressPattern{
	DoneChar:    "|",
	PendingChar: " ",
	BarPadLeft:  " ",
	BarPadRight: " ",
}

var ProgressPatternBlocks = ProgressPattern{
	DoneChar:    "█",
	PendingChar: " ",
	BarPadLeft:  " ",
	BarPadRight: " ",
}

// Progress holds state for an animated progress bar
type Progress struct {
	progressSymbol string
	progressText   string
	helperText     string
	totalSteps     int
	currentStep    int
	barWidth       int
	theme          Theme
	pattern        ProgressPattern

	stop bool
	wg   sync.WaitGroup
	mu   sync.Mutex
}

// Constructor with defaults
func NewProgress() *Progress {
	return &Progress{
		progressSymbol: "[~] ",
		progressText:   "Working...",
		totalSteps:     100,
		currentStep:    0,
		barWidth:       40,
		theme:          ThemeDefault,
		pattern:        ProgressPatternDefault,
	}
}

// Fluent config
func (p *Progress) WithProgressSymbol(s string) *Progress { p.progressSymbol = s; return p }
func (p *Progress) WithWidth(w int) *Progress {
	if w >= 0 {
		p.barWidth = w
	}
	return p
}
func (p *Progress) WithPattern(x ProgressPattern) *Progress { p.pattern = x; return p }
func (p *Progress) WithTheme(h Theme) *Progress             { p.theme = h; return p }
func (p *Progress) WithProgressText(t string) *Progress     { p.progressText = t; return p }
func (p *Progress) WithHelperText(t string) *Progress       { p.helperText = t; return p }
func (p *Progress) WithTotalSteps(t int) *Progress {
	if t > 0 {
		p.totalSteps = t
	}
	return p
}

// Start begins the render loop
func (p *Progress) Start() {
	p.stop = false
	os.Stdout.Write([]byte("\033[s"))    // save cursor
	os.Stdout.Write([]byte("\033[?25l")) // hide cursor
	os.Stdout.Write([]byte("\r\n"))      // clear line

	if p.helperText != "" {
		os.Stdout.WriteString(p.theme.MutedStyle(p.helperText) + "\n")
	}

	// handle Ctrl+C
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		p.Stop(true)
		os.Exit(1)
	}()

	p.wg.Go(func() {
		for !p.stop {
			p.renderOnce()
			time.Sleep(100 * time.Millisecond)
		}
		// final cleanup
		os.Stdout.Write([]byte("\033[u\033[J\033[?25h"))
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

	if p.totalSteps <= 0 {
		return
	}

	p.currentStep++
	if p.currentStep > p.totalSteps {
		p.currentStep = p.totalSteps
	}
}

// internal renderer
func (p *Progress) renderOnce() {
	p.mu.Lock()
	defer p.mu.Unlock()

	ratio := float64(p.currentStep) / float64(p.totalSteps)
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}

	filled := int(ratio * float64(p.barWidth))
	if filled > p.barWidth {
		filled = p.barWidth
	}
	pending := p.barWidth - filled

	bar := strings.Repeat(p.pattern.DoneChar, filled) +
		strings.Repeat(p.pattern.PendingChar, pending)

	percent := strconv.Itoa(int(ratio * 100))
	for len(percent) < 3 {
		percent = " " + percent
	}
	percent += "% "

	os.Stdout.Write([]byte("\r\033[K")) // restore cursor + clear line
	os.Stdout.WriteString(p.theme.SecondaryStyle(p.progressSymbol))
	os.Stdout.WriteString(p.theme.PrimaryStyle(p.progressText))
	os.Stdout.WriteString(p.theme.AccentStyle(p.pattern.BarPadLeft))
	os.Stdout.WriteString(p.theme.AccentStyle(bar))
	os.Stdout.WriteString(p.theme.AccentStyle(p.pattern.BarPadRight))
	os.Stdout.WriteString(p.theme.SecondaryStyle(" " + percent + " "))
}

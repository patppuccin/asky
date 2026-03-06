package asky

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Spinner frame pattern presets.
var (
	SpinnerDefault  = []string{"(⠋)", "(⠙)", "(⠹)", "(⠸)", "(⠼)", "(⠴)", "(⠦)", "(⠧)", "(⠇)", "(⠏)"}
	SpinnerDots     = []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}
	SpinnerDotsMini = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	SpinnerCircles  = []string{"◐", "◓", "◑", "◒"}
	SpinnerSquares  = []string{"▖", "▌", "▘", "▀", "▝", "▐", "▗", "▄"}
	SpinnerLine     = []string{"-", "\\", "|", "/"}
	SpinnerPipes    = []string{"╾", "│", "╸", "┤", "├", "└", "┴", "┬", "┐", "┘"}
	SpinnerMoons    = []string{"🌑", "🌒", "🌓", "🌔", "🌕", "🌖", "🌗", "🌘"}
)

// spinner renders an animated spinner on a single line.
// It saves the cursor before animating and restores it on stop,
// leaving no trace in the terminal output.
// Construct one with [Spinner].
type spinner struct {
	cfg      Config
	frames   []string
	label    string
	interval time.Duration
	stop     bool
	wg       sync.WaitGroup
}

// Spinner returns a spinner builder with sensible defaults.
//
//	sp := asky.Spinner().WithLabel("deploying...")
//	sp.Start()
//	// ... do work ...
//	sp.Stop()
func Spinner() *spinner {
	return &spinner{
		cfg:      pkgConfig,
		frames:   SpinnerDefault,
		label:    "Loading...",
		interval: 100 * time.Millisecond,
	}
}

// WithStyles overrides the [StyleMap] for this spinner.
func (sp *spinner) WithStyles(s *StyleMap) *spinner {
	sp.cfg.Styles = s
	return sp
}

// WithFrames sets a custom frame pattern for the spinner animation.
func (sp *spinner) WithFrames(frames []string) *spinner {
	sp.frames = frames
	return sp
}

// WithLabel sets the label displayed beside the spinner frame.
func (sp *spinner) WithLabel(label string) *spinner {
	sp.label = label
	return sp
}

// WithInterval sets the frame animation interval. Defaults to 100ms.
func (sp *spinner) WithInterval(d time.Duration) *spinner {
	sp.interval = d
	return sp
}

// Start begins the spinner animation in a background goroutine.
// The cursor is saved before animating and restored on [spinner.Stop],
// leaving the terminal clean.
// In accessible mode, prints a single static line instead of animating.
func (sp *spinner) Start() {
	if sp.cfg.Accessible {
		stdOutput.Write([]byte(sp.frames[0] + " " + sp.label + "\n"))
		return
	}

	// Save cursor and hide it before animating
	stdOutput.Write([]byte(ansiSaveCursor + ansiHideCursor))

	// Watch for Ctrl+C & restore terminal before exit
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		sp.Stop()
		os.Exit(1)
	}()

	sp.wg.Go(func() {
		defer stdOutput.Write([]byte(ansiRestoreCursor + ansiClearLine + ansiShowCursor))
		i := 0
		for !sp.stop {
			frame := safeStyle(sp.cfg.Styles.SpinnerPrefix).Sprint(sp.frames[i%len(sp.frames)])
			label := safeStyle(sp.cfg.Styles.SpinnerLabel).Sprint(sp.label)
			stdOutput.Write([]byte(ansiRestoreCursor + ansiClearLine + frame + " " + label))
			i++
			time.Sleep(sp.interval)
		}
	})
}

// Stop halts the spinner, restores the cursor, and clears the spinner line.
// Safe to call multiple times.
func (sp *spinner) Stop() {
	if sp.cfg.Accessible || sp.stop {
		return
	}
	sp.stop = true
	sp.wg.Wait()
}

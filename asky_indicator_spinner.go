package asky

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Patterns ---------------------------------------
var SpinnerPatternDefault = []string{"[⠋] ", "[⠙] ", "[⠹] ", "[⠸] ", "[⠼] ", "[⠴] ", "[⠦] ", "[⠧] ", "[⠇] ", "[⠏] "}
var SpinnerPatternDots = []string{"⠋ ", "⠙ ", "⠹ ", "⠸ ", "⠼ ", "⠴ ", "⠦ ", "⠧ ", "⠇ ", "⠏ "}
var SpinnerPatternCircles = []string{"◐ ", "◓ ", "◑ ", "◒ "}
var SpinnerPatternSquares = []string{"▖ ", "▘ ", "▝ ", "▗ "}
var SpinnerPatternLines = []string{"╾ ", "│ ", "╸ ", "┤ ", "├ ", "└ ", "┴ ", "┬ ", "┐ ", "┘ "}
var SpinnerPatternMoons = []string{"🌑 ", "🌒 ", "🌓 ", "🌔 ", "🌕 ", "🌖 ", "🌗 ", "🌘 "}

// Definition -------------------------------------
type Spinner struct {
	theme       Theme
	frames      []string
	label       string
	description string
	stop        bool
	wg          sync.WaitGroup
}

// Constructors -----------------------------------
func NewSpinner() *Spinner {
	return &Spinner{
		label:  "Loading...",
		frames: SpinnerPatternDefault,
		theme:  ThemeDefault,
	}
}

func (s *Spinner) WithTheme(t Theme) *Spinner          { s.theme = t; return s }
func (s *Spinner) WithFrames(f []string) *Spinner      { s.frames = f; return s }
func (s *Spinner) WithLabel(txt string) *Spinner       { s.label = txt; return s }
func (s *Spinner) WithDescription(txt string) *Spinner { s.description = txt; return s }

// Presentation --------------------------------------------
func (s *Spinner) Start() {
	// Get the style preset
	preset := newPreset(s.theme)

	// Save cursor state before prompt & defer reset
	os.Stdout.WriteString(ansiSaveCursor + ansiHideCursor + ansiClearLineEnd + "\n\r")

	// Print the helper line (no need to redraw on updates)
	if s.description != "" {
		os.Stdout.WriteString(preset.accent.Sprint(s.description) + "\n")
	}

	// Watch for Ctrl+C and set stop flag
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		s.Stop()
		os.Exit(1) // cleanup & quit
	}()

	// Run the spinner render loop until stop (completion or interrupt)
	s.wg.Go(func() {
		defer os.Stdout.WriteString(ansiRestoreCursor + ansiClearScreenEnd + ansiReset + ansiShowCursor)
		i := 0
		for !s.stop {
			currFrame := s.frames[i%len(s.frames)]
			os.Stdout.WriteString(preset.primary.Sprint(currFrame))
			os.Stdout.WriteString(preset.secondary.Sprint(s.label) + ansiClearLineEnd + "\r")
			i++
			time.Sleep(200 * time.Millisecond)
		}
	})
}

func (s *Spinner) Stop() {
	// Stop the spinner & wait for the render loop to exit
	s.stop = true
	s.wg.Wait()
}

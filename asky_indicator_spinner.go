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
	labelText  string
	helperText string
	frames     []string
	theme      Theme
	stop       bool
	wg         sync.WaitGroup
}

// Constructors -----------------------------------
func NewSpinner() *Spinner {
	return &Spinner{
		labelText: "Loading...",
		frames:    SpinnerPatternDefault,
		theme:     ThemeDefault,
	}
}

func (s *Spinner) WithLabelText(txt string) *Spinner  { s.labelText = txt; return s }
func (s *Spinner) WithHelperText(txt string) *Spinner { s.helperText = txt; return s }
func (s *Spinner) WithFrames(f []string) *Spinner     { s.frames = f; return s }
func (s *Spinner) WithTheme(t Theme) *Spinner         { s.theme = t; return s }

// Renderer(s) ------------------------------------
func (s *Spinner) Start() {
	saveCursor()
	hideCursor()
	os.Stdout.Write([]byte("\r\n")) // clear line

	// Print the helper + label
	if s.helperText != "" {
		os.Stdout.WriteString(s.theme.MutedStyle(s.helperText) + "\n")
	}

	// watch for Ctrl+C and set stop flag
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		s.Stop()
		os.Exit(1) // cleanup & quit
	}()

	s.wg.Go(func() {
		i := 0
		for !s.stop {
			thisFrame := s.frames[i%len(s.frames)]
			os.Stdout.Write([]byte(s.theme.SecondaryStyle(thisFrame)))
			os.Stdout.Write([]byte(s.theme.PrimaryStyle(s.labelText)))
			os.Stdout.Write([]byte("\r"))
			time.Sleep(250 * time.Millisecond)
			i++
		}
		// always cleanup here
		restoreCursor()
		clearTillEnd()
		showCursor()
	})
}

func (s *Spinner) Stop() {
	s.stop = true
	s.wg.Wait() // wait until loop exits & cleanup runs
}

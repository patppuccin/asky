package asky

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// --- Presets ---------------------------------------------
var SpinnerPatternDefault = []string{"[⠋] ", "[⠙] ", "[⠹] ", "[⠸] ", "[⠼] ", "[⠴] ", "[⠦] ", "[⠧] ", "[⠇] ", "[⠏] "}
var SpinnerPatternDots = []string{"⣾ ", "⣽ ", "⣻ ", "⢿ ", "⡿ ", "⣟ ", "⣯ ", "⣷ "}
var SpinnerPatternDotsMini = []string{"⠋ ", "⠙ ", "⠹ ", "⠸ ", "⠼ ", "⠴ ", "⠦ ", "⠧ ", "⠇ ", "⠏ "}
var SpinnerPatternCircles = []string{"◐ ", "◓ ", "◑ ", "◒ "}
var SpinnerPatternSquares = []string{"▖ ", "▌ ", "▘ ", "▀ ", "▝ ", "▐ ", "▗ ", "▄ "}
var SpinnerPatternLine = []string{"- ", "\\ ", "| ", "/ "}
var SpinnerPatternPipes = []string{"╾ ", "│ ", "╸ ", "┤ ", "├ ", "└ ", "┴ ", "┬ ", "┐ ", "┘ "}
var SpinnerPatternMoons = []string{"🌑 ", "🌒 ", "🌓 ", "🌔 ", "🌕 ", "🌖 ", "🌗 ", "🌘 "}

// --- Definition ------------------------------------------
type Spinner struct {
	theme       *Theme
	style       *Style
	frames      []string
	label       string
	description string
	stop        bool
	wg          sync.WaitGroup
}

// --- Initiation ------------------------------------------
func NewSpinner() *Spinner {
	return &Spinner{
		label:  "Loading...",
		frames: SpinnerPatternDefault,
	}
}

// Configuration -------------------------------------------
func (sp *Spinner) WithTheme(theme Theme) *Spinner      { sp.theme = &theme; return sp }
func (sp *Spinner) WithStyle(style Style) *Spinner      { sp.style = &style; return sp }
func (sp *Spinner) WithFrames(frames []string) *Spinner { sp.frames = frames; return sp }
func (sp *Spinner) WithLabel(txt string) *Spinner       { sp.label = txt; return sp }
func (sp *Spinner) WithDescription(txt string) *Spinner { sp.description = txt; return sp }

// Presentation --------------------------------------------
func (sp *Spinner) Start() {
	// Sanity check to skip render if both label and prefix are empty
	if sp.label == "" && len(sp.frames) == 0 {
		return
	}

	// Setup theme and style (apply defaults if not set)
	if sp.theme == nil {
		sp.theme = &ThemeDefault
	}
	if sp.style == nil {
		sp.style = StyleDefault(sp.theme)
	}

	// Ensure terminal is large enough for the prompt
	_ = makeSpace(4)

	// Save cursor state before prompt
	stdOutput.Write([]byte(ansiSaveCursor + ansiHideCursor + ansiClearLine + "\n\r"))

	// Print the helper line (no need to redraw on updates)
	if sp.description != "" {
		stdOutput.Write([]byte(sp.style.SpinnerDesc.Sprint(sp.description) + "\n"))
	}

	// Watch for Ctrl+C and set stop flag
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		sp.Stop()
		os.Exit(1) // cleanup & quit
	}()

	// Run the spinner render loop until stop (completion or interrupt)
	sp.wg.Go(func() {
		defer os.Stdout.WriteString(ansiRestoreCursor + ansiClearScreen + ansiReset + ansiShowCursor)
		i := 0
		for !sp.stop {
			currFrame := sp.frames[i%len(sp.frames)]
			stdOutput.Write([]byte(sp.style.SpinnerPrefix.Sprint(currFrame)))
			stdOutput.Write([]byte(sp.style.SpinnerLabel.Sprint(sp.label) + ansiClearLine + "\r"))
			i++
			time.Sleep(200 * time.Millisecond)
		}
	})
}

func (sp *Spinner) Stop() {
	// Stop the spinner & wait for the render loop to exit
	sp.stop = true
	sp.wg.Wait()
}

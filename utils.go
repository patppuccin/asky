package asky

import (
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
	"syscall"

	"github.com/fatih/color"
	"golang.org/x/term"
)

// ErrInterrupted is returned when the user interrupts a prompt (e.g. Ctrl+C).
var ErrInterrupted = errors.New("prompt interrupted")

// ErrTerminalTooSmall is returned when the terminal dimensions are insufficient
// to render a component.
var ErrTerminalTooSmall = errors.New("terminal dimensions too small")

// ErrNoSelectionChoices is returned when a selection prompt is given no choices.
var ErrNoSelectionChoices = errors.New("no choices supplied for selection prompt")

// ErrInvalidSelectionBounds is returned when min count exceeds max count
// in a multi-select prompt configuration.
var ErrInvalidSelectionBounds = errors.New("min count must not exceed max count for multi select prompt")

const (
	ansiHideCursor    = "\033[?25l"
	ansiShowCursor    = "\033[?25h"
	ansiSaveCursor    = "\033[s"
	ansiRestoreCursor = "\033[u"

	ansiReset       = "\033[0m"
	ansiClearLine   = "\033[K"
	ansiClearScreen = "\033[J"
)

// ansiCursorLeft moves the cursor n positions to the left.
func ansiCursorLeft(n int) {
	if n > 0 {
		stdOutput.Write([]byte("\033[" + strconv.Itoa(n) + "D"))
	}
}

// ansiCursorUp moves the cursor n positions up.
func ansiCursorUp(n int) {
	if n > 0 {
		stdOutput.Write([]byte("\033[" + strconv.Itoa(n) + "A"))
	}
}

// termSize returns the current terminal width and height in columns and rows.
func termSize() (int, int, error) {
	return term.GetSize(int(os.Stdout.Fd()))
}

// reserveLines writes n blank lines to stdout then moves the cursor back up,
// reserving vertical space for a component to render into.
// Returns [ErrTerminalTooSmall] if the terminal is too narrow or short.
func reserveLines(lines int) error {
	width, height, _ := termSize()
	if height < lines || width < 50 {
		return ErrTerminalTooSmall
	}
	for range lines {
		os.Stdout.WriteString("\n")
	}
	ansiCursorUp(lines)
	return nil
}

// safeStyle returns s if non-nil, otherwise a no-op Reset style.
// Guards against nil fields on a partially constructed [StyleMap].
func safeStyle(s *color.Color) *color.Color {
	if s != nil {
		return s
	}
	return color.New(color.Reset)
}

// pick returns val if non-empty, otherwise fallback.
func pick(val, fallback string) string {
	if val != "" {
		return val
	}
	return fallback
}

// isInterrupt reports whether err represents a user cancellation.
// Covers io.EOF from bufio, syscall.EINTR from term on Unix,
// and interrupted system call errors on Windows.
func isInterrupt(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, io.EOF) ||
		errors.Is(err, syscall.EINTR) ||
		strings.Contains(err.Error(), "interrupted")
}

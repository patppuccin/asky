package asky

import (
	"errors"
	"io"
	"os"
	"strings"
	"syscall"

	"github.com/fatih/color"
	"golang.org/x/term"
)

// termSize returns the current terminal width and height in columns and rows.
func termSize() (int, int, error) {
	return term.GetSize(int(os.Stdout.Fd()))
}

// reserveLines writes n blank lines to stdout then moves the cursor back up,
// reserving vertical space for a component to render into.
// Returns [ErrTerminalTooSmall] if the terminal has fewer than the
// required number of lines or has width less than 42 characters.
func reserveLines(lines int) error {
	width, height, _ := termSize()
	if height < lines || width < 42 {
		return ErrTerminalTooSmall
	}
	for range lines {
		stdOutput.Write([]byte("\n"))
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

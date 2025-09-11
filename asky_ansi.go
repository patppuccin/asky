package asky

import (
	"strconv"
)

const (
	// ANSI escape codes for cursor visibility and position
	ansiHideCursor    = "\033[?25l"
	ansiShowCursor    = "\033[?25h"
	ansiSaveCursor    = "\033[s"
	ansiRestoreCursor = "\033[u"

	// ANSI escape codes for screen and line control
	ansiReset       = "\033[0m"
	ansiClearLine   = "\033[K"
	ansiClearScreen = "\033[J"
)

// Moves the cursor n positions left.
func ansiCursorLeft(n int) {
	if n > 0 {
		stdOutput.Write([]byte("\033[" + strconv.Itoa(n) + "D"))
	}
}

// Moves the cursor n positions up.
func ansiCursorUp(n int) {
	if n > 0 {
		stdOutput.Write([]byte("\033[" + strconv.Itoa(n) + "A"))
	}
}

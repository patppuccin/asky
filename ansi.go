package asky

import (
	"strconv"
)

const (
	ansiHideCursor    = "\033[?25l"
	ansiShowCursor    = "\033[?25h"
	ansiSaveCursor    = "\0337"
	ansiRestoreCursor = "\0338"

	ansiReset       = "\033[0m\033[0 q"
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

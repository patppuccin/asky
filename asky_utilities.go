package asky

import (
	"errors"
	"os"
	"strconv"
)

func init() {
	// if runtime.GOOS == "windows" {
	// 	EnableANSISupport()
	// }
}

// Custom Errors -------------------------------------------
var ErrInterrupted = errors.New("prompt interrupted")
var ErrNoOptions = errors.New("no options available")
var ErrTerminalTooSmall = errors.New("terminal dimensions too small")

// ANSI Escape Functions -----------------------------------
const (
	ansiHideCursor     = "\033[?25l"
	ansiShowCursor     = "\033[?25h"
	ansiSaveCursor     = "\033[s"
	ansiRestoreCursor  = "\033[u"
	ansiReset          = "\033[0m"
	ansiClearLineEnd   = "\033[K"
	ansiClearScreenEnd = "\033[J"
)

func cursorMoveLeft(n int) {
	if n > 0 {
		os.Stdout.Write([]byte("\033[" + strconv.Itoa(n) + "D"))
	}
}
func cursorMoveRight(n int) {
	if n > 0 {
		os.Stdout.Write([]byte("\033[" + strconv.Itoa(n) + "C"))
	}
}

// Helper Functions ---------------------------------------
// func getTerminalDimensions() (int, int, error) {
// 	width, height, err := term.GetSize(int(os.Stdout.Fd()))
// 	if err != nil {
// 		return 0, 0, err
// 	}
// 	return width, height, nil
// }

// // IsDumbTerminal returns true if output is redirected or TERM=dumb
// func IsDumbTerminal() bool {
// 	// check if stdout is a terminal
// 	fi, err := os.Stdout.Stat()
// 	if err != nil {
// 		return true
// 	}
// 	if (fi.Mode() & os.ModeCharDevice) == 0 {
// 		return true // not a terminal (redirected/piped)
// 	}

// 	// check TERM
// 	term := strings.ToLower(os.Getenv("TERM"))
// 	if term == "dumb" || term == "" {
// 		return true
// 	}

// 	return false
// }

// Common Utilities ----------------------------------------

// EnableANSISupport enables ANSI colors in Windows consoles (Windows 10+)
// func EnableANSISupport() {
// 	if runtime.GOOS != "windows" {
// 		return
// 	}
// 	const ENABLE_VIRTUAL_TERMINAL_PROCESSING = 0x0004

// 	kernel32 := syscall.NewLazyDLL("kernel32.dll")
// 	setConsoleMode := kernel32.NewProc("SetConsoleMode")
// 	getConsoleMode := kernel32.NewProc("GetConsoleMode")

// 	handle := syscall.Handle(os.Stdout.Fd())
// 	var mode uint32
// 	_, _, _ = getConsoleMode.Call(uintptr(handle), uintptr(unsafe.Pointer(&mode)))
// 	mode |= ENABLE_VIRTUAL_TERMINAL_PROCESSING
// 	_, _, _ = setConsoleMode.Call(uintptr(handle), uintptr(mode))
// }

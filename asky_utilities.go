package asky

// HEX Codes to ANSI Colors --------------------------------
// func Hex(hex, s string) string {
// 	var r, g, b int
// 	_, err := fmt.Sscanf(hex, "#%02x%02x%02x", &r, &g, &b)
// 	if err != nil {
// 		return s // fallback: no color if parse fails
// 	}
// 	return fmt.Sprintf("\033[38;2;%d;%d;%dm%s\033[0m", r, g, b, s)
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

// // EnableANSISupport enables ANSI colors in Windows consoles (Windows 10+)
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

// if runtime.GOOS == "windows" {
// 	EnableANSISupport()
// }

package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/patppuccin/asky"
)

// IsDumbTerminal returns true if output is redirected or TERM=dumb
func IsDumbTerminal() bool {
	// check if stdout is a terminal
	fi, err := os.Stdout.Stat()
	if err != nil {
		return true
	}
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		return true // not a terminal (redirected/piped)
	}

	// check TERM
	term := strings.ToLower(os.Getenv("TERM"))
	if term == "dumb" || term == "" {
		return true
	}

	return false
}

// EnableANSISupport enables ANSI colors in Windows consoles (Windows 10+)
func EnableANSISupport() {
	if runtime.GOOS != "windows" {
		return
	}
	const ENABLE_VIRTUAL_TERMINAL_PROCESSING = 0x0004

	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	setConsoleMode := kernel32.NewProc("SetConsoleMode")
	getConsoleMode := kernel32.NewProc("GetConsoleMode")

	handle := syscall.Handle(os.Stdout.Fd())
	var mode uint32
	_, _, _ = getConsoleMode.Call(uintptr(handle), uintptr(unsafe.Pointer(&mode)))
	mode |= ENABLE_VIRTUAL_TERMINAL_PROCESSING
	_, _, _ = setConsoleMode.Call(uintptr(handle), uintptr(mode))
}

func main() {

	// if runtime.GOOS == "windows" {
	// 	EnableANSISupport()
	// }

	// // fmt.Println("Hello, World!")

	// fname, err := asky.NewTextInput().
	// 	WithPromptText("Please enter your first name").
	// 	WithDefault("John").
	// 	WithHelper("Enter the first name you want to use").
	// 	WithSeparator(": ").
	// 	WithTheme(asky.ThemeOsakaJade).
	// 	Render()
	// if err != nil {
	// 	if errors.Is(err, asky.ErrInterrupted) {
	// 		fmt.Println("Input Cancelled")
	// 		return
	// 	}
	// 	fmt.Println("Error: " + err.Error())
	// }

	// lname, err := asky.NewTextInput().
	// 	WithPromptText("Please enter your last name").
	// 	WithHelper("Enter the last name you want to use").
	// 	WithTheme(asky.ThemeOsakaJade).
	// 	Render()
	// if err != nil {
	// 	if errors.Is(err, asky.ErrInterrupted) {
	// 		fmt.Println("Input Cancelled")
	// 		return
	// 	}
	// 	fmt.Println("Error: " + err.Error())
	// }

	// fmt.Println("User's Name: " + fname + " " + lname)

	// sf := asky.NewSpinner().WithFrames(asky.SpinnerPatternDots).WithTheme(asky.ThemeOsakaJade)
	// sf.Start("Telling jokes...")
	// time.Sleep(3 * time.Second) // simulate work
	// sf.Stop()                   // or s.Stop(false, "Failed")

	// pronouns, err := asky.NewTextInput().
	// 	WithPromptText("What are your pronouns?").
	// 	WithDefault("he/him").
	// 	WithHelper("Enter your preferred pronouns").
	// 	WithSeparator(": ").
	// 	WithTheme(asky.ThemeOsakaJade).
	// 	Render()
	// if err != nil {
	// 	if errors.Is(err, asky.ErrInterrupted) {
	// 		fmt.Println("Input Cancelled")
	// 		return
	// 	}
	// 	fmt.Println("Error: " + err.Error())
	// }

	// fmt.Println("Pronouns: " + pronouns)

	// s := asky.NewSpinner().
	// 	WithLabelText("Petting Cats...").
	// 	WithHelperText("This may take a while...").
	// 	WithTheme(asky.ThemeDefault)
	// s.Start()
	// time.Sleep(3 * time.Second) // simulate work
	// s.Stop()                    // or s.Stop(false, "Failed")

	// ok, _ := asky.NewConfirm().
	// 	WithPromptText("Proceed with deployment?").
	// 	WithHelperText("This action is irreversible").
	// 	WithDefaultOption(false).
	// 	WithTheme(asky.ThemeDefault).
	// 	Render()

	// if ok {
	// 	fmt.Println("Proceeding with the deployment...")
	// } else {
	// 	fmt.Println("Deployment cancelled")
	// }

	pb := asky.NewProgress().
		WithWidth(30).
		WithPattern(asky.ProgressPatternMathSymbols).
		WithProgressText("Uploading files").
		WithHelperText("This may take a while...").
		WithTheme(asky.ThemeDefault).
		WithTotalSteps(100)

	pb.Start()

	for range 100 {
		time.Sleep(100 * time.Millisecond)
		pb.Increment()
	}

	pb.Stop(false)
	fmt.Println(asky.ThemeDefault.SuccessStyle("[+]") + " Done!")

}

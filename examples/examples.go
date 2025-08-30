package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/patppuccin/asky"
)

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

func main() {

	fmt.Println(asky.ThemeDefault.InfoStyle("\n# Demo of the asky library ----------------\n"))

	// Text Input
	fname, err := asky.NewTextInput().
		WithPromptText("Please enter your first name").
		WithDefault("James").
		WithHelper("First name is the name that comes first").
		WithSeparator(": ").
		WithTheme(asky.ThemeDefault).
		Render()
	if err != nil {
		if errors.Is(err, asky.ErrInterrupted) {
			fmt.Println("Input Cancelled")
			return
		}
		fmt.Println("Error: " + err.Error())
	}

	lname, err := asky.NewTextInput().
		WithPromptText("Please enter your last name").
		WithHelper("Enter the last name you want to use").
		WithTheme(asky.ThemeDefault).
		Render()
	if err != nil {
		if errors.Is(err, asky.ErrInterrupted) {
			fmt.Println("Input Cancelled")
			return
		}
		fmt.Println("Error: " + err.Error())
	}

	fmt.Println(
		asky.ThemeDefault.SuccessStyle("[+]"),
		asky.ThemeDefault.PrimaryStyle("User's Name:"),
		asky.ThemeDefault.AccentStyle(fname+" "+lname),
	)

	sf := asky.NewSpinner().
		WithLabelText("Asking permission from mom...").
		WithHelperText("Wheew! this could take a while").
		WithFrames(asky.SpinnerPatternDots).
		WithTheme(asky.ThemeDefault)
	sf.Start()
	time.Sleep(3 * time.Second) // simulate work
	sf.Stop()                   // or s.Stop(false, "Failed")
	fmt.Println(
		asky.ThemeDefault.SuccessStyle("[+]"),
		asky.ThemeDefault.PrimaryStyle("Yay! Permission granted"),
	)

	ok, _ := asky.NewConfirm().
		WithPromptText("Ready to go out?").
		WithHelperText("This fun is irreversible").
		WithDefaultOption(false).
		WithTheme(asky.ThemeDefault).
		Render()

	if ok {
		fmt.Println(
			asky.ThemeDefault.SuccessStyle("[+]"),
			asky.ThemeDefault.PrimaryStyle("Let's gooooooooooooooooo!"),
		)

	} else {
		fmt.Println(
			asky.ThemeDefault.ErrorStyle("[-]"),
			asky.ThemeDefault.PrimaryStyle("Why? Why have you forsaken me?"),
		)
	}

	passwordStrengthValidator := func(pw string) (string, bool) {
		switch {
		case len(pw) == 0:
			return "Input cannot be empty", false
		case len(pw) < 4:
			return "weak buddy", false
		case len(pw) < 8:
			return "I mean, it's alriiii", false
		default:
			return "that's why he's the GOAT, the GOAT", true
		}
	}

	pwd, err := asky.NewSecureInput().
		WithPromptText("Please set a ultra-secure password").
		WithHelper("Something that is unique & not about pizza").
		WithTheme(asky.ThemeDefault).
		WithValidator(passwordStrengthValidator).
		Render()
	if err != nil {
		if errors.Is(err, asky.ErrInterrupted) {
			fmt.Println(
				asky.ThemeDefault.ErrorStyle("[-]"),
				asky.ThemeDefault.PrimaryStyle("Password Cancelled"),
			)
			return
		} else {
			fmt.Println(
				asky.ThemeDefault.ErrorStyle("[-]"),
				asky.ThemeDefault.PrimaryStyle("Error: "+err.Error()),
			)
		}
	} else {
		fmt.Println(
			asky.ThemeDefault.SuccessStyle("[+]"),
			asky.ThemeDefault.PrimaryStyle("Ultra-secure Password: "+pwd),
		)
	}

	pb := asky.NewProgress().
		WithWidth(30).
		WithPattern(asky.ProgressPatternMathSymbols).
		WithProgressText("Waiting for the pizza...").
		WithHelperText("Getting Ready to go out").
		WithTheme(asky.ThemeDefault).
		WithTotalSteps(100)

	pb.Start()

	for range 100 {
		time.Sleep(50 * time.Millisecond)
		pb.Increment()
	}

	pb.Stop(false)
	fmt.Println(
		asky.ThemeDefault.SuccessStyle("[+]"),
		asky.ThemeDefault.PrimaryStyle("Yay! (Hot) Pizza is delivered"),
	)

	fmt.Printf(asky.ThemeDefault.InfoStyle("\n# Demo Completed ----------------\n\n\n"))

}

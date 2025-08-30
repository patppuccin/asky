package main

import (
	"errors"
	"fmt"
	"regexp"
	"time"
	"unicode"

	"github.com/patppuccin/asky"
)

func passwordStrengthValidator(pw string) (string, bool) {
	if len(pw) == 0 {
		return "empty", false
	}
	if len(pw) < 8 {
		return "too short", false
	}

	var (
		hasUpper, hasLower, hasNumber, hasSpecial bool
	)

	for _, c := range pw {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsPunct(c), unicode.IsSymbol(c):
			hasSpecial = true
		}
	}

	onlyOneCharType, _ := regexp.MatchString(`^(.)\1+$`, pw)
	if onlyOneCharType {
		return "repetitive", false
	}

	score := 0
	if hasUpper {
		score++
	}
	if hasLower {
		score++
	}
	if hasNumber {
		score++
	}
	if hasSpecial {
		score++
	}

	switch {
	case score < 2:
		return "weak", false
	case score == 2:
		return "ok", true
	case score >= 3:
		return "strong", true
	}

	return "invalid", false
}

func main() {

	fmt.Printf(asky.ThemeDefault.InfoStyle("\n# Showcasing the asky library ---------------\n"))

	// --- Prompt Showcase: Text Input --------------------------
	fname, err := asky.NewTextInput().
		WithPromptText("Enter first name").
		WithDefault("James").
		WithHelper("This name is used as a greet name").
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

	// --- Prompt Showcase: Text Input -------------------------
	lname, err := asky.NewTextInput().
		WithPromptText("Enter last name").
		WithHelper("This is optional").
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

	fmt.Println()
	fmt.Println(
		asky.ThemeDefault.SuccessStyle("[+]"),
		asky.ThemeDefault.PrimaryStyle("User's Name:"),
		asky.ThemeDefault.AccentStyle(fname+" "+lname),
	)

	// --- Prompt Showcase: Secure Input -----------------------

	pwd, err := asky.NewSecureInput().
		WithPromptText("Enter a secure password").
		WithHelper("Use a mix of letters, numbers & symbols").
		WithSeparator(": ").
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
			asky.ThemeDefault.PrimaryStyle(fmt.Sprintf("Password is %d characters long", len(pwd))),
		)
	}

	// --- Prompt Showcase: Confirmation -----------------------
	ok, _ := asky.NewConfirm().
		WithPromptText("Create account with username " + fname + "?").
		WithHelperText("You will be able to change this later").
		WithDefaultOption(false).
		WithTheme(asky.ThemeDefault).
		Render()

	if ok {
		fmt.Println(
			asky.ThemeDefault.SuccessStyle("[+]"),
			asky.ThemeDefault.PrimaryStyle("User created"),
		)

	} else {
		fmt.Println(
			asky.ThemeDefault.ErrorStyle("[-]"),
			asky.ThemeDefault.PrimaryStyle("It's okay, you can try again later"),
		)
	}

	// --- Indicator: Loading Spinner --------------------------
	sp := asky.NewSpinner().
		WithLabelText("Getting things ready...").
		WithHelperText("Fetching data from the deep web").
		WithFrames(asky.SpinnerPatternDots).
		WithTheme(asky.ThemeDefault)
	sp.Start()
	time.Sleep(3 * time.Second) // simulate work
	sp.Stop()                   // or s.Stop(false, "Failed")
	fmt.Println(
		asky.ThemeDefault.SuccessStyle("[+]"),
		asky.ThemeDefault.PrimaryStyle("All things are ready"),
	)

	// --- Indicator: Progress Bar -----------------------------
	pb := asky.NewProgress().
		WithWidth(30).
		WithPattern(asky.ProgressPatternMathSymbols).
		WithProgressText("Running preparations...").
		WithHelperText("This will take a while").
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
		asky.ThemeDefault.PrimaryStyle("All preparations are done"),
	)

	fmt.Printf(asky.ThemeDefault.InfoStyle("\n# Showcase Completed ------------------------\n\n"))
}

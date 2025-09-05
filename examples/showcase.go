package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/patppuccin/asky"
)

func passwordStrengthValidator(pw string) (string, bool) {
	// Step 1: Trivial cases
	if len(pw) == 0 {
		return "empty", false
	}
	if len(pw) < 8 {
		return "too short", false
	}

	// Step 2: Scan character classes
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

	// Step 3: score
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

	// Step 4: map score to labels
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

func firstNameValidator(name string) (string, bool) {
	name = strings.TrimSpace(strings.ToLower(name))
	if len(name) == 0 {
		return "empty", false
	}
	if name == "james" {
		return "James'es are a bad name", false
	}
	return "valid", true
}

func ageValidator(age string) (string, bool) {
	age = strings.TrimSpace(age)
	ageNum, err := strconv.Atoi(age)
	if err != nil {
		return "invalid age", false
	}
	if ageNum < 18 {
		return "too young", false
	}
	if ageNum > 100 {
		return "too old", false
	}
	return "valid", true
}

func main() {

	askyTheme := asky.ThemeCatppuccinMocha

	// asky.NewStatus().WithLabel("Welcome to Asky").WithLevel(asky.StatusLevelInfo).Render()
	// asky.NewStatus().WithLabel("This is a success message").WithLevel(asky.StatusLevelSuccess).Render()
	// asky.NewStatus().WithLabel("This is a debug message").WithLevel(asky.StatusLevelDebug).Render()
	// asky.NewStatus().WithLabel("This is a status message").WithLevel(asky.StatusLevelInfo).Render()
	// asky.NewStatus().WithLabel("This is a warning message").WithLevel(asky.StatusLevelWarn).Render()
	// asky.NewStatus().WithLabel("This is an error message").WithLevel(asky.StatusLevelError).Render()

	// --- Prompt Showcase: Text Input --------------------------
	fname, err := asky.NewTextInput().
		WithLabel("Enter first name:").
		WithPlaceholder("John").
		WithDescription("This name is used as a greet name").
		WithValidator(firstNameValidator).
		WithTheme(askyTheme).
		Render()
	if err != nil {
		if errors.Is(err, asky.ErrInterrupted) {
			asky.NewStatus().WithLabel("User cancelled input").WithLevel(asky.StatusLevelInfo).Render()
		} else {
			asky.NewStatus().WithLabel("Error: " + err.Error()).WithLevel(asky.StatusLevelError).Render()
		}
	}

	// --- Prompt Showcase: Text Input -------------------------
	lname, err := asky.NewTextInput().
		WithLabel("Enter last name").
		WithDescription("This is optional").
		WithTheme(askyTheme).
		Render()
	if err != nil {
		if errors.Is(err, asky.ErrInterrupted) {
			asky.NewStatus().WithLabel("User cancelled input").WithLevel(asky.StatusLevelInfo).Render()
		} else {
			asky.NewStatus().WithLabel("Error: " + err.Error()).WithLevel(asky.StatusLevelError).Render()
		}
	}

	fmt.Println()
	asky.NewStatus().WithLabel("User's Name: " + fname + " " + lname).WithLevel(asky.StatusLevelSuccess).Render()

	// // --- Prompt Showcase: Number Input -----------------------
	// age, err := asky.NewNumberInput().
	// 	WithPromptText("Enter your age").
	// 	WithHelpText("This is a number input").
	// 	WithPromptSeparator(": ").
	// 	WithTheme(askyTheme).
	// 	WithValidator(ageValidator).
	// 	Render()
	// if err != nil {
	// 	if errors.Is(err, asky.ErrInterrupted) {
	// 		fmt.Println("Input Cancelled")
	// 		return
	// 	}
	// 	fmt.Println("Error: " + err.Error())
	// }

	// fmt.Println(
	// 	askyTheme.SuccessStyle("[+]"),
	// 	askyTheme.PrimaryStyle("User's Age:"),
	// 	askyTheme.AccentStyle(age),
	// )

	// // --- Prompt Showcase: Secure Input -----------------------

	// pwd, err := asky.NewSecureInput().
	// 	WithPromptText("Enter a secure password").
	// 	WithHelpText("Use a mix of letters, numbers & symbols").
	// 	WithPromptSeparator(": ").
	// 	WithTheme(askyTheme).
	// 	WithValidator(passwordStrengthValidator).
	// 	Render()
	// if err != nil {
	// 	if errors.Is(err, asky.ErrInterrupted) {
	// 		fmt.Println(
	// 			askyTheme.ErrorStyle("[-]"),
	// 			askyTheme.PrimaryStyle("Password Cancelled"),
	// 		)
	// 		return
	// 	} else {
	// 		fmt.Println(
	// 			askyTheme.ErrorStyle("[-]"),
	// 			askyTheme.PrimaryStyle("Error: "+err.Error()),
	// 		)
	// 	}
	// } else {
	// 	fmt.Println(
	// 		askyTheme.SuccessStyle("[+]"),
	// 		askyTheme.PrimaryStyle(fmt.Sprintf("Password is %d characters long", len(pwd))),
	// 	)
	// }

	// // --- Prompt Showcase: Select -----------------------------
	faveLang, err := asky.NewSingleSelect().
		WithLabel("Pick the favourite language").
		WithDescription("This is used to tailor the recommendations").
		WithPageSize(6).
		WithChoices([]asky.Choice{
			// Systems / Low-level
			{Value: "c", Label: "C"},
			{Value: "cpp", Label: "C++"},
			{Value: "rs", Label: "Rust"},
			{Value: "zig", Label: "Zig"},

			// General-purpose / OO heavyweights
			{Value: "java", Label: "Java", Disabled: true},
			{Value: "cs", Label: "C#"},
			{Value: "kt", Label: "Kotlin", Disabled: true},
			{Value: "swift", Label: "Swift"},
			{Value: "scala", Label: "Scala"},

			// Scripting / Dynamic
			{Value: "py", Label: "Python"},
			{Value: "js", Label: "JavaScript / TypeScript"},
			{Value: "php", Label: "PHP"},
			{Value: "rb", Label: "Ruby"},
			{Value: "perl", Label: "Perl"},

			// Functional / Multi-paradigm
			{Value: "hs", Label: "Haskell"},
			{Value: "clj", Label: "Clojure"},
			{Value: "erl", Label: "Erlang"},
			{Value: "elx", Label: "Elixir"},
			{Value: "fsharp", Label: "F#"},
			{Value: "ml", Label: "OCaml"},

			// Emerging / Modern
			{Value: "go", Label: "Go"},
			{Value: "dart", Label: "Dart"},
			{Value: "nim", Label: "Nim"},
		}).
		WithTheme(askyTheme).
		Render()
	if err != nil {
		if errors.Is(err, asky.ErrInterrupted) {
			asky.NewStatus().WithLabel("User cancelled selection").WithLevel(asky.StatusLevelInfo).Render()
			// return
		}
		asky.NewStatus().WithLabel("Error: " + err.Error()).WithLevel(asky.StatusLevelError).Render()
	} else {
		asky.NewStatus().WithLabel("User's Favourite Languages: " + faveLang.Label).WithLevel(asky.StatusLevelSuccess).Render()
		// return
	}

	// --- Prompt Showcase: Confirmation -----------------------
	ok, _ := asky.NewConfirm().
		WithLabel("Create a brand new account ?").
		WithDescription("You will be able to change this later").
		WithDefaultAnswer(false).
		WithTheme(askyTheme).
		Render()

	if ok {
		asky.NewStatus().WithLabel("User created successfully").WithLevel(asky.StatusLevelSuccess).Render()
	} else {
		asky.NewStatus().WithLabel("Account creation skipped. You can try again later").WithLevel(asky.StatusLevelInfo).Render()
	}

	// --- Indicator: Loading Spinner --------------------------
	sp := asky.NewSpinner().
		WithLabel("Getting things ready...").
		WithDescription("Fetching data from the deep web").
		WithFrames(asky.SpinnerPatternDots).
		WithTheme(askyTheme)
	sp.Start()
	time.Sleep(4 * time.Second) // simulate work
	sp.Stop()                   // or s.Stop(false, "Failed")
	asky.NewStatus().WithLabel("All things are ready").WithLevel(asky.StatusLevelSuccess).Render()

	// // --- Indicator: Progress Bar -----------------------------
	pb := asky.NewProgress().
		WithWidth(30).
		WithPattern(asky.ProgressPatternDefault).
		WithLabel("Running Preparations").
		WithDescription("This will take a while").
		WithSteps(100).
		WithTheme(askyTheme)

	pb.Start()
	for range 100 {
		time.Sleep(50 * time.Millisecond)
		pb.Increment()
	}

	pb.Stop()
	asky.NewStatus().WithLabel("All preparations are done").WithLevel(asky.StatusLevelSuccess).Render()

	fmt.Print(asky.NewStyle().FG(askyTheme.Foreground).Bold().Sprint("\n# Showcase Completed ------------------------\n\n"))
}

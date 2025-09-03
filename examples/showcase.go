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

	faveLang, err := asky.NewSingleSelect().
		WithLabel("Pick the favourite language").
		WithHelp("This is used to tailor the recommendations").
		WithSeparator("").
		WithPageSize(25).
		WithChoices([]asky.Choice{
			// Systems / Low-level
			{Value: "c", Label: "C"},
			{Value: "cpp", Label: "C++"},
			{Value: "rs", Label: "Rust"},
			{Value: "zig", Label: "Zig"},

			// General-purpose / OO heavyweights
			{Value: "java", Label: "Java", Description: "Java is a general-purpose, class-based, object-oriented programming language"},
			{Value: "cs", Label: "C#"},
			{Value: "kt", Label: "Kotlin"},
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
		WithTheme(asky.ThemeDefault).
		Render()
	if err != nil {
		if errors.Is(err, asky.ErrInterrupted) {
			fmt.Println("Input Cancelled")
			return
		}
		fmt.Println("Error: " + err.Error())
	} else {
		fmt.Println(
			asky.ThemeDefault.SuccessStyle("[+]"),
			asky.ThemeDefault.PrimaryStyle("User's Favourite Languages:"),
			asky.ThemeDefault.AccentStyle(faveLang.Label),
		)
		return
	}
	fmt.Print(asky.ThemeDefault.InfoStyle("\n# Showcasing the asky library ---------------\n"))

	// --- Prompt Showcase: Text Input --------------------------
	fname, err := asky.NewTextInput().
		WithPromptText("Enter first name").
		WithDefaultValue("John").
		WithHelpText("This name is used as a greet name").
		WithPromptSeparator(": ").
		WithValidator(firstNameValidator).
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
		WithHelpText("This is optional").
		WithPromptSeparator(": ").
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

	// --- Prompt Showcase: Number Input -----------------------
	age, err := asky.NewNumberInput().
		WithPromptText("Enter your age").
		WithHelpText("This is a number input").
		WithPromptSeparator(": ").
		WithTheme(asky.ThemeDefault).
		WithValidator(ageValidator).
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
		asky.ThemeDefault.PrimaryStyle("User's Age:"),
		asky.ThemeDefault.AccentStyle(age),
	)

	// --- Prompt Showcase: Secure Input -----------------------

	pwd, err := asky.NewSecureInput().
		WithPromptText("Enter a secure password").
		WithHelpText("Use a mix of letters, numbers & symbols").
		WithPromptSeparator(": ").
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

	// --- Prompt Showcase: Select -----------------------------
	// faveLang, err := asky.NewSingleSelect().
	// 	WithLabel("Pick the favourite language").
	// 	WithHelp("This is used to tailor the recommendations").
	// 	WithSeparator("").
	// 	WithPageSize(5).
	// 	WithChoices([]asky.Choice{
	// 		// Systems / Low-level
	// 		{Value: "c", Label: "C"},
	// 		{Value: "cpp", Label: "C++"},
	// 		{Value: "rs", Label: "Rust"},
	// 		{Value: "zig", Label: "Zig"},

	// 		// General-purpose / OO heavyweights
	// 		{Value: "java", Label: "Java"},
	// 		{Value: "cs", Label: "C#"},
	// 		{Value: "kt", Label: "Kotlin"},
	// 		{Value: "swift", Label: "Swift"},
	// 		{Value: "scala", Label: "Scala"},

	// 		// Scripting / Dynamic
	// 		{Value: "py", Label: "Python"},
	// 		{Value: "js", Label: "JavaScript / TypeScript"},
	// 		{Value: "php", Label: "PHP"},
	// 		{Value: "rb", Label: "Ruby"},
	// 		{Value: "perl", Label: "Perl"},

	// 		// Functional / Multi-paradigm
	// 		{Value: "hs", Label: "Haskell"},
	// 		{Value: "clj", Label: "Clojure"},
	// 		{Value: "erl", Label: "Erlang"},
	// 		{Value: "elx", Label: "Elixir"},
	// 		{Value: "fsharp", Label: "F#"},
	// 		{Value: "ml", Label: "OCaml"},

	// 		// Emerging / Modern
	// 		{Value: "go", Label: "Go"},
	// 		{Value: "dart", Label: "Dart"},
	// 		{Value: "nim", Label: "Nim"},
	// 	}).
	// 	WithTheme(asky.ThemeDefault).
	// 	Render()
	// if err != nil {
	// 	if errors.Is(err, asky.ErrInterrupted) {
	// 		fmt.Println("Input Cancelled")
	// 		return
	// 	}
	// 	fmt.Println("Error: " + err.Error())
	// }

	// fmt.Println(
	// 	asky.ThemeDefault.SuccessStyle("[+]"),
	// 	asky.ThemeDefault.PrimaryStyle("User's Favourite Languages:"),
	// 	asky.ThemeDefault.AccentStyle(faveLang.Label),
	// )

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

	fmt.Print(asky.ThemeDefault.InfoStyle("\n# Showcase Completed ------------------------\n\n"))
}

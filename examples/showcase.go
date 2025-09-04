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

	pb2 := asky.NewProgress().
		WithWidth(30).
		WithPattern(asky.ProgressPatternDefault).
		WithLabel("Running Preparations").
		WithHelp("This will take a while").
		WithSteps(100).
		WithTheme(askyTheme)

	pb2.Start()

	for range 100 {
		time.Sleep(50 * time.Millisecond)
		pb2.Increment()
	}

	pb2.Stop(false)
	fmt.Println(
		askyTheme.SuccessStyle("[+]"),
		askyTheme.NeutralStyle("All preparations are done"),
	)

	var err error
	if err == nil {
		return
	}

	// sp2 := asky.NewSpinner().
	// 	WithLabel("Getting things ready...").
	// 	WithHelp("Fetching data from the deep web").
	// 	WithFrames(asky.SpinnerPatternDots).
	// 	WithTheme(askyTheme)
	// sp2.Start()
	// time.Sleep(5 * time.Second) // simulate work
	// sp2.Stop()                  // or s.Stop(false, "Failed")
	// fmt.Println(
	// 	askyTheme.SuccessStyle("[+]"),
	// 	askyTheme.NeutralStyle("All things are ready"),
	// )

	// --- Prompt Showcase: Text Input --------------------------
	fname, err := asky.NewTextInput().
		WithPromptText("Enter first name").
		WithDefaultValue("John").
		WithHelpText("This name is used as a greet name").
		WithPromptSeparator(": ").
		WithValidator(firstNameValidator).
		WithTheme(askyTheme).
		Render()
	if err != nil {
		if errors.Is(err, asky.ErrInterrupted) {
			fmt.Println("Input Cancelled")
		}
		fmt.Println("Error: " + err.Error())
	}

	// --- Prompt Showcase: Text Input -------------------------
	lname, err := asky.NewTextInput().
		WithPromptText("Enter last name").
		WithHelpText("This is optional").
		WithPromptSeparator(": ").
		WithTheme(askyTheme).
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
		askyTheme.SuccessStyle("[+]"),
		askyTheme.PrimaryStyle("User's Name:"),
		askyTheme.AccentStyle(fname+" "+lname),
	)

	// --- Prompt Showcase: Number Input -----------------------
	age, err := asky.NewNumberInput().
		WithPromptText("Enter your age").
		WithHelpText("This is a number input").
		WithPromptSeparator(": ").
		WithTheme(askyTheme).
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
		askyTheme.SuccessStyle("[+]"),
		askyTheme.PrimaryStyle("User's Age:"),
		askyTheme.AccentStyle(age),
	)

	// --- Prompt Showcase: Secure Input -----------------------

	pwd, err := asky.NewSecureInput().
		WithPromptText("Enter a secure password").
		WithHelpText("Use a mix of letters, numbers & symbols").
		WithPromptSeparator(": ").
		WithTheme(askyTheme).
		WithValidator(passwordStrengthValidator).
		Render()
	if err != nil {
		if errors.Is(err, asky.ErrInterrupted) {
			fmt.Println(
				askyTheme.ErrorStyle("[-]"),
				askyTheme.PrimaryStyle("Password Cancelled"),
			)
			return
		} else {
			fmt.Println(
				askyTheme.ErrorStyle("[-]"),
				askyTheme.PrimaryStyle("Error: "+err.Error()),
			)
		}
	} else {
		fmt.Println(
			askyTheme.SuccessStyle("[+]"),
			askyTheme.PrimaryStyle(fmt.Sprintf("Password is %d characters long", len(pwd))),
		)
	}

	// --- Prompt Showcase: Select -----------------------------
	faveLang, err := asky.NewSingleSelect().
		WithLabel("Pick the favourite language").
		WithHelp("This is used to tailor the recommendations").
		WithSeparator("").
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
			fmt.Println("Input Cancelled")
			// return
		}
		fmt.Println("Error: " + err.Error())
	} else {
		fmt.Println(
			askyTheme.SuccessStyle("[+]"),
			askyTheme.PrimaryStyle("User's Favourite Languages:"),
			askyTheme.AccentStyle(faveLang.Label),
		)
		// return
	}

	// --- Prompt Showcase: Confirmation -----------------------
	ok, _ := asky.NewConfirm().
		WithPromptText("Create account with username " + fname + "?").
		WithHelperText("You will be able to change this later").
		WithDefaultOption(false).
		WithTheme(askyTheme).
		Render()

	if ok {
		fmt.Println(
			askyTheme.SuccessStyle("[+]"),
			askyTheme.PrimaryStyle("User created"),
		)

	} else {
		fmt.Println(
			askyTheme.ErrorStyle("[-]"),
			askyTheme.PrimaryStyle("It's okay, you can try again later"),
		)
	}

	// --- Indicator: Loading Spinner --------------------------
	sp := asky.NewSpinner().
		WithLabel("Getting things ready...").
		WithHelp("Fetching data from the deep web").
		WithFrames(asky.SpinnerPatternDots).
		WithTheme(askyTheme)
	sp.Start()
	time.Sleep(3 * time.Second) // simulate work
	sp.Stop()                   // or s.Stop(false, "Failed")
	fmt.Println(
		askyTheme.SuccessStyle("[+]"),
		askyTheme.PrimaryStyle("All things are ready"),
	)

	// --- Indicator: Progress Bar -----------------------------
	pb := asky.NewProgress().
		WithWidth(30).
		WithPattern(asky.ProgressPatternDefault).
		WithLabel("Running preparations...").
		WithHelp("This will take a while").
		WithSteps(100).
		WithTheme(askyTheme)

	pb.Start()

	for range 100 {
		time.Sleep(50 * time.Millisecond)
		pb.Increment()
	}

	pb.Stop(false)
	fmt.Println(
		askyTheme.SuccessStyle("[+]"),
		askyTheme.NeutralStyle("All preparations are done"),
	)

	fmt.Print(askyTheme.InfoStyle("\n# Showcase Completed ------------------------\n\n"))
}

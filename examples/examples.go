package main

import (
	"fmt"

	"github.com/patppuccin/asky"
)

func main() {

	// fname, err := asky.NewTextInput().
	// 	WithPromptText("Please enter your first name").
	// 	WithDefault("John").
	// 	WithHelper("Enter the first name you want to use").
	// 	WithSeparator(": ").
	// 	Render()
	// if err != nil {
	// 	if errors.Is(err, asky.ErrInterrupted) {
	// 		fmt.Println(color.RedString("Input Cancelled"))
	// 		return
	// 	}
	// 	fmt.Println(color.YellowString("Error: ") + err.Error())
	// }

	// lname, err := asky.NewTextInput().
	// 	WithPromptText("Please enter your last name").
	// 	WithHelper("Enter the last name you want to use").
	// 	Render()
	// if err != nil {
	// 	if errors.Is(err, asky.ErrInterrupted) {
	// 		fmt.Println(color.RedString("Input Cancelled"))
	// 		return
	// 	}
	// 	fmt.Println(color.YellowString("Error: ") + err.Error())
	// }

	// fmt.Println("User's Name: " + fname + " " + lname)

	// s := asky.NewSpinner().WithFrames(asky.SpinnerPatternDots)
	// s.Start("Petting Cats...")
	// time.Sleep(3 * time.Second) // simulate work
	// s.Stop()                    // or s.Stop(false, "Failed")

	ok, _ := asky.NewConfirm().
		WithPromptText("Proceed with deployment?").
		WithHelper("This action is irreversible").
		WithDefault(true).
		Render()

	if ok {
		fmt.Println("Proceeding with the deployment...")
	} else {
		fmt.Println("Deployment cancelled")
	}

}

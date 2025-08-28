package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/patppuccin/asky"
)

// func Spinner(ctx context.Context, text string, frames []rune) {
// 	if len(frames) == 0 {
// 		frames = []rune{'⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'}
// 	}

// 	ticker := time.NewTicker(100 * time.Millisecond)
// 	defer ticker.Stop()

// 	i := 0
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			fmt.Print("\r\033[2K")
// 			return
// 		case <-ticker.C:
// 			fmt.Printf("\r%s %s", color.MagentaString("[%c]", frames[i]), text)
// 			i = (i + 1) % len(frames)
// 		}
// 	}
// }

func main() {

	fname, err := asky.NewTextInput().
		WithPromptText("Please enter your first name").
		WithDefault("John").
		WithHelper("Enter the first name you want to use").
		WithSeparator(": ").
		Render()
	if err != nil {
		if errors.Is(err, asky.ErrInterrupted) {
			fmt.Println(color.RedString("Input Cancelled"))
			return
		}
		fmt.Println(color.YellowString("Error: ") + err.Error())
	}

	lname, err := asky.NewTextInput().
		WithPromptText("Please enter your last name").
		WithHelper("Enter the last name you want to use").
		Render()
	if err != nil {
		if errors.Is(err, asky.ErrInterrupted) {
			fmt.Println(color.RedString("Input Cancelled"))
			return
		}
		fmt.Println(color.YellowString("Error: ") + err.Error())
	}

	fmt.Println("User's Name: " + fname + " " + lname)

	s := asky.NewSpinner()
	s.Start("Processing...")
	time.Sleep(3 * time.Second) // simulate work
	s.Stop()                    // or s.Stop(false, "Failed")

}

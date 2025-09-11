package asky

import (
	"os"

	"golang.org/x/term"
)

func getTermDimensions() (int, int, error) {
	termWidth, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || termWidth <= 0 {
		termWidth = 80
	}
	return termWidth, termWidth, err
}

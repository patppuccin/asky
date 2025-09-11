package asky

import (
	"os"

	"golang.org/x/term"
)

func getTermDimensions() (int, int, error) {
	return term.GetSize(int(os.Stdout.Fd()))
}

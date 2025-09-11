package asky

import (
	"os"

	"golang.org/x/term"
)

func getTermDimensions() (int, int, error) {
	return term.GetSize(int(os.Stdout.Fd()))
}

func makeSpace(lines int) error {
	width, height, _ := getTermDimensions()
	if height < lines || width < 50 {
		return ErrTerminalTooSmall
	}
	for range lines {
		os.Stdout.WriteString("\n")
	}
	ansiCursorUp(lines)
	return nil
}

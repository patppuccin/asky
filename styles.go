package asky

import (
	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
)

// stdOutput is the colorable stdout used by all Asky components.
// On Windows, this ensures ANSI escape sequences render correctly.
var stdOutput = colorable.NewColorableStdout()

// StyleMap defines the visual appearance of all Asky TUI components.
// Every field is a [*color.Color] from the fatih/color package — assign
// any value constructed with [color.New] to override a specific style.
//
// As a good practice, always obtain a StyleMap via [NewStyles] — do not
// construct it directly. Unset fields will be unstyled at runtime.
//
//	styles := asky.NewStyles()
//	styles.InputPrefix = color.New(color.FgMagenta, color.Bold)
type StyleMap struct {
	// Log message styles.
	LogSuccessPrefix *color.Color
	LogSuccessLabel  *color.Color
	LogDebugPrefix   *color.Color
	LogDebugLabel    *color.Color
	LogInfoPrefix    *color.Color
	LogInfoLabel     *color.Color
	LogWarnPrefix    *color.Color
	LogWarnLabel     *color.Color
	LogErrorPrefix   *color.Color
	LogErrorLabel    *color.Color
	LogGroupBody     *color.Color

	// Input prompt styles.
	InputPrefix         *color.Color
	InputLabel          *color.Color
	InputPlaceholder    *color.Color
	InputText           *color.Color
	InputValidationPass *color.Color
	InputValidationFail *color.Color
	InputHelp           *color.Color

	// Confirmation prompt styles.
	ConfirmationPrefix *color.Color
	ConfirmationLabel  *color.Color
	ConfirmationHelp   *color.Color

	// Selection prompt styles.
	SelectionPrefix             *color.Color
	SelectionLabel              *color.Color
	SelectionDesc               *color.Color
	SelectionHelp               *color.Color
	SelectionSearchLabel        *color.Color
	SelectionSearchText         *color.Color
	SelectionSearchHint         *color.Color
	SelectionValidationPass     *color.Color
	SelectionValidationFail     *color.Color
	SelectionItemNormalMarker   *color.Color
	SelectionItemNormalLabel    *color.Color
	SelectionItemCurrentMarker  *color.Color
	SelectionItemCurrentLabel   *color.Color
	SelectionItemSelectedMarker *color.Color
	SelectionItemSelectedLabel  *color.Color
	SelectionItemDisabledMarker *color.Color
	SelectionItemDisabledLabel  *color.Color

	// Spinner styles.
	SpinnerPrefix *color.Color
	SpinnerLabel  *color.Color

	// Progress bar styles.
	ProgressPrefix     *color.Color
	ProgressLabel      *color.Color
	ProgressBarPad     *color.Color
	ProgressBarDone    *color.Color
	ProgressBarPending *color.Color
	ProgressBarStatus  *color.Color
}

// NewStyles returns a [StyleMap] with sensible default colors.
//
// The palette uses sharp and distinctive colors with semantic states
// such as green for success, yellow for warnings, red for errors,
// blue for info, and dark gray for muted/dimmed elements.
func NewStyles() *StyleMap {
	return &StyleMap{
		// Log messages
		LogSuccessPrefix: color.New(color.FgGreen),
		LogSuccessLabel:  color.New(color.Reset),
		LogDebugPrefix:   color.New(color.FgHiBlack),
		LogDebugLabel:    color.New(color.Reset),
		LogInfoPrefix:    color.New(color.FgBlue),
		LogInfoLabel:     color.New(color.Reset),
		LogWarnPrefix:    color.New(color.FgYellow),
		LogWarnLabel:     color.New(color.Reset),
		LogErrorPrefix:   color.New(color.FgRed),
		LogErrorLabel:    color.New(color.Reset),
		LogGroupBody:     color.New(color.Reset),

		// Input prompts
		InputPrefix:         color.New(color.FgYellow, color.Bold),
		InputLabel:          color.New(color.Reset),
		InputPlaceholder:    color.New(color.FgHiBlack),
		InputText:           color.New(color.Reset),
		InputValidationPass: color.New(color.FgGreen),
		InputValidationFail: color.New(color.FgRed),
		InputHelp:           color.New(color.FgHiBlack),

		// Confirmation prompts
		ConfirmationPrefix: color.New(color.FgYellow, color.Bold),
		ConfirmationLabel:  color.New(color.Reset),
		ConfirmationHelp:   color.New(color.FgHiBlack),

		// Selection prompts
		SelectionPrefix:             color.New(color.FgYellow, color.Bold),
		SelectionLabel:              color.New(color.Reset),
		SelectionHelp:               color.New(color.FgHiBlack),
		SelectionSearchLabel:        color.New(color.FgYellow, color.Bold),
		SelectionSearchText:         color.New(color.Reset),
		SelectionSearchHint:         color.New(color.FgHiBlack),
		SelectionValidationPass:     color.New(color.FgGreen),
		SelectionValidationFail:     color.New(color.FgRed),
		SelectionItemNormalMarker:   color.New(color.Reset),
		SelectionItemNormalLabel:    color.New(color.Reset),
		SelectionItemCurrentMarker:  color.New(color.FgYellow, color.Bold),
		SelectionItemCurrentLabel:   color.New(color.FgHiYellow),
		SelectionItemSelectedMarker: color.New(color.FgGreen),
		SelectionItemSelectedLabel:  color.New(color.FgGreen),
		SelectionItemDisabledMarker: color.New(color.FgHiBlack),
		SelectionItemDisabledLabel:  color.New(color.FgHiBlack, color.CrossedOut),

		// Spinners
		SpinnerPrefix: color.New(color.FgYellow, color.Bold),
		SpinnerLabel:  color.New(color.Reset),

		// Progress bars
		ProgressPrefix:     color.New(color.FgYellow, color.Bold),
		ProgressLabel:      color.New(color.Reset),
		ProgressBarPad:     color.New(color.FgYellow),
		ProgressBarDone:    color.New(color.FgYellow),
		ProgressBarPending: color.New(color.FgHiBlack),
		ProgressBarStatus:  color.New(color.Reset),
	}
}

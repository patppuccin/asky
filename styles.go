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
	LogBlockBody     *color.Color

	// Input prompt styles.
	InputDesc           *color.Color
	InputPrefix         *color.Color
	InputLabel          *color.Color
	InputPlaceholder    *color.Color
	InputText           *color.Color
	InputValidationPass *color.Color
	InputValidationFail *color.Color
	InputHelp           *color.Color

	// Confirmation prompt styles.
	ConfirmationPrefix         *color.Color
	ConfirmationLabel          *color.Color
	ConfirmationDesc           *color.Color
	ConfirmationHelp           *color.Color
	ConfirmationSelectedItem   *color.Color
	ConfirmationUnselectedItem *color.Color

	// Selection prompt styles.
	SelectionPrefix             *color.Color
	SelectionLabel              *color.Color
	SelectionDesc               *color.Color
	SelectionHelp               *color.Color
	SelectionSearchLabel        *color.Color
	SelectionSearchHint         *color.Color
	SelectionValidationPass     *color.Color
	SelectionValidationFail     *color.Color
	SelectionListItemHeader     *color.Color
	SelectionListItemLabel      *color.Color
	SelectionCurrentItemMarker  *color.Color
	SelectionCurrentItemLabel   *color.Color
	SelectionSelectedItemMarker *color.Color
	SelectionSelectedItemLabel  *color.Color
	SelectionDisabledItemMarker *color.Color
	SelectionDisabledItemLabel  *color.Color

	// Spinner styles.
	SpinnerPrefix *color.Color
	SpinnerLabel  *color.Color
	SpinnerDesc   *color.Color

	// Progress bar styles.
	ProgressPrefix     *color.Color
	ProgressLabel      *color.Color
	ProgressDesc       *color.Color
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
		LogBlockBody:     color.New(color.Reset),

		// Input prompts
		InputDesc:           color.New(color.FgHiBlue),
		InputPrefix:         color.New(color.FgMagenta, color.Bold),
		InputLabel:          color.New(color.FgHiMagenta),
		InputPlaceholder:    color.New(color.FgHiBlack),
		InputText:           color.New(color.Reset),
		InputValidationPass: color.New(color.FgGreen),
		InputValidationFail: color.New(color.FgRed),
		InputHelp:           color.New(color.FgHiBlack),

		// Confirmation prompts
		ConfirmationPrefix:         color.New(color.FgMagenta, color.Bold),
		ConfirmationLabel:          color.New(color.FgHiMagenta),
		ConfirmationDesc:           color.New(color.FgHiBlue),
		ConfirmationHelp:           color.New(color.FgHiBlack),
		ConfirmationSelectedItem:   color.New(color.FgBlack, color.BgMagenta),
		ConfirmationUnselectedItem: color.New(color.FgMagenta),

		// Selection prompts
		SelectionPrefix:             color.New(color.FgMagenta, color.Bold),
		SelectionLabel:              color.New(color.FgHiMagenta),
		SelectionDesc:               color.New(color.FgHiBlue),
		SelectionHelp:               color.New(color.FgHiBlack),
		SelectionSearchLabel:        color.New(color.FgMagenta, color.Bold),
		SelectionSearchHint:         color.New(color.FgHiBlack),
		SelectionValidationPass:     color.New(color.FgGreen),
		SelectionValidationFail:     color.New(color.FgRed),
		SelectionListItemHeader:     color.New(color.FgMagenta, color.Bold),
		SelectionListItemLabel:      color.New(color.Reset),
		SelectionCurrentItemMarker:  color.New(color.FgMagenta, color.Bold),
		SelectionCurrentItemLabel:   color.New(color.FgHiMagenta),
		SelectionSelectedItemMarker: color.New(color.FgGreen),
		SelectionSelectedItemLabel:  color.New(color.FgGreen),
		SelectionDisabledItemMarker: color.New(color.FgHiBlack),
		SelectionDisabledItemLabel:  color.New(color.FgHiBlack, color.CrossedOut),

		// Spinners
		SpinnerPrefix: color.New(color.FgMagenta, color.Bold),
		SpinnerLabel:  color.New(color.FgHiMagenta),
		SpinnerDesc:   color.New(color.FgHiBlue),

		// Progress bars
		ProgressPrefix:     color.New(color.FgMagenta, color.Bold),
		ProgressLabel:      color.New(color.FgHiMagenta),
		ProgressDesc:       color.New(color.FgHiBlue),
		ProgressBarPad:     color.New(color.FgHiMagenta),
		ProgressBarDone:    color.New(color.FgGreen),
		ProgressBarPending: color.New(color.FgYellow),
		ProgressBarStatus:  color.New(color.FgHiMagenta),
	}
}

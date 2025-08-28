package asky

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

type Theme struct {
	PromptSymbolStyle     func(string, ...any) string
	PromptTextStyle       func(string, ...any) string
	HelperTextStyle       func(string, ...any) string
	ErrorTextStyle        func(string, ...any) string
	ConfirmationTextStyle func(string, ...any) string
}

var DefaultTheme = Theme{
	PromptSymbolStyle:     color.YellowString,
	PromptTextStyle:       color.WhiteString,
	HelperTextStyle:       color.HiBlackString,
	ErrorTextStyle:        color.RedString,
	ConfirmationTextStyle: color.YellowString,
}

type TextInput struct {
	PromptSymbol string
	PromptText   string
	defaultValue string
	helperText   string
	Separator    string
	Theme        Theme
}

func NewTextInput() *TextInput {
	return &TextInput{
		PromptSymbol: "[?] ",
		Separator:    ": ",
		Theme:        DefaultTheme,
	}
}

func (ti *TextInput) WithPromptSymbol(p string) *TextInput {
	ti.PromptSymbol = p
	return ti
}

func (ti *TextInput) WithPromptText(p string) *TextInput {
	ti.PromptText = p
	return ti
}

func (ti *TextInput) WithDefault(val string) *TextInput {
	ti.defaultValue = val
	return ti
}

func (ti *TextInput) WithHelper(txt string) *TextInput {
	ti.helperText = txt
	return ti
}

func (ti *TextInput) WithSeparator(sep string) *TextInput {
	ti.Separator = sep
	return ti
}

func (ti *TextInput) WithTheme(th Theme) *TextInput {
	ti.Theme = th
	return ti
}

var ErrInterrupted = errors.New("prompt interrupted")

func (ti *TextInput) Render() (string, error) {
	// Save cursor before printing the prompt
	fmt.Print("\033[s")

	// Print the helper + prompt
	fmt.Println()
	if ti.helperText != "" || ti.defaultValue != "" {
		helper := ti.helperText
		if ti.defaultValue != "" {
			if helper != "" {
				helper += " "
			}
			helper += "(Default: " + ti.defaultValue + ")"
		}
		fmt.Println(ti.Theme.HelperTextStyle(helper))
	}

	fmt.Print(ti.Theme.PromptSymbolStyle(ti.PromptSymbol) + ti.Theme.PromptTextStyle(ti.PromptText+ti.Separator))

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		// restore and clear everything from saved position downwards
		fmt.Print("\033[u\033[J")
		return "", ErrInterrupted
	}

	// Trim newline(s)
	input = strings.TrimRight(input, "\r\n")

	// Apply default if empty
	if input == "" && ti.defaultValue != "" {
		input = ti.defaultValue
	}

	// restore and clear everything from saved position downwards
	fmt.Print("\033[u\033[J")

	return input, nil
}

// SPINNER -------------------------------------------------

var SpinnerPatternDots = []string{"⠋ ", "⠙ ", "⠹ ", "⠸ ", "⠼ ", "⠴ ", "⠦ ", "⠧ ", "⠇ ", "⠏ "}
var SpinnerPatternCircles = []string{"◐ ", "◓ ", "◑ ", "◒ "}
var SpinnerPatternSquares = []string{"▖ ", "▘ ", "▝ ", "▗ "}
var SpinnerPatternLines = []string{"╾ ", "│ ", "╸ ", "┤ ", "├ ", "└ ", "┴ ", "┬ ", "┐ ", "┘ "}
var SpinnerPatternMoons = []string{"🌑 ", "🌒 ", "🌓 ", "🌔 ", "🌕 ", "🌖 ", "🌗 ", "🌘 "}
var SpinnerPatternWave = []string{"▁ ", "▂ ", "▃ ", "▄ ", "▅ ", "▆ ", "▇ ", "█ ", "▇ ", "▆ ", "▅ ", "▄ ", "▃ ", "▂ ", "▁ "}

type Spinner struct {
	frames []string
	theme  Theme
	stop   bool
}

func NewSpinner() *Spinner {
	return &Spinner{
		frames: []string{"[⠋] ", "[⠙] ", "[⠹] ", "[⠸] ", "[⠼] ", "[⠴] ", "[⠦] ", "[⠧] ", "[⠇] ", "[⠏] "},
		theme:  DefaultTheme,
	}
}

func (s *Spinner) WithTheme(t Theme) *Spinner {
	s.theme = t
	return s
}

func (s *Spinner) WithFrames(f []string) *Spinner {
	s.frames = f
	return s
}

func (s *Spinner) Start(text string) {
	// hide cursor
	fmt.Print("\033[?25l")

	go func() {
		i := 0
		for !s.stop {
			thisFrame := s.frames[i%len(s.frames)]
			fmt.Printf("%s%s\r",
				s.theme.PromptSymbolStyle(thisFrame),
				s.theme.PromptTextStyle(text),
			)
			time.Sleep(250 * time.Millisecond)
			i++
		}
		// clear line on stop
		fmt.Print("\r\033[K")
	}()
}

func (s *Spinner) Stop() {
	s.stop = true
	// clear line + show cursor again
	fmt.Print("\r\033[K\033[?25h")
}

// CONFIRM -------------------------------------------------
type Confirm struct {
	promptSymbol string
	promptText   string
	helperText   string
	separator    string
	yesText      string
	noText       string
	defaultValue bool
	theme        Theme
}

func NewConfirm() *Confirm {
	return &Confirm{
		promptSymbol: "[?] ",
		promptText:   "Are you sure?",
		helperText:   "",
		separator:    ": ",
		yesText:      "Yes",
		noText:       "No",
		defaultValue: false,
		theme:        DefaultTheme,
	}
}

func (c *Confirm) WithPromptText(p string) *Confirm {
	c.promptText = p
	return c
}

func (c *Confirm) WithHelper(txt string) *Confirm {
	c.helperText = txt
	return c
}

func (c *Confirm) WithSeparator(sep string) *Confirm {
	c.separator = sep
	return c
}

func (c *Confirm) WithTheme(th Theme) *Confirm {
	c.theme = th
	return c
}

func (c *Confirm) WithYesText(txt string) *Confirm {
	c.yesText = txt
	return c
}

func (c *Confirm) WithNoText(txt string) *Confirm {
	c.noText = txt
	return c
}

func (c *Confirm) WithDefault(val bool) *Confirm {
	c.defaultValue = val
	return c
}

func (c *Confirm) Render() (bool, error) {
	// Save cursor before printing the prompt
	fmt.Print("\033[s")

	fmt.Println()

	// Helper + default
	if c.helperText != "" || c.defaultValue {
		helper := c.helperText
		if helper != "" {
			helper += " "
		}
		defVal := c.noText
		if c.defaultValue {
			defVal = c.yesText
		}
		helper += "(Default: " + defVal + ")"
		fmt.Println(c.theme.HelperTextStyle(helper))
	}

	// Show prompt
	fmt.Printf("%s%s%s%s",
		c.theme.PromptSymbolStyle(c.promptSymbol),
		c.theme.PromptTextStyle(c.promptText),
		c.theme.ConfirmationTextStyle(" ["+c.yesText+"/"+c.noText+"]"),
		c.separator,
	)

	// Read input
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		// restore and clear everything from saved position downwards
		fmt.Print("\033[u\033[J")
		return false, ErrInterrupted
	}

	// Clean input
	input = strings.TrimSpace(strings.ToLower(input))

	// Apply default if empty
	if input == "" {
		// restore and clear everything from saved position downwards
		fmt.Print("\033[u\033[J")
		return c.defaultValue, nil
	}

	// Cleanup
	fmt.Print("\033[u\033[J")

	// Parse yes/no
	switch input {
	case "y", "yes":
		return true, nil
	case "n", "no":
		return false, nil
	default:
		return c.defaultValue, nil
	}
}

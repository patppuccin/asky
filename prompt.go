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
	PromptSymbolStyle func(string, ...any) string
	PromptTextStyle   func(string, ...any) string
	HelperTextStyle   func(string, ...any) string
	ErrorTextStyle    func(string, ...any) string
}

var DefaultTheme = Theme{
	PromptSymbolStyle: color.YellowString,
	PromptTextStyle:   color.WhiteString,
	HelperTextStyle:   color.HiBlackString,
	ErrorTextStyle:    color.RedString,
}

type TextInput struct {
	PromptSymbol string
	PromptText   string
	DefaultValue string
	HelperText   string
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
	ti.DefaultValue = val
	return ti
}

func (ti *TextInput) WithHelper(txt string) *TextInput {
	ti.HelperText = txt
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
	if ti.HelperText != "" || ti.DefaultValue != "" {
		helper := ti.HelperText
		if ti.DefaultValue != "" {
			if helper != "" {
				helper += " "
			}
			helper += "(Default: " + ti.DefaultValue + ")"
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
	if input == "" && ti.DefaultValue != "" {
		input = ti.DefaultValue
	}

	// restore and clear everything from saved position downwards
	fmt.Print("\033[u\033[J")

	return input, nil
}

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

func (s *Spinner) Start(text string) {
	// hide cursor
	fmt.Print("\033[?25l")

	go func() {
		i := 0
		for !s.stop {
			thisFrame := s.frames[i%len(s.frames)]
			fmt.Printf("%s%s\r", s.theme.PromptSymbolStyle(thisFrame), s.theme.PromptTextStyle(text))
			time.Sleep(80 * time.Millisecond)
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

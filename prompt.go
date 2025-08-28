package asky

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

// -----------------------------------------------------------------------------
// Theme
// -----------------------------------------------------------------------------

type Theme struct {
	Prompt func(string, ...any) string
	Helper func(string, ...any) string
	Error  func(string, ...any) string
}

var DefaultTheme = Theme{
	Prompt: color.CyanString,
	Helper: color.HiBlackString,
	Error:  color.RedString,
}

// -----------------------------------------------------------------------------
// TextInput
// -----------------------------------------------------------------------------

type TextInputProps struct {
	Prompt       string
	DefaultValue string
	HelperText   string
	Separator    string
	Theme        Theme
}

type TextInput struct {
	props TextInputProps
}

func TextInputPrompt() *TextInput {
	return &TextInput{
		props: TextInputProps{
			Separator: ": ",
			Theme:     DefaultTheme,
		},
	}
}

func (ti *TextInput) Prompt(p string) *TextInput {
	ti.props.Prompt = p
	return ti
}

func (ti *TextInput) DefaultResponse(val string) *TextInput {
	ti.props.DefaultValue = val
	return ti
}

func (ti *TextInput) HelperText(txt string) *TextInput {
	ti.props.HelperText = txt
	return ti
}

func (ti *TextInput) WithSeparator(sep string) *TextInput {
	ti.props.Separator = sep
	return ti
}

func (ti *TextInput) WithTheme(th Theme) *TextInput {
	ti.props.Theme = th
	return ti
}

func (ti *TextInput) Render() (string, error) {
	if ti.props.HelperText != "" {
		fmt.Println(ti.props.Theme.Helper(ti.props.HelperText))
	}

	fmt.Print(ti.props.Theme.Prompt(ti.props.Prompt) + ti.props.Separator)

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	input = strings.TrimSpace(input)

	if input == "" && ti.props.DefaultValue != "" {
		return ti.props.DefaultValue, nil
	}

	return input, nil
}

package asky

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

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
		Theme:        ThemeDefault,
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
	os.Stdout.Write([]byte("\033[s"))

	// Print the helper + prompt
	// fmt.Println()
	os.Stdout.Write([]byte("\n"))
	if ti.helperText != "" || ti.defaultValue != "" {
		helper := ti.helperText
		if ti.defaultValue != "" {
			if helper != "" {
				helper += " "
			}
			helper += "(Default: " + ti.defaultValue + ")"
		}
		os.Stdout.WriteString(ti.Theme.MutedStyle(helper) + "\n")
	}

	os.Stdout.WriteString(ti.Theme.SecondaryStyle(ti.PromptSymbol))
	os.Stdout.WriteString(ti.Theme.PrimaryStyle(ti.PromptText + ti.Separator))

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')

	os.Stdout.Write([]byte("\033[u\033[J"))

	if err != nil {
		return "", ErrInterrupted
	}

	// Trim newline(s)
	input = strings.TrimRight(input, "\r\n")

	// Apply default if empty
	if input == "" && ti.defaultValue != "" {
		input = ti.defaultValue
	}

	return input, nil
}

package asky

import (
	"bufio"
	"os"
	"strings"
)

type Confirm struct {
	promptSymbol string
	promptText   string
	helperText   string
	separator    string
	defaultValue bool
	theme        Theme
}

func NewConfirm() *Confirm {
	return &Confirm{
		promptSymbol: "[?] ",
		promptText:   "Are you sure?",
		helperText:   "",
		separator:    ": ",
		defaultValue: false,
		theme:        ThemeDefault,
	}
}

func (c *Confirm) WithPromptSymbol(p string) *Confirm {
	c.promptSymbol = p
	return c
}

func (c *Confirm) WithPromptText(p string) *Confirm {
	c.promptText = p
	return c
}

func (c *Confirm) WithHelperText(txt string) *Confirm {
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

func (c *Confirm) WithDefaultOption(val bool) *Confirm {
	c.defaultValue = val
	return c
}

func (c *Confirm) Render() (bool, error) {
	os.Stdout.Write([]byte("\r\033[s"))

	// Helper + default
	os.Stdout.WriteString("\r\n")

	yChar := "y"
	nChar := "N"
	if c.helperText != "" || c.defaultValue {
		helper := c.helperText
		if helper != "" {
			helper += " "
		}
		defVal := "No"
		if c.defaultValue {
			yChar, nChar = "Y", "n"
			defVal = "Yes"
		}
		helper += "(Default: " + defVal + ")"
		os.Stdout.WriteString(c.theme.MutedStyle(helper))
		os.Stdout.Write([]byte("\n"))
	}
	// Show prompt
	os.Stdout.WriteString(c.theme.SecondaryStyle(c.promptSymbol))
	os.Stdout.WriteString(c.theme.PrimaryStyle(c.promptText))
	os.Stdout.WriteString(c.theme.AccentStyle(" ["))
	os.Stdout.WriteString(c.theme.SuccessStyle(yChar))
	os.Stdout.WriteString(c.theme.AccentStyle("/"))
	os.Stdout.WriteString(c.theme.ErrorStyle(nChar))
	os.Stdout.WriteString(c.theme.AccentStyle("]"))
	os.Stdout.WriteString(c.theme.PrimaryStyle(c.separator))

	// Read input
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')

	os.Stdout.Write([]byte("\033[u\033[J"))

	if err != nil {
		return false, ErrInterrupted
	}

	// Parse yes/no
	switch strings.TrimSpace(strings.ToLower(input)) {
	case "y", "yes":
		return true, nil
	case "n", "no":
		return false, nil
	default:
		return c.defaultValue, nil
	}
}

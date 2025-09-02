package asky

import "strconv"

// PrimaryStyle
// SecondaryStyle
// TertiaryStyle
// AccentStyle
// MutedStyle

// SuccessStyle
// InfoStyle
// WarningStyle
// ErrorStyle
// NeutralStyle

// ActiveStyle
// InactiveStyle

type Theme struct {
	Primary   string // Labels
	Secondary string // Symbols
	Accent    string // Updates
	Muted     string
	Success   string
	Info      string
	Warning   string
	Error     string
}

func (t Theme) PrimaryStyle(s string) string   { return hex(t.Primary, s) }
func (t Theme) SecondaryStyle(s string) string { return hex(t.Secondary, s) }
func (t Theme) AccentStyle(s string) string    { return hex(t.Accent, s) }
func (t Theme) MutedStyle(s string) string     { return hex(t.Muted, s) }
func (t Theme) SuccessStyle(s string) string   { return hex(t.Success, s) }
func (t Theme) InfoStyle(s string) string      { return hex(t.Info, s) }
func (t Theme) WarningStyle(s string) string   { return hex(t.Warning, s) }
func (t Theme) ErrorStyle(s string) string     { return hex(t.Error, s) }

// Helper Functions ---------------------------------------

func hex(hex, s string) string {
	// expect "#RRGGBB"
	if len(hex) != 7 || hex[0] != '#' {
		return s // fallback
	}

	r, err1 := strconv.ParseInt(hex[1:3], 16, 0)
	g, err2 := strconv.ParseInt(hex[3:5], 16, 0)
	b, err3 := strconv.ParseInt(hex[5:7], 16, 0)
	if err1 != nil || err2 != nil || err3 != nil {
		return s
	}

	return "\033[38;2;" +
		strconv.Itoa(int(r)) + ";" +
		strconv.Itoa(int(g)) + ";" +
		strconv.Itoa(int(b)) + "m" +
		s + "\033[0m"
}

// Color Themes for the library -------------------------

// Catppuccin Mocha
// ThemeCatppuccinMocha
var ThemeDefault = Theme{
	Primary:   "#f5c2e7", // pink (labels / main text)
	Secondary: "#cba6f7", // mauve (symbols / secondary highlights)
	Accent:    "#89b4fa", // blue (updates / spinners / progress)
	Muted:     "#a6adc8", // overlay0 (subtle / de-emphasized)
	Success:   "#a6e3a1", // green
	Info:      "#94e2d5", // teal
	Warning:   "#f9e2af", // yellow
	Error:     "#f38ba8", // red
}

// // Default (simple ANSI approximations)
// var ThemeDefault = Theme{
// 	PromptSymbolHex:     "\033[33m", // yellow
// 	PromptTextHex:       "\033[37m", // white
// 	HelperTextHex:       "\033[90m", // gray
// 	ErrorTextHex:        "\033[31m", // red
// 	ConfirmationTextHex: "\033[33m", // yellow
// }

// // Catppuccin Mocha
// var ThemeCatppuccinMocha = Theme{
// 	PromptSymbolHex:     "#cba6f7", // mauve
// 	PromptTextHex:       "#89b4fa", // blue
// 	HelperTextHex:       "#a6adc8", // overlay
// 	ErrorTextHex:        "#f38ba8", // red
// 	ConfirmationTextHex: "#94e2d5", // teal
// }

// // Tokyo Night
// var ThemeTokyoNight = Theme{
// 	PromptSymbolHex:     "#7aa2f7", // blue
// 	PromptTextHex:       "#c0caf5", // fg
// 	HelperTextHex:       "#565f89", // comment gray
// 	ErrorTextHex:        "#f7768e", // red
// 	ConfirmationTextHex: "#9ece6a", // green
// }

// // Gruvbox
// var ThemeGruvbox = Theme{
// 	PromptSymbolHex:     "#fabd2f", // yellow
// 	PromptTextHex:       "#ebdbb2", // fg
// 	HelperTextHex:       "#928374", // gray
// 	ErrorTextHex:        "#fb4934", // red
// 	ConfirmationTextHex: "#b8bb26", // green
// }

// // Osaka Jade (custom deep-green palette)
// var ThemeOsakaJade = Theme{
// 	PromptSymbolHex:     "#7ec07a", // jade
// 	PromptTextHex:       "#cbe6c1", // light green
// 	HelperTextHex:       "#6a9955", // muted green
// 	ErrorTextHex:        "#e06c75", // red
// 	ConfirmationTextHex: "#98c379", // confirm green
// }

// // Plain (no colors, safe for CI / dumb terminals)
// var ThemePlain = Theme{
// 	PromptSymbolHex:     "",
// 	PromptTextHex:       "",
// 	HelperTextHex:       "",
// 	ErrorTextHex:        "",
// 	ConfirmationTextHex: "",
// }

// type Theme struct {
// 	PromptSymbolStyle     func(string, ...any) string
// 	PromptTextStyle       func(string, ...any) string
// 	HelperTextStyle       func(string, ...any) string
// 	ErrorTextStyle        func(string, ...any) string
// 	ConfirmationTextStyle func(string, ...any) string
// }

// var ThemeDefault = Theme{
// 	PromptSymbolStyle:     color.YellowString,
// 	PromptTextStyle:       color.WhiteString,
// 	HelperTextStyle:       color.HiBlackString,
// 	ErrorTextStyle:        color.RedString,
// 	ConfirmationTextStyle: color.YellowString,
// }

// PromptSymbolHex string
// 	PromptTextHex   string
// 	HelperTextHex   string

// 	ErrorTextHex        string
// 	ConfirmationTextHex string

// func (t Theme) PromptSymbolStyle(s string) string     { return hex(t.PromptSymbolHex, s) }
// func (t Theme) PromptTextStyle(s string) string       { return hex(t.PromptTextHex, s) }
// func (t Theme) HelperTextStyle(s string) string       { return hex(t.HelperTextHex, s) }
// func (t Theme) ErrorTextStyle(s string) string        { return hex(t.ErrorTextHex, s) }
// func (t Theme) ConfirmationTextStyle(s string) string { return hex(t.ConfirmationTextHex, s) }

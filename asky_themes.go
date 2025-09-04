package asky

import (
	"strconv"
	"strings"
)

type Style struct {
	Fg            string
	Bg            string
	Dim           bool
	Bold          bool
	Italic        bool
	Underline     bool
	Strikethrough bool
}

// --- Positive Style Builders -----------------------------
func (st Style) WithFgHEX(hex string) Style {
	r, g, b := hexParse(hex)
	st.Fg = "38;2;" + strconv.Itoa(r) + ";" + strconv.Itoa(g) + ";" + strconv.Itoa(b)
	return st
}

func (st Style) WithBgHEX(hex string) Style {
	r, g, b := hexParse(hex)
	st.Bg = "48;2;" + strconv.Itoa(r) + ";" + strconv.Itoa(g) + ";" + strconv.Itoa(b)
	return st
}

func (st Style) WithFgRGB(r, g, b int) Style {
	st.Fg = "38;2;" + strconv.Itoa(r) + ";" + strconv.Itoa(g) + ";" + strconv.Itoa(b)
	return st
}

func (st Style) WithBgRGB(r, g, b int) Style {
	st.Bg = "48;2;" + strconv.Itoa(r) + ";" + strconv.Itoa(g) + ";" + strconv.Itoa(b)
	return st
}

func (st Style) WithFgANSI(base int) Style {
	st.Fg = strconv.Itoa(base) // directly assigns 30–37 or 90–97
	return st
}

func (st Style) WithBgANSI(code int) Style {
	st.Bg = strconv.Itoa(code) // directly assigns 40–47 or 100–107
	return st
}

func (st Style) WithDim() Style           { st.Dim = true; return st }
func (st Style) WithBold() Style          { st.Bold = true; return st }
func (st Style) WithItalic() Style        { st.Italic = true; return st }
func (st Style) WithUnderline() Style     { st.Underline = true; return st }
func (st Style) WithStrikethrough() Style { st.Strikethrough = true; return st }

// --- Negative Style Builders -----------------------------
func (st Style) WithoutFGColor() Style       { st.Fg = ""; return st }
func (st Style) WithoutBGColor() Style       { st.Bg = ""; return st }
func (st Style) WithoutDim() Style           { st.Dim = false; return st }
func (st Style) WithoutBold() Style          { st.Bold = false; return st }
func (st Style) WithoutItalic() Style        { st.Italic = false; return st }
func (st Style) WithoutUnderline() Style     { st.Underline = false; return st }
func (st Style) WithoutStrikethrough() Style { st.Strikethrough = false; return st }

func hexParse(hx string) (int, int, int) {
	hx = strings.TrimPrefix(hx, "#")
	if len(hx) != 6 {
		return 0, 0, 0
	}
	r, _ := strconv.ParseUint(hx[0:2], 16, 8)
	g, _ := strconv.ParseUint(hx[2:4], 16, 8)
	b, _ := strconv.ParseUint(hx[4:6], 16, 8)
	return int(r), int(g), int(b)
}

// --- Styled Printer --------------------------------------
func (st Style) Sprint(s string) string {
	if s == "" {
		return ""
	}

	// quick exit: no styling at all
	if !st.Bold && !st.Dim && !st.Italic &&
		!st.Underline && !st.Strikethrough &&
		st.Fg == "" && st.Bg == "" {
		return s
	}

	var codes []string

	// text attributes
	if st.Bold {
		codes = append(codes, "1")
	}
	if st.Dim {
		codes = append(codes, "2")
	}
	if st.Italic {
		codes = append(codes, "3")
	}
	if st.Underline {
		codes = append(codes, "4")
	}
	if st.Strikethrough {
		codes = append(codes, "9")
	}

	// colors (already prebuilt by builders)
	if st.Fg != "" {
		codes = append(codes, st.Fg)
	}
	if st.Bg != "" {
		codes = append(codes, st.Bg)
	}

	var b strings.Builder
	b.Grow(10*len(codes) + len(s) + 4) // allocate once
	b.WriteString("\x1b[")
	b.WriteString(strings.Join(codes, ";"))
	b.WriteString("m")
	b.WriteString(s)
	b.WriteString("\x1b[0m")

	return b.String()
}

type Theme struct {
	Primary   Style
	Secondary Style
	Accent    Style
	Neutral   Style
	Muted     Style
	Success   Style
	Info      Style
	Warning   Style
	Error     Style
	Disabled  Style
	Emphasis  Style
	Outline   Style
}

func (t Theme) PrimaryStyle(s string) string   { return t.Primary.Sprint(s) }
func (t Theme) SecondaryStyle(s string) string { return t.Secondary.Sprint(s) }
func (t Theme) AccentStyle(s string) string    { return t.Accent.Sprint(s) }
func (t Theme) NeutralStyle(s string) string   { return t.Neutral.Sprint(s) }
func (t Theme) MutedStyle(s string) string     { return t.Muted.Sprint(s) }
func (t Theme) SuccessStyle(s string) string   { return t.Success.Sprint(s) }
func (t Theme) InfoStyle(s string) string      { return t.Info.Sprint(s) }
func (t Theme) WarningStyle(s string) string   { return t.Warning.Sprint(s) }
func (t Theme) ErrorStyle(s string) string     { return t.Error.Sprint(s) }
func (t Theme) DisabledStyle(s string) string  { return t.Disabled.Sprint(s) }
func (t Theme) EmphasisStyle(s string) string  { return t.Emphasis.Sprint(s) }
func (t Theme) OutlineStyle(s string) string   { return t.Outline.Sprint(s) }

// Color Themes for the library -------------------------

var ThemeDefault = Theme{
	Primary:   (Style{}).WithFgANSI(35),                 // magenta
	Secondary: (Style{}).WithFgANSI(34),                 // blue
	Accent:    (Style{}).WithFgANSI(36),                 // cyan
	Neutral:   (Style{}).WithFgANSI(37),                 // white
	Muted:     (Style{}).WithFgANSI(90),                 // bright black (gray)
	Success:   (Style{}).WithFgANSI(32),                 // green
	Info:      (Style{}).WithFgANSI(36),                 // cyan
	Warning:   (Style{}).WithFgANSI(33),                 // yellow
	Error:     (Style{}).WithFgANSI(31),                 // red
	Disabled:  (Style{}).WithFgANSI(90).WithDim(),       // gray dim
	Emphasis:  (Style{}).WithFgANSI(35).WithUnderline(), // magenta underline
	Outline:   (Style{}).WithFgANSI(90),                 // gray border
}

var ThemeCatppuccinMocha = Theme{
	Primary:   (Style{}).WithFgHEX("#cba6f7"),                               // Catppuccin Mocha Mauve
	Secondary: (Style{}).WithFgHEX("#89b4fa"),                               // Blue
	Accent:    (Style{}).WithFgHEX("#94e2d5"),                               // Teal
	Neutral:   (Style{}).WithFgHEX("#cdd6f4"),                               // Text
	Muted:     (Style{}).WithFgHEX("#a6adc8"),                               // Subtext0
	Success:   (Style{}).WithFgHEX("#a6e3a1"),                               // Green
	Info:      (Style{}).WithFgHEX("#89dceb"),                               // Sky
	Warning:   (Style{}).WithFgHEX("#f9e2af"),                               // Yellow
	Error:     (Style{}).WithFgHEX("#f38ba8"),                               // Red
	Disabled:  (Style{}).WithFgHEX("#585b70").WithDim().WithStrikethrough(), // Overlay0
	Emphasis:  (Style{}).WithFgHEX("#b4befe").WithUnderline(),               // Lavender, underlined
	Outline:   (Style{}).WithFgHEX("#6c7086"),                               // Overlay1, for borders
}

var ThemeGruvbox = Theme{
	Primary:   (Style{}).WithFgHEX("#d3869b"),                               // mauve/pink
	Secondary: (Style{}).WithFgHEX("#83a598"),                               // aqua
	Accent:    (Style{}).WithFgHEX("#fabd2f"),                               // yellow
	Neutral:   (Style{}).WithFgHEX("#ebdbb2"),                               // fg
	Muted:     (Style{}).WithFgHEX("#a89984"),                               // gray
	Success:   (Style{}).WithFgHEX("#b8bb26"),                               // green
	Info:      (Style{}).WithFgHEX("#83a598"),                               // blue
	Warning:   (Style{}).WithFgHEX("#fabd2f"),                               // yellow
	Error:     (Style{}).WithFgHEX("#fb4934"),                               // red
	Disabled:  (Style{}).WithFgHEX("#665c54").WithDim().WithStrikethrough(), // dark gray
	Emphasis:  (Style{}).WithFgHEX("#d3869b").WithUnderline(),               // pink underlined
	Outline:   (Style{}).WithFgHEX("#3c3836"),                               // bg2
}

var ThemeKanagawa = Theme{
	Primary:   (Style{}).WithFgHEX("#957fb8"),                               // purple
	Secondary: (Style{}).WithFgHEX("#7e9cd8"),                               // blue
	Accent:    (Style{}).WithFgHEX("#7aa89f"),                               // teal
	Neutral:   (Style{}).WithFgHEX("#dcd7ba"),                               // fg
	Muted:     (Style{}).WithFgHEX("#938aa9"),                               // muted purple-gray
	Success:   (Style{}).WithFgHEX("#98bb6c"),                               // green
	Info:      (Style{}).WithFgHEX("#7fb4ca"),                               // cyan
	Warning:   (Style{}).WithFgHEX("#e6c384"),                               // yellow
	Error:     (Style{}).WithFgHEX("#e46876"),                               // red
	Disabled:  (Style{}).WithFgHEX("#727169").WithDim().WithStrikethrough(), // gray
	Emphasis:  (Style{}).WithFgHEX("#957fb8").WithUnderline(),               // purple underline
	Outline:   (Style{}).WithFgHEX("#54546d"),                               // bg edge
}

var ThemeTokyoNight = Theme{
	Primary:   (Style{}).WithFgHEX("#bb9af7"),                               // purple
	Secondary: (Style{}).WithFgHEX("#7aa2f7"),                               // blue
	Accent:    (Style{}).WithFgHEX("#7dcfff"),                               // cyan
	Neutral:   (Style{}).WithFgHEX("#c0caf5"),                               // fg
	Muted:     (Style{}).WithFgHEX("#565f89"),                               // muted gray
	Success:   (Style{}).WithFgHEX("#9ece6a"),                               // green
	Info:      (Style{}).WithFgHEX("#2ac3de"),                               // aqua
	Warning:   (Style{}).WithFgHEX("#e0af68"),                               // yellow
	Error:     (Style{}).WithFgHEX("#f7768e"),                               // red
	Disabled:  (Style{}).WithFgHEX("#414868").WithDim().WithStrikethrough(), // dark gray
	Emphasis:  (Style{}).WithFgHEX("#bb9af7").WithUnderline(),               // purple underline
	Outline:   (Style{}).WithFgHEX("#3b4261"),                               // border
}

var ThemeDracula = Theme{
	Primary:   (Style{}).WithFgHEX("#bd93f9"),                               // purple
	Secondary: (Style{}).WithFgHEX("#8be9fd"),                               // cyan
	Accent:    (Style{}).WithFgHEX("#ffb86c"),                               // orange
	Neutral:   (Style{}).WithFgHEX("#f8f8f2"),                               // fg
	Muted:     (Style{}).WithFgHEX("#6272a4"),                               // muted blue-gray
	Success:   (Style{}).WithFgHEX("#50fa7b"),                               // green
	Info:      (Style{}).WithFgHEX("#8be9fd"),                               // cyan
	Warning:   (Style{}).WithFgHEX("#f1fa8c"),                               // yellow
	Error:     (Style{}).WithFgHEX("#ff5555"),                               // red
	Disabled:  (Style{}).WithFgHEX("#44475a").WithDim().WithStrikethrough(), // gray
	Emphasis:  (Style{}).WithFgHEX("#bd93f9").WithUnderline(),               // purple underline
	Outline:   (Style{}).WithFgHEX("#282a36"),                               // border
}

var ThemeOsakaJade = Theme{
	Primary:   (Style{}).WithFgHEX("#00a37a"),                               // jade
	Secondary: (Style{}).WithFgHEX("#3dbd93"),                               // lighter jade
	Accent:    (Style{}).WithFgHEX("#41c7b9"),                               // aqua
	Neutral:   (Style{}).WithFgHEX("#d0f0e0"),                               // light fg
	Muted:     (Style{}).WithFgHEX("#5a7d73"),                               // gray-green
	Success:   (Style{}).WithFgHEX("#7ed09e"),                               // pastel green
	Info:      (Style{}).WithFgHEX("#4fb0c6"),                               // blue-green
	Warning:   (Style{}).WithFgHEX("#f0c674"),                               // amber
	Error:     (Style{}).WithFgHEX("#d9534f"),                               // red
	Disabled:  (Style{}).WithFgHEX("#4a605a").WithDim().WithStrikethrough(), // muted
	Emphasis:  (Style{}).WithFgHEX("#00a37a").WithUnderline(),               // jade underline
	Outline:   (Style{}).WithFgHEX("#355e4d"),                               // dark green
}

var ThemeSolarizedLight = Theme{
	Primary:   (Style{}).WithFgHEX("#268bd2"), // blue
	Secondary: (Style{}).WithFgHEX("#2aa198"), // cyan
	Accent:    (Style{}).WithFgHEX("#6c71c4"), // violet
	Neutral:   (Style{}).WithFgHEX("#657b83"), // base00
	Muted:     (Style{}).WithFgHEX("#93a1a1"), // base1
	Success:   (Style{}).WithFgHEX("#859900"), // green
	Info:      (Style{}).WithFgHEX("#268bd2"), // blue
	Warning:   (Style{}).WithFgHEX("#b58900"), // yellow
	Error:     (Style{}).WithFgHEX("#dc322f"), // red
	Disabled:  (Style{}).WithFgHEX("#93a1a1").WithDim().WithStrikethrough(),
	Emphasis:  (Style{}).WithFgHEX("#6c71c4").WithUnderline(),
	Outline:   (Style{}).WithFgHEX("#eee8d5"), // base2
}

var ThemeCatppuccinLatte = Theme{
	Primary:   (Style{}).WithFgHEX("#8839ef"), // mauve
	Secondary: (Style{}).WithFgHEX("#1e66f5"), // blue
	Accent:    (Style{}).WithFgHEX("#04a5e5"), // sky
	Neutral:   (Style{}).WithFgHEX("#4c4f69"), // text
	Muted:     (Style{}).WithFgHEX("#6c6f85"), // subtext0
	Success:   (Style{}).WithFgHEX("#40a02b"), // green
	Info:      (Style{}).WithFgHEX("#209fb5"), // teal
	Warning:   (Style{}).WithFgHEX("#df8e1d"), // yellow
	Error:     (Style{}).WithFgHEX("#d20f39"), // red
	Disabled:  (Style{}).WithFgHEX("#9ca0b0").WithDim().WithStrikethrough(),
	Emphasis:  (Style{}).WithFgHEX("#8839ef").WithUnderline(),
	Outline:   (Style{}).WithFgHEX("#e6e9ef"), // crust
}

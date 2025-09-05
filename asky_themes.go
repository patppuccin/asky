package asky

import (
	"strconv"
	"strings"
)

// --- Color Definition -----------------------------------
type color string

func ColorFromHex(hx string) color {
	hx = strings.TrimPrefix(strings.TrimSpace(hx), "#")
	if len(hx) != 6 {
		return ""
	}
	var result strings.Builder
	result.WriteString("rgb:")
	for i := 0; i < 6; i += 2 {
		val, err := strconv.ParseUint(hx[i:i+2], 16, 8)
		if err != nil {
			return ""
		}
		if i > 0 {
			result.WriteString(",")
		}
		result.WriteString(strconv.Itoa(int(val)))
	}
	return color(result.String())
}
func ColorFromRGB(r, g, b int) color {
	if r < 0 || r > 255 || g < 0 || g > 255 || b < 0 || b > 255 {
		return ""
	}
	return color("rgb:" +
		strconv.Itoa(r) + "," +
		strconv.Itoa(g) + "," +
		strconv.Itoa(b))
}
func ColorFromANSI(i int) color {
	if i < 0 || i > 255 {
		return ""
	}
	return color("ansi:" + strconv.Itoa(i))
}

// --- Style Definition -----------------------------------
type style struct {
	fg, bg                   color
	dim, bold, italic        bool
	underline, strikethrough bool
}

func NewStyle() *style { return &style{} }

func (st *style) FG(c color) *style     { st.fg = c; return st }
func (st *style) BG(c color) *style     { st.bg = c; return st }
func (st *style) Dim() *style           { st.dim = true; return st }
func (st *style) Bold() *style          { st.bold = true; return st }
func (st *style) Italic() *style        { st.italic = true; return st }
func (st *style) Underline() *style     { st.underline = true; return st }
func (st *style) Strikethrough() *style { st.strikethrough = true; return st }
func (st *style) isEmpty() bool {
	return !st.bold && !st.dim && !st.italic &&
		!st.underline && !st.strikethrough &&
		st.fg == "" && st.bg == ""
}
func (st *style) parseColor(repr string, bg bool) (string, bool) {
	if repr == "" {
		return "", false
	}
	// rgb:r,g,b
	if strings.HasPrefix(repr, "rgb:") {
		s := repr[4:] // after "rgb:"
		// find two commas in one pass
		c1, c2 := -1, -1
		for i := 0; i < len(s); i++ {
			if s[i] == ',' {
				if c1 == -1 {
					c1 = i
				} else {
					c2 = i
					break
				}
			}
		}
		if c1 == -1 || c2 == -1 {
			return "", false
		}
		r := s[:c1]
		g := s[c1+1 : c2]
		b := s[c2+1:]
		if bg {
			return "48;2;" + r + ";" + g + ";" + b, true
		}
		return "38;2;" + r + ";" + g + ";" + b, true
	}

	// ansi:n
	if strings.HasPrefix(repr, "ansi:") {
		n := repr[5:]
		if n == "" {
			return "", false
		}
		if bg {
			return "48;5;" + n, true
		}
		return "38;5;" + n, true
	}
	return "", false
}
func (st *style) Sprint(text string) string {

	if text == "" || st.isEmpty() {
		return text
	}

	// Build the escape sequence
	var b strings.Builder
	b.Grow(len(text) + 24)
	b.WriteString("\x1b[")

	first := true // flag to check code presence
	write := func(code string) {
		if code == "" {
			return
		}
		if !first {
			b.WriteByte(';')
		}
		first = false
		b.WriteString(code)
	}

	if st.bold {
		write("1")
	}
	if st.dim {
		write("2")
	}
	if st.italic {
		write("3")
	}
	if st.underline {
		write("4")
	}
	if st.strikethrough {
		write("9")
	}
	if code, ok := st.parseColor(string(st.fg), false); ok {
		write(code)
	}
	if code, ok := st.parseColor(string(st.bg), true); ok {
		write(code)
	}
	if first {
		return text
	}

	b.WriteString("m")
	b.WriteString(text)
	b.WriteString("\x1b[0m")
	return b.String()
}

// --- Style Presets --------------------------------------
type preset struct {
	theme     *Theme
	success   *style
	debug     *style
	info      *style
	warn      *style
	err       *style
	neutral   *style
	muted     *style
	disabled  *style
	primary   *style
	secondary *style
	accent    *style
	highlight *style
}

func newPreset(theme Theme) *preset {
	return &preset{
		theme:     &theme,
		success:   NewStyle().FG(theme.Green).Bold(),
		debug:     NewStyle().FG(theme.Muted).Bold(),
		info:      NewStyle().FG(theme.Blue).Bold(),
		warn:      NewStyle().FG(theme.Yellow).Bold(),
		err:       NewStyle().FG(theme.Red).Bold(),
		neutral:   NewStyle().FG(theme.Foreground),
		muted:     NewStyle().FG(theme.Muted),
		disabled:  NewStyle().FG(theme.Foreground).Dim().Strikethrough(),
		highlight: NewStyle().FG(theme.Background).BG(theme.Highlight).Bold(),
		primary:   NewStyle().FG(theme.Primary),
		secondary: NewStyle().FG(theme.Secondary),
		accent:    NewStyle().FG(theme.Accent),
	}
}

// --- Theme Definition -----------------------------------
type Theme struct {
	Background    color
	BackgroundAlt color
	Foreground    color
	ForegroundAlt color
	Highlight     color
	Cursor        color
	Muted         color
	Outline       color
	Red           color
	Yellow        color
	Green         color
	Blue          color
	Teal          color
	Primary       color
	Secondary     color
	Accent        color
}

// --- Theme Presets --------------------------------------

// Default ANSI (light & dark)
var ThemeDefault = Theme{
	Background:    ColorFromANSI(0),  // black
	BackgroundAlt: ColorFromANSI(8),  // bright black (gray)
	Foreground:    ColorFromANSI(15), // bright white
	ForegroundAlt: ColorFromANSI(7),  // white
	Highlight:     ColorFromANSI(8),  // bright black (as bg highlight)
	Cursor:        ColorFromANSI(15), // bright white
	Muted:         ColorFromANSI(8),  // bright black
	Outline:       ColorFromANSI(8),  // bright black
	Red:           ColorFromANSI(9),  // bright red
	Yellow:        ColorFromANSI(11), // bright yellow
	Green:         ColorFromANSI(10), // bright green
	Blue:          ColorFromANSI(12), // bright blue
	Teal:          ColorFromANSI(14), // bright cyan
	Primary:       ColorFromANSI(5),  // magenta
	Secondary:     ColorFromANSI(13), // bright magenta
	Accent:        ColorFromANSI(3),  // yellow (closest to Accent in ANSI)
}

// Catppuccin Mocha (dark)
var ThemeCatppuccinMocha = Theme{
	Background:    ColorFromHex("#1e1e2e"), // base
	BackgroundAlt: ColorFromHex("#181825"), // mantle
	Foreground:    ColorFromHex("#cdd6f4"), // text
	ForegroundAlt: ColorFromHex("#bac2de"), // subtext1
	Highlight:     ColorFromHex("#cba6f7"), // mauve (used sparingly)
	Cursor:        ColorFromHex("#f5e0dc"), // rosewater
	Muted:         ColorFromHex("#6c7086"), // overlay0
	Outline:       ColorFromHex("#45475a"), // surface0
	Red:           ColorFromHex("#f38ba8"), // red (error/danger)
	Yellow:        ColorFromHex("#f9e2af"), // yellow (warning)
	Green:         ColorFromHex("#a6e3a1"), // green (success)
	Blue:          ColorFromHex("#89b4fa"), // blue (info)
	Primary:       ColorFromHex("#cba6f7"), // mauve (primary brand/selection)
	Secondary:     ColorFromHex("#f2cdcd"), // flamingo (secondary actions)
	Accent:        ColorFromHex("#fab387"), // peach (highlighted actions/callouts)
}

// Gruvbox (dark)
var ThemeGruvbox = Theme{
	Background:    ColorFromHex("#282828"),
	BackgroundAlt: ColorFromHex("#3c3836"),
	Foreground:    ColorFromHex("#ebdbb2"),
	ForegroundAlt: ColorFromHex("#d5c4a1"),
	Highlight:     ColorFromHex("#3c3836"),
	Cursor:        ColorFromHex("#ebdbb2"),
	Muted:         ColorFromHex("#a89984"),
	Outline:       ColorFromHex("#504945"),
	Red:           ColorFromHex("#fb4934"),
	Yellow:        ColorFromHex("#fabd2f"),
	Green:         ColorFromHex("#b8bb26"),
	Blue:          ColorFromHex("#83a598"),
	Teal:          ColorFromHex("#8ec07c"),
	Primary:       ColorFromHex("#b16286"),
	Secondary:     ColorFromHex("#d3869b"),
	Accent:        ColorFromHex("#fe8019"),
}

// Kanagawa (dark)
var ThemeKanagawa = Theme{
	Background:    ColorFromHex("#1f1f28"),
	BackgroundAlt: ColorFromHex("#2a2a37"),
	Foreground:    ColorFromHex("#dcd7ba"),
	ForegroundAlt: ColorFromHex("#c8c093"),
	Highlight:     ColorFromHex("#2a2a37"),
	Cursor:        ColorFromHex("#dcd7ba"),
	Muted:         ColorFromHex("#727169"),
	Outline:       ColorFromHex("#363646"),
	Red:           ColorFromHex("#c34043"),
	Yellow:        ColorFromHex("#c0a36e"),
	Green:         ColorFromHex("#98bb6c"),
	Blue:          ColorFromHex("#7e9cd8"),
	Teal:          ColorFromHex("#6a9589"),
	Primary:       ColorFromHex("#957fb8"),
	Secondary:     ColorFromHex("#d27e99"),
	Accent:        ColorFromHex("#ffa066"),
}

// Tokyo Night (dark)
var ThemeTokyoNight = Theme{
	Background:    ColorFromHex("#1a1b26"),
	BackgroundAlt: ColorFromHex("#24283b"),
	Foreground:    ColorFromHex("#c0caf5"),
	ForegroundAlt: ColorFromHex("#a9b1d6"),
	Highlight:     ColorFromHex("#2f334d"),
	Cursor:        ColorFromHex("#c0caf5"),
	Muted:         ColorFromHex("#565f89"),
	Outline:       ColorFromHex("#2a2e3f"),
	Red:           ColorFromHex("#f7768e"),
	Yellow:        ColorFromHex("#e0af68"),
	Green:         ColorFromHex("#9ece6a"),
	Blue:          ColorFromHex("#7aa2f7"),
	Teal:          ColorFromHex("#73daca"),
	Primary:       ColorFromHex("#bb9af7"),
	Secondary:     ColorFromHex("#ff7a93"),
	Accent:        ColorFromHex("#ff9e64"),
}

// Dracula (dark)
var ThemeDracula = Theme{
	Background:    ColorFromHex("#282a36"),
	BackgroundAlt: ColorFromHex("#44475a"),
	Foreground:    ColorFromHex("#f8f8f2"),
	ForegroundAlt: ColorFromHex("#e2e2dc"),
	Highlight:     ColorFromHex("#44475a"),
	Cursor:        ColorFromHex("#f8f8f2"),
	Muted:         ColorFromHex("#6272a4"),
	Outline:       ColorFromHex("#3a3d4a"),
	Red:           ColorFromHex("#ff5555"),
	Yellow:        ColorFromHex("#f1fa8c"),
	Green:         ColorFromHex("#50fa7b"),
	Blue:          ColorFromHex("#8be9fd"),
	Teal:          ColorFromHex("#8be9fd"),
	Primary:       ColorFromHex("#bd93f9"),
	Secondary:     ColorFromHex("#ff79c6"),
	Accent:        ColorFromHex("#ffb86c"),
}

// Osaka Jade (custom green-forward dark)
var ThemeOsakaJade = Theme{
	Background:    ColorFromHex("#0b1d13"),
	BackgroundAlt: ColorFromHex("#10251a"),
	Foreground:    ColorFromHex("#d6f1dd"),
	ForegroundAlt: ColorFromHex("#bfe6cc"),
	Highlight:     ColorFromHex("#153524"),
	Cursor:        ColorFromHex("#d6f1dd"),
	Muted:         ColorFromHex("#6b8f80"),
	Outline:       ColorFromHex("#1e3a29"),
	Red:           ColorFromHex("#ef6a6a"),
	Yellow:        ColorFromHex("#e9d66b"),
	Green:         ColorFromHex("#34d399"),
	Blue:          ColorFromHex("#4fb3ff"),
	Teal:          ColorFromHex("#00b3a4"),
	Primary:       ColorFromHex("#a78bfa"),
	Secondary:     ColorFromHex("#f472b6"),
	Accent:        ColorFromHex("#f4a261"),
}

// Solarized Light
var ThemeSolarizedLight = Theme{
	Background:    ColorFromHex("#fdf6e3"),
	BackgroundAlt: ColorFromHex("#eee8d5"),
	Foreground:    ColorFromHex("#657b83"),
	ForegroundAlt: ColorFromHex("#586e75"),
	Highlight:     ColorFromHex("#eee8d5"),
	Cursor:        ColorFromHex("#586e75"),
	Muted:         ColorFromHex("#93a1a1"),
	Outline:       ColorFromHex("#839496"),
	Red:           ColorFromHex("#dc322f"),
	Yellow:        ColorFromHex("#b58900"),
	Green:         ColorFromHex("#859900"),
	Blue:          ColorFromHex("#268bd2"),
	Teal:          ColorFromHex("#2aa198"),
	Primary:       ColorFromHex("#6c71c4"),
	Secondary:     ColorFromHex("#d33682"),
	Accent:        ColorFromHex("#cb4b16"),
}

// Catppuccin Latte (light)
var CatppuccinLatte = Theme{
	Background:    ColorFromHex("#eff1f5"),
	BackgroundAlt: ColorFromHex("#e6e9ef"),
	Foreground:    ColorFromHex("#4c4f69"),
	ForegroundAlt: ColorFromHex("#5c5f77"),
	Highlight:     ColorFromHex("#ccd0da"),
	Cursor:        ColorFromHex("#4c4f69"),
	Muted:         ColorFromHex("#9ca0b0"),
	Outline:       ColorFromHex("#bcc0cc"),
	Red:           ColorFromHex("#d20f39"),
	Yellow:        ColorFromHex("#df8e1d"),
	Green:         ColorFromHex("#40a02b"),
	Blue:          ColorFromHex("#1e66f5"),
	Teal:          ColorFromHex("#179299"),
	Primary:       ColorFromHex("#8839ef"),
	Secondary:     ColorFromHex("#ea76cb"),
	Accent:        ColorFromHex("#fe640b"),
}

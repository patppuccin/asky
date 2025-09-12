package asky

import (
	"os"
	"strconv"
	"strings"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
)

// --- TTY Standardization --------------------------------
var (
	stdOutput = colorable.NewColorableStdout()
	stdError  = colorable.NewColorableStderr()
	noTTY     = os.Getenv("TERM") == "dumb" ||
		(!isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()))
	noColor = os.Getenv("NO_COLOR") != ""
)

// --- Color Definition ------------------------------------
type color string

func (c color) toSGR(bg bool) (string, bool) {
	if c == "" {
		return "", false
	}
	s := string(c)

	if len(s) > 4 && s[:4] == "rgb:" {
		s := s[4:] // after "rgb:"
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

	if len(s) > 5 && s[:5] == "ansi:" {
		n := strings.TrimPrefix(s, "ansi:")
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

// --- Color Conversion Helpers ----------------------------
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

// --- attribs Definition ------------------------------------
type attribs struct {
	fg, bg                   color
	dim, bold, italic        bool
	underline, strikethrough bool
}

// attribs returns a new attribs with the given attributes.
func (st *attribs) FG(c color) *attribs     { st.fg = c; return st }
func (st *attribs) BG(c color) *attribs     { st.bg = c; return st }
func (st *attribs) Dim() *attribs           { st.dim = true; return st }
func (st *attribs) Bold() *attribs          { st.bold = true; return st }
func (st *attribs) Italic() *attribs        { st.italic = true; return st }
func (st *attribs) Underline() *attribs     { st.underline = true; return st }
func (st *attribs) Strikethrough() *attribs { st.strikethrough = true; return st }
func (st *attribs) isEmpty() bool {
	return !st.bold && !st.dim && !st.italic &&
		!st.underline && !st.strikethrough &&
		st.fg == "" && st.bg == ""
}

// Sprint returns the text with the attribs applied.
func (st *attribs) Sprint(text string) string {
	// If no content, no color, no TTY, or no attributes, return the text as-is
	if text == "" || noColor || noTTY || st.isEmpty() {
		return text
	}

	// Build the ANSI escape sequence
	var b strings.Builder
	b.Grow(len(text) + 24)
	b.WriteString("\x1b[")

	first := true // flag to check escape code presence
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
	if code, ok := st.fg.toSGR(false); ok {
		write(code)
	}
	if code, ok := st.bg.toSGR(true); ok {
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

func NewAttrib() *attribs { return &attribs{} }

// --- Style Definition ------------------------------------
type Style struct {
	theme *Theme

	// Styles for Status Messages
	StatusSuccessPrefix *attribs
	StatusSuccessLabel  *attribs
	StatusDebugPrefix   *attribs
	StatusDebugLabel    *attribs
	StatusInfoPrefix    *attribs
	StatusInfoLabel     *attribs
	StatusWarnPrefix    *attribs
	StatusWarnLabel     *attribs
	StatusErrorPrefix   *attribs
	StatusErrorLabel    *attribs

	// Styles for Banners
	BannerLabel           *attribs
	BannerLabelPadChar    *attribs
	BannerSubLabel        *attribs
	BannerSubLabelPadChar *attribs

	// Styles for Text & Secure Input Prompts
	InputDesc           *attribs
	InputPrefix         *attribs
	InputLabel          *attribs
	InputPlaceholder    *attribs
	InputText           *attribs
	InputValidationPass *attribs
	InputValidationFail *attribs
	InputHelp           *attribs

	// Styles for Confirmation Prompts
	ConfirmationPrefix         *attribs
	ConfirmationLabel          *attribs
	ConfirmationDesc           *attribs
	ConfirmationHelp           *attribs
	ConfirmationSelectedItem   *attribs
	ConfirmationUnselectedItem *attribs

	// Styles for Selection Prompts
	Selectionprefix             *attribs
	SelectionLabel              *attribs
	SelectionDesc               *attribs
	SelectionHelp               *attribs
	SelectionSearchHint         *attribs
	SelectionValidationPass     *attribs
	SelectionValidationFail     *attribs
	SelectionListItemHeader     *attribs
	SelectionListItemLabel      *attribs
	SelectionCurrentItemMarker  *attribs
	SelectionCurrentItemLabel   *attribs
	SelectionSelectedItemMarker *attribs
	SelectionSelectedItemLabel  *attribs
	SelectionDisabledItemMarker *attribs
	SelectionDisabledItemLabel  *attribs

	// Styles for Spinners
	SpinnerPrefix *attribs
	SpinnerLabel  *attribs
	SpinnerDesc   *attribs

	// Styles for Progress Bars
	ProgressPrefix     *attribs
	ProgressLabel      *attribs
	ProgressDesc       *attribs
	ProgressBarPad     *attribs
	ProgressBarDone    *attribs
	ProgressBarPending *attribs
	ProgressBarStatus  *attribs
}

func StyleDefault(theme *Theme) *Style {
	return &Style{
		theme: theme,

		// Default Styles for Status Messages
		StatusSuccessPrefix: NewAttrib().FG(theme.Green),
		StatusSuccessLabel:  NewAttrib().FG(theme.Foreground),
		StatusDebugPrefix:   NewAttrib().FG(theme.Muted),
		StatusDebugLabel:    NewAttrib().FG(theme.Foreground),
		StatusInfoPrefix:    NewAttrib().FG(theme.Blue),
		StatusInfoLabel:     NewAttrib().FG(theme.Foreground),
		StatusWarnPrefix:    NewAttrib().FG(theme.Yellow),
		StatusWarnLabel:     NewAttrib().FG(theme.Foreground),
		StatusErrorPrefix:   NewAttrib().FG(theme.Red),
		StatusErrorLabel:    NewAttrib().FG(theme.Foreground),

		// Default Styles for Banners
		BannerLabel:           NewAttrib().FG(theme.Primary),
		BannerLabelPadChar:    NewAttrib().FG(theme.Accent),
		BannerSubLabel:        NewAttrib().FG(theme.Secondary),
		BannerSubLabelPadChar: NewAttrib().FG(theme.Accent),

		// Default Styles for Text & Secure Input Prompts
		InputDesc:           NewAttrib().FG(theme.Accent),
		InputPrefix:         NewAttrib().FG(theme.Primary),
		InputLabel:          NewAttrib().FG(theme.Secondary),
		InputPlaceholder:    NewAttrib().FG(theme.Muted),
		InputText:           NewAttrib().FG(theme.Foreground),
		InputValidationPass: NewAttrib().FG(theme.Green),
		InputValidationFail: NewAttrib().FG(theme.Red),
		InputHelp:           NewAttrib().FG(theme.Muted),

		// Default Styles for Confirmation Prompts
		ConfirmationPrefix:         NewAttrib().FG(theme.Primary),
		ConfirmationLabel:          NewAttrib().FG(theme.Secondary),
		ConfirmationDesc:           NewAttrib().FG(theme.Accent),
		ConfirmationHelp:           NewAttrib().FG(theme.Muted),
		ConfirmationSelectedItem:   NewAttrib().FG(theme.Background).BG(theme.Primary),
		ConfirmationUnselectedItem: NewAttrib().FG(theme.Primary),

		// Default Styles for Selection Prompts
		Selectionprefix:             NewAttrib().FG(theme.Primary),
		SelectionLabel:              NewAttrib().FG(theme.Secondary),
		SelectionDesc:               NewAttrib().FG(theme.Accent),
		SelectionHelp:               NewAttrib().FG(theme.Muted),
		SelectionSearchHint:         NewAttrib().FG(theme.Muted),
		SelectionValidationPass:     NewAttrib().FG(theme.Green),
		SelectionValidationFail:     NewAttrib().FG(theme.Red),
		SelectionListItemHeader:     NewAttrib().FG(theme.Primary),
		SelectionListItemLabel:      NewAttrib().FG(theme.Foreground),
		SelectionCurrentItemMarker:  NewAttrib().FG(theme.Primary),
		SelectionCurrentItemLabel:   NewAttrib().FG(theme.Primary),
		SelectionSelectedItemMarker: NewAttrib().FG(theme.Green),
		SelectionSelectedItemLabel:  NewAttrib().FG(theme.Green),
		SelectionDisabledItemMarker: NewAttrib().FG(theme.Muted),
		SelectionDisabledItemLabel:  NewAttrib().FG(theme.Muted).Strikethrough(),

		// Default Styles for Spinners
		SpinnerPrefix: NewAttrib().FG(theme.Primary),
		SpinnerLabel:  NewAttrib().FG(theme.Secondary),
		SpinnerDesc:   NewAttrib().FG(theme.Accent),

		// Default Styles for Progress Bars
		ProgressPrefix:     NewAttrib().FG(theme.Primary),
		ProgressLabel:      NewAttrib().FG(theme.Secondary),
		ProgressDesc:       NewAttrib().FG(theme.Accent),
		ProgressBarPad:     NewAttrib().FG(theme.Secondary),
		ProgressBarDone:    NewAttrib().FG(theme.Green),
		ProgressBarPending: NewAttrib().FG(theme.Yellow),
		ProgressBarStatus:  NewAttrib().FG(theme.Secondary),
	}
}

// --- Theme Definition ------------------------------------
type Theme struct {
	// Base
	Background    color
	BackgroundAlt color
	Foreground    color
	ForegroundAlt color

	// Emphasis
	Primary   color
	Secondary color
	Accent    color
	Highlight color
	Muted     color

	// Semantic
	Red    color
	Green  color
	Yellow color
	Blue   color
	Purple color
	Orange color
}

// --- Theme Presets ---------------------------------------
var ThemeDefault = Theme{
	// Base
	Background:    ColorFromANSI(0),  // black
	BackgroundAlt: ColorFromANSI(8),  // bright black (gray)
	Foreground:    ColorFromANSI(15), // bright white
	ForegroundAlt: ColorFromANSI(7),  // white

	// Emphasis
	Primary:   ColorFromANSI(4), // blue
	Secondary: ColorFromANSI(5), // magenta/purple
	Accent:    ColorFromANSI(6), // cyan (close enough for accent)
	Highlight: ColorFromANSI(3), // yellow
	Muted:     ColorFromANSI(8), // gray

	// Semantic
	Red:    ColorFromANSI(9),  // bright red
	Green:  ColorFromANSI(10), // bright green
	Yellow: ColorFromANSI(11), // bright yellow
	Blue:   ColorFromANSI(12), // bright blue
	Purple: ColorFromANSI(13), // bright magenta
	Orange: ColorFromANSI(3),  // reuse yellow as orange (ANSI has no native orange)
}

var ThemeCatppuccinMocha = Theme{
	Background:    ColorFromHex("#1e1e2e"), // base
	BackgroundAlt: ColorFromHex("#181825"), // mantle
	Foreground:    ColorFromHex("#cdd6f4"), // text
	ForegroundAlt: ColorFromHex("#bac2de"), // subtext1

	Primary:   ColorFromHex("#cba6f7"), // mauve
	Secondary: ColorFromHex("#f2cdcd"), // flamingo
	Accent:    ColorFromHex("#fab387"), // peach
	Highlight: ColorFromHex("#f5e0dc"), // rosewater
	Muted:     ColorFromHex("#6c7086"), // overlay0

	Red:    ColorFromHex("#f38ba8"),
	Green:  ColorFromHex("#a6e3a1"),
	Yellow: ColorFromHex("#f9e2af"),
	Blue:   ColorFromHex("#89b4fa"),
	Purple: ColorFromHex("#cba6f7"),
	Orange: ColorFromHex("#fab387"),
}

var ThemeCatppuccinLatte = Theme{
	Background:    ColorFromHex("#eff1f5"), // base
	BackgroundAlt: ColorFromHex("#e6e9ef"), // mantle
	Foreground:    ColorFromHex("#4c4f69"), // text
	ForegroundAlt: ColorFromHex("#5c5f77"), // subtext1

	Primary:   ColorFromHex("#8839ef"), // mauve
	Secondary: ColorFromHex("#ea76cb"), // flamingo
	Accent:    ColorFromHex("#fe640b"), // peach
	Highlight: ColorFromHex("#dc8a78"), // rosewater
	Muted:     ColorFromHex("#9ca0b0"), // overlay0

	Red:    ColorFromHex("#d20f39"),
	Green:  ColorFromHex("#40a02b"),
	Yellow: ColorFromHex("#df8e1d"),
	Blue:   ColorFromHex("#1e66f5"),
	Purple: ColorFromHex("#8839ef"),
	Orange: ColorFromHex("#fe640b"),
}

var ThemeGruvboxDark = Theme{
	Background:    ColorFromHex("#282828"),
	BackgroundAlt: ColorFromHex("#3c3836"),
	Foreground:    ColorFromHex("#ebdbb2"),
	ForegroundAlt: ColorFromHex("#d5c4a1"),

	Primary:   ColorFromHex("#b16286"), // purple
	Secondary: ColorFromHex("#d3869b"), // pinkish
	Accent:    ColorFromHex("#fe8019"), // orange
	Highlight: ColorFromHex("#fabd2f"), // yellow
	Muted:     ColorFromHex("#a89984"),

	Red:    ColorFromHex("#fb4934"),
	Green:  ColorFromHex("#b8bb26"),
	Yellow: ColorFromHex("#fabd2f"),
	Blue:   ColorFromHex("#83a598"),
	Purple: ColorFromHex("#d3869b"),
	Orange: ColorFromHex("#fe8019"),
}

var ThemeTokyoNight = Theme{
	Background:    ColorFromHex("#1a1b26"),
	BackgroundAlt: ColorFromHex("#24283b"),
	Foreground:    ColorFromHex("#c0caf5"),
	ForegroundAlt: ColorFromHex("#a9b1d6"),

	Primary:   ColorFromHex("#7aa2f7"), // blue
	Secondary: ColorFromHex("#bb9af7"), // purple
	Accent:    ColorFromHex("#7dcfff"), // cyan
	Highlight: ColorFromHex("#e0af68"), // yellow/orange
	Muted:     ColorFromHex("#565f89"),

	Red:    ColorFromHex("#f7768e"),
	Green:  ColorFromHex("#9ece6a"),
	Yellow: ColorFromHex("#e0af68"),
	Blue:   ColorFromHex("#7aa2f7"),
	Purple: ColorFromHex("#bb9af7"),
	Orange: ColorFromHex("#ff9e64"),
}

var ThemeKanagawa = Theme{
	Background:    ColorFromHex("#1f1f28"),
	BackgroundAlt: ColorFromHex("#2a2a37"),
	Foreground:    ColorFromHex("#dcd7ba"),
	ForegroundAlt: ColorFromHex("#c8c093"),

	Primary:   ColorFromHex("#7e9cd8"), // blue
	Secondary: ColorFromHex("#957fb8"), // purple
	Accent:    ColorFromHex("#ffa066"), // orange
	Highlight: ColorFromHex("#e6c384"), // yellow
	Muted:     ColorFromHex("#727169"),

	Red:    ColorFromHex("#c34043"),
	Green:  ColorFromHex("#98bb6c"),
	Yellow: ColorFromHex("#c0a36e"),
	Blue:   ColorFromHex("#7e9cd8"),
	Purple: ColorFromHex("#957fb8"),
	Orange: ColorFromHex("#ffa066"),
}

var ThemeDracula = Theme{
	Background:    ColorFromHex("#282a36"),
	BackgroundAlt: ColorFromHex("#44475a"),
	Foreground:    ColorFromHex("#f8f8f2"),
	ForegroundAlt: ColorFromHex("#e2e2dc"),

	Primary:   ColorFromHex("#bd93f9"), // purple
	Secondary: ColorFromHex("#ff79c6"), // pink
	Accent:    ColorFromHex("#50fa7b"), // green
	Highlight: ColorFromHex("#f1fa8c"), // yellow
	Muted:     ColorFromHex("#6272a4"),

	Red:    ColorFromHex("#ff5555"),
	Green:  ColorFromHex("#50fa7b"),
	Yellow: ColorFromHex("#f1fa8c"),
	Blue:   ColorFromHex("#8be9fd"),
	Purple: ColorFromHex("#bd93f9"),
	Orange: ColorFromHex("#ffb86c"),
}

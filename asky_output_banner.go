package asky

type Alignment int

const (
	AlignLeft Alignment = iota
	AlignCenter
	AlignRight
)

type FrameStyle int

const (
	FrameStyleRounded FrameStyle = iota
	FrameStyleSquare
	FrameStyleNone
)

type Banner struct {
	Theme      Theme
	Label      string
	SubLabel   string
	Alignment  Alignment
	FrameStyle FrameStyle
}

func NewBanner() Banner {
	return Banner{
		Theme:      ThemeDefault,
		Label:      "",
		SubLabel:   "",
		Alignment:  AlignLeft,
		FrameStyle: FrameStyleRounded,
	}
}

func (b Banner) WithTheme(theme Theme) Banner                { b.Theme = theme; return b }
func (b Banner) WithLabel(label string) Banner               { b.Label = label; return b }
func (b Banner) WithSubLabel(subLabel string) Banner         { b.SubLabel = subLabel; return b }
func (b Banner) WithAlignment(alignment Alignment) Banner    { b.Alignment = alignment; return b }
func (b Banner) WithFrameStyle(frameStyle FrameStyle) Banner { b.FrameStyle = frameStyle; return b }

func (b Banner) Render() {
	if b.Label == "" {
		return
	}
}

package asky

import (
	"strings"

	"github.com/mattn/go-runewidth"
)

// --- Definitions -----------------------------------------
type Alignment int

const (
	AlignLeft Alignment = iota
	AlignCenter
	AlignRight
)

type Banner struct {
	theme           *Theme
	style           *Style
	label           string
	labelOffset     int
	labelPadChar    string
	subLabel        string
	subLabelOffset  int
	subLabelPadChar string
	width           int
	alignment       Alignment
}

// --- Initialization --------------------------------------
func NewBanner() Banner {
	return Banner{
		label:           "",
		labelOffset:     0,
		labelPadChar:    " ",
		subLabel:        "",
		subLabelOffset:  0,
		subLabelPadChar: " ",
		alignment:       AlignLeft,
	}
}

// --- Configuration ---------------------------------------
func (bn Banner) WithTheme(theme Theme) Banner      { bn.theme = &theme; return bn }
func (bn Banner) WithStyle(style Style) Banner      { bn.style = &style; return bn }
func (bn Banner) WithLabel(label string) Banner     { bn.label = label; return bn }
func (bn Banner) WithLabelOffset(offset int) Banner { bn.labelOffset = max(0, offset); return bn }
func (bn Banner) WithLabelPadChar(padChar string) Banner {
	if runewidth.StringWidth(padChar) < 1 {
		bn.labelPadChar = " "
	} else {
		bn.labelPadChar = padChar
	}
	return bn
}
func (bn Banner) WithSubLabel(subLabel string) Banner  { bn.subLabel = subLabel; return bn }
func (bn Banner) WithSubLabelOffset(offset int) Banner { bn.subLabelOffset = max(0, offset); return bn }
func (bn Banner) WithSubLabelPadChar(padChar string) Banner {
	if runewidth.StringWidth(padChar) < 1 {
		bn.subLabelPadChar = " "
	} else {
		bn.subLabelPadChar = padChar
	}
	return bn
}
func (bn Banner) WithWidth(width int) Banner               { bn.width = min(0, width); return bn }
func (bn Banner) WithAlignment(alignment Alignment) Banner { bn.alignment = alignment; return bn }

// --- Presentation ----------------------------------------
func (bn Banner) Render() {
	// Sanity check to skip render if both label and subLabel are empty
	if bn.label == "" && bn.subLabel == "" {
		return
	}

	// Setup theme and style (apply defaults if not set)
	if bn.theme == nil {
		bn.theme = &ThemeDefault
	}
	if bn.style == nil {
		bn.style = StyleDefault(bn.theme)
	}

	// Render the banner with the configured label and subLabel
	if bn.label != "" {
		line := padLine(bn.style.BannerLabelPadChar, bn.style.BannerLabel, bn.label, bn.alignment, bn.labelPadChar, bn.labelOffset)
		stdOutput.Write([]byte(line + "\n"))
	}
	if bn.subLabel != "" {
		line := padLine(bn.style.BannerSubLabelPadChar, bn.style.BannerSubLabel, bn.subLabel, bn.alignment, bn.subLabelPadChar, bn.subLabelOffset)
		stdOutput.Write([]byte(line + "\n"))
	}
}

// --- Helpers ---------------------------------------------
func repeatPadChar(padChar string, padWidth int) string {
	// Fallback to space if padChar has no display width
	if runewidth.StringWidth(padChar) < 1 {
		padChar = " "
	}

	var b strings.Builder
	curWidth := 0

	// Repeat padChar until we reach the target display width
	for curWidth < padWidth {
		b.WriteString(padChar)
		curWidth += runewidth.StringWidth(padChar)
	}

	// Trim excess if the last padChar overshot the width (CJK, emojis, ligatures, etc.)
	result := b.String()
	if curWidth > padWidth {
		result = runewidth.Truncate(result, padWidth, "")
	}

	return result
}

func padLine(padStyle *attribs, contentStyle *attribs, content string, alignment Alignment, padChar string, offset int) string {
	// Get the terminal width
	termWidth, _, err := getTermDimensions()
	if err != nil || termWidth <= 0 {
		termWidth = 80
	}

	spacedContent := " " + content + " "
	spacedContentWidth := runewidth.StringWidth(spacedContent)

	if spacedContentWidth+offset*2 > termWidth {
		avail := termWidth - offset*2 - 2
		avail = max(0, avail)
		trunc := runewidth.Truncate(content, avail, "...")
		spacedContent = " " + trunc + " "
		spacedContentWidth = runewidth.StringWidth(spacedContent)
	}

	space := max(0, termWidth-spacedContentWidth)
	leftWidth := 0
	rightWidth := 0
	switch alignment {
	case AlignCenter:
		leftWidth = space / 2
		rightWidth = space - leftWidth
	case AlignRight:
		rightWidth = offset
		leftWidth = space - rightWidth
	default: // defaults to left alignment
		leftWidth = offset
		rightWidth = space - leftWidth
	}

	return padStyle.Sprint(repeatPadChar(padChar, leftWidth)) +
		contentStyle.Sprint(spacedContent) +
		padStyle.Sprint(repeatPadChar(padChar, rightWidth))
}

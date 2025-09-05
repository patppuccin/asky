package asky

import "os"

type StatusLevel int

const (
	StatusLevelDefault StatusLevel = iota
	StatusLevelSuccess
	StatusLevelDebug
	StatusLevelInfo
	StatusLevelWarn
	StatusLevelError
)

type Status struct {
	Theme  Theme
	Prefix string
	Label  string
	Level  StatusLevel
}

func NewStatus() Status {
	return Status{
		Theme:  ThemeDefault,
		Prefix: "",
		Label:  "",
		Level:  StatusLevelDefault,
	}
}

func (s Status) WithTheme(theme Theme) Status       { s.Theme = theme; return s }
func (s Status) WithPrefix(prefix string) Status    { s.Prefix = prefix; return s }
func (s Status) WithLabel(label string) Status      { s.Label = label; return s }
func (s Status) WithLevel(level StatusLevel) Status { s.Level = level; return s }

func (s Status) Render() {
	if s.Label == "" {
		return
	}

	var styledPrefix string
	switch s.Level {
	case StatusLevelSuccess:
		if s.Prefix == "" {
			s.Prefix = "[✓] "
		}
		styledPrefix = s.Theme.SuccessStyle(s.Prefix)
	case StatusLevelDebug:
		if s.Prefix == "" {
			s.Prefix = "[-] "
		}
		styledPrefix = s.Theme.MutedStyle(s.Prefix)
	case StatusLevelInfo:
		if s.Prefix == "" {
			s.Prefix = "[i] "
		}
		styledPrefix = s.Theme.InfoStyle(s.Prefix)
	case StatusLevelWarn:
		if s.Prefix == "" {
			s.Prefix = "[!] "
		}
		styledPrefix = s.Theme.WarningStyle(s.Prefix)
	case StatusLevelError:
		if s.Prefix == "" {
			s.Prefix = "[x] "
		}
		styledPrefix = s.Theme.ErrorStyle(s.Prefix)
	default:
		if s.Prefix == "" {
			s.Prefix = "[~] "
		}
		styledPrefix = s.Theme.PrimaryStyle(s.Prefix)
	}

	styledLabel := s.Theme.NeutralStyle(s.Label)
	os.Stdout.WriteString(styledPrefix + styledLabel + "\n")
}

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
	p := newPreset(s.Theme)
	var styledPrefix string
	if s.Label == "" {
		return
	}

	getPrefix := func(px string) string {
		if s.Prefix == "" {
			return px
		} else {
			return s.Prefix
		}
	}

	switch s.Level {
	case StatusLevelSuccess:
		styledPrefix = p.success.Sprint(getPrefix("[✓] "))
	case StatusLevelDebug:
		styledPrefix = p.debug.Sprint(getPrefix("[-] "))
	case StatusLevelInfo:
		styledPrefix = p.info.Sprint(getPrefix("[i] "))
	case StatusLevelWarn:
		styledPrefix = p.warn.Sprint(getPrefix("[!] "))
	case StatusLevelError:
		styledPrefix = p.err.Sprint(getPrefix("[x] "))
	default:
		styledPrefix = p.neutral.Sprint(getPrefix("[~] "))
	}

	os.Stdout.WriteString(styledPrefix + p.neutral.Sprint(s.Label) + "\n")
}

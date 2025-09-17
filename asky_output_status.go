package asky

// --- Definition -----------------------------------------
type statusLevel int

const (
	StatusLevelDebug statusLevel = iota
	StatusLevelSuccess
	StatusLevelInfo
	StatusLevelWarn
	StatusLevelError
)

type status struct {
	theme  *Theme
	style  *Style
	prefix string
	label  string
	level  statusLevel
}

// --- Initialization --------------------------------------
func NewStatus() *status {
	return &status{
		prefix: "",
		label:  "",
		level:  StatusLevelDebug,
	}
}

// --- Configuration ---------------------------------------
func (st status) WithTheme(theme Theme) status       { st.theme = &theme; return st }
func (st status) WithStyle(style Style) status       { st.style = &style; return st }
func (st status) WithPrefix(prefix string) status    { st.prefix = prefix; return st }
func (st status) WithLabel(label string) status      { st.label = label; return st }
func (st status) WithLevel(level statusLevel) status { st.level = level; return st }

// --- Presentation ----------------------------------------
func (st status) getPrefix(px string) string {
	if st.prefix == "" {
		return px
	}
	return st.prefix
}

func (st status) Render() {
	// Sanity check to skip render if both label and prefix are empty
	if st.label == "" && st.prefix == "" {
		return
	}

	// Setup theme and style (apply defaults if not set)
	if st.theme == nil {
		st.theme = &ThemeDefault
	}
	if st.style == nil {
		st.style = StyleDefault(st.theme)
	}

	// Construct the styled prefix and label (as per the status level)
	var styledPrefix string
	var styledLabel string
	switch st.level {
	case StatusLevelSuccess:
		styledPrefix = st.style.StatusSuccessPrefix.Sprint(st.getPrefix("[✓] "))
		styledLabel = st.style.StatusSuccessLabel.Sprint(st.label)
	case StatusLevelInfo:
		styledPrefix = st.style.StatusInfoPrefix.Sprint(st.getPrefix("[i] "))
		styledLabel = st.style.StatusInfoLabel.Sprint(st.label)
	case StatusLevelWarn:
		styledPrefix = st.style.StatusWarnPrefix.Sprint(st.getPrefix("[!] "))
		styledLabel = st.style.StatusWarnLabel.Sprint(st.label)
	case StatusLevelError:
		styledPrefix = st.style.StatusErrorPrefix.Sprint(st.getPrefix("[x] "))
		styledLabel = st.style.StatusErrorLabel.Sprint(st.label)
	default:
		styledPrefix = st.style.StatusDebugPrefix.Sprint(st.getPrefix("[-] "))
		styledLabel = st.style.StatusDebugLabel.Sprint(st.label)
	}

	// Render the styled prefix and label
	stdOutput.Write([]byte(styledPrefix + styledLabel + "\n"))
}

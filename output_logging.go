package asky

import "github.com/fatih/color"

// ==== Log Message ===========================================================+

// logMessage prints a single styled log line with a level prefix.
// Construct one with [LogMessage].
type logMessage struct {
	cfg    Config
	prefix string
}

// LogMessage returns a builder for printing a single styled log line.
//
//	asky.LogMessage().Info("server started")
//	asky.LogMessage().WithPrefix("(done)").Success("deployment complete")
func LogMessage() *logMessage {
	return &logMessage{cfg: pkgConfig}
}

// WithStyles overrides the [StyleMap] for this message.
func (l *logMessage) WithStyles(s *StyleMap) *logMessage {
	l.cfg.Styles = s
	return l
}

// WithPrefix overrides the default level prefix symbol.
func (l *logMessage) WithPrefix(p string) *logMessage {
	l.prefix = p
	return l
}

// Success prints a success message.
func (l *logMessage) Success(msg string) {
	l.render(l.cfg.Styles.LogSuccessPrefix, l.cfg.Styles.LogSuccessLabel, "(✓)", msg)
}

// Debug prints a debug message.
func (l *logMessage) Debug(msg string) {
	l.render(l.cfg.Styles.LogDebugPrefix, l.cfg.Styles.LogDebugLabel, "(~)", msg)
}

// Info prints an info message.
func (l *logMessage) Info(msg string) {
	l.render(l.cfg.Styles.LogInfoPrefix, l.cfg.Styles.LogInfoLabel, "(i)", msg)
}

// Warn prints a warning message.
func (l *logMessage) Warn(msg string) {
	l.render(l.cfg.Styles.LogWarnPrefix, l.cfg.Styles.LogWarnLabel, "(!)", msg)
}

// Error prints an error message.
func (l *logMessage) Error(msg string) {
	l.render(l.cfg.Styles.LogErrorPrefix, l.cfg.Styles.LogErrorLabel, "(✗)", msg)
}

func (l *logMessage) render(pfxStyle, labelStyle *color.Color, defaultPfx, msg string) {
	pfx := safeStyle(pfxStyle).Sprint(pick(l.prefix, defaultPfx))
	label := safeStyle(labelStyle).Sprint(msg)
	stdOutput.Write([]byte(pfx + " " + label + "\n"))
}

// ==== Log Group ==============================================================

// logGroup prints a styled title line followed by indented message lines.
// Construct one with [LogGroup].
type logGroup struct {
	cfg    Config
	prefix string
}

// LogGroup returns a builder for printing a styled title with indented body lines.
//
//	asky.LogGroup().Info("config loaded", "host: localhost", "port: 8080")
//	asky.LogGroup().WithPrefix("DONE:").Success("deploy finished", "3 services restarted")
func LogGroup() *logGroup {
	return &logGroup{cfg: pkgConfig}
}

// WithStyles overrides the [StyleMap] for this group.
func (l *logGroup) WithStyles(s *StyleMap) *logGroup {
	l.cfg.Styles = s
	return l
}

// WithPrefix overrides the default level prefix label.
func (l *logGroup) WithPrefix(p string) *logGroup {
	l.prefix = p
	return l
}

// Success prints a success group.
func (l *logGroup) Success(title string, msgs ...string) {
	l.render(l.cfg.Styles.LogSuccessPrefix, l.cfg.Styles.LogSuccessLabel, "SUCCESS:", title, msgs...)
}

// Debug prints a debug group.
func (l *logGroup) Debug(title string, msgs ...string) {
	l.render(l.cfg.Styles.LogDebugPrefix, l.cfg.Styles.LogDebugLabel, "DEBUG:", title, msgs...)
}

// Info prints an info group.
func (l *logGroup) Info(title string, msgs ...string) {
	l.render(l.cfg.Styles.LogInfoPrefix, l.cfg.Styles.LogInfoLabel, "INFO:", title, msgs...)
}

// Warn prints a warning group.
func (l *logGroup) Warn(title string, msgs ...string) {
	l.render(l.cfg.Styles.LogWarnPrefix, l.cfg.Styles.LogWarnLabel, "WARN:", title, msgs...)
}

// Error prints an error group.
func (l *logGroup) Error(title string, msgs ...string) {
	l.render(l.cfg.Styles.LogErrorPrefix, l.cfg.Styles.LogErrorLabel, "ERROR:", title, msgs...)
}

func (l *logGroup) render(pfxStyle, labelStyle *color.Color, defaultPfx, title string, msgs ...string) {
	pfx := safeStyle(pfxStyle).Sprint(pick(l.prefix, defaultPfx))
	titleStr := safeStyle(labelStyle).Sprint(title)
	stdOutput.Write([]byte(pfx + " " + titleStr + "\n"))
	for _, msg := range msgs {
		stdOutput.Write([]byte("  " + safeStyle(l.cfg.Styles.LogBlockBody).Sprint(msg) + "\n"))
	}
}

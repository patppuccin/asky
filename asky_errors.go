package asky

import "errors"

var ErrInterrupted = errors.New("prompt interrupted")
var ErrTerminalTooSmall = errors.New("terminal dimensions too small")
var ErrNoSelectionChoices = errors.New("no choices supplied for selection prompt")
var ErrInvalidSelectionCount = errors.New("max count > min count for multi select prompt")

package asky

import "errors"

var ErrInterrupted = errors.New("prompt interrupted")
var ErrTerminalTooSmall = errors.New("terminal dimensions too small")

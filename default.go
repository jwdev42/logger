//This file is part of logger. ©2020-2023 Jörg Walter.

package logger

import "io"

// Stores the default logger
var defaultLogger *Logger

// Returns the default logger, will panic if SetupDefaultLogger() has not been called yet.
func Default() *Logger {
	if defaultLogger == nil {
		panic("Programming error: Default logger is not set, please set it up first by calling SetupDefaultLogger()")
	}
	return defaultLogger
}

// Creates a default logger that can be used by calling Default().
func SetupDefaultLogger(w io.Writer, level Level, delimiter string) {
	defaultLogger = New(w, level, delimiter)
}

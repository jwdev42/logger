//This file is part of logger. ©2020-2023 Jörg Walter.

package logger

import (
	"fmt"
	"strings"
)

const (
	LevelInvalid  Level = iota //Not a loglevel, represents an invalid loglevel internally.
	LevelPanic                 //Reports a crash or a program termination that was necessary to prevent something more severe.
	LevelAlert                 //Reports a condition that needs immediate action by the user.
	LevelCritical              //Reports a hard error.
	LevelError                 //Reports a normal error condition.
	LevelWarning               //Reports something that isn't an error but might cause problems later.
	LevelNotice                //Reports conditions that are no errors but may require special handling.
	LevelInfo                  //Reports normal user information about expected conditions.
	LevelDebug                 //Reports information that is only interesting for debugging or development.
)

// Represents a loglevel.
type Level int

// Panics if the loglevel does not exist.
func assertLoglevel(lvl Level) {
	if lvl < LevelPanic || lvl > LevelDebug {
		panic(fmt.Sprintf("Log level %d is not defined", lvl))
	}
}

// Tries to associate the input string with a specific loglevel.
// Returns that loglevel on success, on failure LevelInvalid and an
// error is returned.
func ParseLevel(input string) (Level, error) {
	switch lvl := strings.ToLower(input); lvl {
	case "panic":
		return LevelPanic, nil
	case "alert":
		return LevelAlert, nil
	case "critical":
		return LevelCritical, nil
	case "error":
		return LevelError, nil
	case "warning":
		return LevelWarning, nil
	case "notice":
		return LevelNotice, nil
	case "info":
		return LevelInfo, nil
	case "debug":
		return LevelDebug, nil
	}
	return LevelInvalid, fmt.Errorf("Input sequence %q cannot be associated with a defined loglevel", input)
}

// String returns the string representation of a Level. If the Level is
// not defined, String returns "Undefined".
func (r Level) String() string {
	switch r {
	case LevelPanic:
		return "Panic"
	case LevelAlert:
		return "Alert"
	case LevelCritical:
		return "Critical"
	case LevelError:
		return "Error"
	case LevelWarning:
		return "Warning"
	case LevelNotice:
		return "Notice"
	case LevelInfo:
		return "Info"
	case LevelDebug:
		return "Debug"
	}
	return "Undefined"
}

// Loglevels returns map with the string representations of all
// available loglevels.
func Loglevels() map[Level]string {
	levels := make(map[Level]string)
	for lvl := LevelPanic; lvl <= LevelDebug; lvl++ {
		levels[lvl] = lvl.String()
	}
	return levels
}

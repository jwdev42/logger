//This file is part of logger. ©2020-2023 Jörg Walter.

// The package logger provides a Logger type that can receive log records. Logger will only display a record if it is as severe as
// or more severe than the loglevel set for the Logger. Users of Logger therefore can control how much logging output they will see
// while running their program. A Logger can be used by multiple goroutines.
//
// For usable loglevels see CONSTANTS.
package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	LevelPanic = iota
	LevelAlert
	LevelCritical
	LevelError
	LevelWarning
	LevelNotice
	LevelInfo
	LevelDebug
)

var level2str = map[int]string{
	LevelPanic:    "[Panic]",
	LevelAlert:    "[Alert]",
	LevelCritical: "[Critical]",
	LevelError:    "[Error]",
	LevelWarning:  "[Warning]",
	LevelNotice:   "[Notice]",
	LevelInfo:     "[Info]",
	LevelDebug:    "[Debug]",
}

// Panics if the loglevel does not exist.
func assertLoglevel(level int) {
	if level < LevelPanic || level > LevelDebug {
		panic(fmt.Sprintf("Log level %d is not defined", level))
	}
}

// Logger is the data type used for sending log records to.
type Logger struct {
	mu         *sync.Mutex
	delimiter  string
	timeFormat string
	level      int
	out        io.Writer
}

// New constructs a new Logger. It will print a log record to its given writer if it fulfills the
// Logger's designated loglevel or a more severe one. If you set the level to LevelCritical, the Logger
// will print all messages of LevelPanic or LevelAlert or LevelCritical.
func New(w io.Writer, level int, delimiter string) *Logger {
	if len(delimiter) < 1 {
		panic("delimiter must have one or more characters")
	}
	assertLoglevel(level)
	return &Logger{
		delimiter: delimiter,
		level:     level,
		mu:        new(sync.Mutex),
		out:       w,
	}
}

// Alert sends a message of loglevel LevelAlert to the Logger.
func (l *Logger) Alert(v ...any) (n int, err error) {
	return l.Println(LevelAlert, v...)
}

// Alertf sends a formatted message of loglevel LevelAlert to the Logger.
func (l *Logger) Alertf(format string, a ...any) (n int, err error) {
	return l.Printf(LevelAlert, format, a...)
}

// Critical sends a message of loglevel LevelCritical to the Logger.
func (l *Logger) Critical(v ...any) (n int, err error) {
	return l.Println(LevelCritical, v...)
}

// Criticalf sends a formatted message of loglevel LevelCritical to the Logger.
func (l *Logger) Criticalf(format string, a ...any) (n int, err error) {
	return l.Printf(LevelCritical, format, a...)
}

// Die sends a message of loglevel LevelPanic to the Logger, then exits with code 1.
func (l *Logger) Die(v ...any) {
	l.Panic(v...)
	os.Exit(1)
}

// Dief sends a formatted message of loglevel LevelPanic to the Logger, then exits with code 1.
func (l *Logger) Dief(format string, a ...any) {
	l.Panicf(format, a...)
	os.Exit(1)
}

// Debug sends a message of loglevel LevelDebug to the Logger.
func (l *Logger) Debug(v ...any) (n int, err error) {
	return l.Println(LevelDebug, v...)
}

// Debugf sends a formatted message of loglevel LevelDebug to the Logger.
func (l *Logger) Debugf(format string, a ...any) (n int, err error) {
	return l.Printf(LevelDebug, format, a...)
}

// Error sends a message of loglevel LevelError to the Logger.
func (l *Logger) Error(v ...any) (n int, err error) {
	return l.Println(LevelError, v...)
}

// Errorf sends a formatted message of loglevel LevelError to the Logger.
func (l *Logger) Errorf(format string, a ...any) (n int, err error) {
	return l.Printf(LevelError, format, a...)
}

// Info sends a message of loglevel LevelInfo to the Logger.
func (l *Logger) Info(v ...any) (n int, err error) {
	return l.Println(LevelInfo, v...)
}

// Infof sends a formatted message of loglevel LevelInfo to the Logger.
func (l *Logger) Infof(format string, a ...any) (n int, err error) {
	return l.Printf(LevelInfo, format, a...)
}

// Level returns the Logger's current loglevel as an integer.
func (l *Logger) Level() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.level
}

// LevelStr returns the string representation of the Logger's current loglevel.
func (l *Logger) LevelStr() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return level2str[l.level]
}

// Notice sends a message of loglevel LevelNotice to the Logger.
func (l *Logger) Notice(v ...any) (n int, err error) {
	return l.Println(LevelNotice, v...)
}

// Noticef sends a formatted message of loglevel LevelNotice to the Logger.
func (l *Logger) Noticef(format string, a ...any) (n int, err error) {
	return l.Printf(LevelNotice, format, a...)
}

// Panic sends a message of loglevel LevelPanic to the Logger.
// Please note that it does NOT call panic()!
func (l *Logger) Panic(v ...any) (n int, err error) {
	return l.Println(LevelPanic, v...)
}

// Panicf sends a formatted message of loglevel LevelPanic to the Logger.
// Please note that it does NOT call panic()!
func (l *Logger) Panicf(format string, a ...any) (n int, err error) {
	return l.Printf(LevelPanic, format, a...)
}

// Println writes the log message if its log level is equally severe or more severe than that set for the Logger.
func (l *Logger) Println(level int, v ...any) (n int, err error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if !l.trigger(level) {
		return 0, nil
	}
	if len(l.timeFormat) > 0 {
		return fmt.Fprintf(l.out,
			"%s%s%s%s%s\n",
			level2str[level],
			l.delimiter,
			time.Now().Format(l.timeFormat),
			l.delimiter,
			fmt.Sprint(v...))
	} else {
		return fmt.Fprintf(l.out,
			"%s%s%s\n",
			level2str[level],
			l.delimiter,
			fmt.Sprint(v...))
	}
}

// Printf writes a formatted log message if the logger was configured to print the given level.
func (l *Logger) Printf(level int, format string, a ...any) (n int, err error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if !l.trigger(level) {
		return 0, nil
	}
	if len(l.timeFormat) > 0 {
		return fmt.Fprintf(l.out,
			"%s%s%s%s%s",
			level2str[level],
			l.delimiter,
			time.Now().Format(l.timeFormat),
			l.delimiter,
			l.autoAppendLF(fmt.Sprintf(format, a...)))
	} else {
		return fmt.Fprintf(l.out,
			"%s%s%s",
			level2str[level],
			l.delimiter,
			l.autoAppendLF(fmt.Sprintf(format, a...)))
	}
}

// SetLevel sets a new loglevel for the Logger. Setting an invalid loglevel will cause a panic.
func (l *Logger) SetLevel(level int) {
	assertLoglevel(level)
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// SetOutput changes the writer the Logger will write its messages to.
func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out = w
}

// SetTimeFormat takes a format string as defined in the "(t Time) Format" function of go's "time" module.
// If such a string is set, log records will display a timestamp formatted like specified by the format string.
// To remove timestamps from future log records, set the format string to "".
func (l *Logger) SetTimeFormat(format string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.timeFormat = format
}

// TimeFormat returns the current format string for the timestamp. If it returns "", log records will have no timestamp.
func (l *Logger) TimeFormat() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.timeFormat
}

// Warning sends a message of loglevel LevelWarning to the Logger.
func (l *Logger) Warning(v ...any) (n int, err error) {
	return l.Println(LevelWarning, v...)
}

// Warningf sends a formatted message of loglevel LevelWarning to the Logger.
func (l *Logger) Warningf(format string, a ...any) (n int, err error) {
	return l.Printf(LevelWarning, format, a...)
}

// trigger returns true if the Logger should print a message of loglevel
// level, otherwise it returns false.
func (l *Logger) trigger(level int) bool {
	assertLoglevel(level)
	if level <= l.level {
		return true
	}
	return false
}

// autoAppendLF appends one newline character at the end of input and returns
// a new string if input doesn't already end with a newline character.
func (l *Logger) autoAppendLF(input string) string {
	if strings.HasSuffix(input, "\n") {
		return input
	}
	return input + "\n"
}

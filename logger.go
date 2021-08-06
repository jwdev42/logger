//This file is part of logger. ©2020 Jörg Walter.

//The package logger provides a Logger type that can receive log records. Logger will only display a record if it is as severe as
//or more severe than the loglevel set for the Logger. Users of Logger therefore can control how much logging output they will see
//while running their program. A Logger can be used by multiple goroutines.
//
//For usable loglevels see CONSTANTS.
package logger

import (
	"fmt"
	"io"
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

//Panics if the loglevel does not exist.
func assertLoglevel(level int) {
	if level < LevelPanic || level > LevelDebug {
		panic(fmt.Sprintf("Log level %d is not defined", level))
	}
}

//Logger is the data type used for sending log records to.
type Logger struct {
	mu, muLevel, muTime *sync.Mutex
	delimiter           string
	timeFormat          string
	level               int
	out                 io.Writer
}

//New constructs a new Logger. It will print a log record to its given writer if it fulfills the
//Logger's designated loglevel or a more severe one. If you set the level to LevelCritical, the Logger
//will print all messages of LevelPanic or LevelAlert or LevelCritical.
func New(w io.Writer, level int, delimiter string) *Logger {
	if len(delimiter) < 1 {
		panic("delimiter must have one or more characters")
	}
	assertLoglevel(level)
	return &Logger{
		delimiter: delimiter,
		level:     level,
		mu:        new(sync.Mutex),
		muLevel:   new(sync.Mutex),
		muTime:    new(sync.Mutex),
		out:       w,
	}
}

func (l *Logger) trigger(level int) bool {
	assertLoglevel(level)
	l.muLevel.Lock()
	defer l.muLevel.Unlock()
	if level <= l.level {
		return true
	}
	return false
}

//Alert sends a message of loglevel LevelAlert to the Logger.
func (l *Logger) Alert(v ...interface{}) (n int, err error) {
	return l.Println(LevelAlert, v...)
}

//Critical sends a message of loglevel LevelCritical to the Logger.
func (l *Logger) Critical(v ...interface{}) (n int, err error) {
	return l.Println(LevelCritical, v...)
}

//Debug sends a message of loglevel LevelDebug to the Logger.
func (l *Logger) Debug(v ...interface{}) (n int, err error) {
	return l.Println(LevelDebug, v...)
}

//Error sends a message of loglevel LevelError to the Logger.
func (l *Logger) Error(v ...interface{}) (n int, err error) {
	return l.Println(LevelError, v...)
}

//Info sends a message of loglevel LevelInfo to the Logger.
func (l *Logger) Info(v ...interface{}) (n int, err error) {
	return l.Println(LevelInfo, v...)
}

//Level returns the Logger's current loglevel as an integer.
func (l *Logger) Level() int {
	l.muLevel.Lock()
	defer l.muLevel.Unlock()
	return l.level
}

//LevelStr returns the string representation of the Logger's current loglevel.
func (l *Logger) LevelStr() string {
	l.muLevel.Lock()
	defer l.muLevel.Unlock()
	return level2str[l.level]
}

//Notice sends a message of loglevel LevelNotice to the Logger.
func (l *Logger) Notice(v ...interface{}) (n int, err error) {
	return l.Println(LevelNotice, v...)
}

//Panic sends a message of loglevel LevelPanic to the Logger. Please note that it does NOT call panic()!
func (l *Logger) Panic(v ...interface{}) (n int, err error) {
	return l.Println(LevelPanic, v...)
}

//Println writes the log message if its log level is equally severe or more severe than that set for the Logger.
func (l *Logger) Println(level int, v ...interface{}) (n int, err error) {
	if !l.trigger(level) {
		return 0, nil
	}
	record := strings.TrimRight(fmt.Sprint(v...), "\n")
	dateTimeFormat := l.TimeFormat()
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(dateTimeFormat) <= 0 {
		return fmt.Fprintf(l.out, "%s%s%s\n", level2str[level], l.delimiter, record)
	} else {
		return fmt.Fprintf(l.out, "%s%s%s%s%s\n", level2str[level], l.delimiter, time.Now().Format(dateTimeFormat), l.delimiter, record)
	}
}

//SetLevel sets a new loglevel for the Logger. Setting an invalid loglevel will cause a panic.
func (l *Logger) SetLevel(level int) {
	assertLoglevel(level)
	l.muLevel.Lock()
	defer l.muLevel.Unlock()
	l.level = level
}

//SetOutput changes the writer the Logger will write its messages to.
func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out = w
}

//SetTimeFormat takes a format string as defined in the "(t Time) Format" function of go's "time" module.
//If such a string is set, log records will display a timestamp formatted like specified by the format string.
//To remove timestamps from future log records, set the format string to "".
func (l *Logger) SetTimeFormat(format string) {
	l.muTime.Lock()
	defer l.muTime.Unlock()
	l.timeFormat = format
}

//TimeFormat returns the current format string for the timestamp. If it returns "", log records will have no timestamp.
func (l *Logger) TimeFormat() string {
	l.muTime.Lock()
	defer l.muTime.Unlock()
	return l.timeFormat
}

//Warning sends a message of loglevel LevelWarning to the Logger.
func (l *Logger) Warning(v ...interface{}) (n int, err error) {
	return l.Println(LevelWarning, v...)
}

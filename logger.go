//This file is part of logger.
//©2020 Jörg Walter, License: LGPL-3

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
	Level_Panic = iota
	Level_Alert
	Level_Critical
	Level_Error
	Level_Warning
	Level_Notice
	Level_Info
	Level_Debug
)

var level2str = map[int]string{
	Level_Panic:    "[Panic]",
	Level_Alert:    "[Alert]",
	Level_Critical: "[Critical]",
	Level_Error:    "[Error]",
	Level_Warning:  "[Warning]",
	Level_Notice:   "[Notice]",
	Level_Info:     "[Info]",
	Level_Debug:    "[Debug]",
}

func level_must_be_in_range(level int) {
	if level < Level_Panic || level > Level_Debug {
		panic(fmt.Sprintf("Log level %d is not defined", level))
	}
}

//Logger is the data type used for sending log records to.
type Logger struct {
	mu, mu_level, mu_time *sync.Mutex
	delimiter             string
	time_format           string
	level                 int
	out                   io.Writer
}

//New constructs a new Logger. It will print a log record to its given writer if it fulfills the
//Logger's designated loglevel or a more severe one. If you set the level to Level_Critical, the Logger
//will print all messages of Level_Panic or Level_Alert or Level_Critical.
func New(w io.Writer, level int, delimiter string) *Logger {
	if len(delimiter) < 1 {
		panic("delimiter must have one or more characters")
	}
	level_must_be_in_range(level)
	return &Logger{
		delimiter: delimiter,
		level:     level,
		mu:        new(sync.Mutex),
		mu_level:  new(sync.Mutex),
		mu_time:   new(sync.Mutex),
		out:       w,
	}
}

func (l *Logger) trigger(level int) bool {
	level_must_be_in_range(level)
	l.mu_level.Lock()
	defer l.mu_level.Unlock()
	if level <= l.level {
		return true
	}
	return false
}

//Alert sends a message of loglevel Level_Alert to the Logger.
func (l *Logger) Alert(v ...interface{}) (n int, err error) {
	return l.Println(Level_Alert, v...)
}

//Critical sends a message of loglevel Level_Critical to the Logger.
func (l *Logger) Critical(v ...interface{}) (n int, err error) {
	return l.Println(Level_Critical, v...)
}

//Debug sends a message of loglevel Level_Debug to the Logger.
func (l *Logger) Debug(v ...interface{}) (n int, err error) {
	return l.Println(Level_Debug, v...)
}

//Error sends a message of loglevel Level_Error to the Logger.
func (l *Logger) Error(v ...interface{}) (n int, err error) {
	return l.Println(Level_Error, v...)
}

//Info sends a message of loglevel Level_Info to the Logger.
func (l *Logger) Info(v ...interface{}) (n int, err error) {
	return l.Println(Level_Info, v...)
}

//Level returns the Logger's current loglevel as an integer.
func (l *Logger) Level() int {
	l.mu_level.Lock()
	defer l.mu_level.Unlock()
	return l.level
}

//LevelStr returns the string representation of the Logger's current loglevel.
func (l *Logger) LevelStr() string {
	l.mu_level.Lock()
	defer l.mu_level.Unlock()
	return level2str[l.level]
}

//Notice sends a message of loglevel Level_Notice to the Logger.
func (l *Logger) Notice(v ...interface{}) (n int, err error) {
	return l.Println(Level_Notice, v...)
}

//Panic sends a message of loglevel Level_Panic to the Logger. Please note that it does NOT call panic()!
func (l *Logger) Panic(v ...interface{}) (n int, err error) {
	return l.Println(Level_Panic, v...)
}

//Println writes the log message if its log level is equally severe or more severe than that set for the Logger.
func (l *Logger) Println(level int, v ...interface{}) (n int, err error) {
	if !l.trigger(level) {
		return 0, nil
	}
	record := strings.TrimRight(fmt.Sprint(v...), "\n")
	datetime_f := l.TimeFormat()
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(datetime_f) <= 0 {
		return fmt.Fprintf(l.out, "%s%s%s\n", level2str[level], l.delimiter, record)
	} else {
		return fmt.Fprintf(l.out, "%s%s%s%s%s\n", level2str[level], l.delimiter, time.Now().Format(datetime_f), l.delimiter, record)
	}
}

//SetLevel sets a new loglevel for the Logger. Setting an invalid loglevel will cause a panic.
func (l *Logger) SetLevel(level int) {
	level_must_be_in_range(level)
	l.mu_level.Lock()
	defer l.mu_level.Unlock()
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
	l.mu_time.Lock()
	defer l.mu_time.Unlock()
	l.time_format = format
}

//TimeFormat returns the current format string for the timestamp. If it returns "", log records will have no timestamp.
func (l *Logger) TimeFormat() string {
	l.mu_time.Lock()
	defer l.mu_time.Unlock()
	return l.time_format
}

//Warning sends a message of loglevel Level_Warning to the Logger.
func (l *Logger) Warning(v ...interface{}) (n int, err error) {
	return l.Println(Level_Warning, v...)
}

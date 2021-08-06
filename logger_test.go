//This file is part of logger.
//©2020 Jörg Walter, License: LGPL-3

package logger

import (
	"fmt"
	"strings"
	"testing"
)

const loglevelDelimiter = " - "

func doeslog(t *testing.T, loglevel int, messageLevel int, msg string) {
	b := new(strings.Builder)
	l := New(b, loglevel, loglevelDelimiter)
	l.Println(messageLevel, msg)
	expected := fmt.Sprintf("%s%s%s\n", level2str[messageLevel], loglevelDelimiter, msg)
	if b.String() != expected {
		t.Errorf("Logstring %q is not equal to expected string %q", b.String(), expected)
	}
}

func doesnotlog(t *testing.T, loglevel int, messageLevel int, msg string) {
	b := new(strings.Builder)
	l := New(b, loglevel, loglevelDelimiter)
	l.Println(messageLevel, msg)
	if b.String() != "" {
		t.Errorf("Log message %q of level %s should not have been printed via a logger at level %s", msg, level2str[messageLevel], level2str[loglevel])
	}
}

func TestLoglevels(t *testing.T) {
	const msg = "Test Message!"
	for loglevel := LevelDebug; loglevel >= LevelPanic; loglevel-- {
		for messageLevel := loglevel; messageLevel >= LevelPanic; messageLevel-- {
			doeslog(t, loglevel, messageLevel, msg)
		}
	}

	for loglevel := LevelPanic; loglevel <= LevelInfo; loglevel++ {
		for messageLevel := loglevel + 1; messageLevel <= LevelDebug; messageLevel++ {
			doesnotlog(t, loglevel, messageLevel, msg)
		}
	}
}

func TestPrintShortcuts(t *testing.T) {
	const msg = "This is a sample log record"
	const format = "%s%s%s\n"
	const errformat = "Expected %q. Got %q"
	var expect string
	b := new(strings.Builder)
	l := New(b, LevelDebug, loglevelDelimiter)

	//Panic
	l.Panic(msg)
	expect = fmt.Sprintf(format, level2str[LevelPanic], loglevelDelimiter, msg)
	if b.String() != expect {
		t.Errorf(errformat, expect, b.String())
	}
	b.Reset()

	//Alert
	l.Alert(msg)
	expect = fmt.Sprintf(format, level2str[LevelAlert], loglevelDelimiter, msg)
	if b.String() != expect {
		t.Errorf(errformat, expect, b.String())
	}
	b.Reset()

	//Critical
	l.Critical(msg)
	expect = fmt.Sprintf(format, level2str[LevelCritical], loglevelDelimiter, msg)
	if b.String() != expect {
		t.Errorf(errformat, expect, b.String())
	}
	b.Reset()

	//Error
	l.Error(msg)
	expect = fmt.Sprintf(format, level2str[LevelError], loglevelDelimiter, msg)
	if b.String() != expect {
		t.Errorf(errformat, expect, b.String())
	}
	b.Reset()

	//Warning
	l.Warning(msg)
	expect = fmt.Sprintf(format, level2str[LevelWarning], loglevelDelimiter, msg)
	if b.String() != expect {
		t.Errorf(errformat, expect, b.String())
	}
	b.Reset()

	//Notice
	l.Notice(msg)
	expect = fmt.Sprintf(format, level2str[LevelNotice], loglevelDelimiter, msg)
	if b.String() != expect {
		t.Errorf(errformat, expect, b.String())
	}
	b.Reset()

	//Info
	l.Info(msg)
	expect = fmt.Sprintf(format, level2str[LevelInfo], loglevelDelimiter, msg)
	if b.String() != expect {
		t.Errorf(errformat, expect, b.String())
	}
	b.Reset()

	//Debug
	l.Debug(msg)
	expect = fmt.Sprintf(format, level2str[LevelDebug], loglevelDelimiter, msg)
	if b.String() != expect {
		t.Errorf(errformat, expect, b.String())
	}
	b.Reset()
}

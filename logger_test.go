//This file is part of logger. ©2020 Jörg Walter.

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
	var expect string
	b := new(strings.Builder)
	l := New(b, LevelDebug, loglevelDelimiter)

	for logLevel := LevelPanic; logLevel <= LevelDebug; logLevel++ {
		switch logLevel {
		case LevelPanic:
			l.Panic(msg)
		case LevelAlert:
			l.Alert(msg)
		case LevelCritical:
			l.Critical(msg)
		case LevelError:
			l.Error(msg)
		case LevelWarning:
			l.Warning(msg)
		case LevelNotice:
			l.Notice(msg)
		case LevelInfo:
			l.Info(msg)
		case LevelDebug:
			l.Debug(msg)
		default:
			t.Fatalf("Unknown loglevel %d", logLevel)
		}
		expect = fmt.Sprintf("%s%s%s\n", level2str[logLevel], loglevelDelimiter, msg)
		if b.String() != expect {
			t.Errorf("Expected %q. Got %q", expect, b.String())
		}
		b.Reset()
	}
}

func TestPrintfShortcuts(t *testing.T) {
	const msg = "This is a sample log record for"
	var expect string
	b := new(strings.Builder)
	l := New(b, LevelDebug, loglevelDelimiter)

	for logLevel := LevelPanic; logLevel <= LevelDebug; logLevel++ {
		switch logLevel {
		case LevelPanic:
			l.Panicf("%s %s", msg, level2str[logLevel])
		case LevelAlert:
			l.Alertf("%s %s", msg, level2str[logLevel])
		case LevelCritical:
			l.Criticalf("%s %s", msg, level2str[logLevel])
		case LevelError:
			l.Errorf("%s %s", msg, level2str[logLevel])
		case LevelWarning:
			l.Warningf("%s %s", msg, level2str[logLevel])
		case LevelNotice:
			l.Noticef("%s %s", msg, level2str[logLevel])
		case LevelInfo:
			l.Infof("%s %s", msg, level2str[logLevel])
		case LevelDebug:
			l.Debugf("%s %s", msg, level2str[logLevel])
		default:
			t.Fatalf("Unknown loglevel %d", logLevel)
		}
		expect = fmt.Sprintf(
			"%s%s%s %s\n",
			level2str[logLevel],
			loglevelDelimiter,
			msg,
			level2str[logLevel])
		if b.String() != expect {
			t.Errorf("Expected %q. Got %q", expect, b.String())
		}
		b.Reset()
	}
}

func TestPrintfAutoLF(t *testing.T) {
	const msg = "5 + 4 = 9"
	expect := fmt.Sprintf("%s%s%s\n", level2str[LevelDebug], loglevelDelimiter, msg)
	b := new(strings.Builder)
	l := New(b, LevelDebug, loglevelDelimiter)
	//without LF
	l.Debugf("%d + %d = %d", 5, 4, 9)
	if b.String() != expect {
		t.Errorf("Expected %q. Got %q", expect, b.String())
	}
	b.Reset()
	//with LF
	l.Debugf("%d + %d = %d\n", 5, 4, 9)
	if b.String() != expect {
		t.Errorf("Expected %q. Got %q", expect, b.String())
	}
	b.Reset()
}

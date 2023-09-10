//This file is part of logger. ©2020-2023 Jörg Walter.

package logger

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"sync"
	"testing"
)

const loglevelDelimiter = " - "

func doeslog(t *testing.T, loglevel Level, messageLevel Level, msg string) {
	b := new(strings.Builder)
	l := New(b, loglevel, loglevelDelimiter)
	l.Println(messageLevel, msg)
	expected := fmt.Sprintf("[%s]%s%s\n", messageLevel.String(), loglevelDelimiter, msg)
	if b.String() != expected {
		t.Errorf("Logstring %q is not equal to expected string %q", b.String(), expected)
	}
}

func doesnotlog(t *testing.T, loglevel Level, messageLevel Level, msg string) {
	b := new(strings.Builder)
	l := New(b, loglevel, loglevelDelimiter)
	l.Println(messageLevel, msg)
	if b.String() != "" {
		t.Errorf("Log message %q of level %s should not have been printed via a logger at level %s", msg, messageLevel.String(), loglevel.String())
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
		expect = fmt.Sprintf("[%s]%s%s\n", logLevel.String(), loglevelDelimiter, msg)
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
			l.Panicf("%s %s", msg, logLevel.String())
		case LevelAlert:
			l.Alertf("%s %s", msg, logLevel.String())
		case LevelCritical:
			l.Criticalf("%s %s", msg, logLevel.String())
		case LevelError:
			l.Errorf("%s %s", msg, logLevel.String())
		case LevelWarning:
			l.Warningf("%s %s", msg, logLevel.String())
		case LevelNotice:
			l.Noticef("%s %s", msg, logLevel.String())
		case LevelInfo:
			l.Infof("%s %s", msg, logLevel.String())
		case LevelDebug:
			l.Debugf("%s %s", msg, logLevel.String())
		default:
			t.Fatalf("Unknown loglevel %d", logLevel)
		}
		expect = fmt.Sprintf(
			"[%s]%s%s %s\n",
			logLevel.String(),
			loglevelDelimiter,
			msg,
			logLevel.String())
		if b.String() != expect {
			t.Errorf("Expected %q. Got %q", expect, b.String())
		}
		b.Reset()
	}
}

func TestPrintfAutoLF(t *testing.T) {
	const msg = "5 + 4 = 9"
	expect := fmt.Sprintf("[%s]%s%s\n", LevelDebug.String(), loglevelDelimiter, msg)
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

// Checks if the mutex lock works correctly
func TestMutex(t *testing.T) {
	const messageLength = 36
	const goroutines = 16
	b := new(bytes.Buffer)
	l := New(b, LevelDebug, loglevelDelimiter)
	//Create output
	wg := new(sync.WaitGroup)
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go testMutexWorker(wg, l, i)
	}
	wg.Wait()
	//Check output
	scanner := bufio.NewScanner(b)
	var text string
	for line := 1; scanner.Scan(); line++ {
		text = scanner.Text()
		if len(text) != messageLength {
			t.Errorf("Malformed text at line %d: %s", line, text)
		}
	}
	if err := scanner.Err(); err != nil {
		t.Errorf("Scanner returned error: %s", err)
	}
}

func testMutexWorker(wg *sync.WaitGroup, l *Logger, id int) {
	for i := 0; i < 10000; i++ {
		l.Debugf("Goroutine %02d: Message %04d", id, i)
	}
	wg.Done()
}

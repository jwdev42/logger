//This file is part of logger.
//©2020 Jörg Walter, License: LGPL-3

package logger

import (
	"fmt"
	"strings"
	"testing"
)

const loglevel_delimiter = " - "

func doeslog(t *testing.T, logger_level int, message_level int, msg string) {
	b := new(strings.Builder)
	l := New(b, logger_level, loglevel_delimiter)
	l.Println(message_level, msg)
	expected := fmt.Sprintf("%s%s%s\n", level2str[message_level], loglevel_delimiter, msg)
	if b.String() != expected {
		t.Errorf("Logstring %q is not equal to expected string %q", b.String(), expected)
	}
}

func doesnotlog(t *testing.T, logger_level int, message_level int, msg string) {
	b := new(strings.Builder)
	l := New(b, logger_level, loglevel_delimiter)
	l.Println(message_level, msg)
	if b.String() != "" {
		t.Errorf("Log message %q of level %s should not have been printed via a logger at level %s", msg, level2str[message_level], level2str[logger_level])
	}
}

func TestLoglevels(t *testing.T) {
	const msg = "Test Message!"
	for logger_level := Level_Debug; logger_level >= Level_Panic; logger_level-- {
		for message_level := logger_level; message_level >= Level_Panic; message_level-- {
			doeslog(t, logger_level, message_level, msg)
		}
	}

	for logger_level := Level_Panic; logger_level <= Level_Info; logger_level++ {
		for message_level := logger_level + 1; message_level <= Level_Debug; message_level++ {
			doesnotlog(t, logger_level, message_level, msg)
		}
	}
}

func TestPrintShortcuts(t *testing.T) {
	const msg = "This is a sample log record"
	const format = "%s%s%s\n"
	const errformat = "Expected %q. Got %q"
	var expect string
	b := new(strings.Builder)
	l := New(b, Level_Debug, loglevel_delimiter)

	//Panic
	l.Panic(msg)
	expect = fmt.Sprintf(format, level2str[Level_Panic], loglevel_delimiter, msg)
	if b.String() != expect {
		t.Errorf(errformat, expect, b.String())
	}
	b.Reset()

	//Alert
	l.Alert(msg)
	expect = fmt.Sprintf(format, level2str[Level_Alert], loglevel_delimiter, msg)
	if b.String() != expect {
		t.Errorf(errformat, expect, b.String())
	}
	b.Reset()

	//Critical
	l.Critical(msg)
	expect = fmt.Sprintf(format, level2str[Level_Critical], loglevel_delimiter, msg)
	if b.String() != expect {
		t.Errorf(errformat, expect, b.String())
	}
	b.Reset()

	//Error
	l.Error(msg)
	expect = fmt.Sprintf(format, level2str[Level_Error], loglevel_delimiter, msg)
	if b.String() != expect {
		t.Errorf(errformat, expect, b.String())
	}
	b.Reset()

	//Warning
	l.Warning(msg)
	expect = fmt.Sprintf(format, level2str[Level_Warning], loglevel_delimiter, msg)
	if b.String() != expect {
		t.Errorf(errformat, expect, b.String())
	}
	b.Reset()

	//Notice
	l.Notice(msg)
	expect = fmt.Sprintf(format, level2str[Level_Notice], loglevel_delimiter, msg)
	if b.String() != expect {
		t.Errorf(errformat, expect, b.String())
	}
	b.Reset()

	//Info
	l.Info(msg)
	expect = fmt.Sprintf(format, level2str[Level_Info], loglevel_delimiter, msg)
	if b.String() != expect {
		t.Errorf(errformat, expect, b.String())
	}
	b.Reset()

	//Debug
	l.Debug(msg)
	expect = fmt.Sprintf(format, level2str[Level_Debug], loglevel_delimiter, msg)
	if b.String() != expect {
		t.Errorf(errformat, expect, b.String())
	}
	b.Reset()
}

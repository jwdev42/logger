//This file is part of logger.
//©2020 Jörg Walter, License: LGPL-3

package logger

import (
	"testing"
)

func TestStr2arg(t *testing.T) {
	teststr := "[Test]"
	expected := "test"
	arg, _ := str2arg(teststr)
	if arg != expected {
		t.Errorf("Expected result: %q. Actual result: %q", expected, arg)
	}
}

func TestLevelFlag_String(t *testing.T) {
	for k, v := range level2str {
		arg, err := str2arg(v)
		if err != nil {
			t.Log(err)
			t.FailNow()
		}
		flag := new(LevelFlag)
		if err := flag.Set(arg); err != nil {
			t.Error(err)
		}
		if int(*flag) != k {
			t.Errorf("Expected flag %d, got flag %d", k, int(*flag))
		}
	}
}

//This file is part of logger. ©2020-2023 Jörg Walter.

package logger

import (
	"strings"
	"testing"
)

func TestLevelFlagString(t *testing.T) {
	for k, v := range Loglevels() {
		arg := strings.ToLower(v)
		flag := new(LevelFlag)
		if err := flag.Set(arg); err != nil {
			t.Error(err)
		}
		if Level(*flag) != k {
			t.Errorf("Expected flag %d, got flag %d", k, Level(*flag))
		}
	}
}

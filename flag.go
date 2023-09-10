//This file is part of logger. ©2020-2023 Jörg Walter.

package logger

import (
	"fmt"
)

// LevelFlag implements flag.Getter
type LevelFlag Level

func (f *LevelFlag) Set(arg string) error {
	lvl, err := ParseLevel(arg)
	if err != nil {
		return fmt.Errorf("%q does not represent a valid loglevel", arg)
	}
	*f = LevelFlag(lvl)
	return nil
}

func (f *LevelFlag) String() string {
	if f == nil {
		return ""
	}
	return Level(*f).String()
}

func (f *LevelFlag) Get() any {
	if f == nil {
		return nil
	}
	return Level(*f)
}

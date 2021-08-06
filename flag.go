//This file is part of logger. ©2020 Jörg Walter.

package logger

import (
	"fmt"
	"strings"
)

type LevelFlag int

func (f *LevelFlag) Set(arg string) error {
	for k, v := range level2str {
		level, err := str2arg(v)
		if err != nil {
			return err
		}
		if arg == level {
			*f = LevelFlag(k)
			return nil
		}
	}
	return fmt.Errorf("%q is not a valid loglevel string representation", arg)
}

func (f *LevelFlag) String() string {
	if f == nil {
		return ""
	}
	if strrep, ok := level2str[int(*f)]; ok {
		s, _ := str2arg(strrep)
		return s
	}
	return ""
}

//Converts a string representation of a log level to a command line argument representation of itself.
func str2arg(input string) (string, error) {
	if input[0] != '[' || input[len(input)-1] != ']' || len(input) < 3 {
		return "", fmt.Errorf("%q is not a valid loglevel string representation", input)
	}
	return strings.ToLower(input[1 : len(input)-1]), nil
}

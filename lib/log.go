package lib

import (
	"flag"
	"fmt"
)

var flagVerbose int

func init() {
	flag.IntVar(&flagVerbose, "v", 0, "verbosity level (0=Info; 1=Debug; 2=Trace)")
}

const (
	LogInfo  = 0
	LogDebug = 1
	LogTrace = 2
)

func Log(v int, msg string) {
	if v <= flagVerbose {
		fmt.Println(msg)
	}
}

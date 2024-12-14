package lib

import (
	"flag"
	"fmt"
)

var LogLevel int

func init() {
	flag.IntVar(&LogLevel, "v", 0, "log level (0=Info; 1=Debug; 2=Trace)")
}

const (
	LogInfo  = 0
	LogDebug = 1
	LogTrace = 2
)

func Log(v int, a ...any) {
	if v <= LogLevel {
		fmt.Println(a...)
	}
}

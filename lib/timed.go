package lib

import (
	"fmt"
	"time"
)

func Timed(name string, f func(string), input string) {
	t0 := time.Now()
	f(input)
	elapsed := time.Since(t0)

	//fmt.Printf("%s (%q): %v\n\n", name, input, elapsed)
	fmt.Printf("%s: %v\n\n", name, elapsed)
}

func TimedFunc(name string, f func()) {
	t0 := time.Now()
	f()
	elapsed := time.Since(t0)

	//fmt.Printf("%s (%q): %v\n\n", name, input, elapsed)
	fmt.Printf("%s: %v\n\n", name, elapsed)
}

package lib

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/pprof"
	"time"

	"github.com/phicode/challenges/lib/assert"
	"github.com/phicode/challenges/lib/assets"
)

func Profile(repeat int, pprofFile, name string, f func(string), input string) {
	inputAbsolute := assets.MustFind(input)
	inputDir := filepath.Dir(inputAbsolute)
	pprofAbsolute := filepath.Join(inputDir, pprofFile)

	file, err := os.Create(pprofAbsolute)
	assert.NoErr(err)
	defer file.Close()

	assert.NoErr(pprof.StartCPUProfile(file))
	t0 := time.Now()
	for i := 0; i < repeat; i++ {
		f(input)
	}
	elapsed := time.Since(t0)
	pprof.StopCPUProfile()
	fmt.Printf("Profiled %s: %v -> %s\n\n", name, elapsed, pprofAbsolute)
	fmt.Printf("Run: go tool pprof -http localhost:8000 %s %s\n", os.Args[0], pprofAbsolute)
}

package assets

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/phicode/challenges/lib/assert"
)

func Find(name string) (string, error) {
	baseDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		testPath := filepath.Join(baseDir, name)
		s, err := os.Stat(testPath)
		if err != nil && !os.IsNotExist(err) {
			return "", err
		}
		if err == nil && s.Mode().IsRegular() {
			return testPath, nil
		}
		parentBaseDir := filepath.Dir(baseDir)
		if parentBaseDir == baseDir {
			return "", errors.New("file not found:" + name)
		}
		baseDir = parentBaseDir
	}

}

func MustFind(name string) string {
	p, err := Find(name)
	assert.NoErr(err)
	return p
}

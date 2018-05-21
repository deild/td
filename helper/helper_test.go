package helper

import (
	"os"
	"path"
	"testing"
)

func TestCheck(t *testing.T) {
	cwd, _ := os.Getwd()
	extra := "/TODOtestingFILE"
	_, err := os.Stat(path.Join(cwd, extra))
	if os.IsNotExist(err) {
		var w *os.File
		w, _ = os.Create(path.Join(cwd, extra))
		Check(w.Close)
	}
	os.Remove(path.Join(cwd, extra))
}

func TestCheckError(t *testing.T) {
	var w *os.File
	Check(w.Close)
}

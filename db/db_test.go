package db

import (
	"fmt"
	"os"
	"path"
	"testing"
	"time"
)

func TestWhenFileDoesntExist(t *testing.T) {
	cwd, _ := os.Getwd()
	extra := fmt.Sprint("/TODOtestingFOLDER/", time.Now().Format("20060102150405"))
	os.Setenv(envDBPath, path.Join(cwd, extra))
	db, _ := NewDataStore()
	if db.Check() == nil {
		t.Errorf("Expected database check to return error, but it didn't.")
	}
	os.Unsetenv(envDBPath)
}

func TestWhenDirExist(t *testing.T) {
	cwd, _ := os.Getwd()
	extra := fmt.Sprint("/TODOtestingFOLDER/", time.Now().Format("20060102150405"))
	os.Setenv(envDBPath, path.Join(cwd, extra))
	os.MkdirAll(path.Join(cwd, extra), 0700)
	db, _ := NewDataStore()
	if db.Check() != nil {
		t.Errorf("Expected database check to return nil, but it didn't.")
	}
	os.RemoveAll(path.Join(cwd, "/TODOtestingFOLDER/"))
	os.Unsetenv(envDBPath)
}

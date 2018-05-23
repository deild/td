// +build mage

package main

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var (
	binary    = "td"
	version   = "snapshot"
	buildDate = time.Now().UTC().Format(time.RFC3339)
	commit    = "none"
)

// Start by installing vgo,
func Init() error { // nolint: deadcode
	if err := sh.RunV("go", "get", "-u", "gopkg.in/alecthomas/gometalinter.v2"); err != nil {
		return err
	}
	if err := sh.Run("gometalinter.v2", "--install"); err != nil {
		return err
	}
	if err := sh.RunV("go", "get", "-u", "golang.org/x/vgo"); err != nil {
		return err
	}
	return sh.RunV("vgo", "vendor")
}

// A build step that requires additional params,
func Build() error { // nolint: deadcode
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}
	return sh.RunV("vgo", "build", "-ldflags", ldflags(), "-o", binary, "github.com/deild/td")
}

// Remove the temporarily generated files
func Clean() { // nolint: deadcode
	err := sh.Rm(binary)
	if err != nil {
		fmt.Println(err)
	}
	err = sh.Rm("dist")
	if err != nil {
		fmt.Println(err)
	}
	err = sh.Rm("vendor")
	if err != nil {
		fmt.Println(err)
	}
}

// Test, Lint and Build binary
func All() { // nolint: deadcode
	mg.Deps(Build, Lint, Test)
}

// Run Go Meta Linter
func Lint() error { // nolint: deadcode
	return sh.RunV("gometalinter.v2", "--vendor", "./...")
}

// Run tests
func Test() error { // nolint: deadcode
	return sh.RunV("vgo", "test", "-race", "-coverprofile=coverage.txt", "-covermode=atomic", "./...")
}

func ldflags() string {
	commit, err := sh.Output("git", "rev-parse", "--short", "HEAD")
	if err != nil {
		fmt.Print("WARNING: git rev-parse --short HEAD error:", err)
	}

	version, err := sh.Output("git", "describe", "--tags")
	if err != nil {
		fmt.Print("WARNING: git describe --tags error:", err)
	}

	return fmt.Sprintf("-s -w -X main.date=%s -X main.commit=%s -X main.version=%s", buildDate, commit, version)
}

// Generates a new release. Expects the TAG environment variable to be set,
// which will create a new tag with that name.
func Release() (err error) { // nolint: deadcode
	if os.Getenv("TAG") == "" {
		return errors.New("MSG and TAG environment variables are required")
	}
	if err := sh.RunV("git", "tag", "-a", "$TAG"); err != nil {
		return err
	}
	if err := sh.RunV("git", "push", "origin", "$TAG"); err != nil {
		return err
	}
	defer func() {
		if err != nil {
			err = sh.RunV("git", "tag", "--delete", "$TAG")
			err = sh.RunV("git", "push", "--delete", "origin", "$TAG")
		}
	}()
	return sh.RunV("goreleaser")
}

// Create a snapshot with goreleaser
func Snapshot() error { // nolint: deadcode
	return sh.RunV("goreleaser", "--rm-dist", "--snapshot")
}

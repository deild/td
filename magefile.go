// +build mage

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var (
	Default   = Build
	version   = "snapshot"
	buildDate = time.Now().UTC().Format(time.RFC3339)
	commit    = "none"
)

// Start by installing vgo
func Getvgo() error {
	fmt.Println("+ get vgo")
	return sh.Run("go", "get", "-u", "golang.org/x/vgo")
}

// A build step that requires additional params,
func Build() error {
	fmt.Println("+ build")
	return sh.Run("vgo", "build", "-ldflags", ldflags(), "-o", "td", ".")
}

// Clean up after yourself
func Clean() {
	fmt.Println("+ clean")
	err := os.RemoveAll("td")
	if err != nil {
		fmt.Println(err)
	}
	err = os.RemoveAll("build")
	if err != nil {
		fmt.Println(err)
	}
	err = os.RemoveAll("dist")
	if err != nil {
		fmt.Println(err)
	}
}

func getgox() error {
	return sh.Run("go", "get", "-u", "github.com/mitchellh/gox")
}

// Build binary for all os
func All() error {
	mg.Deps(getgox)
	fmt.Println("+ all")
	return sh.Run("gox", "-os=linux", "-os=windows", "-os=darwin", "-arch=amd64", "-output=build/{{.Dir}}_{{.OS}}_{{.Arch}}", "-ldflags", ldflags())
}

func getLint() error {
	return sh.Run("go", "get", "-u", "github.com/golang/lint/golint")
}

// Run Go Linter
func Lint() error {
	mg.Deps(getLint)
	fmt.Println("+ lint")
	out, err := sh.Output("golint", "./...")
	fmt.Println(out)
	return err
}

// Run tests
func Test() error {
	fmt.Println("+ test")
	return sh.Run("vgo", "test", "./...")
}

// Install myapp binary
func Install() error {
	fmt.Println("+ install")
	return sh.Run("vgo", "install", "-ldflags", ldflags(), ".")
}

func ldflags() string {
	commit, err := sh.Output("git", "rev-parse", "--short", "HEAD")
	if err != nil {
		fmt.Printf("WARNING: git rev-parse error")
	}
	hashtag, err := sh.Output("git", "rev-list", "--tags", "--max-count=1")
	if err != nil {
		fmt.Printf("WARNING: git rev-list error")
	}
	if hashtag != "" {
		tag, err := sh.Output("git", "describe", "--tags", hashtag)
		if err != nil {
			fmt.Printf("WARNING: git describe error")
		}
		if tag != "" {
			version = tag
		}
	}
	return fmt.Sprintf("-s -w -X main.date=%s -X main.commit=%s -X main.version=%s", buildDate, commit, version)
}

// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var Default = Build

// A build step that requires additional params,
func Build() error {
	fmt.Println("+ build")
	cmd := exec.Command("vgo", "build", "-o", "td", ".")
	return cmd.Run()
}

// Start by installing vgo
func Getvgo() error {
	fmt.Println("+ get vgo")
	cmd := exec.Command("go", "get", "-u", "golang.org/x/vgo")
	return cmd.Run()
}

// Clean up after yourself
func Clean() {
	fmt.Println("+ clean")
	os.RemoveAll("td")
}

func getgox() error {
	return sh.Run("go", "get", "-u", "github.com/mitchellh/gox")
}

// Build binary for all os
func All() error {
	mg.Deps(getgox)
	fmt.Println("+ all")
	cmd := exec.Command("gox", "-os=linux", "-os=windows", "-os=darwin", "-arch=amd64", "-output=./build/{{.Dir}}_{{.OS}}_{{.Arch}}", "-ldflags=\"-s -w -X main.version=1.3.0\"")
	return cmd.Run()
}

func getLint() error {
	return sh.Run("go", "get", "-u", "github.com/golang/lint/golint")
}

// Run Go Meta Linter
func Lint() error {
	mg.Deps(getLint)
	fmt.Println("+ lint")
	out, err := sh.Output("golint", "./...")
	fmt.Println(out)
	return err
}

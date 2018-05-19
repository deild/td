// +build mage

package main

import (
	"fmt"
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var Default = Build

// A build step that requires additional params,
func Build() error {
	fmt.Println("+ build")
	return sh.Run("vgo", "build", "-o", "td", ".")

}

// Start by installing vgo
func Getvgo() error {
	fmt.Println("+ get vgo")
	return sh.Run("go", "get", "-u", "golang.org/x/vgo")
}

// Clean up after yourself
func Clean() {
	fmt.Println("+ clean")
	err := os.RemoveAll("td")
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
	return sh.Run("gox", "-os=linux", "-os=windows", "-os=darwin", "-arch=amd64", "-output=build/{{.Dir}}_{{.OS}}_{{.Arch}}", "-ldflags=\"-s -w -X main.version=1.4.0\"")
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

// Run tests
func Test() error {
	fmt.Println("+ test")
	return sh.Run("vgo", "test", "-test.v", "./...")
}

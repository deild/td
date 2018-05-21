package helper

import "fmt"

// Check use to check return error of a deferred function
func Check(f func() error) {
	if err := f(); err != nil {
		fmt.Println("Received error:", err)
	}
}

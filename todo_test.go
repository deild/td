package main

import "testing"

func ExampleTodo() {
	todo := Todo{
		ID:       0,
		Desc:     "Test td",
		Status:   "pending",
		Modified: "",
	}
	todo.MakeOutput(false)
	// Output: 0 | ✕ Test td
}
func TestWithoutColor(t *testing.T) {
	todo := NewTodo()
	todo.ID = 0
	todo.Desc = "Test without color"
	todo.Modified = ""
	todo.MakeOutput(false)
	// Output: 0 | ✕ Test td
}

func TestWithColor(t *testing.T) {
	todo := NewTodo()
	todo.ID = 0
	todo.Desc = "Test color"
	todo.Modified = ""
	todo.MakeOutput(true)
	// Output: 0 | ✕ Test td
}

func TestStatusWIP(t *testing.T) {
	todo := NewTodo()
	todo.Status = WIP
	todo.MakeOutput(false)
	// Output: 0 | ✕ Test td
}

func TestStatusDone(t *testing.T) {
	todo := NewTodo()
	todo.Status = DONE
	todo.MakeOutput(false)
	// Output: 0 | ✕ Test td
}

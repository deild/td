package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/daviddengcn/go-colortext"
	p "github.com/deild/td/printer"
)

// Todo todo's structure
type Todo struct {
	ID       int64  `json:"id"`
	Desc     string `json:"desc"`
	Status   string `json:"status"`
	Modified string `json:"modified"`
}

// NewTodo create a pending todo
func NewTodo() *Todo {
	var todo = new(Todo)
	todo.Status = PENDING
	return todo
}

// MakeOutput print todo
func (t *Todo) MakeOutput(useColor bool) {
	var symbole string
	var color ct.Color

	switch t.Status {
	case DONE:
		color = ct.Green
		symbole = p.OkSign
	case WIP:
		color = ct.Blue
		symbole = p.WpSign
	default:
		color = ct.Red
		symbole = p.KoSign
	}

	hashtagReg := regexp.MustCompile(`#[^\\s]*`)

	spaceCount := 6 - len(strconv.FormatInt(t.ID, 10))

	fmt.Print(strings.Repeat(" ", spaceCount), t.ID, " | ")
	if useColor {
		ct.ChangeColor(color, false, ct.None, false)
	}
	fmt.Print(symbole)
	if useColor {
		ct.ResetColor()
	}
	fmt.Print(" ")
	pos := 0
	for _, token := range hashtagReg.FindAllStringIndex(t.Desc, -1) {
		fmt.Print(t.Desc[pos:token[0]])
		if useColor {
			ct.ChangeColor(ct.Yellow, false, ct.None, false)
		}
		fmt.Print(t.Desc[token[0]:token[1]])
		if useColor {
			ct.ResetColor()
		}
		pos = token[1]
	}
	fmt.Println(t.Desc[pos:])
}

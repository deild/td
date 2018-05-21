package main

import (
	"fmt"
	"strconv"

	"github.com/daviddengcn/go-colortext"
	"github.com/urfave/cli"
)

// Initialize a collection of todos
func initialize(c *cli.Context) error {
	db, err := NewDataStore()
	if err != nil {
		return exitError(err)
	}

	if err := db.Initialize(); err != nil {
		return exitError(err)
	}

	printSucces("Initialized empty to-do file as \"%s\".\n", db.Path)
	return nil
}

// Add a new todo
func add(c *cli.Context) error {

	if len(c.Args()) != 1 {
		return exitError(
			fmt.Errorf("You must provide a name to your todo.\nUsage: %s", c.Command.UsageText))
	}

	collection, err := NewCollection()
	if err != nil {
		return exitError(err)
	}

	todo := NewTodo()
	todo.Desc = c.Args()[0]

	id, err := collection.CreateTodo(todo)
	if err != nil {
		return exitError(err)
	}

	if err := collection.WriteTodos(); err != nil {
		return exitError(err)
	}

	printSucces("#%d \"%s\" is now added to your todos.\n", id, c.Args()[0])
	return nil
}

func modify(c *cli.Context) error {

	if len(c.Args()) != 2 {
		return exitError(
			fmt.Errorf("You must provide the id and the new text for your todo.\nUsage: %s", c.Command.UsageText))
	}

	collection, err := NewCollection()
	if err != nil {
		return exitError(err)
	}

	id, err := strconv.ParseInt(c.Args()[0], 10, 32)
	if err != nil {
		return exitError(err)
	}

	_, err = collection.Modify(id, c.Args()[1])
	if err != nil {
		return exitError(err)
	}

	if err := collection.WriteTodos(); err != nil {
		return exitError(err)
	}

	printSucces("\"%s\" has now a new description: %s\n", c.Args()[0], c.Args()[1])
	return nil
}

func toggle(c *cli.Context) error {

	if len(c.Args()) != 1 {
		return exitError(
			fmt.Errorf("You must provide the position of the item you want to change.\nUsage: %s", c.Command.UsageText))
	}

	collection, err := NewCollection()
	if err != nil {
		return exitError(err)
	}

	id, err := strconv.ParseInt(c.Args()[0], 10, 32)
	if err != nil {
		return exitError(err)
	}

	todo, err := collection.Toggle(id)
	if err != nil {
		return exitError(err)
	}

	if err := collection.WriteTodos(); err != nil {
		return exitError(err)
	}

	var status string

	switch todo.Status {
	case "wip":
		status = "marked as work in progress"
	default:
		status = "marked as " + todo.Status
	}

	printSucces("Your todo %d is now %s.\n", id, status)
	return nil
}

func search(c *cli.Context) error {
	if len(c.Args()) != 1 {
		return exitError(
			fmt.Errorf("You must provide a string search.\nUsage: %s", c.Command.UsageText))
	}

	collection, err := NewCollection()
	if err != nil {
		return exitError(err)
	}

	collection.Search(c.Args()[0])

	if len(collection.Todos) == 0 {
		ct.ChangeColor(ct.Cyan, false, ct.None, false)
		fmt.Printf("Sorry, there's no todos containing \"%s\".\n", c.Args()[0])
		ct.ResetColor()
		return nil
	}

	if len(collection.Todos) > 0 {
		fmt.Println()
		for _, todo := range collection.Todos {
			todo.MakeOutput(true)
		}
		fmt.Println()
	} else {
		printSucces("There's no todo to show.")
	}
	return nil
}

func reorder(c *cli.Context) error {
	collection, err := NewCollection()
	if err != nil {
		return exitError(err)
	}

	exact := c.Bool("exact")

	if exact {
		ids := make([]int64, len(c.Args()))
		for i, sid := range c.Args() {
			id, err := strconv.ParseInt(sid, 10, 32)
			if err != nil {
				return exitError(err)
			}
			ids[i] = id
		}
		if err := collection.ReorderByIDs(ids); err != nil {
			return exitError(err)
		}
	}

	if err := collection.Reorder(); err != nil {
		return exitError(err)
	}

	if err := collection.WriteTodos(); err != nil {
		return exitError(err)
	}

	printSucces("Your list is now reordered.")
	return nil

}

func swap(c *cli.Context) error {

	if len(c.Args()) != 2 {
		return exitError(
			fmt.Errorf("You must provide two position if you want to swap todos.\nUsage: %s", c.Command.UsageText))
	}

	collection, err := NewCollection()
	if err != nil {
		return exitError(err)
	}

	idA, err := strconv.ParseInt(c.Args()[0], 10, 32)
	if err != nil {
		return exitError(err)
	}

	idB, err := strconv.ParseInt(c.Args()[1], 10, 32)
	if err != nil {
		return exitError(err)
	}

	_, err = collection.Find(idA)
	if err != nil {
		return exitError(err)
	}

	_, err = collection.Find(idB)
	if err != nil {
		return exitError(err)
	}

	if err := collection.Swap(idA, idB); err != nil {
		return exitError(err)
	}

	if err := collection.Reorder(); err != nil {
		return exitError(err)
	}

	if err := collection.WriteTodos(); err != nil {
		return exitError(err)
	}

	printSucces("\"%s\" and \"%s\" has been swapped\n", c.Args()[0], c.Args()[1])

	return nil
}

func wip(c *cli.Context) error {

	if len(c.Args()) != 1 {
		return exitError(
			fmt.Errorf("You must provide the position of the item you want to change.\nUsage: %s", c.Command.UsageText))
	}

	collection, err := NewCollection()
	if err != nil {
		return exitError(err)
	}

	id, err := strconv.ParseInt(c.Args()[0], 10, 32)
	if err != nil {
		return exitError(err)
	}

	todo, err := collection.SetStatus(id, WIP)
	if err != nil {
		return exitError(err)
	}

	if err := collection.WriteTodos(); err != nil {
		return exitError(err)
	}

	var status string
	switch todo.Status {
	case WIP:
		status = "marked as work in progress"
	default:
		status = todo.Status
	}

	printSucces("Your todo %d is now %s.\n", id, status)
	return nil
}

func clean(c *cli.Context) error {

	collection, err := NewCollection()

	if err != nil {
		return exitError(err)
	}

	collection.RemoveFinishedTodos()

	if err := collection.WriteTodos(); err != nil {
		return exitError(err)
	}

	printSucces("Your list is now flushed of finished todos.")
	return nil
}

func noSubcommands(c *cli.Context) error {

	collection, err := NewCollection()
	if err != nil {
		return exitError(err)
	}

	if !c.IsSet("all") {
		switch {
		case c.IsSet("done"):
			collection.ListDoneTodos()
		case c.IsSet("wip"):
			collection.ListWorkInProgressTodos()
		default:
			collection.ListUndoneTodos()
		}
	}

	if len(collection.Todos) > 0 {
		fmt.Println()
		for _, todo := range collection.Todos {
			todo.MakeOutput(true)
		}
		fmt.Println()

	} else {
		ct.ChangeColor(ct.Yellow, false, ct.None, false)
		fmt.Println("There's no todo to show.")
		ct.ResetColor()
	}
	return nil
}

func printSucces(format string, a ...interface{}) {
	ct.ChangeColor(ct.Cyan, false, ct.None, false)
	fmt.Println(format, a)
	ct.ResetColor()
}

func exitError(message error) *cli.ExitError {
	ct.ChangeColor(ct.Red, false, ct.None, false)
	return cli.NewExitError(message, 1)
}

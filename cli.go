package main

import (
	"fmt"
	"os"
	"time"

	"github.com/daviddengcn/go-colortext"
	"github.com/urfave/cli"
)

const version = "dev"

var (
	// Hold subcommands
	cmds []cli.Command
	// Hold flags
	flags []cli.Flag
)

func init() {
	flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "done, d",
			Usage: "print done todos",
		},
		cli.BoolFlag{
			Name:  "wip, w",
			Usage: "print working in progress todos",
		},
		cli.BoolFlag{
			Name:  "all, a",
			Usage: "print all todos",
		},
	}
	cmds = []cli.Command{
		{
			Name:      "init",
			ShortName: "i",
			Usage:     "Initialize a collection of todos. If not path defined, it will create a file named .todos in the current directory.",
			UsageText: "td init",
			Action:    initialize,
		},
		{
			Name:      "add",
			ShortName: "a",
			Usage:     "Add a new todo",
			UsageText: "td add \"call mum\"",
			Action:    add,
		},
		{
			Name:      "modify",
			ShortName: "m",
			Usage:     "Modify the text of an existing todo",
			UsageText: "td modify 2 \"call dad\"",
			Action:    modify,
		},
		{
			Name:      "toggle",
			ShortName: "t",
			Usage:     "Toggle the status of a todo by giving his id",
			UsageText: "td toggle 1",
			Action:    toggle,
		},
		{
			Name:      "wip",
			ShortName: "w",
			Usage:     "Change the status of a todo to \"Work In Progress\" by giving its id",
			UsageText: "td wip 1",
			Action:    wip,
		},
		{
			Name:      "clean",
			ShortName: "c",
			Usage:     "Remove finished todos from the list",
			Action:    clean,
		},
		{
			Name:      "reorder",
			ShortName: "r",
			Usage:     "Reset ids of todo",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "exact, e",
					Usage: "Reorder the list according to the specified IDs and append the remaining items to the end of the list",
				},
			},
			Action: reorder,
		},
		{
			Name:      "swap",
			ShortName: "sw",
			Usage:     "Swap the position of two todos",
			UsageText: "td swap 9 3",
			Action:    swap,
		},
		{
			Name:      "search",
			ShortName: "s",
			Usage:     "Search a string in all todos",
			UsageText: "td search \"project-1\"",
			Action:    search,
		},
	}

}

func main() {
	app := cli.NewApp()
	app.Name = "td"
	app.Usage = "Your todos manager"
	app.Version = version
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Gaël Gillard",
			Email: "gael@gaelgillard.com",
		},
	}
	app.Flags = flags
	app.Commands = cmds
	app.Action = func(c *cli.Context) error {

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
			ct.ChangeColor(ct.Cyan, false, ct.None, false)
			fmt.Println("There's no todo to show.")
			ct.ResetColor()
		}
		return nil
	}
	app.After = func(c *cli.Context) error {
		db, _ := NewDataStore()
		ct.ChangeColor(ct.Magenta, false, ct.None, false)
		// fmt.Println(fmt.Sprintf("%-25s %s", `to-do file:`, db.Path))
		fmt.Println(db.Path)
		ct.ResetColor()
		return nil
	}

	app.Before = func(c *cli.Context) error {

		if len(c.Args()) == 1 {
			exceptions := []string{"init", "i", "help", "h"}
			for _, x := range exceptions {
				if c.Args()[0] == x {
					return nil
				}
			}
		}

		db, err := NewDataStore()
		if err != nil {
			return err
		}

		if err := db.Check(); err != nil {
			errDS := fmt.Errorf(`
===============================================================================

ERROR:

  File to store your todos could be found. Your current file location is:
  %s

  Run 'td init' to start a new to-do list or set/update the environment
  variable named 'TODO_DB_PATH' with the correct location of your to-dos file.

  Example 'export TODO_DB_PATH=$HOME/Dropbox/todo.json'

  If 'TODO_DB_PATH' is blank, it will reference to a file named '.todos' in the
  current working folder, and if there's no file, it will create one.

===============================================================================

				 `, db.Path)

			return cli.NewExitError(errDS, 1)
		}

		return nil
	}

	app.Run(os.Args)
}

func exitError(message error) *cli.ExitError {
	ct.ChangeColor(ct.Red, false, ct.None, false)
	return cli.NewExitError(message, 1)
}

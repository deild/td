package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/daviddengcn/go-colortext"
	"github.com/deild/td/db"
	"github.com/urfave/cli"
)

var (
	// Hold subcommands
	cmds []cli.Command
	// Hold flags
	flags []cli.Flag
	// Hold authors
	authors []cli.Author
	version = "dev"
	date    = "unknown"
	commit  = "none"
)

func init() {
	flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "done, d",
			Usage: "print done todos",
		},
		cli.BoolFlag{
			Name:  "wip, w",
			Usage: "print work in progress todos",
		},
		cli.BoolFlag{
			Name:  "all, a",
			Usage: "print all todos",
		},
	}

	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("%s %s\nBuild date: %s\nCommit: %s\n", c.App.Name, c.App.Version, date, commit)
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
			UsageText: "td clean",
			Action:    clean,
		},
		{
			Name:      "reorder",
			ShortName: "r",
			Usage:     "Reset ids of todo",
			UsageText: "td reset [-exact id2 id1]",
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
	authors = []cli.Author{
		cli.Author{
			Name:  "Tolvä",
			Email: "tolva@tuta.io",
		},
		cli.Author{
			Name:  "Gaël Gillard",
			Email: "gael@gaelgillard.com",
		},
		cli.Author{
			Name:  "Tarcísio Gruppi",
			Email: "txgruppi@gmail.com",
		},
		cli.Author{
			Name:  "Victor Alves",
			Email: "victor.alves@sentia.com",
		},
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "td"
	app.Usage = "Your todos manager"
	app.Version = version
	app.Compiled = time.Now().UTC()
	app.Authors = authors
	app.Copyright = "Copyright (c) 2018 Tolvä"
	app.Flags = flags
	app.Commands = cmds
	app.Action = noSubcommands
	app.After = func(c *cli.Context) error {
		data, _ := db.NewDataStore()
		ct.ChangeColor(ct.Magenta, false, ct.None, false)
		fmt.Println(data.Path)
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

		data, err := db.NewDataStore()
		if err != nil {
			return err
		}

		if err := data.Check(); err != nil {
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

				 `, data.Path)

			return cli.NewExitError(errDS, 1)
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

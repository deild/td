# td

[![GitHub Release](https://img.shields.io/github/release/deild/td.svg)](https://github.com/deild/td/releases/latest)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/deild/td)
[![Software License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![SemVer](https://img.shields.io/badge/semver-2.0.0-blue.svg)](https://semver.org/)
[![Powered By: GoReleaser](https://img.shields.io/badge/powered%20by-goreleaser-blue.svg)](https://github.com/goreleaser)

[![Travis](https://travis-ci.org/deild/td.svg?branch=master)](https://travis-ci.org/deild/td)
[![Coverage Status](https://coveralls.io/repos/github/deild/td/badge.svg?branch=master)](https://coveralls.io/github/deild/td?branch=master)
[![codecov](https://codecov.io/gh/deild/td/branch/master/graph/badge.svg)](https://codecov.io/gh/deild/td)
[![CodeFactor](https://www.codefactor.io/repository/github/deild/td/badge)](https://www.codefactor.io/repository/github/deild/td)
[![Go Report Card](https://goreportcard.com/badge/github.com/deild/td)](https://goreportcard.com/report/github.com/deild/td)

> Your todo list in your terminal.
>
> ![Screenshot](screenshot.png)

## Usage

### Installation

- From *binary*: go to the [release page](https://github.com/deild/td/releases)
- From *source*: `go get github.com/deild/td`

### Information

*td* will look at a `.todos` files to store your todos (like Git does: it will try recursively in each parent folder). This permit to have different list of todos per folder.

If it doesn't find a `.todos`, *td* use an environment variable to store your todos: `TODO_DB_PATH` where you define the path to the JSON file. If the file doesn't exist, the program will create it for you.

### CLI

```sh
NAME:
   td - Your todos manager

USAGE:
   td [global options] command [command options] [arguments...]

VERSION:
   1.5.0

AUTHORS:
   Tolvä <tolva@tuta.io>
   Gaël Gillard <gael@gaelgillard.com>
   Tarcísio Gruppi <txgruppi@gmail.com>
   Victor Alves <victor.alves@sentia.com>

COMMANDS:
     init, i     Initialize a collection of todos. If not path defined, it will create a file named .todos in the current directory.
     add, a      Add a new todo
     modify, m   Modify the text of an existing todo
     toggle, t   Toggle the status of a todo by giving his id
     wip, w      Change the status of a todo to "Work In Progress" by giving its id
     clean, c    Remove finished todos from the list
     reorder, r  Reset ids of todo
     swap, sw    Swap the position of two todos
     search, s   Search a string in all todos
     help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --done, -d     print done todos
   --wip, -w      print work in progress todos
   --all, -a      print all todos
   --help, -h     show help
   --version, -v  print the version

COPYRIGHT:
   Copyright (c) 2018 Tolvä
```

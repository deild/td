package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
)

const WIP = "wip"
const DONE = "done"
const PENDING = "pending"

//Collection of todo
type Collection struct {
	Todos []*Todo
}

func NewCollection() (*Collection, error) {
	var collection *Collection
	collection = new(Collection)

	if err := collection.RetrieveTodos(); err != nil {
		return nil, err
	}

	return collection, nil
}

// RemoveAtIndex remove one todo with its index
func (c *Collection) RemoveAtIndex(item int) {
	s := *c
	s.Todos = append(s.Todos[:item], s.Todos[item+1:]...)
	*c = s
}

// RetrieveTodos load todo from disk
func (c *Collection) RetrieveTodos() error {
	db, err := NewDataStore()
	if err != nil {
		return err
	}
	if err := db.Check(); err != nil {
		return err
	}

	file, err := os.OpenFile(db.Path, os.O_RDONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(&c.Todos)
	return err
}

// WriteTodos write the collection on disk
func (c *Collection) WriteTodos() error {
	db, err := NewDataStore()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(db.Path, os.O_RDWR|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := json.MarshalIndent(&c.Todos, "", "  ")
	if err != nil {
		return err
	}
	if _, err = file.Write(data); err != nil {
		return err
	}
	return file.Sync()
}

func (c *Collection) ListPendingTodos() error {
	for i := len(c.Todos) - 1; i >= 0; i-- {
		if c.Todos[i].Status != PENDING {
			c.RemoveAtIndex(i)
		}
	}
	return nil
}

func (c *Collection) ListUndoneTodos() error {
	for i := len(c.Todos) - 1; i >= 0; i-- {
		if c.Todos[i].Status == DONE {
			c.RemoveAtIndex(i)
		}
	}
	return nil
}

func (c *Collection) ListDoneTodos() error {
	for i := len(c.Todos) - 1; i >= 0; i-- {
		if c.Todos[i].Status != DONE {
			c.RemoveAtIndex(i)
		}
	}
	return nil
}

func (c *Collection) ListWorkInProgressTodos() {
	for i := len(c.Todos) - 1; i >= 0; i-- {
		if c.Todos[i].Status != WIP {
			c.RemoveAtIndex(i)
		}
	}
}

// CreateTodo new todo in the collection
func (c *Collection) CreateTodo(newTodo *Todo) (int64, error) {
	var err error
	var highestID int64 = 0
	for _, todo := range c.Todos {
		if todo.ID > highestID {
			highestID = todo.ID
		}
	}

	newTodo.ID = (highestID + 1)
	newTodo.Modified = time.Now().Local().String()
	c.Todos = append(c.Todos, newTodo)

	return newTodo.ID, err
}

// Find a todo for an id or an error
func (c *Collection) Find(id int64) (foundedTodo *Todo, err error) {
	founded := false
	for _, todo := range c.Todos {
		if id == todo.ID {
			foundedTodo = todo
			founded = true
		}
	}
	if !founded {
		err = errors.New("The todo with the id " + strconv.FormatInt(id, 10) + " was not found.")
	}
	return
}

func (c *Collection) SetStatus(id int64, status string) (*Todo, error) {

	todo, err := c.Find(id)

	if err != nil {
		return todo, err
	}

	todo.Status = status

	todo.Modified = time.Now().Local().String()

	return todo, err
}

func (c *Collection) Toggle(id int64) (*Todo, error) {
	var status string

	todo, err := c.Find(id)

	if err != nil {
		return todo, err
	}

	switch todo.Status {
	case PENDING:
		status = WIP
	case WIP:
		status = DONE
	default:
		status = PENDING
	}

	return c.SetStatus(id, status)
}

// Modify the desc of the todo find by id
func (c *Collection) Modify(id int64, desc string) (*Todo, error) {
	todo, err := c.Find(id)

	if err != nil {
		return todo, err
	}

	todo.Desc = desc
	todo.Modified = time.Now().Local().String()

	return todo, err
}

// RemoveFinishedTodos delete completed todos
func (c *Collection) RemoveFinishedTodos() error {
	return c.ListUndoneTodos()

}

// Reorder the collection
func (c *Collection) Reorder() error {
	for i, todo := range c.Todos {
		todo.ID = int64(i + 1)
	}
	return nil

}

// Swap inverse two elements of the collection
func (c *Collection) Swap(idA int64, idB int64) error {
	var positionA int
	var positionB int

	for i, todo := range c.Todos {
		switch todo.ID {
		case idA:
			positionA = i
			todo.ID = idB
		case idB:
			positionB = i
			todo.ID = idA
		}
	}

	c.Todos[positionA], c.Todos[positionB] = c.Todos[positionB], c.Todos[positionA]
	return nil

}

// Search retains only the elements that matches a sentence
func (c *Collection) Search(sentence string) {
	sentence = regexp.QuoteMeta(sentence)
	re := regexp.MustCompile("(?i)" + sentence)
	for i := len(c.Todos) - 1; i >= 0; i-- {
		if !re.MatchString(c.Todos[i].Desc) {
			c.RemoveAtIndex(i)
		}
	}
}

func (c *Collection) ReorderByIDs(ids []int64) error {
	idsMap := map[int64]int{}
	for index, id := range ids {
		if _, ok := idsMap[id]; ok {
			return fmt.Errorf("The ID %d is already in the list", id)
		}
		idsMap[id] = index
	}

	ordered := make([]*Todo, len(ids))
	rest := []*Todo{}

	for _, todo := range c.Todos {
		if index, ok := idsMap[todo.ID]; ok {
			ordered[index] = todo
			continue
		}
		rest = append(rest, todo)
	}

	newTodos := make([]*Todo, len(c.Todos))
	index := 0
	var idCounter int64 = 1
	for _, todo := range ordered {
		if todo == nil {
			continue
		}
		todo.ID = idCounter
		newTodos[index] = todo
		index++
		idCounter++
	}
	for _, todo := range rest {
		todo.ID = idCounter
		newTodos[index] = todo
		index++
		idCounter++
	}

	c.Todos = newTodos

	return nil
}

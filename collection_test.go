package main

import (
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	"github.com/deild/td/db"
)

var statuses = []string{PENDING, DONE, WIP}

func createTmpTodo(t *testing.T) {
	cwd, _ := os.Getwd()
	extra := fmt.Sprint("/TODOtestingFOLDER/", time.Now().Format("20060102150405"))
	os.Setenv(db.EnvDBPath, path.Join(path.Join(cwd, extra), ".todos"))
	os.MkdirAll(path.Join(cwd, extra), 0700)
	ds, err := db.NewDataStore()
	if err != nil {
		t.Error("error on NewCollection", err)
	}
	fmt.Println(os.Getenv(db.EnvDBPath))
	if err := ds.Initialize(); err != nil {
		t.Error("error on NewCollection", err)
	}
}

func removeTmpTodo() {
	cwd, _ := os.Getwd()
	os.RemoveAll(path.Join(cwd, "/TODOtestingFOLDER/"))
	os.Unsetenv(db.EnvDBPath)
}
func TestNewCollection(t *testing.T) {
	createTmpTodo(t)
	_, err := NewCollection()
	if err != nil {
		t.Error("error on NewCollection", err)
	}
	removeTmpTodo()
}

func TestWriteTodos(t *testing.T) {
	createTmpTodo(t)
	var collection Collection
	var todos []*Todo

	task := NewTodo()
	task.ID = 1
	task.Status = PENDING
	todos = append(todos, task)
	collection.Todos = todos
	err := collection.WriteTodos()
	if err != nil {
		t.Error("error can't WriteTodos", err)
	}
	removeTmpTodo()
}

func TestCreateTodo(t *testing.T) {
	createTmpTodo(t)
	collection, _ := NewCollection()
	task := NewTodo()
	task.ID = 1
	task.Status = PENDING
	id, err := collection.CreateTodo(task)
	if err != nil {
		t.Error("error can't WriteTodos", err)
	}
	if id != 1 {
		t.Errorf("expected 1, got %d.", id)
	}
	task = NewTodo()
	task.ID = 2
	task.Status = PENDING
	id, err = collection.CreateTodo(task)
	if err != nil {
		t.Error("error can't WriteTodos", err)
	}
	if id != 2 {
		t.Errorf("expected 2, got %d.", id)
	}
	removeTmpTodo()
}

func TestListPendingTodos(t *testing.T) {
	var collection Collection
	var todos []*Todo

	for id := range make([]int, len(statuses)) {
		task := NewTodo()
		task.ID = int64(id)
		task.Status = statuses[id]
		task.Desc = fmt.Sprintf("This is task number %d as %s", task.ID, task.Status)
		todos = append(todos, task)
	}
	collection.Todos = todos
	collection.ListPendingTodos()

	if len(collection.Todos) != 1 {
		t.Error("Expected only one task, got", len(collection.Todos))
	}

	for _, td := range collection.Todos {
		if td.Status != PENDING {
			t.Errorf("Expected status of tasks to be only \"pending\", got at least one as \"%s\"", td.Status)
		}
	}
}

func TestListUndoneTodos(t *testing.T) {
	var collection Collection
	var todos []*Todo

	for id := range make([]int, len(statuses)) {
		task := NewTodo()
		task.ID = int64(id)
		task.Status = statuses[id]
		todos = append(todos, task)
	}
	collection.Todos = todos
	collection.ListUndoneTodos()

	if len(collection.Todos) != 2 {
		t.Error("Expected only one task, got", len(collection.Todos))
	}

	for _, td := range collection.Todos {
		if td.Status != PENDING && td.Status != WIP {
			t.Errorf("Expected status of tasks to be only \"pending\" or \"work in progress\", got at least one as \"%s\"", td.Status)
		}
	}
}

func TestListWorkInProgressTodos(t *testing.T) {
	var collection Collection
	var todos []*Todo

	for id := range make([]int, len(statuses)) {
		task := NewTodo()
		task.ID = int64(id)
		task.Status = statuses[id]
		task.Desc = fmt.Sprintf("This is task number %d as %s", task.ID, task.Status)
		todos = append(todos, task)
	}
	collection.Todos = todos
	collection.ListWorkInProgressTodos()

	if len(collection.Todos) != 1 {
		t.Error("Expected only one task, got", len(collection.Todos))
	}

	for _, td := range collection.Todos {
		if td.Status != WIP {
			t.Errorf("Expected status of tasks to be only \"work in progress\", got at least one as \"%s\"", td.Status)
		}
	}
}

func TestListDoneTodos(t *testing.T) {
	var collection Collection
	var todos []*Todo

	for id := range make([]int, len(statuses)) {
		task := NewTodo()
		task.ID = int64(id)
		task.Status = statuses[id]
		todos = append(todos, task)
	}
	collection.Todos = todos
	collection.ListDoneTodos()

	if len(collection.Todos) != 1 {
		t.Error("Expected only one task, got", len(collection.Todos))
	}

	for _, td := range collection.Todos {
		if td.Status != DONE {
			t.Errorf("Expected status of tasks to be only \"done\", got at least one as \"%s\"", td.Status)
		}
	}
}

func TestToggleStatus(t *testing.T) {
	// PENDING > WIP > DONE > PENDING
	var collection Collection
	var todos []*Todo

	task := NewTodo()
	task.ID = 1
	task.Status = PENDING
	todos = append(todos, task)
	collection.Todos = todos

	collection.Toggle(1)
	if collection.Todos[0].Status != WIP {
		t.Errorf("Expected status to go from \"pending\" to \"work in progress\", changed to \"%s\" instead.", collection.Todos[0].Status)
	}

	collection.Toggle(1)
	if collection.Todos[0].Status != DONE {
		t.Errorf("Expected status to go from \"work in progress\" to \"done\", changed to \"%s\" instead.", collection.Todos[0].Status)
	}

	collection.Toggle(1)
	if collection.Todos[0].Status != PENDING {
		t.Errorf("Expected status to go from \"done\" to \"pending\", changed to \"%s\" instead.", collection.Todos[0].Status)
	}
}

func TestSetStatus(t *testing.T) {
	var collection Collection
	var todos []*Todo

	task := NewTodo()
	task.ID = 1
	task.Status = PENDING
	todos = append(todos, task)
	collection.Todos = todos

	collection.SetStatus(1, DONE)
	if collection.Todos[0].Status != DONE {
		t.Errorf("Expected to set status to \"done\", got \"%s\" instead.", collection.Todos[0].Status)
	}
}

func TestCantSetStatus(t *testing.T) {
	var collection Collection
	var todos []*Todo

	task := NewTodo()
	task.ID = 1
	task.Status = PENDING
	todos = append(todos, task)
	collection.Todos = todos

	_, err := collection.SetStatus(2, DONE)
	if err == nil {
		t.Errorf("Expected can't set status, got set status instead.")
	}
}

func TestTodoModifyDescription(t *testing.T) {
	var collection Collection
	var todos []*Todo

	oldDesc := []string{
		"Test 1",
		"Test 2",
		"Test 3",
		"Test 4",
	}

	newDesc := []string{
		"New test 1",
		"New test 2",
		"New test 3",
		"New test 4",
	}

	for id := range make([]int, len(oldDesc)) {
		task := NewTodo()
		task.ID = int64(id)
		task.Desc = oldDesc[id]
		todos = append(todos, task)
	}
	collection.Todos = todos

	for id, td := range collection.Todos {
		if td.Desc != oldDesc[id] {
			t.Error("Something is wrong with the test, description should be", oldDesc[id])
			t.FailNow()
		}
		collection.Modify(int64(id), newDesc[id])
		if td.Desc != newDesc[id] {
			t.Error("Something is wrong with the test, description should be", newDesc[id])
			t.FailNow()
		}
	}
}

func collectionFromTaskDesk(taskDesc []string) (Collection, []*Todo) {
	var collection Collection
	var todos []*Todo

	for id := range make([]int, len(taskDesc)) {
		task := NewTodo()
		task.ID = int64(id + 1)
		task.Desc = taskDesc[id]
		todos = append(todos, task)
	}
	collection.Todos = todos
	return collection, todos
}
func TestTodoSwap(t *testing.T) {

	taskDesc := []string{
		"Test 1",
		"Test 2",
		"Test 3",
		"Test 4",
	}

	collection, _ := collectionFromTaskDesk(taskDesc)
	collection.Swap(2, 4)

	secondTodo, err := collection.Find(2)
	if err != nil {
		t.Errorf("Expect to find item 2, but this happened: %s", err)
		t.FailNow()
	}
	if secondTodo.Desc != taskDesc[3] {
		t.Errorf("Expected the description from second todo item to have be \"%s\", but it didn't", taskDesc[3])
		t.FailNow()
	}

	lastTodo, err := collection.Find(4)
	if err != nil {
		t.Errorf("Expect to find item 4, but this happened: %s", err)
		t.FailNow()
	}
	if lastTodo.Desc != taskDesc[1] {
		t.Errorf("Expected the description from last todo item to have be \"%s\", but it didn't", taskDesc[1])
	}

}

func TestRemoveAtIndex(t *testing.T) {
	taskDesc := []string{
		"Test 1",
		"Test 2",
		"Test 3",
		"Test 4",
	}
	collection, todos := collectionFromTaskDesk(taskDesc)
	collection.RemoveAtIndex(2)
	if len(collection.Todos) != len(todos)-1 {
		t.Errorf("Expected size of current to-dos to one less then original slice, but it was %d and the other %d", len(collection.Todos), len(todos))
	}
}

func TestDontFind(t *testing.T) {
	var collection Collection
	var todos []*Todo

	task := NewTodo()
	task.ID = 1
	task.Status = PENDING
	todos = append(todos, task)
	collection.Todos = todos

	_, err := collection.Find(2)
	if err == nil {
		t.Error("Expect don't find item 2, but it find")
		t.FailNow()
	}

	_, err = collection.Find(1)
	if err != nil {
		t.Errorf("Expect find item 1, but this happened: %s", err)
		t.FailNow()
	}
}

func TestRemoveFinishedTodos(t *testing.T) {
	var collection Collection
	var todos []*Todo

	for id := range make([]int, len(statuses)) {
		task := NewTodo()
		task.ID = int64(id)
		task.Status = statuses[id]
		todos = append(todos, task)
	}
	collection.Todos = todos
	collection.RemoveFinishedTodos()

	if len(collection.Todos) != 2 {
		t.Error("Expected only one task, got", len(collection.Todos))
	}

	for _, td := range collection.Todos {
		if td.Status != PENDING && td.Status != WIP {
			t.Errorf("Expected status of tasks to be only \"pending\" or \"work in progress\", got at least one as \"%s\"", td.Status)
		}
	}
}

func TestSearchMatch(t *testing.T) {
	var collection Collection
	var todos []*Todo

	task := NewTodo()
	task.ID = 1
	task.Status = PENDING
	task.Desc = "Match"
	todos = append(todos, task)
	collection.Todos = todos

	collection.Search("Ma")
	if len(collection.Todos) != 1 {
		t.Error("Expected only one task, got", len(collection.Todos))
	}
}

func TestReorder(t *testing.T) {
	var collection Collection
	var todos []*Todo

	for id := range make([]int, len(statuses)) {
		task := NewTodo()
		task.ID = int64(id)
		task.Status = statuses[id]
		todos = append(todos, task)
	}
	collection.Todos = todos
	collection.Reorder()

	if len(collection.Todos) != 3 {
		t.Error("Expected 3 tasks, got", len(collection.Todos))
	}

	for l := range make([]int, len(statuses)) {
		if collection.Todos[l].ID != int64(l+1) {
			t.Errorf("Expected ID of tasks %d, got %d", l, collection.Todos[l].ID)
		}
	}

}

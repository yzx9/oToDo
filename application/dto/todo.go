package dto

import (
	"time"

	"github.com/devfeel/mapper"
	"github.com/yzx9/otodo/domain/todo"
)

func init() {
	mapper.Register(&NewTodo{})
	mapper.Register(&Todo{})
	mapper.Register(&todo.Todo{})
}

type NewTodo struct {
	Title      string     `json:"title"`
	Memo       string     `json:"momo"`
	Importance bool       `json:"importance"`
	Deadline   *time.Time `json:"deadline"`
	Notified   bool       `json:"notified"`
	NotifyAt   *time.Time `json:"notifyAt"`
	Done       bool       `json:"done"`
	DoneAt     *time.Time `json:"dontAt"`

	UserID           int64   `json:"userID"`
	TodoListID       int64   `json:"todoListID"`
	Files            []int64 `json:"files"`
	Steps            []int64 `json:"steps"`
	TodoRepeatPlanID int64   `json:"todoRepeatPlanID"` // TODO
}

func (t NewTodo) ToEntity() todo.Todo {
	var to todo.Todo
	if err := mapper.Mapper(t, to); err != nil {
		panic(err)
	}
	return to
}

type Todo struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Title      string     `json:"title"`
	Memo       string     `json:"momo"`
	Importance bool       `json:"importance"`
	Deadline   *time.Time `json:"deadline"`
	Notified   bool       `json:"notified"`
	NotifyAt   *time.Time `json:"notifyAt"`
	Done       bool       `json:"done"`
	DoneAt     *time.Time `json:"dontAt"`

	UserID           int64   `json:"userID"`
	TodoListID       int64   `json:"todoListID"`
	Files            []int64 `json:"files"`
	Steps            []int64 `json:"steps"`
	TodoRepeatPlanID int64   `json:"todoRepeatPlanID"` // TODO
}

func (t Todo) ToEntity() todo.Todo {
	var to todo.Todo
	if err := mapper.Mapper(t, to); err != nil {
		panic(err)
	}
	return to
}

func (Todo) FromEntity(t todo.Todo) Todo {
	var to Todo
	if err := mapper.Mapper(t, to); err != nil {
		panic(err)
	}
	return to
}

type NewTodoStep struct {
	Name string `json:"name"`
	Done bool   `json:"done"`

	TodoID int64 `json:"todoID"`
}

func (step NewTodoStep) AssembleTo(entity *todo.TodoStep) {
	entity.Name = step.Name
	entity.Done = step.Done
}

type TodoStep struct {
	ID int64 `json:"id"`

	Name   string     `json:"name"`
	Done   bool       `json:"done"`
	DoneAt *time.Time `json:"doneAt"`

	TodoID int64 `json:"todoID"`
}

func (s TodoStep) ToEntity() todo.TodoStep {
	return todo.TodoStep{
		ID: s.ID,

		Name:   s.Name,
		Done:   s.Done,
		DoneAt: s.DoneAt,

		TodoID: s.TodoID,
	}
}

func (TodoStep) FromEntity(entity todo.TodoStep) TodoStep {
	return TodoStep{
		ID: entity.ID,

		Name:   entity.Name,
		Done:   entity.Done,
		DoneAt: entity.DoneAt,

		TodoID: entity.TodoID,
	}
}

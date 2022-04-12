package dto

import (
	"time"

	"github.com/yzx9/otodo/domain/todo"
)

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
	return todo.Todo{
		Title:      t.Title,
		Memo:       t.Memo,
		Importance: t.Importance,
		Deadline:   t.Deadline,
		Notified:   t.Notified,
		NotifyAt:   t.NotifyAt,
		Done:       t.Done,
		DoneAt:     t.DoneAt,

		UserID:           t.UserID,
		TodoListID:       t.TodoListID,
		Files:            t.Files,
		Steps:            t.Steps,
		TodoRepeatPlanID: t.TodoRepeatPlanID,
	}
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
	return todo.Todo{
		ID:        t.ID,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,

		Title:      t.Title,
		Memo:       t.Memo,
		Importance: t.Importance,
		Deadline:   t.Deadline,
		Notified:   t.Notified,
		NotifyAt:   t.NotifyAt,
		Done:       t.Done,
		DoneAt:     t.DoneAt,

		UserID:           t.UserID,
		TodoListID:       t.TodoListID,
		Files:            t.Files,
		Steps:            t.Steps,
		TodoRepeatPlanID: t.TodoRepeatPlanID,
	}
}

func (Todo) FromEntity(t todo.Todo) Todo {
	return Todo{
		ID:        t.ID,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,

		Title:      t.Title,
		Memo:       t.Memo,
		Importance: t.Importance,
		Deadline:   t.Deadline,
		Notified:   t.Notified,
		NotifyAt:   t.NotifyAt,
		Done:       t.Done,
		DoneAt:     t.DoneAt,

		UserID:           t.UserID,
		TodoListID:       t.TodoListID,
		Files:            t.Files,
		Steps:            t.Steps,
		TodoRepeatPlanID: t.TodoRepeatPlanID,
	}
}

type NewTodoStep struct {
	Name string `json:"name"`
	Done bool   `json:"done"`

	TodoID int64 `json:"todoID"`
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

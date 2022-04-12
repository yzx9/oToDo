package dto

import (
	"time"

	"github.com/yzx9/otodo/domain/todolist"
)

type NewTodoList struct {
	Name      string `json:"name"`
	IsBasic   bool   `json:"isBasic"`
	IsSharing bool   `json:"isSharing"`

	UserID           int64 `json:"userID"`
	TodoListFolderID int64 `json:"todoListFolderID"`
}

func (dto NewTodoList) ToEntity() todolist.TodoList {
	return todolist.TodoList{
		Name:      dto.Name,
		IsBasic:   dto.IsBasic,
		IsSharing: dto.IsSharing,

		UserID:           dto.UserID,
		TodoListFolderID: dto.TodoListFolderID,
	}
}

type TodoList struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Name      string `json:"name"`
	IsBasic   bool   `json:"isBasic"`
	IsSharing bool   `json:"isSharing"`

	UserID           int64 `json:"userID"`
	TodoListFolderID int64 `json:"todoListFolderID"`
}

func (TodoList) FromEntity(entity todolist.TodoList) TodoList {
	return TodoList{
		ID:        entity.ID,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,

		Name:      entity.Name,
		IsBasic:   entity.IsBasic,
		IsSharing: entity.IsSharing,

		UserID:           entity.UserID,
		TodoListFolderID: entity.TodoListFolderID,
	}
}

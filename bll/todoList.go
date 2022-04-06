package bll

import (
	"fmt"

	"github.com/yzx9/otodo/infrastructure/repository"
	"github.com/yzx9/otodo/infrastructure/util"
)

func CreateTodoList(userID int64, todoList *repository.TodoList) error {
	todoList.IsBasic = false
	todoList.UserID = userID
	todoList.TodoListFolderID = 0
	if err := repository.TodoListRepo.Insert(todoList); err != nil {
		return fmt.Errorf("fails to create todo list: %w", err)
	}

	return nil
}

func GetTodoList(userID, todoListID int64) (repository.TodoList, error) {
	return OwnOrSharedTodoList(userID, todoListID)
}

func ForceGetTodoList(todoListID int64) (repository.TodoList, error) {
	list, err := repository.TodoListRepo.Find(todoListID)
	if err != nil {
		return repository.TodoList{}, fmt.Errorf("fails to get todo list: %w", err)
	}

	return list, nil
}

func GetTodoLists(userID int64) ([]repository.TodoList, error) {
	vec, err := repository.TodoListRepo.FindByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("fails to get user todo lists: %w", err)
	}

	shared, err := repository.TodoListRepo.FindSharedOnesByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("fails to get user shared todo lists: %w", err)
	}

	vec = append(vec, shared...)
	return vec, nil
}

func UpdateTodoList(userID int64, todoList *repository.TodoList) error {
	oldTodoList, err := OwnTodoList(userID, todoList.ID)
	if err != nil {
		return err
	}

	if oldTodoList.IsBasic {
		return util.NewErrorWithForbidden("unable to update basic todo list")
	}

	if err := repository.TodoListRepo.Save(todoList); err != nil {
		return fmt.Errorf("fails to update todo list: %w", err)
	}

	return nil
}

func DeleteTodoList(userID, todoListID int64) (repository.TodoList, error) {
	// only allow delete by owner, not shared users
	todoList, err := OwnTodoList(userID, todoListID)
	if err != nil {
		return repository.TodoList{}, err
	}

	// disable to delete basic todo list
	if todoList.IsBasic {
		return repository.TodoList{}, util.NewErrorWithPreconditionFailed("unable to delete basic todo list: %v", todoListID)
	}

	// cascade delete todos
	if _, err = repository.TodoRepo.DeleteAllByTodoList(todoListID); err != nil {
		return repository.TodoList{}, fmt.Errorf("fails to cascade delete todos: %w", err)
	}

	if err = repository.TodoListRepo.Delete(todoListID); err != nil {
		return repository.TodoList{}, fmt.Errorf("fails to delete todo list: %w", err)
	}

	return todoList, nil
}

// owner
func OwnTodoList(userID, todoListID int64) (repository.TodoList, error) {
	todoList, err := repository.TodoListRepo.Find(todoListID)
	if err != nil {
		return repository.TodoList{}, fmt.Errorf("fails to get todo list: %v", todoListID)
	}

	if todoList.UserID != userID {
		return repository.TodoList{}, util.NewErrorWithForbidden("unable to handle non-owned todo list: %v", todoListID)
	}

	return todoList, nil
}

package service

import (
	"fmt"

	"github.com/yzx9/otodo/domain/todolist"
	"github.com/yzx9/otodo/infrastructure/repository"
)

func GetActiveTodoListSharings(userID, todoListID int64) ([]repository.Sharing, error) {
	sharings, err := repository.SharingRepo.FindAllActiveOnes(userID, repository.SharingTypeTodoList)
	if err != nil {
		return nil, fmt.Errorf("fails to get sharing tokens: %w", err)
	}

	vec := make([]repository.Sharing, 0)
	for i := range sharings {
		if sharings[i].RelatedID == todoListID {
			vec = append(vec, sharings[i])
		}
	}

	return vec, nil
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

func GetTodoListFolder(userID, todoListFolderID int64) (repository.TodoListFolder, error) {
	return todolist.OwnTodoListFolder(userID, todoListFolderID)
}

func GetTodoListFolders(userID int64) ([]repository.TodoListFolder, error) {
	vec, err := repository.TodoListFolderRepo.FindAllByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("fails to get todo list folder: %w", err)
	}

	return vec, nil
}

func GetTodoListSharedUsers(userID, todoListID int64) ([]repository.User, error) {
	_, err := todolist.OwnOrSharedTodoList(userID, todoListID)
	if err != nil {
		return nil, err
	}

	users, err := repository.TodoListRepo.FindAllSharedUsers(todoListID)
	if err != nil {
		return nil, fmt.Errorf("fails to get todo list shared users: %w", err)
	}

	return users, nil
}

package service

import (
	"fmt"

	"github.com/yzx9/otodo/application/dto"
	"github.com/yzx9/otodo/domain/todolist"
	"github.com/yzx9/otodo/domain/user"
)

func GetActiveTodoListSharings(userID, todoListID int64) ([]todolist.Sharing, error) {
	sharings, err := SharingRepository.FindAllActive(userID, todolist.SharingTypeTodoList)
	if err != nil {
		return nil, fmt.Errorf("fails to get sharing tokens: %w", err)
	}

	vec := make([]todolist.Sharing, 0)
	for i := range sharings {
		if sharings[i].RelatedID == todoListID {
			vec = append(vec, sharings[i])
		}
	}

	return vec, nil
}

func GetTodoLists(userID int64) ([]todolist.TodoList, error) {
	vec, err := TodoListRepository.FindAllByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("fails to get user todo lists: %w", err)
	}

	shared, err := TodoListRepository.FindAllSharedByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("fails to get user shared todo lists: %w", err)
	}

	vec = append(vec, shared...)
	return vec, nil
}

func GetTodoListFolder(userID, todoListFolderID int64) (todolist.TodoListFolder, error) {
	return todolist.OwnTodoListFolder(userID, todoListFolderID)
}

func GetTodoListFolders(userID int64) ([]todolist.TodoListFolder, error) {
	vec, err := TodoListFolderRepository.FindAllByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("fails to get todo list folder: %w", err)
	}

	return vec, nil
}

func GetTodoListSharedUsers(userID, todoListID int64) ([]user.User, error) {
	if _, err := todolist.OwnOrSharedTodoList(userID, todoListID); err != nil {
		return nil, err
	}

	users, err := TodoListSharingRepository.FindAllSharedUsers(todoListID)
	if err != nil {
		return nil, fmt.Errorf("fails to get todo list shared users: %w", err)
	}

	return users, nil
}

func GetSharingInfo(token string) (dto.SharingToken, error) {
	sharing, err := todolist.GetSharing(token)
	if err != nil {
		return dto.SharingToken{}, err
	}

	return dto.SharingToken{
		Token:     sharing.Token,
		Type:      sharing.Type,
		CreatedAt: sharing.CreatedAt,
	}, nil
}

func GetSharingTodoListInfo(token string) (dto.SharingTodoList, error) {
	sharing, err := todolist.GetSharing(token)
	if err != nil {
		return dto.SharingTodoList{}, err
	}

	user, err := GetUser(sharing.UserID)
	if err != nil {
		return dto.SharingTodoList{}, err
	}

	list, err := TodoListRepository.Find(sharing.RelatedID)
	if err != nil {
		return dto.SharingTodoList{}, fmt.Errorf("fails to get todo list: %w", err)
	}

	return dto.SharingTodoList{
		UserNickname: user.Nickname,
		TodoListName: list.Name,
	}, nil
}

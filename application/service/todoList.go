package service

import (
	"fmt"

	"github.com/yzx9/otodo/application/dto"
	"github.com/yzx9/otodo/domain/identity"
	"github.com/yzx9/otodo/domain/sharing"
	"github.com/yzx9/otodo/domain/todo"
)

func CreateTodoList(userID int64, todoList dto.NewTodoList) (dto.TodoList, error) {
	entity := todoList.ToEntity()
	if err := todo.CreateTodoList(userID, &entity); err != nil {
		return dto.TodoList{}, err
	}

	return dto.TodoList{}.FromEntity(entity), nil
}

func DeleteTodoList(userID, id int64) (todo.TodoList, error) {
	return todo.DeleteTodoList(userID, id)
}

func CreateTodoListSharing(userID, todoListID int64) (sharing.Sharing, error) {
	return todo.CreateTodoListSharing(userID, todoListID)
}

func DeleteTodoListSharing(userID int64, token string) error {
	return todo.DeleteTodoListSharing(userID, token)
}

func CreateTodoListSharedUser(userID int64, token string) error {
	return todo.CreateTodoListSharedUser(userID, token)
}

func DeleteTodoListSharedUser(operatorID int64, userID int64, todoListID int64) error {
	return todo.DeleteTodoListSharedUser(operatorID, userID, todoListID)
}

func GetActiveTodoListSharings(userID, todoListID int64) ([]sharing.Sharing, error) {
	sharings, err := SharingRepository.FindAllActive(userID, sharing.SharingTypeTodoList)
	if err != nil {
		return nil, fmt.Errorf("fails to get sharing tokens: %w", err)
	}

	vec := make([]sharing.Sharing, 0)
	for i := range sharings {
		if sharings[i].RelatedID == todoListID {
			vec = append(vec, sharings[i])
		}
	}

	return vec, nil
}

func GetTodoList(userID, todoListID int64) (todo.TodoList, error) {
	return todo.GetTodoList(userID, todoListID)
}

func GetTodoLists(userID int64) ([]todo.TodoList, error) {
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

func GetTodoListFolder(userID, todoListFolderID int64) (todo.TodoListFolder, error) {
	return todo.OwnTodoListFolder(userID, todoListFolderID)
}

func GetTodoListFolders(userID int64) ([]todo.TodoListFolder, error) {
	vec, err := TodoListFolderRepository.FindAllByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("fails to get todo list folder: %w", err)
	}

	return vec, nil
}

func GetTodoListSharedUsers(userID, todoListID int64) ([]identity.User, error) {
	if _, err := todo.OwnOrSharedTodoList(userID, todoListID); err != nil {
		return nil, err
	}

	users, err := TodoListSharingRepository.FindAllSharedUsers(todoListID)
	if err != nil {
		return nil, fmt.Errorf("fails to get todo list shared users: %w", err)
	}

	return users, nil
}

func GetSharingInfo(token string) (dto.SharingToken, error) {
	sharing, err := sharing.GetSharing(token)
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
	sharing, err := sharing.GetSharing(token)
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

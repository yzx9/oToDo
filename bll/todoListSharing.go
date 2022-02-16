package bll

import (
	"fmt"

	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/entity"
	"github.com/yzx9/otodo/utils"
)

func CreateTodoListSharedUser(userID, todoListID string) error {
	exist, err := ExistTodoListSharing(userID, todoListID)
	if err != nil {
		return err
	}

	if !exist {
		err = dal.InsertTodoListSharedUser(userID, todoListID)
		if err != nil {
			return fmt.Errorf("fails to create todo list shared user: %w", err)
		}
	}

	return nil
}

func GetTodoListSharedUsers(todoListID string) ([]entity.User, error) {
	users, err := dal.SelectTodoListSharedUsers(todoListID)
	if err != nil {
		return nil, fmt.Errorf("fails to get todo list shared users: %w", err)
	}

	return users, nil
}

func DeleteTodoListSharedUser(userID, todoListID string) error {
	if err := dal.DeleteTodoListSharedUser(userID, todoListID); err != nil {
		return fmt.Errorf("fails to delete todo list shared users: %w", err)
	}

	return nil
}

func ExistTodoListSharing(userID, todoListID string) (bool, error) {
	exist, err := dal.ExistTodoListSharing(userID, todoListID)
	if err != nil {
		return false, fmt.Errorf("fails to valid sharing: %w", err)
	}

	return exist, nil
}

// owner or shared user
func OwnOrSharedTodoList(userID, todoListID string) (entity.TodoList, error) {
	todoList, err := dal.SelectTodoList(todoListID)
	if err != nil {
		return entity.TodoList{}, fmt.Errorf("fails to get todo list: %v", todoListID)
	}

	if todoList.UserID != userID {
		sharing, err := ExistTodoListSharing(userID, todoListID)
		if err != nil {
			return entity.TodoList{}, fmt.Errorf("fails to get todo list sharing: %v", todoListID)
		}

		if !sharing {
			return entity.TodoList{}, utils.NewErrorWithForbidden("unable to handle unauthorized todo list: %v", todoListID)
		}
	}

	return todoList, nil
}

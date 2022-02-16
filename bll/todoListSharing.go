package bll

import (
	"fmt"

	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/model/entity"
	"github.com/yzx9/otodo/util"
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

func GetTodoListSharedUsers(userID string, todoListID string) ([]entity.User, error) {
	_, err := OwnOrSharedTodoList(userID, todoListID)
	if err != nil {
		return nil, err
	}

	users, err := dal.SelectTodoListSharedUsers(todoListID)
	if err != nil {
		return nil, fmt.Errorf("fails to get todo list shared users: %w", err)
	}

	return users, nil
}

// Delete shared user from todo list,
// can be called by owner to delete anyone,
// or called by shared user to delete themselves
func DeleteTodoListSharedUser(operatorID, userID, todoListID string) error {
	todoList, err := OwnOrSharedTodoList(operatorID, todoListID)
	if err != nil {
		return err
	}

	if todoList.UserID != operatorID && userID != operatorID {
		return util.NewErrorWithForbidden("unable to delete shared user")
	}

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
			return entity.TodoList{}, util.NewErrorWithForbidden("unable to handle unauthorized todo list: %v", todoListID)
		}
	}

	return todoList, nil
}

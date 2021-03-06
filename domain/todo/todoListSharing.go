package todo

import (
	"fmt"

	"github.com/yzx9/otodo/domain/sharing"
	"github.com/yzx9/otodo/util"
)

func CreateTodoListSharing(userID, todoListID int64) (sharing.Sharing, error) {
	todoList, err := OwnTodoList(userID, todoListID)
	if err != nil {
		return sharing.Sharing{}, err
	}

	if todoList.IsBasic {
		return sharing.Sharing{}, fmt.Errorf("unable to share basic todo list: %v", todoListID)
	}

	record, err := sharing.CreateSharing(userID, todoListID, sharing.SharingTypeTodoList)
	if err != nil {
		return sharing.Sharing{}, fmt.Errorf("fails to create sharing token: %w", err)
	}

	return record, nil
}

func DeleteTodoListSharing(userID int64, token string) error {
	record, err := sharing.GetSharing(token)
	if err != nil {
		return err
	}

	if record.Type != sharing.SharingTypeTodoList {
		return util.NewErrorWithForbidden("invalid sharing token: %v")
	}

	if record.UserID != userID {
		return util.NewErrorWithForbidden("unable to delete non-own sharing token")
	}

	if err := record.Delete(); err != nil {
		return fmt.Errorf("fails to delete sharing: %w", err)
	}

	return nil
}

/**
 * oTodo List Shared Users
 */
func CreateTodoListSharedUser(userID int64, token string) error {
	sharing, err := sharing.GetSharing(token)
	if err != nil {
		return err
	}

	exist, err := ExistTodoListSharing(userID, sharing.RelatedID)
	if err != nil {
		return err
	}

	if exist {
		return nil
	}

	err = TodoListSharingRepository.SaveSharedUser(userID, sharing.RelatedID)
	if err != nil {
		return fmt.Errorf("fails to create todo list shared user: %w", err)
	}

	return nil
}

// Delete shared user from todo list,
// can be called by owner to delete anyone,
// or called by shared user to delete themselves
func DeleteTodoListSharedUser(operatorID, userID, todoListID int64) error {
	todoList, err := OwnOrSharedTodoList(operatorID, todoListID)
	if err != nil {
		return err
	}

	if todoList.UserID != operatorID && userID != operatorID {
		return util.NewErrorWithForbidden("unable to delete shared user")
	}

	if err := TodoListSharingRepository.DeleteSharedUser(userID, todoListID); err != nil {
		return fmt.Errorf("fails to delete todo list shared users: %w", err)
	}

	return nil
}

func ExistTodoListSharing(userID, todoListID int64) (bool, error) {
	exist, err := TodoListSharingRepository.ExistSharing(userID, todoListID)
	if err != nil {
		return false, fmt.Errorf("fails to valid sharing: %w", err)
	}

	return exist, nil
}

// owner or shared user
func OwnOrSharedTodoList(userID, todoListID int64) (TodoList, error) {
	todoList, err := TodoListRepository.Find(todoListID)
	if err != nil {
		return TodoList{}, fmt.Errorf("fails to get todo list: %v", todoListID)
	}

	if todoList.UserID != userID {
		sharing, err := ExistTodoListSharing(userID, todoListID)
		if err != nil {
			return TodoList{}, fmt.Errorf("fails to get todo list sharing: %v", todoListID)
		}

		if !sharing {
			return TodoList{}, util.NewErrorWithForbidden("unable to handle unauthorized todo list: %v", todoListID)
		}
	}

	return todoList, nil
}

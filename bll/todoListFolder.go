package bll

import (
	"fmt"

	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/model/dto"
	"github.com/yzx9/otodo/model/entity"
	"github.com/yzx9/otodo/util"
)

func CreateTodoListFolder(userID int64, todoListFolderName string) (entity.TodoListFolder, error) {
	folder := entity.TodoListFolder{
		Name:   todoListFolderName,
		UserID: userID,
	}
	if err := dal.InsertTodoListFolder(&folder); err != nil {
		return entity.TodoListFolder{}, err
	}

	return folder, nil
}

func GetTodoListFolder(userID, todoListFolderID int64) (entity.TodoListFolder, error) {
	return OwnTodoListFolder(userID, todoListFolderID)
}

func GetTodoListFolders(userID int64) ([]entity.TodoListFolder, error) {
	vec, err := dal.SelectTodoListFolders(userID)
	return vec, fmt.Errorf("fails to get user: %w", err)
}

func DeleteTodoListFolder(userID, todoListFolderID int64) (entity.TodoListFolder, error) {
	write := func(err error) (entity.TodoListFolder, error) {
		return entity.TodoListFolder{}, err
	}

	folder, err := OwnTodoListFolder(userID, todoListFolderID)
	if err != nil {
		return write(err)
	}

	// TODO[feat] Whether to cascade delete todo lists
	// Cascade delete todo lists
	if _, err = dal.DeleteTodoListsByFolder(todoListFolderID); err != nil {
		return write(fmt.Errorf("fails to cascade delete todo lists: %w", err))
	}

	if err = dal.DeleteTodoListFolder(todoListFolderID); err != nil {
		return write(fmt.Errorf("fails to delete todo list folder: %w", err))
	}

	return folder, nil
}

// Verify permission
func OwnTodoListFolder(userID, todoListFolderID int64) (entity.TodoListFolder, error) {
	todoListFolder, err := dal.SelectTodoListFolder(todoListFolderID)
	if err != nil {
		return entity.TodoListFolder{}, fmt.Errorf("fails to get todo list folder: %v", todoListFolderID)
	}

	if todoListFolder.UserID != userID {
		return entity.TodoListFolder{}, util.NewErrorWithForbidden("unable to handle non-owned todo list folder: %v", todoListFolderID)
	}

	return todoListFolder, nil
}

// Get Menu, folder+list tree
func GetTodoListMenu(userID int64) ([]dto.TodoListMenu, error) {
	folders, err := GetTodoListFolders(userID)
	if err != nil {
		return nil, fmt.Errorf("fails to get user menu: %w", err)
	}

	lists, err := GetTodoLists(userID)
	if err != nil {
		return nil, fmt.Errorf("fails to get user menu: %w", err)
	}

	// TODO[feat]: Sortable
	menu := make([]dto.TodoListMenu, len(folders))
	for i := range folders {
		menu = append(menu, dto.TodoListMenu{
			ID:       folders[i].ID,
			Name:     folders[i].Name,
			IsLeaf:   false,
			Children: make([]dto.TodoListMenu, 0),
		})
	}

	for i := range lists {
		item := dto.TodoListMenu{
			ID:     lists[i].ID,
			Name:   lists[i].Name,
			IsLeaf: true,
		}

		if lists[i].TodoListFolderID == 0 {
			menu = append(menu, item)
			continue
		}

		for j := range menu {
			if menu[j].ID == lists[i].TodoListFolderID {
				menu[i].Children = append(menu[i].Children, item)
			}
		}
		// TODO[bug]: need log if data inconsistency
	}

	return menu, nil
}

package bll

import (
	"fmt"

	"github.com/yzx9/otodo/infrastructure/repository"
	"github.com/yzx9/otodo/model/dto"
)

// Get Menu, folder+list tree
func GetTodoListMenu(userID int64) ([]dto.TodoListMenuItem, error) {
	folders, err := GetTodoListFolders(userID)
	if err != nil {
		return nil, fmt.Errorf("fails to get user menu: %w", err)
	}

	lists, err := repository.SelectTodoListsWithMenuFormat(userID)
	if err != nil {
		return nil, fmt.Errorf("fails to get user menu: %w", err)
	}

	// TODO[feat]: Sortable
	menu := make([]dto.TodoListMenuItem, 0)
	for i := range folders {
		menu = append(menu, dto.TodoListMenuItem{
			TodoListMenuItemRaw: dto.TodoListMenuItemRaw{
				ID:    folders[i].ID,
				Name:  folders[i].Name,
				Count: 0,
			},
			IsLeaf:   false,
			Children: make([]dto.TodoListMenuItem, 0),
		})
	}

	for i := range lists {
		item := dto.TodoListMenuItem{
			TodoListMenuItemRaw: lists[i],
			IsLeaf:              true,
		}

		if lists[i].TodoListFolderID == 0 {
			menu = append(menu, item)
			continue
		}

		for j := range menu {
			if menu[j].ID == lists[i].TodoListFolderID {
				menu[j].Count += lists[i].Count
				menu[j].Children = append(menu[j].Children, item)
			}
		}
		// TODO[bug]: need log if data inconsistency
	}

	return menu, nil
}

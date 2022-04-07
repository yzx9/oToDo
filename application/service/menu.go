package service

import (
	"fmt"

	"github.com/yzx9/otodo/application/dto"
	"github.com/yzx9/otodo/infrastructure/repository"
)

// Get Menu, folder+list tree
func GetMenu(userID int64) ([]dto.TodoListMenuItem, error) {
	folders, err := GetTodoListFolders(userID)
	if err != nil {
		return nil, fmt.Errorf("fails to get user menu: %w", err)
	}

	lists, err := repository.TodoListRepo.FindByUserWithMenuFormat(userID)
	if err != nil {
		return nil, fmt.Errorf("fails to get user menu: %w", err)
	}

	// TODO[feat]: Sortable
	menu := make([]dto.TodoListMenuItem, 0)
	for i := range folders {
		menu = append(menu, dto.TodoListMenuItem{
			TodoListMenuItem: repository.TodoListMenuItem{
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
			TodoListMenuItem: lists[i],
			IsLeaf:           true,
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

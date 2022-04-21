package todo

import (
	"fmt"
)

type MenuItem struct {
	ID               int64
	Name             string
	Count            int
	TodoListFolderID int64

	IsLeaf   bool
	Children []MenuItem
}

// Get Menu, folder+list tree
func GetMenu(userID int64) ([]MenuItem, error) {
	folders, err := TodoListFolderRepository.FindAllByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("fails to get user menu: %w", err)
	}

	lists, err := TodoListRepository.FindByUserOnMenuFormat(userID)
	if err != nil {
		return nil, fmt.Errorf("fails to get user menu: %w", err)
	}

	// TODO[feat]: Sortable
	menu := make([]MenuItem, 0)
	for i := range folders {
		menu = append(menu, MenuItem{
			ID:       folders[i].ID,
			Name:     folders[i].Name,
			Count:    0,
			IsLeaf:   false,
			Children: make([]MenuItem, 0),
		})
	}

	for i := range lists {
		item := MenuItem{
			ID:     lists[i].ID,
			Name:   lists[i].Name,
			Count:  lists[i].Count,
			IsLeaf: true,
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

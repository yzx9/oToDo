package service

import (
	"github.com/yzx9/otodo/application/dto"
	"github.com/yzx9/otodo/domain/todolist"
	"github.com/yzx9/otodo/util"
)

// Get Menu, folder+list tree
func GetMenu(userID int64) ([]dto.MenuItem, error) {
	menu, err := todolist.GetMenu(userID)
	if err != nil {
		return nil, err
	}

	items := make([]dto.MenuItem, 0, len(menu))

	var assembler func(a todolist.MenuItem) dto.MenuItem
	assembler = func(a todolist.MenuItem) dto.MenuItem {
		return dto.MenuItem{
			ID:    a.ID,
			Name:  a.Name,
			Count: a.Count,

			IsLeaf:   a.IsLeaf,
			Children: util.Map(assembler, a.Children),
		}
	}
	for _, a := range menu {
		items = append(items, assembler(a))
	}

	return items, nil
}

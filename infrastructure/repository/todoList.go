package repository

import (
	"github.com/yzx9/otodo/model/dto"
	"github.com/yzx9/otodo/model/entity"
	"github.com/yzx9/otodo/util"
)

func InsertTodoList(todoList *entity.TodoList) error {
	err := db.
		Create(todoList).
		Error

	return util.WrapGormErr(err, "todo list")
}

func SelectTodoList(id int64) (entity.TodoList, error) {
	var list entity.TodoList
	err := db.
		Where(&entity.TodoList{
			Entity: entity.Entity{
				ID: id,
			},
		}).
		First(&list).
		Error

	return list, util.WrapGormErr(err, "todo list")
}

func SelectTodoLists(userId int64) ([]entity.TodoList, error) {
	var lists []entity.TodoList
	err := db.
		Where(entity.TodoList{
			UserID: userId,
		}).
		Find(&lists).
		Error

	return lists, util.WrapGormErr(err, "todo list")
}

func SelectTodoListsWithMenuFormat(userID int64) ([]dto.TodoListMenuItemRaw, error) {
	var lists []dto.TodoListMenuItemRaw
	err := db.
		Model(entity.TodoList{}).
		Where(entity.TodoList{
			UserID: userID,
		}).
		Not(entity.TodoList{
			IsBasic: true, // Skip basic todo list
		}).
		Select("id", "name", "todo_list_folder_id", "(SELECT count(todos.id) FROM todos WHERE todos.todo_list_id = todo_lists.id) as count").
		Find(&lists).
		Error

	return lists, util.WrapGormErr(err, "todo list")
}

func SaveTodoList(todoList *entity.TodoList) error {
	re := db.Save(&todoList)
	return util.WrapGormErr(re.Error, "todo list")
}

func DeleteTodoList(id int64) error {
	err := db.
		Delete(&entity.Todo{
			Entity: entity.Entity{
				ID: id,
			},
		}).
		Error

	return util.WrapGormErr(err, "todo list")
}

func DeleteTodoListsByFolder(todoListFolderID int64) (int64, error) {
	re := db.
		Where(entity.TodoList{
			TodoListFolderID: todoListFolderID,
		}).
		Delete(entity.TodoList{})

	return re.RowsAffected, util.WrapGormErr(re.Error, "todo list")
}

func ExistTodoList(id int64) (bool, error) {
	var count int64
	err := db.
		Model(&entity.TodoList{}).
		Where(&entity.TodoList{
			Entity: entity.Entity{
				ID: id,
			},
		}).
		Count(&count).
		Error

	return count != 0, util.WrapGormErr(err, "todo list")
}

/**
 * Sharing
 */

func InsertTodoListSharedUser(userID, todoListID int64) error {
	err := db.
		Model(&entity.User{
			Entity: entity.Entity{
				ID: userID,
			},
		}).
		Association("SharedTodoLists").
		Append(&entity.TodoList{
			Entity: entity.Entity{
				ID: todoListID,
			},
		})

	return util.WrapGormErr(err, "todo list shared user")
}

func SelectSharedTodoLists(userID int64) ([]entity.TodoList, error) {
	var lists []entity.TodoList
	err := db.
		Model(&entity.User{
			Entity: entity.Entity{
				ID: userID,
			},
		}).
		Association("SharedTodoLists").
		Find(&lists)

	return lists, util.WrapGormErr(err, "user shared todo list")
}

func SelectTodoListSharedUsers(todoListID int64) ([]entity.User, error) {
	var users []entity.User
	err := db.
		Model(&entity.TodoList{
			Entity: entity.Entity{
				ID: todoListID,
			},
		}).
		Association("SharedUsers").
		Find(&users)

	return users, util.WrapGormErr(err, "todo list shared users")
}

func DeleteTodoListSharedUser(userID, todoListID int64) error {
	err := db.
		Model(&entity.TodoList{
			Entity: entity.Entity{
				ID: todoListID,
			},
		}).
		Association("SharedUsers").
		Delete(&entity.User{
			Entity: entity.Entity{
				ID: userID,
			},
		})

	return util.WrapGormErr(err, "todo list shared users")
}

func ExistTodoListSharing(userID, todoListID int64) (bool, error) {
	// TODO[pref]: count in db
	var lists []entity.TodoList
	err := db.
		Model(&entity.User{
			Entity: entity.Entity{
				ID: userID,
			},
		}).
		Association("SharedTodoLists").
		Find(&lists, &entity.TodoList{
			Entity: entity.Entity{
				ID: todoListID,
			},
		})

	if err != nil {
		return false, util.WrapGormErr(err, "todo list sharing")
	}

	return len(lists) != 0, nil
}

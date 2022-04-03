package repository

import (
	"github.com/yzx9/otodo/infrastructure/util"
	"github.com/yzx9/otodo/model/dto"
)

type TodoList struct {
	Entity

	Name      string `json:"name" gorm:"size:128"`
	IsBasic   bool   `json:"-"`
	IsSharing bool   `json:"isSharing"`

	UserID int64 `json:"userID"`
	User   User  `json:"-"`

	TodoListFolderID int64          `json:"todoListFolderID"`
	TodoListFolder   TodoListFolder `json:"-"`

	SharedUsers []*User `json:"-" gorm:"many2many:todo_list_shared_users"`
}

func InsertTodoList(todoList *TodoList) error {
	err := db.
		Create(todoList).
		Error

	return util.WrapGormErr(err, "todo list")
}

func SelectTodoList(id int64) (TodoList, error) {
	var list TodoList
	err := db.
		Where(&TodoList{
			Entity: Entity{
				ID: id,
			},
		}).
		First(&list).
		Error

	return list, util.WrapGormErr(err, "todo list")
}

func SelectTodoLists(userId int64) ([]TodoList, error) {
	var lists []TodoList
	err := db.
		Where(TodoList{
			UserID: userId,
		}).
		Find(&lists).
		Error

	return lists, util.WrapGormErr(err, "todo list")
}

func SelectTodoListsWithMenuFormat(userID int64) ([]dto.TodoListMenuItemRaw, error) {
	var lists []dto.TodoListMenuItemRaw
	err := db.
		Model(TodoList{}).
		Where(TodoList{
			UserID: userID,
		}).
		Not(TodoList{
			IsBasic: true, // Skip basic todo list
		}).
		Select("id", "name", "todo_list_folder_id", "(SELECT count(todos.id) FROM todos WHERE todos.todo_list_id = todo_lists.id) as count").
		Find(&lists).
		Error

	return lists, util.WrapGormErr(err, "todo list")
}

func SaveTodoList(todoList *TodoList) error {
	re := db.Save(&todoList)
	return util.WrapGormErr(re.Error, "todo list")
}

func DeleteTodoList(id int64) error {
	err := db.
		Delete(&Todo{
			Entity: Entity{
				ID: id,
			},
		}).
		Error

	return util.WrapGormErr(err, "todo list")
}

func DeleteTodoListsByFolder(todoListFolderID int64) (int64, error) {
	re := db.
		Where(TodoList{
			TodoListFolderID: todoListFolderID,
		}).
		Delete(TodoList{})

	return re.RowsAffected, util.WrapGormErr(re.Error, "todo list")
}

func ExistTodoList(id int64) (bool, error) {
	var count int64
	err := db.
		Model(&TodoList{}).
		Where(&TodoList{
			Entity: Entity{
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
		Model(&User{
			Entity: Entity{
				ID: userID,
			},
		}).
		Association("SharedTodoLists").
		Append(&TodoList{
			Entity: Entity{
				ID: todoListID,
			},
		})

	return util.WrapGormErr(err, "todo list shared user")
}

func SelectSharedTodoLists(userID int64) ([]TodoList, error) {
	var lists []TodoList
	err := db.
		Model(&User{
			Entity: Entity{
				ID: userID,
			},
		}).
		Association("SharedTodoLists").
		Find(&lists)

	return lists, util.WrapGormErr(err, "user shared todo list")
}

func SelectTodoListSharedUsers(todoListID int64) ([]User, error) {
	var users []User
	err := db.
		Model(&TodoList{
			Entity: Entity{
				ID: todoListID,
			},
		}).
		Association("SharedUsers").
		Find(&users)

	return users, util.WrapGormErr(err, "todo list shared users")
}

func DeleteTodoListSharedUser(userID, todoListID int64) error {
	err := db.
		Model(&TodoList{
			Entity: Entity{
				ID: todoListID,
			},
		}).
		Association("SharedUsers").
		Delete(&User{
			Entity: Entity{
				ID: userID,
			},
		})

	return util.WrapGormErr(err, "todo list shared users")
}

func ExistTodoListSharing(userID, todoListID int64) (bool, error) {
	// TODO[pref]: count in db
	var lists []TodoList
	err := db.
		Model(&User{
			Entity: Entity{
				ID: userID,
			},
		}).
		Association("SharedTodoLists").
		Find(&lists, &TodoList{
			Entity: Entity{
				ID: todoListID,
			},
		})

	if err != nil {
		return false, util.WrapGormErr(err, "todo list sharing")
	}

	return len(lists) != 0, nil
}

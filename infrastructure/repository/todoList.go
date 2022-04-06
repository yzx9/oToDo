package repository

import (
	"github.com/yzx9/otodo/infrastructure/util"
	"gorm.io/gorm"
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

type TodoListMenuItem struct {
	ID               int64  `json:"id"`
	Name             string `json:"name"`
	Count            int    `json:"count"`
	TodoListFolderID int64  `json:"-"`
}

var TodoListRepo TodoListRepository

type TodoListRepository struct {
	db *gorm.DB
}

func (r *TodoListRepository) Insert(todoList *TodoList) error {
	err := r.db.
		Create(todoList).
		Error

	return util.WrapGormErr(err, "todo list")
}

func (r *TodoListRepository) Find(id int64) (TodoList, error) {
	var list TodoList
	err := r.db.
		Where(&TodoList{
			Entity: Entity{
				ID: id,
			},
		}).
		First(&list).
		Error

	return list, util.WrapGormErr(err, "todo list")
}

func (r *TodoListRepository) FindByUser(userId int64) ([]TodoList, error) {
	var lists []TodoList
	err := r.db.
		Where(TodoList{
			UserID: userId,
		}).
		Find(&lists).
		Error

	return lists, util.WrapGormErr(err, "todo list")
}

func (r *TodoListRepository) FindByUserWithMenuFormat(userID int64) ([]TodoListMenuItem, error) {
	var lists []TodoListMenuItem
	err := r.db.
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

func (r *TodoListRepository) Save(todoList *TodoList) error {
	re := r.db.Save(&todoList)
	return util.WrapGormErr(re.Error, "todo list")
}

func (r *TodoListRepository) Delete(id int64) error {
	err := r.db.
		Delete(&Todo{
			Entity: Entity{
				ID: id,
			},
		}).
		Error

	return util.WrapGormErr(err, "todo list")
}

func (r *TodoListRepository) DeleteAllByFolder(todoListFolderID int64) (int64, error) {
	re := r.db.
		Where(TodoList{
			TodoListFolderID: todoListFolderID,
		}).
		Delete(TodoList{})

	return re.RowsAffected, util.WrapGormErr(re.Error, "todo list")
}

func (r *TodoListRepository) Exist(id int64) (bool, error) {
	var count int64
	err := r.db.
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

func (r *TodoListRepository) InsertSharedUser(userID, todoListID int64) error {
	err := r.db.
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

func (r *TodoListRepository) FindSharedOnesByUser(userID int64) ([]TodoList, error) {
	var lists []TodoList
	err := r.db.
		Model(&User{
			Entity: Entity{
				ID: userID,
			},
		}).
		Association("SharedTodoLists").
		Find(&lists)

	return lists, util.WrapGormErr(err, "user shared todo list")
}

func (r *TodoListRepository) SelectTodoListSharedUsers(todoListID int64) ([]User, error) {
	var users []User
	err := r.db.
		Model(&TodoList{
			Entity: Entity{
				ID: todoListID,
			},
		}).
		Association("SharedUsers").
		Find(&users)

	return users, util.WrapGormErr(err, "todo list shared users")
}

func (r *TodoListRepository) DeleteSharedUser(userID, todoListID int64) error {
	err := r.db.
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

func (r *TodoListRepository) ExistSharing(userID, todoListID int64) (bool, error) {
	// TODO[pref]: count in db
	var lists []TodoList
	err := r.db.
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

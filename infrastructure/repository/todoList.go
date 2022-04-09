package repository

import (
	"github.com/yzx9/otodo/domain/todolist"
	"github.com/yzx9/otodo/infrastructure/util"
	"gorm.io/gorm"
)

type TodoList struct {
	Entity

	Name      string `gorm:"size:128"`
	IsBasic   bool
	IsSharing bool

	UserID int64
	User   User

	TodoListFolderID int64
	TodoListFolder   TodoListFolder

	SharedUsers []*User `gorm:"many2many:todo_list_shared_users"`
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

func NewTodoListRepository(db *gorm.DB) TodoListRepository {
	return TodoListRepository{db: db}
}

func (r TodoListRepository) Save(entity *todolist.TodoList) error {
	po := r.convertToPO(entity)
	re := r.db.Save(&po)
	entity.ID = po.ID
	return util.WrapGormErr(re.Error, "todo list")
}

func (r TodoListRepository) Delete(id int64) error {
	err := r.db.
		Delete(&Todo{
			Entity: Entity{
				ID: id,
			},
		}).
		Error

	return util.WrapGormErr(err, "todo list")
}

func (r TodoListRepository) DeleteAllByFolder(todoListFolderID int64) (int64, error) {
	re := r.db.
		Where(TodoList{
			TodoListFolderID: todoListFolderID,
		}).
		Delete(TodoList{})

	return re.RowsAffected, util.WrapGormErr(re.Error, "todo list")
}

func (r TodoListRepository) Find(id int64) (todolist.TodoList, error) {
	var PO TodoList
	err := r.db.
		Where(&TodoList{
			Entity: Entity{
				ID: id,
			},
		}).
		First(&PO).
		Error

	return r.convertToEntity(PO), util.WrapGormErr(err, "todo list")
}

func (r TodoListRepository) FindByUser(userId int64) ([]todolist.TodoList, error) {
	var POs []TodoList
	err := r.db.
		Where(TodoList{
			UserID: userId,
		}).
		Find(&POs).
		Error

	if err != nil {
		return nil, util.WrapGormErr(err, "todo list")
	}

	entities := make([]todolist.TodoList, len(POs))
	for i := range POs {
		entities = append(entities, r.convertToEntity(POs[i]))
	}

	return entities, nil
}

func (r TodoListRepository) FindByUserWithMenuFormat(userID int64) ([]TodoListMenuItem, error) {
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

func (r TodoListRepository) Exist(id int64) (bool, error) {
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

func (r TodoListRepository) convertToPO(entity *todolist.TodoList) TodoList {
	return TodoList{
		Entity: Entity{
			ID:        entity.ID,
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt,
		},

		UserID: entity.UserID,

		TodoListFolderID: entity.TodoListFolderID,
		TodoListFolder:   TodoListFolder{}, // TODO

		SharedUsers: nil, // TODO
	}
}

func (r TodoListRepository) convertToEntity(po TodoList) todolist.TodoList {
	return todolist.TodoList{
		ID:        po.ID,
		CreatedAt: po.CreatedAt,
		UpdatedAt: po.UpdatedAt,

		UserID: po.UserID,

		TodoListFolderID: po.TodoListFolderID,

		SharedUsers: nil, // TODO
	}
}

/**
 * Sharing
 */

var TodoListSharingRepo TodoListSharingRepository

type TodoListSharingRepository struct {
	db *gorm.DB
}

func NewTodoListSharingRepository(db *gorm.DB) TodoListSharingRepository {
	return TodoListSharingRepository{db: db}
}

func (r TodoListSharingRepository) SaveSharedUser(userID, todoListID int64) error {
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

func (r TodoListSharingRepository) FindSharedOnesByUser(userID int64) ([]TodoList, error) {
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

func (r TodoListSharingRepository) FindAllSharedUsers(todoListID int64) ([]User, error) {
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

func (r TodoListSharingRepository) DeleteSharedUser(userID, todoListID int64) error {
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

func (r TodoListSharingRepository) ExistSharing(userID, todoListID int64) (bool, error) {
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

func (r TodoListSharingRepository) convertToPO(entity *todolist.TodoList) TodoList {
	return TodoList{
		Entity: Entity{
			ID:        entity.ID,
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt,
		},

		UserID: entity.UserID,

		TodoListFolderID: entity.TodoListFolderID,
		TodoListFolder:   TodoListFolder{}, // TODO

		SharedUsers: nil, // TODO
	}
}

func (r TodoListSharingRepository) convertToEntity(po TodoList) todolist.TodoList {
	return todolist.TodoList{
		ID:        po.ID,
		CreatedAt: po.CreatedAt,
		UpdatedAt: po.UpdatedAt,

		UserID: po.UserID,

		TodoListFolderID: po.TodoListFolderID,

		SharedUsers: nil, // TODO
	}
}

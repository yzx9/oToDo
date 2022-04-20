package repository

import (
	"github.com/yzx9/otodo/domain/todolist"
	"github.com/yzx9/otodo/domain/user"
	"github.com/yzx9/otodo/util"
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

func (r TodoListRepository) FindAllByUser(userID int64) ([]todolist.TodoList, error) {
	var POs []TodoList
	err := r.db.
		Where(TodoList{
			UserID: userID,
		}).
		Find(&POs).
		Error

	if err != nil {
		return nil, util.WrapGormErr(err, "todo list")
	}

	return r.convertToEntities(POs), nil
}

func (r TodoListRepository) FindAllSharedByUser(userID int64) ([]todolist.TodoList, error) {
	var POs []TodoList
	err := r.db.
		Model(&User{
			Entity: Entity{
				ID: userID,
			},
		}).
		Association("SharedTodoLists").
		Find(&POs)

	return r.convertToEntities(POs), util.WrapGormErr(err, "user shared todo list")
}

func (r TodoListRepository) FindByUserOnMenuFormat(userID int64) ([]todolist.MenuItem, error) {
	var items []todolist.MenuItem
	err := r.db.
		Model(TodoList{}).
		Where(TodoList{
			UserID: userID,
		}).
		Not(TodoList{
			IsBasic: true, // Skip basic todo list
		}).
		Select("id", "name", "todo_list_folder_id", "(SELECT count(todos.id) FROM todos WHERE todos.todo_list_id = todo_lists.id) as count").
		Find(&items).
		Error

	return items, util.WrapGormErr(err, "todo list")
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

func (r TodoListRepository) convertToEntities(POs []TodoList) []todolist.TodoList {
	return util.Map(r.convertToEntity, POs)
}

/**
 * Sharing
 */

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

func (r TodoListSharingRepository) FindAllSharedUsers(todoListID int64) ([]user.User, error) {
	var POs []User
	err := r.db.
		Model(&TodoList{
			Entity: Entity{
				ID: todoListID,
			},
		}).
		Association("SharedUsers").
		Find(&POs)

	// TODO: change to dep
	userRepo := NewUserRepository(r.db)
	return userRepo.convertToEntities(POs), util.WrapGormErr(err, "todo list shared users")
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

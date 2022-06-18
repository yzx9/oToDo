package repository

import (
	"github.com/yzx9/otodo/domain/identity"
	"github.com/yzx9/otodo/domain/todo"
	"github.com/yzx9/otodo/infrastructure/repository"
	"github.com/yzx9/otodo/util"
	"gorm.io/gorm"
)

type TodoList struct {
	repository.Entity

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

func (r TodoListRepository) Save(entity *todo.TodoList) error {
	po := r.convertToPO(entity)
	re := r.db.Save(&po)
	entity.ID = po.ID
	return util.WrapGormErr(re.Error, "todo list")
}

func (r TodoListRepository) Delete(id int64) error {
	err := r.db.
		Delete(&Todo{
			Entity: repository.Entity{
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

func (r TodoListRepository) Find(id int64) (todo.TodoList, error) {
	var PO TodoList
	err := r.db.
		Where(&TodoList{
			Entity: repository.Entity{
				ID: id,
			},
		}).
		First(&PO).
		Error

	return r.convertToEntity(PO), util.WrapGormErr(err, "todo list")
}

func (r TodoListRepository) FindAllByUser(userID int64) ([]todo.TodoList, error) {
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

func (r TodoListRepository) FindAllSharedByUser(userID int64) ([]todo.TodoList, error) {
	var POs []TodoList
	err := r.db.
		Model(&User{
			Entity: repository.Entity{
				ID: userID,
			},
		}).
		Association("SharedTodoLists").
		Find(&POs)

	return r.convertToEntities(POs), util.WrapGormErr(err, "user shared todo list")
}

func (r TodoListRepository) FindByUserOnMenuFormat(userID int64) ([]todo.MenuItem, error) {
	var items []todo.MenuItem
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
			Entity: repository.Entity{
				ID: id,
			},
		}).
		Count(&count).
		Error

	return count != 0, util.WrapGormErr(err, "todo list")
}

func (r TodoListRepository) convertToPO(entity *todo.TodoList) TodoList {
	return TodoList{
		Entity: repository.Entity{
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

func (r TodoListRepository) convertToEntity(po TodoList) todo.TodoList {
	return todo.TodoList{
		ID:        po.ID,
		CreatedAt: po.CreatedAt,
		UpdatedAt: po.UpdatedAt,

		UserID: po.UserID,

		TodoListFolderID: po.TodoListFolderID,

		SharedUsers: nil, // TODO
	}
}

func (r TodoListRepository) convertToEntities(POs []TodoList) []todo.TodoList {
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
			Entity: repository.Entity{
				ID: userID,
			},
		}).
		Association("SharedTodoLists").
		Append(&TodoList{
			Entity: repository.Entity{
				ID: todoListID,
			},
		})

	return util.WrapGormErr(err, "todo list shared user")
}

func (r TodoListSharingRepository) FindAllSharedUsers(todoListID int64) ([]identity.User, error) {
	var POs []User
	err := r.db.
		Model(&TodoList{
			Entity: repository.Entity{
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
			Entity: repository.Entity{
				ID: todoListID,
			},
		}).
		Association("SharedUsers").
		Delete(&User{
			Entity: repository.Entity{
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
			Entity: repository.Entity{
				ID: userID,
			},
		}).
		Association("SharedTodoLists").
		Find(&lists, &TodoList{
			Entity: repository.Entity{
				ID: todoListID,
			},
		})

	if err != nil {
		return false, util.WrapGormErr(err, "todo list sharing")
	}

	return len(lists) != 0, nil
}

func (r TodoListSharingRepository) convertToPO(entity *todo.TodoList) TodoList {
	return TodoList{
		Entity: repository.Entity{
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

func (r TodoListSharingRepository) convertToEntity(po TodoList) todo.TodoList {
	return todo.TodoList{
		ID:        po.ID,
		CreatedAt: po.CreatedAt,
		UpdatedAt: po.UpdatedAt,

		UserID: po.UserID,

		TodoListFolderID: po.TodoListFolderID,

		SharedUsers: nil, // TODO
	}
}

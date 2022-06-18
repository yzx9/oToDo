package repository

import (
	"time"

	"github.com/yzx9/otodo/domain/todo"
	"github.com/yzx9/otodo/infrastructure/repository"
	"github.com/yzx9/otodo/util"
	"gorm.io/gorm"
)

type Todo struct {
	repository.Entity

	Title      string `gorm:"size:128"`
	Memo       string
	Importance bool
	Deadline   *time.Time
	Notified   bool
	NotifyAt   *time.Time
	Done       bool
	DoneAt     *time.Time

	UserID int64
	User   User

	TodoListID int64
	TodoList   TodoList

	Files []File `gorm:"many2many:todo_files"`

	Steps []TodoStep

	TodoRepeatPlanID int64
	TodoRepeatPlan   TodoRepeatPlan

	NextID *int64 // next todo id if repeat
	Next   *Todo
}

var TodoRepo TodoRepository

type TodoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) TodoRepository {
	return TodoRepository{db: db}
}

func (r TodoRepository) Save(entity *todo.Todo) error {
	po := r.convertToPO(entity)
	err := r.db.Save(&po).Error
	entity.ID = po.ID
	return util.WrapGormErr(err, "todo")
}

func (r TodoRepository) Delete(id int64) error {
	err := r.db.
		Delete(&Todo{
			Entity: repository.Entity{
				ID: id,
			},
		}).
		Error

	return util.WrapGormErr(err, "todo")
}

func (r TodoRepository) DeleteAllByTodoList(todoListID int64) (int64, error) {
	re := r.db.
		Where(Todo{
			TodoListID: todoListID,
		}).
		Delete(Todo{})

	return re.RowsAffected, util.WrapGormErr(re.Error, "todo")
}

func (r TodoRepository) Find(id int64) (todo.Todo, error) {
	var po Todo
	err := r.db.
		Scopes(preloadTodoInfo).
		Where(&Todo{
			Entity: repository.Entity{
				ID: id,
			},
		}).
		First(&po).
		Error

	return r.convertToEntity(po), util.WrapGormErr(err, "todo")
}

func (r TodoRepository) FindAllByTodoList(todoListID int64) ([]todo.Todo, error) {
	var POs []Todo
	err := r.db.
		Scopes(preloadTodoInfo).
		Where(Todo{
			TodoListID: todoListID,
		}).
		Find(&POs).
		Error

	return r.convertToEntities(POs), util.WrapGormErr(err, "todos")
}

func (r TodoRepository) FindAllByUser(userID int64) ([]todo.Todo, error) {
	var POs []Todo
	err := r.db.
		Scopes(filterTodoUser(userID)).
		Find(&POs).
		Error

	return r.convertToEntities(POs), util.WrapGormErr(err, "all todos")
}

func (r TodoRepository) FindAllInBasicTodoList(userID int64) ([]todo.Todo, error) {
	var basicTodoList TodoList
	err := r.db.
		Model(&TodoList{}).
		Where(TodoList{UserID: userID, IsBasic: true}).
		Find(&basicTodoList).
		Error
	if err != nil {
		return nil, util.WrapGormErr(err, "basic todolist")
	}

	var POs []Todo
	err = r.db.
		Scopes(filterTodoUser(userID)).
		Find(&POs).
		Error

	return r.convertToEntities(POs), util.WrapGormErr(err, "basic todolist")
}

func (r TodoRepository) FindAllImportantOnesByUser(userID int64) ([]todo.Todo, error) {
	var POs []Todo
	err := r.db.
		Scopes(filterTodoUser(userID)).
		Where("Importance", true).
		Find(&POs).
		Error

	return r.convertToEntities(POs), util.WrapGormErr(err, "important todos")
}

func (r TodoRepository) FindAllPlanedOnesByUser(userID int64) ([]todo.Todo, error) {
	var POs []Todo
	err := r.db.
		Scopes(filterTodoUser(userID)).
		Not("deadline", nil).
		Order("deadline").
		Find(&POs).
		Error

	if err != nil {
		return nil, util.WrapGormErr(err, "not notified todos")
	}

	return r.convertToEntities(POs), util.WrapGormErr(err, "planed todos")
}

func (r TodoRepository) FindAllNotNotifiedOnesByUser(userID int64) ([]todo.Todo, error) {
	var POs []Todo
	err := r.db.
		Scopes(filterTodoUser(userID)).
		Not("notified", false).
		Order("notify_at").
		Find(&POs).
		Error

	if err != nil {
		return nil, util.WrapGormErr(err, "not notified todos")
	}

	return r.convertToEntities(POs), nil
}

func (r TodoRepository) convertToPO(entity *todo.Todo) Todo {
	return Todo{
		Entity: repository.Entity{
			ID:        entity.ID,
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt,
		},

		Title:      entity.Title,
		Memo:       entity.Memo,
		Importance: entity.Importance,
		Deadline:   entity.Deadline,
		Notified:   entity.Notified,
		NotifyAt:   entity.NotifyAt,
		Done:       entity.Done,
		DoneAt:     entity.DoneAt,

		UserID:           entity.UserID,
		TodoListID:       entity.TodoListID,
		Files:            nil, // TODO
		Steps:            nil, // TODO
		TodoRepeatPlanID: entity.TodoRepeatPlanID,
		NextID:           entity.NextID,
	}
}

func (r TodoRepository) convertToEntity(po Todo) todo.Todo {
	return todo.Todo{
		ID:        po.ID,
		CreatedAt: po.CreatedAt,
		UpdatedAt: po.UpdatedAt,

		Title:      po.Title,
		Memo:       po.Memo,
		Importance: po.Importance,
		Deadline:   po.Deadline,
		Notified:   po.Notified,
		NotifyAt:   po.NotifyAt,
		Done:       po.Done,
		DoneAt:     po.DoneAt,

		UserID:           po.UserID,
		TodoListID:       po.TodoListID,
		Files:            nil, // TODO
		Steps:            nil, // TODO
		TodoRepeatPlanID: po.TodoRepeatPlanID,
		NextID:           po.NextID,
	}
}

func (r TodoRepository) convertToEntities(POs []Todo) []todo.Todo {
	return util.Map(r.convertToEntity, POs)
}

/**
 * File
 */

type TodoFileRepository struct {
	db *gorm.DB
}

func NewTodoFileRepository(db *gorm.DB) TodoFileRepository {
	return TodoFileRepository{db: db}
}

func (r TodoFileRepository) Save(todoID, fileID int64) error {
	err := r.db.
		Where(Todo{
			Entity: repository.Entity{
				ID: todoID,
			},
		}).
		Association("Files").
		Append(&File{
			Entity: repository.Entity{
				ID: fileID,
			},
		})

	return util.WrapGormErr(err, "todo file")
}

func (r TodoFileRepository) FindAllByTodo(todoID int64) ([]File, error) {
	var files []File
	err := r.db.
		Where(Todo{
			Entity: repository.Entity{
				ID: todoID,
			},
		}).
		Find(&files).
		Error

	return files, util.WrapGormErr(err, "todo file")
}

/**
 * Helpers
 */

func filterTodoUser(userID int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.
			Where(Todo{UserID: userID}).
			Scopes(preloadTodoInfo)
	}
}

func preloadTodoInfo(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Files").
		Preload("Steps").
		Preload("TodoRepeatPlan")
}

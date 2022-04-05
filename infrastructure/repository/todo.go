package repository

import (
	"time"

	"github.com/yzx9/otodo/infrastructure/util"
	"gorm.io/gorm"
)

type Todo struct {
	Entity

	Title      string     `json:"title" gorm:"size:128"`
	Memo       string     `json:"memo"`
	Importance bool       `json:"importance"`
	Deadline   *time.Time `json:"deadline"`
	Notified   bool       `json:"notified"`
	NotifyAt   *time.Time `json:"notifyAt"`
	Done       bool       `json:"done"`
	DoneAt     *time.Time `json:"doneAt"`

	UserID int64 `json:"userID"`
	User   User  `json:"-"`

	TodoListID int64    `json:"todolistID"`
	TodoList   TodoList `json:"-"`

	Files []File `json:"files" gorm:"many2many:todo_files"`

	Steps []TodoStep `json:"steps"`

	TodoRepeatPlanID int64          `json:"-"`
	TodoRepeatPlan   TodoRepeatPlan `json:"todoRepeatPlan"`

	NextID *int64 `json:"nextID"` // next todo id if repeat
	Next   *Todo  `json:"-"`
}

var TodoRepo TodoRepository

type TodoRepository struct {
	db *gorm.DB
}

func (r *TodoRepository) InsertTodo(todo *Todo) error {
	err := r.db.Create(todo).Error
	return util.WrapGormErr(err, "todo")
}

func (r *TodoRepository) SelectTodo(id int64) (Todo, error) {
	var todo Todo
	err := r.db.
		Scopes(todoPreload).
		Where(&Todo{
			Entity: Entity{
				ID: id,
			},
		}).
		First(&todo).
		Error

	return todo, util.WrapGormErr(err, "todo")
}

func (r *TodoRepository) SelectTodos(todoListID int64) ([]Todo, error) {
	var todos []Todo
	err := r.db.
		Scopes(todoPreload).
		Where(Todo{
			TodoListID: todoListID,
		}).
		Find(&todos).
		Error

	return todos, util.WrapGormErr(err, "todos")
}

func (r *TodoRepository) SelectAllTodos(userID int64) ([]Todo, error) {
	var todos []Todo
	err := r.db.
		Scopes(todoUser(userID)).
		Find(&todos).
		Error

	return todos, util.WrapGormErr(err, "all todos")
}

func (r *TodoRepository) SelectImportantTodos(userID int64) ([]Todo, error) {
	var todos []Todo
	err := r.db.
		Scopes(todoUser(userID)).
		Where("Importance", true).
		Find(&todos).
		Error

	return todos, util.WrapGormErr(err, "important todos")
}

func (r *TodoRepository) SelectPlanedTodos(userID int64) ([]Todo, error) {
	var todos []Todo
	err := r.db.
		Scopes(todoUser(userID)).
		Not("deadline", nil).
		Order("deadline").
		Find(&todos).
		Error

	return todos, util.WrapGormErr(err, "planed todos")
}

func (r *TodoRepository) SelectNotNotifiedTodos(userID int64) ([]Todo, error) {
	var todos []Todo
	err := r.db.
		Scopes(todoUser(userID)).
		Not("notified", false).
		Order("notify_at").
		Find(&todos).
		Error

	return todos, util.WrapGormErr(err, "not notified todos")
}

func (r *TodoRepository) SaveTodo(todo *Todo) error {
	err := r.db.Save(&todo).Error
	return util.WrapGormErr(err, "todo")
}

func (r *TodoRepository) DeleteTodo(id int64) error {
	err := r.db.
		Delete(&Todo{
			Entity: Entity{
				ID: id,
			},
		}).
		Error

	return util.WrapGormErr(err, "todo")
}

func (r *TodoRepository) DeleteTodos(todoListID int64) (int64, error) {
	re := r.db.
		Where(Todo{
			TodoListID: todoListID,
		}).
		Delete(Todo{})

	return re.RowsAffected, util.WrapGormErr(re.Error, "todo")
}

/**
 * oTodo File
 */

func (r *TodoRepository) InsertTodoFile(todoID, fileID int64) error {
	err := r.db.
		Where(Todo{
			Entity: Entity{
				ID: todoID,
			},
		}).
		Association("Files").
		Append(&File{
			Entity: Entity{
				ID: fileID,
			},
		})

	return util.WrapGormErr(err, "todo file")
}

func (r *TodoRepository) SelectTodoFiles(todoID int64) ([]File, error) {
	var files []File
	err := r.db.
		Where(Todo{
			Entity: Entity{
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

func todoUser(userID int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.
			Where(Todo{UserID: userID}).
			Scopes(todoPreload)
	}
}

func todoPreload(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Files").
		Preload("Steps").
		Preload("TodoRepeatPlan")
}

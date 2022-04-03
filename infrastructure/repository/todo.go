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

func InsertTodo(todo *Todo) error {
	err := db.Create(todo).Error
	return util.WrapGormErr(err, "todo")
}

func SelectTodo(id int64) (Todo, error) {
	var todo Todo
	err := db.
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

func SelectTodos(todoListID int64) ([]Todo, error) {
	var todos []Todo
	err := db.
		Scopes(todoPreload).
		Where(Todo{
			TodoListID: todoListID,
		}).
		Find(&todos).
		Error

	return todos, util.WrapGormErr(err, "todos")
}

func SelectAllTodos(userID int64) ([]Todo, error) {
	var todos []Todo
	err := db.
		Scopes(todoUser(userID)).
		Find(&todos).
		Error

	return todos, util.WrapGormErr(err, "all todos")
}

func SelectImportantTodos(userID int64) ([]Todo, error) {
	var todos []Todo
	err := db.
		Scopes(todoUser(userID)).
		Where("Importance", true).
		Find(&todos).
		Error

	return todos, util.WrapGormErr(err, "important todos")
}

func SelectPlanedTodos(userID int64) ([]Todo, error) {
	var todos []Todo
	err := db.
		Scopes(todoUser(userID)).
		Not("deadline", nil).
		Order("deadline").
		Find(&todos).
		Error

	return todos, util.WrapGormErr(err, "planed todos")
}

func SelectNotNotifiedTodos(userID int64) ([]Todo, error) {
	var todos []Todo
	err := db.
		Scopes(todoUser(userID)).
		Not("notified", false).
		Order("notify_at").
		Find(&todos).
		Error

	return todos, util.WrapGormErr(err, "not notified todos")
}

func SaveTodo(todo *Todo) error {
	err := db.Save(&todo).Error
	return util.WrapGormErr(err, "todo")
}

func DeleteTodo(id int64) error {
	err := db.
		Delete(&Todo{
			Entity: Entity{
				ID: id,
			},
		}).
		Error

	return util.WrapGormErr(err, "todo")
}

func DeleteTodos(todoListID int64) (int64, error) {
	re := db.
		Where(Todo{
			TodoListID: todoListID,
		}).
		Delete(Todo{})

	return re.RowsAffected, util.WrapGormErr(re.Error, "todo")
}

/**
 * oTodo File
 */

func InsertTodoFile(todoID, fileID int64) error {
	err := db.
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

func SelectTodoFiles(todoID int64) ([]File, error) {
	var files []File
	err := db.
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

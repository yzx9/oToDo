package dal

import (
	"github.com/yzx9/otodo/model/entity"
	"github.com/yzx9/otodo/util"
	"gorm.io/gorm"
)

func InsertTodo(todo *entity.Todo) error {
	err := db.Create(todo).Error
	return util.WrapGormErr(err, "todo")
}

func SelectTodo(id int64) (entity.Todo, error) {
	var todo entity.Todo
	err := db.
		Scopes(todoPreload).
		Where(&entity.Todo{
			Entity: entity.Entity{
				ID: id,
			},
		}).
		First(&todo).
		Error

	return todo, util.WrapGormErr(err, "todo")
}

func SelectTodos(todoListID int64) ([]entity.Todo, error) {
	var todos []entity.Todo
	err := db.
		Scopes(todoPreload).
		Where(entity.Todo{
			TodoListID: todoListID,
		}).
		Find(&todos).
		Error

	return todos, util.WrapGormErr(err, "todos")
}

func SelectAllTodos(userID int64) ([]entity.Todo, error) {
	var todos []entity.Todo
	err := db.
		Scopes(todoUser(userID)).
		Find(&todos).
		Error

	return todos, util.WrapGormErr(err, "all todos")
}

func SelectImportantTodos(userID int64) ([]entity.Todo, error) {
	var todos []entity.Todo
	err := db.
		Scopes(todoUser(userID)).
		Where("Importance", true).
		Find(&todos).
		Error

	return todos, util.WrapGormErr(err, "important todos")
}

func SelectPlanedTodos(userID int64) ([]entity.Todo, error) {
	var todos []entity.Todo
	err := db.
		Scopes(todoUser(userID)).
		Not("deadline", nil).
		Order("deadline").
		Find(&todos).
		Error

	return todos, util.WrapGormErr(err, "planed todos")
}

func SelectNotNotifiedTodos(userID int64) ([]entity.Todo, error) {
	var todos []entity.Todo
	err := db.
		Scopes(todoUser(userID)).
		Not("notified", false).
		Order("notify_at").
		Find(&todos).
		Error

	return todos, util.WrapGormErr(err, "not notified todos")
}

func SaveTodo(todo *entity.Todo) error {
	err := db.Save(&todo).Error
	return util.WrapGormErr(err, "todo")
}

func DeleteTodo(id int64) error {
	err := db.
		Delete(&entity.Todo{
			Entity: entity.Entity{
				ID: id,
			},
		}).
		Error

	return util.WrapGormErr(err, "todo")
}

func DeleteTodos(todoListID int64) (int64, error) {
	re := db.
		Where(entity.Todo{
			TodoListID: todoListID,
		}).
		Delete(entity.Todo{})

	return re.RowsAffected, util.WrapGormErr(re.Error, "todo")
}

/**
 * oTodo File
 */

func InsertTodoFile(todoID, fileID int64) error {
	err := db.
		Where(entity.Todo{
			Entity: entity.Entity{
				ID: todoID,
			},
		}).
		Association("Files").
		Append(&entity.File{
			Entity: entity.Entity{
				ID: fileID,
			},
		})

	return util.WrapGormErr(err, "todo file")
}

func SelectTodoFiles(todoID int64) ([]entity.File, error) {
	var files []entity.File
	err := db.
		Where(entity.Todo{Entity: entity.Entity{ID: todoID}}).
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
			Where(entity.Todo{UserID: userID}).
			Scopes(todoPreload)
	}
}

func todoPreload(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Files").
		Preload("Steps").
		Preload("TodoRepeatPlan")
}

package dal

import (
	"github.com/yzx9/otodo/entity"
	"github.com/yzx9/otodo/utils"
)

func InsertTodo(todo *entity.Todo) error {
	re := db.Create(todo)
	return utils.WrapGormErr(re.Error, "todo")
}

func SelectTodo(id string) (entity.Todo, error) {
	var todo entity.Todo
	re := db.Where("ID = ?", id).First(&todo)
	return todo, utils.WrapGormErr(re.Error, "todo")
}

func SelectTodos(todoListID string) ([]entity.Todo, error) {
	var todos []entity.Todo
	re := db.Where("TodoListID = ?", todoListID).Find(&todos)
	return todos, utils.WrapGormErr(re.Error, "todos")
}

func SelectImportantTodos(userID string) ([]entity.Todo, error) {
	var todos []entity.Todo
	re := db.Where("UserID = ?", userID).Where("Importance", true).Find(&todos)
	return todos, utils.WrapGormErr(re.Error, "important todos")
}

func SelectPlanedTodos(userID string) ([]entity.Todo, error) {
	var todos []entity.Todo
	re := db.Where("UserID = ?", userID).Not("Deadline", nil).Order("Deadline").Find(&todos)
	return todos, utils.WrapGormErr(re.Error, "planed todos")
}

func SelectNotNotifiedTodos(userID string) ([]entity.Todo, error) {
	var todos []entity.Todo
	re := db.Where("UserID = ?", userID).Not("Notified", false).Order("Deadline").Find(&todos)
	return todos, utils.WrapGormErr(re.Error, "not notified todos")
}

func UpdateTodo(todo *entity.Todo) error {
	re := db.Save(&todo)
	return utils.WrapGormErr(re.Error, "todo")
}

func DeleteTodo(id string) error {
	re := db.Delete(&entity.Todo{
		Entity: entity.Entity{
			ID: id,
		},
	})
	return utils.WrapGormErr(re.Error, "todo")
}

func DeleteTodos(todoListID string) (int64, error) {
	re := db.Delete(entity.Todo{}, "TodoListID = ?", todoListID)
	return re.RowsAffected, utils.WrapGormErr(re.Error, "todo")
}

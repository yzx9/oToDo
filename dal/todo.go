package dal

import (
	"github.com/yzx9/otodo/model/entity"
	"github.com/yzx9/otodo/utils"
	"gorm.io/gorm"
)

func InsertTodo(todo *entity.Todo) error {
	re := db.Create(todo)
	return utils.WrapGormErr(re.Error, "todo")
}

func SelectTodo(id string) (entity.Todo, error) {
	var todo entity.Todo
	re := getTodosWithPreload().Where("id = ?", id).First(&todo)
	return todo, utils.WrapGormErr(re.Error, "todo")
}

func SelectTodos(todoListID string) ([]entity.Todo, error) {
	var todos []entity.Todo
	re := getTodosWithPreload().Where(entity.Todo{TodoListID: todoListID}).Find(&todos)
	return todos, utils.WrapGormErr(re.Error, "todos")
}

func SelectImportantTodos(userID string) ([]entity.Todo, error) {
	var todos []entity.Todo
	re := getTodosWithPreload().Where(entity.Todo{UserID: userID}).Where("Importance", true).Find(&todos)
	return todos, utils.WrapGormErr(re.Error, "important todos")
}

func SelectPlanedTodos(userID string) ([]entity.Todo, error) {
	var todos []entity.Todo
	re := getTodosWithPreload().Where(entity.Todo{UserID: userID}).Not("deadline", nil).Order("deadline").Find(&todos)
	return todos, utils.WrapGormErr(re.Error, "planed todos")
}

func SelectNotNotifiedTodos(userID string) ([]entity.Todo, error) {
	var todos []entity.Todo
	re := getTodosWithPreload().Where(entity.Todo{UserID: userID}).Not("notified", false).Order("notify_at").Find(&todos)
	return todos, utils.WrapGormErr(re.Error, "not notified todos")
}

func getTodosWithPreload() *gorm.DB {
	return db.Preload("Files").Preload("Steps").Preload("TodoRepeatPlan")
}

func SaveTodo(todo *entity.Todo) error {
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
	re := db.Where(entity.Todo{TodoListID: todoListID}).Delete(entity.Todo{})
	return re.RowsAffected, utils.WrapGormErr(re.Error, "todo")
}

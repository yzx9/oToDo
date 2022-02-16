package dal

import (
	"github.com/yzx9/otodo/model/entity"
	"github.com/yzx9/otodo/util"
	"gorm.io/gorm"
)

func InsertTodo(todo *entity.Todo) error {
	re := db.Create(todo)
	return util.WrapGormErr(re.Error, "todo")
}

func SelectTodo(id int64) (entity.Todo, error) {
	var todo entity.Todo
	re := getTodosWithPreload().Where("id = ?", id).First(&todo)
	return todo, util.WrapGormErr(re.Error, "todo")
}

func SelectTodos(todoListID int64) ([]entity.Todo, error) {
	var todos []entity.Todo
	re := getTodosWithPreload().Where(entity.Todo{TodoListID: todoListID}).Find(&todos)
	return todos, util.WrapGormErr(re.Error, "todos")
}

func SelectImportantTodos(userID int64) ([]entity.Todo, error) {
	var todos []entity.Todo
	re := getTodosWithPreload().Where(entity.Todo{UserID: userID}).Where("Importance", true).Find(&todos)
	return todos, util.WrapGormErr(re.Error, "important todos")
}

func SelectPlanedTodos(userID int64) ([]entity.Todo, error) {
	var todos []entity.Todo
	re := getTodosWithPreload().Where(entity.Todo{UserID: userID}).Not("deadline", nil).Order("deadline").Find(&todos)
	return todos, util.WrapGormErr(re.Error, "planed todos")
}

func SelectNotNotifiedTodos(userID int64) ([]entity.Todo, error) {
	var todos []entity.Todo
	re := getTodosWithPreload().Where(entity.Todo{UserID: userID}).Not("notified", false).Order("notify_at").Find(&todos)
	return todos, util.WrapGormErr(re.Error, "not notified todos")
}

func getTodosWithPreload() *gorm.DB {
	return db.Preload("Files").Preload("Steps").Preload("TodoRepeatPlan")
}

func SaveTodo(todo *entity.Todo) error {
	re := db.Save(&todo)
	return util.WrapGormErr(re.Error, "todo")
}

func DeleteTodo(id int64) error {
	re := db.Delete(&entity.Todo{
		Entity: entity.Entity{
			ID: id,
		},
	})
	return util.WrapGormErr(re.Error, "todo")
}

func DeleteTodos(todoListID int64) (int64, error) {
	re := db.Where(entity.Todo{TodoListID: todoListID}).Delete(entity.Todo{})
	return re.RowsAffected, util.WrapGormErr(re.Error, "todo")
}

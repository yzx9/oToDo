package dal

import (
	"github.com/yzx9/otodo/entity"
	"github.com/yzx9/otodo/utils"
)

func InsertTodoList(todoList *entity.TodoList) error {
	re := db.Create(todoList)
	return utils.WrapGormErr(re.Error, "todo list")
}

func SelectTodoList(id string) (entity.TodoList, error) {
	var list entity.TodoList
	re := db.Where("ID = ?", id).First(&list)
	return list, utils.WrapGormErr(re.Error, "todo list")
}

func SelectTodoLists(userId string) ([]entity.TodoList, error) {
	var lists []entity.TodoList
	re := db.Where("UserID = ?", userId).Find(&lists)
	return lists, utils.WrapGormErr(re.Error, "todo list")
}

func DeleteTodoList(id string) error {
	re := db.Delete(&entity.Todo{
		Entity: entity.Entity{
			ID: id,
		},
	})
	return utils.WrapGormErr(re.Error, "todo list")
}

func DeleteTodoListsByFolder(todoListFolderID string) (int64, error) {
	re := db.Delete(entity.TodoList{}, "TodoListFolderID = ?", todoListFolderID)
	return re.RowsAffected, utils.WrapGormErr(re.Error, "todo list")
}

func ExistTodoList(id string) (bool, error) {
	var count int64
	re := db.Model(&entity.TodoList{}).Where("ID = ?", id).Count(&count)
	return count != 0, utils.WrapGormErr(re.Error, "todo list")
}

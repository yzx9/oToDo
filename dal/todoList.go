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
	re := db.Where("id = ?", id).First(&list)
	return list, utils.WrapGormErr(re.Error, "todo list")
}

func SelectTodoLists(userId string) ([]entity.TodoList, error) {
	var lists []entity.TodoList
	re := db.Where(entity.TodoList{UserID: userId}).Find(&lists)
	return lists, utils.WrapGormErr(re.Error, "todo list")
}

func SelectSharedTodoLists(userId string) ([]entity.TodoList, error) {
	var lists []entity.TodoList
	err := db.Model(&entity.User{}).Where(entity.TodoList{UserID: userId}).Association("SharedTodoLists").Find(&lists)
	return lists, utils.WrapGormErr(err, "todo list")
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
	re := db.Where(entity.TodoList{TodoListFolderID: todoListFolderID}).Delete(entity.TodoList{})
	return re.RowsAffected, utils.WrapGormErr(re.Error, "todo list")
}

func ExistTodoList(id string) (bool, error) {
	var count int64
	re := db.Model(&entity.TodoList{}).Where("id = ?", id).Count(&count)
	return count != 0, utils.WrapGormErr(re.Error, "todo list")
}

func ExistTodoListSharing(userID, todoListID string) (bool, error) {
	user := entity.User{}
	user.ID = userID
	var lists []entity.TodoList
	if err := db.Model(&user).Association("SharedTodoLists").Find(&lists, "id = ?", todoListID); err != nil {
		return false, utils.WrapGormErr(err, "todo list")
	}

	return len(lists) != 0, nil
}

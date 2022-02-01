package bll

import (
	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/entity"
)

func GetTodos() []entity.Todo {
	return dal.GetTodos()
}

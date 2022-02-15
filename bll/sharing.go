package bll

import "github.com/yzx9/otodo/dal"

func HasSharing(userID, todoListID string) (bool, error) {
	return dal.ExistTodoListSharing(userID, todoListID)
}

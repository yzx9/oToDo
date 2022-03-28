package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/api/common"
	"github.com/yzx9/otodo/bll"
	"github.com/yzx9/otodo/model/dto"
	"github.com/yzx9/otodo/model/entity"
)

// Create todo list
func PostTodoListHandler(c *gin.Context) {
	list := entity.TodoList{}
	if err := c.ShouldBind(&list); err != nil {
		common.AbortWithError(c, err)
		return
	}

	userID := common.MustGetAccessUserID(c)
	if err := bll.CreateTodoList(userID, &list); err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, list)
}

// Get todo list
func GetTodoListHandler(c *gin.Context) {
	todoListID, err := common.GetRequiredParamID(c, "id")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	userID := common.MustGetAccessUserID(c)
	todoList, err := bll.GetTodoList(userID, todoListID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, todoList)
}

// Get todos in todo list
func GetTodoListTodosHandler(c *gin.Context) {
	todoListID, err := common.GetRequiredParamID(c, "id")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	userID := common.MustGetAccessUserID(c)
	todos, err := bll.GetTodos(userID, todoListID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, todos)
}

// Delete todo list
func DeleteTodoListHandler(c *gin.Context) {
	id, err := common.GetRequiredParamID(c, "id")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	userID := common.MustGetAccessUserID(c)
	todo, err := bll.DeleteTodoList(userID, id)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, todo)
}

/**
 * oTodo List Sharing
 */

// Get shared users in todo list
func GetTodoListSharedUsersHandler(c *gin.Context) {
	userID := common.MustGetAccessUserID(c)
	todoListID, err := common.GetRequiredParamID(c, "id")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	users, err := bll.GetTodoListSharedUsers(userID, todoListID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	// TODO[bug]: filter user fields

	c.JSON(http.StatusOK, users)
}

// Delete shared user from todo list
func DeleteTodoListSharedUserHandler(c *gin.Context) {
	operatorID := common.MustGetAccessUserID(c)
	todoListID, err := common.GetRequiredParamID(c, "id")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	userID, err := common.GetRequiredParamID(c, "user-id")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	if err := bll.DeleteTodoListSharedUser(operatorID, userID, todoListID); err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

// Create share link for todo list
func PostTodoListSharingsHandler(c *gin.Context) {
	userID := common.MustGetAccessUserID(c)
	todoListID, err := common.GetRequiredParamID(c, "id")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	sharing, err := bll.CreateTodoListSharing(userID, todoListID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.SharingToken{
		Token:     sharing.Token,
		Type:      sharing.Type,
		CreatedAt: sharing.CreatedAt,
	})
}

// Get current share link
func GetTodoListSharingsHandler(c *gin.Context) {
	userID := common.MustGetAccessUserID(c)
	todoListID, err := common.GetRequiredParamID(c, "id")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	sharings, err := bll.GetActiveTodoListSharings(userID, todoListID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	vec := make([]dto.SharingToken, 0)
	for i := range sharings {
		vec = append(vec, dto.SharingToken{
			Token:     sharings[i].Token,
			Type:      sharings[i].Type,
			CreatedAt: sharings[i].CreatedAt,
		})
	}

	c.JSON(http.StatusOK, vec)
}

// Close todo list sharing
func DeleteTodoListISSharingsHandler(c *gin.Context) {
	// TODO: Implement this
	operatorID := common.MustGetAccessUserID(c)
	todoListID, err := common.GetRequiredParamID(c, "id")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	userID, err := common.GetRequiredParamID(c, "user-id")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	if err := bll.DeleteTodoListSharedUser(operatorID, userID, todoListID); err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

// Join todo list by share token
func PostTodoListSharingHandler(c *gin.Context) {
	userID := common.MustGetAccessUserID(c)
	token, err := common.GetRequiredParam(c, "token")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	if err := bll.CreateTodoListSharedUser(userID, token); err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

// Inactive share link
func DeleteTodoListSharingHandler(c *gin.Context) {
	userID := common.MustGetAccessUserID(c)
	token, err := common.GetRequiredParam(c, "token")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	if err := bll.DeleteTodoListSharing(userID, token); err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/bll"
	"github.com/yzx9/otodo/entity"
	"github.com/yzx9/otodo/otodo"
	"github.com/yzx9/otodo/utils"
	"github.com/yzx9/otodo/web/common"
)

// Create todo step
func PostTodoStepHandler(c *gin.Context) {
	todoID, err := common.GetRequiredParam(c, "id")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	payload := struct {
		Name string `json:"name"`
	}{}
	if c.ShouldBind(&payload) != nil {
		common.AbortWithError(c, utils.NewError(otodo.ErrPreconditionRequired, "name required"))
		return
	}

	userID := common.MustGetAccessUserID(c)
	step, err := bll.CreateTodoStep(userID, todoID, payload.Name)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, step)
}

// Update todo step
func PutTodoStepHandler(c *gin.Context) {
	step := entity.TodoStep{}
	err := c.ShouldBind(&step)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	userID := common.MustGetAccessUserID(c)
	newStep, err := bll.UpdateTodoStep(userID, step)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, newStep)
}

// Delete todo step
func DeleteTodoStepHandler(c *gin.Context) {
	todoID, err := common.GetRequiredParam(c, "id")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	stepID, ok := c.Params.Get("step-id")
	if !ok {
		common.AbortWithError(c, fmt.Errorf("step-id required"))
		return
	}

	userID := common.MustGetAccessUserID(c)
	newStep, err := bll.DeleteTodoStep(userID, todoID, stepID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, newStep)
}

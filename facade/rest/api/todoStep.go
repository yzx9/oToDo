package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/application/dto"
	"github.com/yzx9/otodo/domain/todo"
	"github.com/yzx9/otodo/facade/rest/common"
	"github.com/yzx9/otodo/infrastructure/errors"
	"github.com/yzx9/otodo/infrastructure/repository"
	"github.com/yzx9/otodo/infrastructure/util"
)

// Create todo step
func PostTodoStepHandler(c *gin.Context) {
	todoID, err := common.GetRequiredParamID(c, "id")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	payload := dto.TodoStepDTO{}
	if c.ShouldBind(&payload) != nil {
		common.AbortWithError(c, util.NewError(errors.ErrPreconditionRequired, "name required"))
		return
	}

	userID := common.MustGetAccessUserID(c)
	step, err := todo.CreateTodoStep(userID, todoID, payload.Name)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, step)
}

// Update todo step
func PutTodoStepHandler(c *gin.Context) {
	step := repository.TodoStep{}
	if err := c.ShouldBind(&step); err != nil {
		common.AbortWithError(c, err)
		return
	}

	stepID, err := common.GetRequiredParamID(c, "step-id")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	userID := common.MustGetAccessUserID(c)
	step.ID = stepID
	if err := todo.UpdateTodoStep(userID, &step); err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, step)
}

// Delete todo step
func DeleteTodoStepHandler(c *gin.Context) {
	todoID, err := common.GetRequiredParamID(c, "id")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	stepID, err := common.GetRequiredParamID(c, "step-id")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	userID := common.MustGetAccessUserID(c)
	step, err := todo.DeleteTodoStep(userID, todoID, stepID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, step)
}

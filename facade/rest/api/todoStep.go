package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/application/dto"
	"github.com/yzx9/otodo/application/service"
	"github.com/yzx9/otodo/facade/rest/common"
	"github.com/yzx9/otodo/infrastructure/errors"
	"github.com/yzx9/otodo/infrastructure/util"
)

// Create todo step
func PostTodoStepHandler(c *gin.Context) {
	dto := dto.NewTodoStep{}
	if c.ShouldBind(&dto) != nil {
		common.AbortWithError(c, util.NewError(errors.ErrPreconditionRequired, "name required"))
		return
	}

	todoID, err := common.GetRequiredParamID(c, "id")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}
	dto.TodoID = todoID

	userID := common.MustGetAccessUserID(c)
	step, err := service.CreateTodoStep(userID, dto)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, step)
}

// Update todo step
func PutTodoStepHandler(c *gin.Context) {
	step := dto.TodoStep{}
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
	newStep, err := service.UpdateTodoStep(userID, step)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, newStep)
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
	step, err := service.DeleteTodoStep(userID, todoID, stepID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, step)
}

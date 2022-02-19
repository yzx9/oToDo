package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/api/common"
	"github.com/yzx9/otodo/bll"
	"github.com/yzx9/otodo/model/dto"
	"github.com/yzx9/otodo/model/entity"
	"github.com/yzx9/otodo/otodo"
	"github.com/yzx9/otodo/util"
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
		common.AbortWithError(c, util.NewError(otodo.ErrPreconditionRequired, "name required"))
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
	if err := bll.UpdateTodoStep(userID, &step); err != nil {
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
	step, err := bll.DeleteTodoStep(userID, todoID, stepID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, step)
}

package common

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetParamUUID(c *gin.Context, key string) (uuid.UUID, error) {
	idRaw, ok := c.Params.Get(key)
	if !ok {
		return uuid.UUID{}, fmt.Errorf("%v required", key)
	}

	id, err := uuid.Parse(idRaw)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("invalid %v", key)
	}

	return id, nil
}

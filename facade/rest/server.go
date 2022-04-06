package rest

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/facade/rest/middleware"
	"github.com/yzx9/otodo/infrastructure/config"
)

func Run() error {
	r := gin.New()
	r.Use(
		gin.Logger(),
		gin.Recovery(),
		middleware.CORSMiddleware(),
		middleware.ErrorMiddleware())

	setupRouter(r)

	port := config.Server.Port
	if port == 0 {
		port = 8080
	}
	host := config.Server.Host
	addr := fmt.Sprintf("%v:%v", host, port)

	r.Run(addr)
	return nil
}

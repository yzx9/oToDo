package web

import (
	"github.com/gin-gonic/gin"

	"github.com/yzx9/otodo/web/handlers"
)

func setupRouter(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		// Ping test
		v1.GET("/ping", handlers.PingHandler)

		// Todo
		v1.GET("/todo", handlers.GetTodosHandler)
	}
}

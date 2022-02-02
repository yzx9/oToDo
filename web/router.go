package web

import (
	"github.com/gin-gonic/gin"

	"github.com/yzx9/otodo/web/handlers"
)

func setupRouter(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		// Auth
		v1.POST("/login", handlers.LoginHandler)
		v1.POST("/logout", handlers.LogoutHandler)
	}
	{
		// Ping test
		v1.GET("/ping", handlers.PingHandler)

		// Todo
		v1.GET("/todo/:id", handlers.GetTodosHandler)

		// Todo List
		v1.GET("/todolist", handlers.GetTodoListsHandler)
		v1.GET("/todolist/:id", handlers.GetTodoListHandler)
	}
}

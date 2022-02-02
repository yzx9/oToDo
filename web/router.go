package web

import (
	"github.com/gin-gonic/gin"

	"github.com/yzx9/otodo/web/handlers"
	"github.com/yzx9/otodo/web/middlewares"
)

func setupRouter(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	setupRouterV1(v1)
}

func setupRouterV1(r *gin.RouterGroup) {
	// Public routes
	{
		// Auth
		r.POST("/login", handlers.LoginHandler)
		r.POST("/logout", handlers.LogoutHandler)

		// Ping test
		r.GET("/ping", handlers.PingHandler)
	}

	// Authorized routes
	r = r.Group("/", middlewares.JwtAuthMiddleware)
	{
		// Ping test
		r.GET("/hello", handlers.HelloHandler)

		// Todo
		r.GET("/todo/:id", handlers.GetTodosHandler)

		// Todo List
		r.GET("/todolist", handlers.GetTodoListsHandler)
		r.GET("/todolist/:id", handlers.GetTodoListHandler)
	}
}

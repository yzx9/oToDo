package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/web/handlers"
	"github.com/yzx9/otodo/web/middlewares"
)

func setupRouter(e *gin.Engine) {
	r := e.Group("/api")

	// Public routes
	{
		// Ping test
		r.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })

		// Auth
		r.POST("/session", handlers.PostSessionHandler)
		r.DELETE("/session", handlers.DeleteSessionHandler)
	}

	// Authorized routes
	r = r.Group("/", middlewares.JwtAuthMiddleware)
	{
		// Auth
		r.GET("/session", handlers.GetSessionHandler)
		r.GET("/session/token", handlers.PostAccessTokenHandler)

		// Todo
		r.GET("/todo/:id", handlers.GetTodosHandler)

		// Todo List
		r.GET("/todolist", handlers.GetTodoListsHandler)
		r.GET("/todolist/:id", handlers.GetTodoListHandler)
	}
}

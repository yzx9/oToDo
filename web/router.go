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
		r.GET("/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		})

		// Session
		r.POST("/sessions", handlers.PostSessionHandler)

		r.POST("/sessions/current/tokens", handlers.PostSessionTokenHandler)

		// User
		r.POST("/users", handlers.PostUserHandler)
	}

	// Authorized routes
	r = r.Group("/", middlewares.JwtAuthMiddleware())
	{
		// Session
		r.GET("/sessions", handlers.GetSessionHandler)
		r.DELETE("/sessions", handlers.DeleteSessionHandler)

		// File
		// r.MaxMultipartMemory = MaxFileSize // 限制 Gin 上传文件时最大内存 (默认 32 MiB)
		r.GET("/files/:id", handlers.GetFileHandler)

		// User
		r.GET("/users/current", handlers.GetCurrentUserHandler)

		r.GET("/users/current/todo-lists", handlers.GetCurrentUserTodoListsHandler)
		r.GET("/users/current/todo-lists/basic", handlers.GetCurrentUserBasicTodoListHandler)

		r.GET("/users/current/todos/daily", handlers.GetCurrentUserDailyTodosHandler)
		r.GET("/users/current/todos/planned", handlers.GetCurrentUserPlannedTodosHandler)
		r.GET("/users/current/todos/important", handlers.GetCurrentUserImportantTodosHandler)

		r.GET("/users/current/todo-list-folders", handlers.GetCurrentUserTodoListFoldersHandler)

		// Todo
		r.POST("/todos", handlers.PostTodoHandler)
		r.PUT("/todos/:id", handlers.PutTodoHandler)
		r.PATCH("/todos/:id", handlers.PatchTodoHandler)
		r.GET("/todos/:id", handlers.GetTodoHandler)
		r.DELETE("/todos/:id", handlers.DeleteTodoHandler)

		r.POST("/todos/:id/files", handlers.PostTodoFileHandler)

		r.POST("/todos/:id/steps", handlers.PostTodoStepHandler)
		r.PUT("/todos/:id/steps", handlers.PutTodoStepHandler)

		r.POST("/todos/:id/steps/:step-id", handlers.DeleteTodoStepHandler)

		// Todo List
		r.GET("/todo-lists/:id", handlers.GetTodoListHandler)
		r.DELETE("/todo-lists/:id", handlers.DeleteTodoListHandler)

		r.GET("/todo-lists/:id/todos", handlers.GetTodoListTodosHandler)

		// Todo List Folder
		r.GET("/todo-list-folders/:id", handlers.GetTodoListFolderHandler)
		r.DELETE("/todo-list-folders/:id", handlers.DeleteTodoListFolderHandler)
	}
}

package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/web/handlers"
	"github.com/yzx9/otodo/web/middlewares"
)

func (s *Server) setupRouter() {
	r := s.engine.Group("/api")

	// Public routes
	{
		// Ping test
		r.GET("/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		})

		// File
		// r.MaxMultipartMemory = MaxFileSize // 限制 Gin 上传文件时最大内存 (默认 32 MiB)
		r.POST("/files", handlers.PostFileHandler)
		r.GET("/files/:id", handlers.GetFileHandler)
		r.GET("/files/presigned/:id", handlers.GetFilePresignedHandler) // TODO[feat]: 或许不需要独立API

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
		r.POST("/files/:id/presign", handlers.PostFilePresignHandler) // TODO[feat]: 调整API, 允许管理PreSign，允许设置有效时长

		// Current User
		r.GET("/users/current", handlers.GetCurrentUserHandler)

		r.GET("/users/current/todo-lists", handlers.GetCurrentUserTodoListsHandler)
		r.GET("/users/current/todo-lists/basic", handlers.GetCurrentUserBasicTodoListHandler)

		r.GET("/users/current/todos/daily", handlers.GetCurrentUserDailyTodosHandler)
		r.GET("/users/current/todos/planned", handlers.GetCurrentUserPlannedTodosHandler)
		r.GET("/users/current/todos/important", handlers.GetCurrentUserImportantTodosHandler)
		r.GET("/users/current/todos/not-notified", handlers.GetCurrentUserNotNotifiedTodosHandler)

		r.GET("/users/current/todo-list-folders", handlers.GetCurrentUserTodoListFoldersHandler)

		// Todo
		r.POST("/todos", handlers.PostTodoHandler) // TODO[feat]: done 属性是否需要独立API？否则无法返回重复产生的Todo
		r.PUT("/todos/:id", handlers.PutTodoHandler)
		r.PATCH("/todos/:id", handlers.PatchTodoHandler)
		r.GET("/todos/:id", handlers.GetTodoHandler)
		r.DELETE("/todos/:id", handlers.DeleteTodoHandler)

		r.POST("/todos/:id/files", handlers.PostTodoFileHandler)

		r.POST("/todos/:id/steps", handlers.PostTodoStepHandler)
		r.PUT("/todos/:id/steps", handlers.PutTodoStepHandler)

		r.POST("/todos/:id/steps/:step-id", handlers.DeleteTodoStepHandler)

		// Todo List
		r.POST("/todo-lists", handlers.PostTodoListHandler)
		r.GET("/todo-lists/:id", handlers.GetTodoListHandler)
		r.DELETE("/todo-lists/:id", handlers.DeleteTodoListHandler)

		r.GET("/todo-lists/:id/todos", handlers.GetTodoListTodosHandler)

		r.GET("/todo-lists/:id/shared-users", handlers.GetTodoListSharedUsersHandler)
		r.DELETE("/todo-lists/:id/shared-users/:user-id", handlers.DeleteTodoListSharedUserHandler)

		r.POST("/todo-lists/:id/sharings", handlers.PostTodoListSharingsHandler)
		r.GET("/todo-lists/:id/share-links", handlers.GetTodoListSharingsHandler)

		r.POST("/todo-lists/:id/sharings/:token", handlers.PostTodoListSharingHandler)
		r.DELETE("/todo-lists/:id/sharings/:token", handlers.DeleteTodoListSharingHandler)

		// Todo List Folder
		r.GET("/todo-list-folders/:id", handlers.GetTodoListFolderHandler)
		r.DELETE("/todo-list-folders/:id", handlers.DeleteTodoListFolderHandler)

		// Sharing
		r.GET("/sharings/:token", handlers.GetSharingHandler)
	}
}

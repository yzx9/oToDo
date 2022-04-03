package user_interface

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/user_interface/api"
	"github.com/yzx9/otodo/user_interface/middleware"
)

func (s *Server) setupRouter() {
	r := s.engine.Group("/api")

	// Public routes
	{
		// Ping test
		r.GET("/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		})

		// Session
		r.POST("/sessions", api.PostSessionHandler)

		r.GET("/sessions/oauth/github", api.GetSessionOAuthGithub)
		r.POST("/sessions/oauth/github", api.PostSessionOAuthGithub)

		r.POST("/sessions/current/tokens", api.PostSessionTokenHandler)

		// File
		// r.MaxMultipartMemory = MaxFileSize // 限制 Gin 上传文件时最大内存 (默认 32 MiB)
		r.POST("/files", api.PostFileHandler)
		r.GET("/files/:id", api.GetFileHandler)

		// User
		r.POST("/users", api.PostUserHandler)

		// Sharing
		r.GET("/sharings/:token", api.GetSharingHandler)
		r.GET("/sharings/:token/todo-list", api.GetSharingTodoListHandler)
	}

	// Authorized routes
	r = r.Group("/", middleware.JwtAuthMiddleware())
	{
		// Session
		r.GET("/sessions", api.GetSessionHandler)
		r.DELETE("/sessions", api.DeleteSessionHandler)

		// File
		r.POST("/files/:id/pre-sign", api.PostFilePreSignHandler) // TODO[feat]: 调整API, 允许管理PreSign，允许设置有效时长

		// Current User
		r.GET("/users/current", api.GetCurrentUserHandler)

		r.GET("/users/current/menu", api.GetCurrentUserMenu)

		r.GET("/users/current/todo-lists", api.GetCurrentUserTodoListsHandler)

		r.GET("/users/current/todos/basic", api.GetCurrentUserBasicTodoListTodosHandler)
		r.GET("/users/current/todos/daily", api.GetCurrentUserDailyTodosHandler)
		r.GET("/users/current/todos/planned", api.GetCurrentUserPlannedTodosHandler)
		r.GET("/users/current/todos/important", api.GetCurrentUserImportantTodosHandler)
		r.GET("/users/current/todos/not-notified", api.GetCurrentUserNotNotifiedTodosHandler)

		r.GET("/users/current/todo-list-folders", api.GetCurrentUserTodoListFoldersHandler)

		// Todo
		r.POST("/todos", api.PostTodoHandler) // TODO[feat]: done 属性是否需要独立API？否则无法返回重复产生的Todo
		r.PUT("/todos/:id", api.PutTodoHandler)
		r.PATCH("/todos/:id", api.PatchTodoHandler)
		r.GET("/todos/:id", api.GetTodoHandler)
		r.DELETE("/todos/:id", api.DeleteTodoHandler)

		r.POST("/todos/:id/files", api.PostTodoFileHandler)

		r.POST("/todos/:id/steps", api.PostTodoStepHandler)
		r.PUT("/todos/:id/steps/:step-id", api.PutTodoStepHandler)
		r.DELETE("/todos/:id/steps/:step-id", api.DeleteTodoStepHandler)

		// Todo List
		r.POST("/todo-lists", api.PostTodoListHandler)
		r.GET("/todo-lists/:id", api.GetTodoListHandler)
		r.DELETE("/todo-lists/:id", api.DeleteTodoListHandler)

		r.GET("/todo-lists/:id/todos", api.GetTodoListTodosHandler)

		r.GET("/todo-lists/:id/shared-users", api.GetTodoListSharedUsersHandler)
		r.DELETE("/todo-lists/:id/shared-users/:user-id", api.DeleteTodoListSharedUserHandler)

		r.POST("/todo-lists/:id/sharings", api.PostTodoListSharingsHandler)
		r.GET("/todo-lists/:id/sharings", api.GetTodoListSharingsHandler)

		r.POST("/todo-lists/:id/sharings/:token", api.PostTodoListSharingHandler)
		r.DELETE("/todo-lists/:id/sharings/:token", api.DeleteTodoListSharingHandler)

		// Todo List Folder
		r.POST("/todo-list-folders", api.PostTodoListFolderHandler)
		r.GET("/todo-list-folders/:id", api.GetTodoListFolderHandler)
		r.DELETE("/todo-list-folders/:id", api.DeleteTodoListFolderHandler)
	}
}

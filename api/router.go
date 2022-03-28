package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/api/handler"
	"github.com/yzx9/otodo/api/middleware"
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
		r.POST("/sessions", handler.PostSessionHandler)

		r.GET("/sessions/oauth/github", handler.GetSessionOAuthGithub)
		r.POST("/sessions/oauth/github", handler.PostSessionOAuthGithub)

		r.POST("/sessions/current/tokens", handler.PostSessionTokenHandler)

		// File
		// r.MaxMultipartMemory = MaxFileSize // 限制 Gin 上传文件时最大内存 (默认 32 MiB)
		r.POST("/files", handler.PostFileHandler)
		r.GET("/files/:id", handler.GetFileHandler)

		// User
		r.POST("/users", handler.PostUserHandler)

		// Sharing
		r.GET("/sharings/:token", handler.GetSharingHandler)
		r.GET("/sharings/:token/todo-list", handler.GetSharingTodoListHandler)

		//SmsCode
		r.GET("/sendcode", handler.SendSmsCode)
		r.POST("/login_sms", handler.SmsLogin)
	}

	// Authorized routes
	r = r.Group("/", middleware.JwtAuthMiddleware())
	{
		// Session
		r.GET("/sessions", handler.GetSessionHandler)
		r.DELETE("/sessions", handler.DeleteSessionHandler)

		// File
		r.POST("/files/:id/pre-sign", handler.PostFilePreSignHandler) // TODO[feat]: 调整API, 允许管理PreSign，允许设置有效时长

		// Current User
		r.GET("/users/current", handler.GetCurrentUserHandler)

		r.GET("/users/current/menu", handler.GetCurrentUserMenu)

		r.GET("/users/current/todo-lists", handler.GetCurrentUserTodoListsHandler)

		r.GET("/users/current/todos/basic", handler.GetCurrentUserBasicTodoListTodosHandler)
		r.GET("/users/current/todos/daily", handler.GetCurrentUserDailyTodosHandler)
		r.GET("/users/current/todos/planned", handler.GetCurrentUserPlannedTodosHandler)
		r.GET("/users/current/todos/important", handler.GetCurrentUserImportantTodosHandler)
		r.GET("/users/current/todos/not-notified", handler.GetCurrentUserNotNotifiedTodosHandler)

		r.GET("/users/current/todo-list-folders", handler.GetCurrentUserTodoListFoldersHandler)

		// Todo
		r.POST("/todos", handler.PostTodoHandler) // TODO[feat]: done 属性是否需要独立API？否则无法返回重复产生的Todo
		r.PUT("/todos/:id", handler.PutTodoHandler)
		r.PATCH("/todos/:id", handler.PatchTodoHandler)
		r.GET("/todos/:id", handler.GetTodoHandler)
		r.DELETE("/todos/:id", handler.DeleteTodoHandler)

		r.POST("/todos/:id/files", handler.PostTodoFileHandler)

		r.POST("/todos/:id/steps", handler.PostTodoStepHandler)
		r.PUT("/todos/:id/steps/:step-id", handler.PutTodoStepHandler)
		r.DELETE("/todos/:id/steps/:step-id", handler.DeleteTodoStepHandler)

		// Todo List
		r.POST("/todo-lists", handler.PostTodoListHandler)
		r.GET("/todo-lists/:id", handler.GetTodoListHandler)
		r.DELETE("/todo-lists/:id", handler.DeleteTodoListHandler)

		r.GET("/todo-lists/:id/todos", handler.GetTodoListTodosHandler)

		r.GET("/todo-lists/:id/shared-users", handler.GetTodoListSharedUsersHandler)
		r.DELETE("/todo-lists/:id/shared-users/:user-id", handler.DeleteTodoListSharedUserHandler)

		r.POST("/todo-lists/:id/sharings", handler.PostTodoListSharingsHandler)
		r.GET("/todo-lists/:id/sharings", handler.GetTodoListSharingsHandler)
		r.DELETE("/todo-lists/:id/sharings", handler.DeleteTodoListSharingHandler)

		r.POST("/todo-lists/:id/sharings/:token", handler.PostTodoListSharingHandler)
		r.DELETE("/todo-lists/:id/sharings/:token", handler.DeleteTodoListSharingHandler)

		// Todo List Folder
		r.POST("/todo-list-folders", handler.PostTodoListFolderHandler)
		r.GET("/todo-list-folders/:id", handler.GetTodoListFolderHandler)
		r.DELETE("/todo-list-folders/:id", handler.DeleteTodoListFolderHandler)
	}
}

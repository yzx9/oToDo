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

		// Auth
		r.POST("/session", handlers.PostSessionHandler)
		r.DELETE("/session", handlers.DeleteSessionHandler)

		r.POST("/session/token", handlers.PostSessionTokenHandler)
	}

	// Authorized routes
	r = r.Group("/", middlewares.JwtAuthMiddleware())
	{
		// Auth
		r.GET("/session", handlers.GetSessionHandler)

		// File
		r.POST("/file", handlers.PostFileHandler)
		r.GET("/file/:id", handlers.GetFileHandler)

		// Todo
		r.POST("/todo", handlers.PostTodoHandler)
		r.PUT("/todo/:id", handlers.PutTodoHandler)
		r.PATCH("/todo/:id", handlers.PatchTodoHandler)
		r.GET("/todo/:id", handlers.GetTodoHandler)
		r.DELETE("/todo/:id", handlers.DeleteTodoHanlder)

		// Todo List
		// r.MaxMultipartMemory = MaxFileSize // 限制 Gin 上传文件时最大内存 (默认 32 MiB)
		r.GET("/todolist", handlers.GetTodoListsHandler)
		r.GET("/todolist/:id", handlers.GetTodoListHandler)
	}
}

package web

import (
	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/web/middlewares"
)

type Server struct {
	engine *gin.Engine
	addr   string
}

func CreateServer() *Server {
	r := gin.New()
	r.Use(
		gin.Logger(),
		gin.Recovery(),
		middlewares.ErrorMiddleware())

	setupRouter(r)

	return &Server{
		engine: r,
		addr:   ":8080",
	}
}

func (s *Server) Listen(addr string) *Server {
	s.addr = addr
	return s
}

func (s *Server) Run() *Server {
	s.engine.Run(s.addr)
	return s
}

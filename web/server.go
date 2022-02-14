package web

import (
	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/bll"
	"github.com/yzx9/otodo/web/middlewares"
)

type Server struct {
	Error  error
	engine *gin.Engine
	addr   string
}

func CreateServer() *Server {
	if err := bll.Init(); err != nil {
		return &Server{
			Error: err,
		}
	}

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
	if s.Error != nil {
		return s
	}

	s.addr = addr
	return s
}

func (s *Server) Run() *Server {
	if s.Error != nil {
		return s
	}

	s.engine.Run(s.addr)
	return s
}

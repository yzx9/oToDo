package user_interface

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/yzx9/otodo/bll"
	"github.com/yzx9/otodo/crosscutting"
	"github.com/yzx9/otodo/infrastructure/config"
	"github.com/yzx9/otodo/user_interface/middleware"
)

type Server struct {
	Error error

	config *viper.Viper
	engine *gin.Engine
	addr   string
}

func NewServer() *Server {
	r := gin.New()
	r.Use(
		gin.Logger(),
		gin.Recovery(),
		middleware.CORSMiddleware(),
		middleware.ErrorMiddleware())

	s := Server{
		engine: r,
		addr:   ":8080",
	}

	s.setupRouter()
	return &s
}

func (s *Server) LoadConfig(dir string) *Server {
	if s.Error != nil {
		return s
	}

	config, err := crosscutting.LoadConfig(dir)
	if err != nil {
		s.Error = fmt.Errorf("fails to load config: %w", err)
		return s
	}
	s.config = config

	return s
}

func (s *Server) LoadAndWatchConfig(dir string) *Server {
	if s.Error != nil {
		return s
	}

	config, err := crosscutting.LoadAndWatchConfig(dir)
	if err != nil {
		s.Error = fmt.Errorf("fails to load config: %w", err)
		return s
	}
	s.config = config

	return s
}

func (s *Server) Run() *Server {
	if s.Error != nil {
		return s
	}

	if err := crosscutting.StartUp(); err != nil {
		s.Error = err
		return s
	}

	if err := bll.StartUp(); err != nil {
		s.Error = err
		return s
	}

	port := config.Server.Port
	if port == 0 {
		port = 8080
	}
	host := config.Server.Host

	s.addr = fmt.Sprintf("%v:%v", host, port)

	s.engine.Run(s.addr)
	return s
}

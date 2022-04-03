package user_interface

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/yzx9/otodo/application"
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

	s.config = viper.New()
	s.config.SetConfigType("yaml")
	s.config.AddConfigPath(dir)

	s.config.SetConfigName("config.yaml")
	if err := s.config.ReadInConfig(); err != nil {
		s.Error = fmt.Errorf("fails to load config.yaml: %w", err)
		return s
	}

	s.config.SetConfigName("secret.yaml")
	if err := s.config.MergeInConfig(); err != nil {
		s.Error = fmt.Errorf("fails to load secret.yaml: %w", err)
		return s
	}

	config.SetConfig(s.config)

	return s
}

func (s *Server) LoadAndWatchConfig(dir string) *Server {
	if s.Error != nil {
		return s
	}

	s.LoadConfig(dir)

	s.config.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed: ", e.Name)
		config.SetConfig(s.config)
	})

	s.config.WatchConfig()

	return s
}

func (s *Server) Run() *Server {
	if s.Error != nil {
		return s
	}

	if err := application.StartUp(); err != nil {
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

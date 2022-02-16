package web

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/yzx9/otodo/bll"
	"github.com/yzx9/otodo/web/middleware"
)

type Server struct {
	Error error

	config *viper.Viper
	engine *gin.Engine
	addr   string
}

func CreateServer() *Server {
	r := gin.New()
	r.Use(
		gin.Logger(),
		gin.Recovery(),
		middleware.ErrorMiddleware())

	s := Server{
		engine: r,
		addr:   ":8080",
	}

	s.setupRouter()
	return &s
}

func (s *Server) LoadConfig(dir string) {
	if s.Error != nil {
		return
	}

	s.config = viper.New()
	s.config.SetConfigType("yaml")
	s.config.AddConfigPath(dir)

	s.config.SetConfigName("config.yaml")
	if err := s.config.ReadInConfig(); err != nil {
		s.Error = fmt.Errorf("fails to load config.yaml: %w", err)
		return
	}

	s.config.SetConfigName("secret.yaml")
	if err := s.config.MergeInConfig(); err != nil {
		s.Error = fmt.Errorf("fails to load secret.yaml: %w", err)
		return
	}

	SetConfig(s.config)
}

func (s *Server) LoadAndWatchConfig(dir string) {
	if s.Error != nil {
		return
	}

	s.LoadConfig(dir)

	s.config.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed: ", e.Name)
		SetConfig(s.config)
	})

	s.config.WatchConfig()
}

func (s *Server) Listen(addr string) {
	if s.Error != nil {
		return
	}

	s.addr = addr
}

func (s *Server) Run() {
	if s.Error != nil {
		return
	}

	if err := bll.Init(); err != nil {
		s.Error = err
		return
	}

	s.engine.Run(s.addr)
}

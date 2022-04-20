package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/config"
	"github.com/yzx9/otodo/driver/rest/middleware"
)

type Server struct {
	server      *http.Server
	shutdown    bool
	errorStream chan error
}

func Run() (s *Server) {
	r := gin.New()
	r.Use(
		gin.Logger(),
		gin.Recovery(),
		middleware.CORSMiddleware(),
		middleware.ErrorMiddleware())

	setupRouter(r)

	port := config.Server.Port
	if port == 0 {
		port = 8080
	}
	host := config.Server.Host
	addr := fmt.Sprintf("%v:%v", host, port)

	s = new(Server)
	s.server = &http.Server{
		Addr:    addr,
		Handler: r,
	}

	s.errorStream = make(chan error)
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.errorStream <- fmt.Errorf("rest server: %w", err)
			s.Shutdown(context.Background())
		}
	}()

	return s
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.shutdown {
		return nil
	}

	close(s.errorStream)
	s.shutdown = true

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}

	return nil
}

func (s *Server) ErrorStream() <-chan error {
	return s.errorStream
}

package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/config"
	"github.com/yzx9/otodo/driver/rest/middleware"
)

func Run() (shutdown func(ctx context.Context) error, errStream <-chan error) {
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

	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	stream := make(chan error)
	errStream = stream
	shutdown = func(ctx context.Context) error {
		if err := server.Shutdown(ctx); err != nil {
			return fmt.Errorf("server shutdown: %w", err)
		}

		return nil
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			stream <- fmt.Errorf("rest server: %w", err)
		}
	}()

	return
}

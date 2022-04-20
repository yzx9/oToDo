package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/yzx9/otodo/config"
	"github.com/yzx9/otodo/driver/rest"
)

func main() {
	// init
	log.Println("load and watch config...")
	onConfigChange, err := config.LoadAndWatch(".")
	if err != nil {
		log.Fatalf("fails to load and watch config: %s", err.Error())
	}

	log.Println("start up application...")
	if err := startUp(); err != nil {
		log.Fatalf("fails to start up application: %s", err.Error())
	}

	// run server
	log.Println("run rest server...")
	shutdownRestServer, restServerErrorStream := rest.Run()
	log.Println("serving...")

	// listen events
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case <-onConfigChange:
				fmt.Println("config file changed.")

			case <-restServerErrorStream:
				fmt.Println("[REST] ", err.Error())
			}
		}
	}()

	// wait quit
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	signal := <-quit
	log.Println("receive signal: ", signal.String())

	log.Println("shutdown server...")
	ctxShutdown, cancelShutdown := context.WithTimeout(ctx, 10*time.Second)
	defer cancelShutdown()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := shutdownRestServer(ctxShutdown); err != nil {
			log.Println("fails to shutdown rest server: %w", err)
		}
	}()

	wg.Wait()
	log.Println("server shutdown.")
}

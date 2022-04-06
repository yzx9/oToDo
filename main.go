package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/yzx9/otodo/crosscutting"
	"github.com/yzx9/otodo/facade/rest"
)

func main() {
	if _, err := crosscutting.LoadAndWatchConfig("."); err != nil {
		err = fmt.Errorf("fails to load and watch config: %w", err)
		log.Fatal(err)
	}

	if err := crosscutting.StartUp(); err != nil {
		err = fmt.Errorf("fails to start up: %w", err)
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	runRestServer := func() {
		defer wg.Done()

		if err := rest.Run(); err != nil {
			err = fmt.Errorf("fails to run rest server: %w", err)
			fmt.Println(err)
		}
	}

	wg.Add(1)
	go runRestServer()

	wg.Wait()
	log.Fatal("all server down.")
}

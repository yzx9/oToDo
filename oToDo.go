package main

import (
	"log"

	"github.com/yzx9/otodo/web"
)

func main() {
	s := web.CreateServer()
	s.LoadAndWatchConfig(".")
	s.Listen("localhost:8080")
	s.Run()

	if s.Error != nil {
		log.Fatal(s.Error)
	}
}

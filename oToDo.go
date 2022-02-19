package main

import (
	"log"

	"github.com/yzx9/otodo/api"
)

func main() {
	s := api.CreateServer().
		LoadAndWatchConfig(".").
		Listen("localhost:8080").
		Run()

	if s.Error != nil {
		log.Fatal(s.Error)
	}
}

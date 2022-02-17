package main

import (
	"log"

	"github.com/yzx9/otodo/web"
)

func main() {
	s := web.CreateServer().
		LoadAndWatchConfig(".").
		Listen("localhost:8080").
		Run()

	if s.Error != nil {
		log.Fatal(s.Error)
	}
}

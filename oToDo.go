package main

import (
	"log"

	"github.com/yzx9/otodo/api"
)

func main() {
	s := api.NewServer().
		LoadAndWatchConfig(".").
		Run()

	if s.Error != nil {
		log.Fatal(s.Error)
	}
}

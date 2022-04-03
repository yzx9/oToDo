package main

import (
	"log"

	ui "github.com/yzx9/otodo/user_interface"
)

func main() {
	s := ui.NewServer().
		LoadAndWatchConfig(".").
		Run()

	if s.Error != nil {
		log.Fatal(s.Error)
	}
}

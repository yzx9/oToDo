package main

import (
	"fmt"

	"github.com/yzx9/otodo/web"
)

func main() {
	err := web.CreateServer().Listen("localhost:8080").Run().Error
	if err != nil {
		fmt.Println(err)
	}
}

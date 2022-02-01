package main

import "github.com/yzx9/otodo/web"

func main() {
	web.CreateServer().Listen("localhost:8080").Run()
}

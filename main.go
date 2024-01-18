package main

import (
	"doubleblind/server"
	"fmt"
)

func main() {
	fmt.Println("Starting server on localhost:8090")
	server.Run()
}
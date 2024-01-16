package main

import (
	"fmt"
	"doubleblind/server"
)

func main() {
	fmt.Println("Starting server on localhost:8090")
	server.Run()
}
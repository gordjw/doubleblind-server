package main

import (
	"doubleblind/server"
	"flag"
	"fmt"
)

type flags struct {
	host string
	port int
}

func main() {
	f := flags{}
	flag.StringVar( &f.host, "host", "localhost", "Listen on this port" )
	flag.IntVar( &f.port, "port", 8080, "Listen on this port" )
	flag.Parse()

	fmt.Printf("Starting server on %s:%d\n", f.host, f.port)
	server.Run(f.host, f.port)
}
package server

import (
	"fmt"
	"net/http"
)

func Hello() string {
	return "HELLO!"
}

func Run() {
	// Default response for now
	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Welcome to my website")
	})

	// Serve static files
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Start the server and listen on port 80
	http.ListenAndServe(":8090", nil)
}
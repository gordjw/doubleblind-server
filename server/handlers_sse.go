package server

import (
	"context"
	"net/http"
	"time"
)

func (env *Env) getSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)

	// Create a context for handling client disconnection
	_, cancel := context.WithCancel(r.Context())
	defer cancel()

	// Register this http.ResponseWriter in the client list to receive SSE events
	// TODO - This won't scale, need to be able to remove w from env.clients when this connection closes
	env.clients = append(env.clients, w)

	// Keep the connection open
	// TODO - Is there a better way to keep the connection open?
	for {
		time.Sleep(1 * time.Second)
	}
}

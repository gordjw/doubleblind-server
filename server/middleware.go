package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
)

type ContextKey string

const ContextKeyAuthToken ContextKey = "authToken"
const ContextJsonResponseKey ContextKey = "jsonResponse"

func middlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), ContextKeyAuthToken, "1")
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func middlewareCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		next.ServeHTTP(w, r)
	})
}

func middlewareJSONResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}

func middlewareHTMLResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		next.ServeHTTP(w, r)
	})
}

func middlewareLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("req: %s\n", r.RequestURI)

		next.ServeHTTP(w, r)
	})
}

func TrackRoute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		rctx := chi.RouteContext(r.Context())
		fmt.Println(rctx.RoutePatterns)
		routePattern := strings.Join(rctx.RoutePatterns, "")
		fmt.Println("route:", routePattern)
	})
}

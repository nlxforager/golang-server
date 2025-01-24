// Package main
package mux

import (
	"golang-server/cmd/servers/auth/handlers"
	"net/http"
)

// The closes entry point sans sockets.
func NewMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", handlers.Hello())

	return mux
}

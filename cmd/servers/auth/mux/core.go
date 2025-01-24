// Package main
package mux

import (
	"net/http"

	"golang-server/cmd/servers/auth/handlers"
)

// The closes entry point sans sockets.
func NewMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", handlers.Hello())

	return mux
}

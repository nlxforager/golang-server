// Package main
package main

import (
	"golang-server/cmd/servers/auth/handlers"
	"net/http"
)

// The closes entry point sans sockets.
func NewMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /", handlers.Hello())

	return mux
}

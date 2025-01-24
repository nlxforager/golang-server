// Package main
package mux

import (
	"net/http"

	"golang-server/cmd/servers/auth/handlers"
	"golang-server/src/domain/auth"
)

// The closes entry point sans sockets.
func NewMux(auth auth.AuthService) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", handlers.Hello())

	authMux := http.NewServeMux()
	authMux.HandleFunc("GET /password/", handlers.AuthByUsername(auth))

	mux.Handle("/", authMux)
	return mux
}

// Package main
package mux

import (
	"fmt"
	"net/http"

	"golang-server/cmd/servers/auth/handlers"
	"golang-server/src/domain/auth"
	"golang-server/src/infrastructure/messaging/email"
)

type AuthMuxOpts struct {
	Auth auth.AuthService
	Mail email.EmailService
}

type MuxOpts struct {
	*AuthMuxOpts
}

// The closes entry point sans sockets.
func NewMux(opts *MuxOpts) *http.ServeMux {
	if opts == nil {
		opts = &MuxOpts{
			AuthMuxOpts: nil,
		}
	}

	mux := http.NewServeMux()

	{

		helloMux := http.NewServeMux()

		hello := handlers.Hello()
		helloMux.HandleFunc("GET /", func(writer http.ResponseWriter, request *http.Request) {

			fmt.Printf("request.URL.Path %s", request.URL.Path)
			if request.URL.Path != "/" {

				http.NotFound(writer, request)
			} else {
				hello(writer, request)
			}
		})

		mux.Handle("/", helloMux)
	}

	if opts.AuthMuxOpts != nil {
		authHandlers, err := handlers.NewAuthHandler(opts.Auth, opts.Mail)
		if err != nil {
			panic(err)
		}
		mux.HandleFunc("POST /password/", authHandlers.AuthByUsernamePassword())
		mux.HandleFunc("POST /otp/", authHandlers.SubmitOtp())
	}

	return mux
}

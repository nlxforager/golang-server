// Package main
package mux

import (
	"fmt"
	"golang-server/cmd/servers/multifact/handlers"
	"golang-server/cmd/servers/multifact/mux/middlewares"
	"golang-server/src/domain/auth"
	"golang-server/src/domain/email"
	"net/http"
	"reflect"
)

type AuthMuxOpts struct {
	Auth auth.AuthService
	Mail email.OTPEmailer
}

type MuxOpts struct {
	*AuthMuxOpts
}

// The closes entry point sans sockets.
func NewMux(opts *MuxOpts) *http.ServeMux {
	if opts == nil {
		opts = &MuxOpts{
			AuthMuxOpts: nil, // FIXME
		}
	}

	mux := http.NewServeMux()

	{
		helloMux := http.NewServeMux()

		hello := handlers.Hello()
		helloMux.HandleFunc("GET /", middlewares.LogMiddleware(func(writer http.ResponseWriter, request *http.Request) {
			if request.URL.Path != "/" {
				http.NotFound(writer, request)
			} else {
				hello(writer, request)
			}
		}))

		mux.Handle("GET /", helloMux)
	}
	// Routes that should be authenticated or part of authentication flow.
	if opts.AuthMuxOpts != nil {
		authMw := middlewares.BearerAuthMiddleware(opts.Auth)

		authHandlers, err := handlers.NewAuthHandler(opts.Auth, opts.Mail)
		if err != nil {
			panic(err)
		}

		mux.HandleFunc("POST /register/", middlewares.LogMiddleware(authHandlers.RegisterUsernamePassword()))
		mux.HandleFunc("POST /token/", middlewares.LogMiddleware(authHandlers.AuthByUsernamePassword()))
		mux.HandleFunc("POST /otp/", middlewares.LogMiddleware(authHandlers.SubmitOtp()))

		mux.HandleFunc("PATCH /user/", middlewares.LogMiddleware.Wrap(authMw)(authHandlers.PatchUser()))
	}

	{
		func() {
			defer func() {
				_ = recover()
				//err := recover()
				//fmt.Printf("cannot log routes %#v\n", err)
			}()

			return
			httpMux := reflect.ValueOf(mux).Elem()
			routes := httpMux.FieldByName("patterns") // This is the map of routes
			for i := 0; i < routes.Len(); i++ {
				fmt.Println(routes.Index(i).Elem().FieldByName("str").String())
			}
		}()
	}

	return mux
}

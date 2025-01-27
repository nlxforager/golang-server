// Package main
package mux

import (
	"fmt"
	"net/http"
	"reflect"

	"golang-server/cmd/servers/multifact/handlers"
	"golang-server/cmd/servers/multifact/mux/docs"
	"golang-server/cmd/servers/multifact/mux/middlewares"
	"golang-server/src/domain/auth"
	"golang-server/src/domain/email"

	swagger "github.com/swaggo/http-swagger"
)

type AuthMuxOpts struct {
	Auth auth.Authenticator
	Mail email.OTPEmailer
}

type MuxOpts struct {
	*AuthMuxOpts
}

// ShowAccount godoc
// @Summary      swagger api
// @Description  swagger docs
// @Tags         accounts
// @Produce      json
// @Param        id   query      string  true  "file name" 	Enums(index.html, doc.json)
// @Success 200 {string} string "ok, html or json"
// @Header       200              {string}  Content-Type  "content type"
// @Failure      404  {object}  int
// @Failure      500  {object}  int
// @Router       /swagger [get]
func swagHandler(w http.ResponseWriter, r *http.Request) {
	swagger.Handler()(w, r)
}

// NewMux The closes entry point sans sockets.
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

		docs.SwaggerInfo.Description = "This is a ? server."
		docs.SwaggerInfo.Title = "multifact"

		helloMux.HandleFunc("GET /swagger", swagHandler)

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
		// ShowAccount godoc
		// @Summary      Show an account
		// @Description  get string by ID
		// @Tags         accounts
		// @Accept       json
		// @Produce      json
		// @Param        id   path      int  true  "Account ID"
		// @Success      200  {object}  model.Account
		// @Failure      400  {object}  httputil.HTTPError
		// @Failure      404  {object}  httputil.HTTPError
		// @Failure      500  {object}  httputil.HTTPError
		// @Router       /accounts/{id} [get]
		mux.HandleFunc("PATCH /user/", middlewares.LogMiddleware.Wrap(authMw)(authHandlers.PatchUser()))

		/*

			GET /
				2025-01-27 07:57:04.939................................REQUEST................................
				GET / HTTP/1.1
				Header: Accept: application/json

				2025-01-27 07:57:04.939................................RESPONSE................................
				HTTP/?.? 200 OK
				Content-Type: text/plain; charset=utf-8
				{"data": {"message": "Hello World"}}

			POST /register/
			POST /token/
			POST /otp/
			PATCH /user/
		*/
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

//GET /swagger
//GET /
//POST /register/
//POST /token/
//POST /otp/
//PATCH /user/

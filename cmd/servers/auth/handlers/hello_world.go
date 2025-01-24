package handlers

import (
	"html/template"
	"log/slog"
	"net/http"

	"golang-server/src/domain/hello"
	"golang-server/src/log"
)

type AcceptFunc func(w http.ResponseWriter, r *http.Request)
type AcceptFuncsOpts struct {
	AcceptFuncs map[string]AcceptFunc
	DefaultFunc func(w http.ResponseWriter, r *http.Request)
}

type Options struct {
	AcceptFuncsOpts
}

func (o *Options) GetAcceptFunc(accepts []string) AcceptFunc {
	for _, accept := range accepts {
		if f, ok := o.AcceptFuncs[accept]; ok {
			return f
		}
	}
	return o.DefaultFunc
}

func Hello() func(w http.ResponseWriter, r *http.Request) {
	l := log.Logger.With(slog.String("handler", "hello_world"))
	l.Info("hello world")

	helloService := hello.NewHelloService()
	options := Options{
		AcceptFuncsOpts: AcceptFuncsOpts{
			AcceptFuncs: map[string]AcceptFunc{
				"application/json": func(w http.ResponseWriter, r *http.Request) {
					message := helloService.Hello()
					w.Write([]byte("{\"data\": {\"message\": \"" + message + "\"}}"))
				},
				"text/html": func(w http.ResponseWriter, r *http.Request) {
					message := helloService.Hello()
					tmp, err := template.New("hello_world").Parse("<html><body><div>{{.}}</div></body></html>")
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					tmp.Execute(w, message)
				},
			},
			DefaultFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotAcceptable)
				w.Write([]byte("{\"error\":\"Invalid Accept Header\"}"))
			},
		},
	}
	return func(w http.ResponseWriter, r *http.Request) {
		accepts := r.Header["Accept"]
		options.GetAcceptFunc(accepts)(w, r)
	}
}

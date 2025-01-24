package handlers

import (
	"log/slog"
	"net/http"

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
	options := Options{
		AcceptFuncsOpts: AcceptFuncsOpts{
			AcceptFuncs: map[string]AcceptFunc{
				"application/json": func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("{\"data\": {\"message\": \"helloworld\"}}"))
				},
			},
			DefaultFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("{\"error\":\"Invalid Accept Header\"}"))
				w.WriteHeader(http.StatusNotAcceptable)
			},
		},
	}
	return func(w http.ResponseWriter, r *http.Request) {
		accepts := r.Header["Accept"]
		options.GetAcceptFunc(accepts)(w, r)
	}
}

package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"

	"golang-server/src/log"
)

func Hello() func(w http.ResponseWriter, r *http.Request) {
	l := log.Logger.With(slog.String("handler", "hello_world"))
	l.Info("hello world")

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l.Info("form::parse::before", slog.Any("Form", r.Form))
		r.ParseForm()
		l.Info("form::parse::after", slog.Any("Form", r.Form))

		some(ctx, r.Form)

		w.Write([]byte("Hello World!"))
	}
}

func some(ctx context.Context, form url.Values) {

}

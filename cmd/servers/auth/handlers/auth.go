package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"golang-server/src/domain/auth"
	"golang-server/src/log"
)

type UserNamePassword struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
}

func AuthByUsername(authService auth.AuthService) func(w http.ResponseWriter, r *http.Request) {
	l := log.Logger.With(slog.String("handler", "hello_world"))
	l.Info("hello world")

	options := Options{
		AcceptFuncsOpts: AcceptFuncsOpts{
			AcceptFuncs: map[string]AcceptFunc{
				"application/json": func(w http.ResponseWriter, r *http.Request) {
					form := &UserNamePassword{}
					json.NewDecoder(r.Body).Decode(form)

					var err error
					if form.Username != nil || form.Password != nil {
						err = fmt.Errorf("insufficent username or password")
					} else if err = authService.ByPasswordAndUsername(*form.Username, *form.Password); err != nil {
						err = fmt.Errorf("username and password invalid")
					}
					if err != nil {
						w.WriteHeader(http.StatusUnauthorized)
					}
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

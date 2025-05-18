package response_types

import (
	"encoding/json"
	"net/http"
)

type Response[T any] struct {
	Data  T      `json:"data"`
	Error string `json:"error"`
}

func Error[T any](w http.ResponseWriter, httpCode int, err error, body T) {
	w.WriteHeader(httpCode)
	var r Response[T]

	r.Data = body
	if err != nil {
		r.Error = err.Error()
	}

	b, _ := json.Marshal(r)
	w.Write(b)
}

func ErrorNoBody(w http.ResponseWriter, httpCode int, err error) {
	Error[any](w, httpCode, err, nil)
}

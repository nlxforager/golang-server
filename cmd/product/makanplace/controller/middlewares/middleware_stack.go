package middlewares

import (
	"net/http"
)

type Middleware = func(handler http.Handler) http.Handler
type MiddewareStack struct {
	h []Middleware
}

func (ms MiddewareStack) Wrap(next func(handler http.Handler) http.Handler) MiddewareStack {
	ms.h = append(ms.h, next)
	return ms
}

func (ms MiddewareStack) Finalize(final http.Handler) http.Handler {
	for _, m := range ms.h {
		final = m(final)
	}
	return final
}

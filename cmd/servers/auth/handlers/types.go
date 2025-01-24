package handlers

import "net/http"

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

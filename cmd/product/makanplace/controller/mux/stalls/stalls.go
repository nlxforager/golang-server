package stalls

import (
	"golang-server/cmd/product/makanplace/controller/middlewares"
	log2 "golang-server/cmd/product/makanplace/log"
	"golang-server/cmd/product/makanplace/service/mkusersessionservice"
	"io"
	"log"
	"net/http"
)

type Response struct {
	LoginUrls map[string]string `json:"login_urls"`

	UserInfo *mkusersessionservice.UserInfo `json:"user_info"`
}

func Register(mux *http.ServeMux, makanTokenCookieKey string, mkService *mkusersessionservice.Service, middlewares middlewares.MiddewareStack) {
	mux.Handle("POST /stall/", middlewares.Finalize(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s [POST /stall]\n", log2.HttpRequestPrefix(r.Context()))
		cookie, _ := r.Cookie(makanTokenCookieKey)
		session := mkService.GetSession(cookie.Value, false)
		if session == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		log.Printf("session found %#v\n", session)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		log.Println(string(body))
		w.WriteHeader(http.StatusNotImplemented)
	})))

	mux.Handle("GET /stall/", middlewares.Finalize(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("stalls getting")
		cookie, _ := r.Cookie(makanTokenCookieKey)
		session := mkService.GetSession(cookie.Value, false)
		if session == nil {
			w.WriteHeader(http.StatusOK)
			return
		}

		log.Printf("session found %#v\n", session)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		log.Println(string(body))
		w.WriteHeader(http.StatusOK)
	})))
}

package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang-server/cmd/product/makanplace/config"
	goauthmux "golang-server/cmd/product/makanplace/mux/oauth_google"
	goauthservice "golang-server/cmd/product/makanplace/service/oauth/google"
)

var Config config.Config

const COOKIE_NAME_GOOGLE_AUTHED_BEFORE = "loginoncebefore"

func main() {
	if err := Init(); err != nil {
		log.Fatal(err)
	}
	eCh := make(chan os.Signal, 1)
	signal.Notify(eCh, syscall.SIGTERM, syscall.SIGINT)

	goauthService := goauthservice.NewService(Config.GoogleAuthConfig)
	mux := http.NewServeMux()
	goauthmux.Register(mux, &goauthService)

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		cookie, _ := r.Cookie(COOKIE_NAME_GOOGLE_AUTHED_BEFORE)
		var cookieVal string
		if cookie != nil {
			cookieVal = cookie.Value
		}
		w.Write([]byte("pong. the client browser has cookie. " + COOKIE_NAME_GOOGLE_AUTHED_BEFORE + "=" + cookieVal))
	})

	go func() {
		log.Println("Listening on " + Config.ServerConfig.Port)
		http.ListenAndServe(Config.ServerConfig.Port, mux)
	}()
	recvSig := <-eCh
	log.Println("Received signal: " + recvSig.String() + " ; exiting...")
}

func Init() (err error) {
	Config, err = config.InitConfig()
	if err != nil {
		return
	}
	return
}

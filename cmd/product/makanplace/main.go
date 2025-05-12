package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang-server/cmd/product/makanplace/config"
	"golang-server/cmd/product/makanplace/service/mkusersessionservice"

	goauthmux "golang-server/cmd/product/makanplace/mux/oauth_google"
	goauthservice "golang-server/cmd/product/makanplace/service/oauth/google"
)

var Config config.Config

func makanTokenCookieKey() string { return "makantoken" }

func main() {
	if err := Init(); err != nil {
		log.Fatal(err)
	}
	eCh := make(chan os.Signal, 1)
	signal.Notify(eCh, syscall.SIGTERM, syscall.SIGINT)

	goauthService := goauthservice.NewService(Config.GoogleAuthConfig)
	mux := http.NewServeMux()

	mkUserSessionService := mkusersessionservice.New()
	goauthmux.Register(mux, makanTokenCookieKey(), &goauthService, mkUserSessionService)

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		cookie, _ := r.Cookie(makanTokenCookieKey())
		var sessionId string
		if cookie != nil {
			sessionId = cookie.Value
		}

		session := mkUserSessionService.GetSession(sessionId)

		sessionB, _ := json.Marshal(session)

		w.Write([]byte("pong. the client browser has cookie. " + makanTokenCookieKey() + "=" + string(sessionB)))
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

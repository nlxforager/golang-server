package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang-server/cmd/product/makanplace/config"
	"golang-server/cmd/product/makanplace/controller/middlewares"
	mkauthmux "golang-server/cmd/product/makanplace/controller/mux/auth"
	goauthmux "golang-server/cmd/product/makanplace/controller/mux/oauth_google"
	"golang-server/cmd/product/makanplace/controller/mux/ping"
	"golang-server/cmd/product/makanplace/service/mkusersessionservice"
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
	makanTokenCookieKey := makanTokenCookieKey()
	goauthloginurl := "/auth/google/login"

	defaultMiddlewares := middlewares.MiddewareStack{}.Wrap(middlewares.WithCORS)

	goauthmux.Register(mux, makanTokenCookieKey, &goauthService, mkUserSessionService, goauthloginurl)
	ping.Register(mux, makanTokenCookieKey, mkUserSessionService, goauthloginurl, defaultMiddlewares)
	mkauthmux.Register(mux, makanTokenCookieKey, mkUserSessionService, goauthloginurl, defaultMiddlewares)

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

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang-server/cmd/product/makanplace/config"
	"golang-server/cmd/product/makanplace/controller/middlewares"

	authrepo "golang-server/cmd/product/makanplace/repositories/auth"

	mkauthmux "golang-server/cmd/product/makanplace/controller/mux/auth"
	goauthmux "golang-server/cmd/product/makanplace/controller/mux/oauth_google"
	"golang-server/cmd/product/makanplace/controller/mux/ping"
	"golang-server/cmd/product/makanplace/service/mkusersessionservice"
	goauthservice "golang-server/cmd/product/makanplace/service/oauth/google"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Config config.Config
var DbConnPool *pgxpool.Pool

func makanTokenCookieKey() string { return "makantoken" }

func main() {
	if err := Init(); err != nil {
		log.Fatal(err)
	}
	eCh := make(chan os.Signal, 1)
	signal.Notify(eCh, syscall.SIGTERM, syscall.SIGINT)

	goauthService := goauthservice.NewService(Config.GoogleAuthConfig)
	mux := http.NewServeMux()
	mkAuthRepository := authrepo.New(DbConnPool)
	mkUserSessionService := mkusersessionservice.New(mkAuthRepository)
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

	// Db Connection Pool

	config, err := pgxpool.ParseConfig(Config.ConnString)
	if err != nil {
		return err
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnIdleTime = 5 * time.Minute
	ctx := context.Background()

	// Create the pool
	DbConnPool, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return err
	}

	err = DbConnPool.Ping(ctx)
	if err != nil {
		panic(err)
	}
	return nil
}

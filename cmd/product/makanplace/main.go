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
	mklog "golang-server/cmd/product/makanplace/log"

	authrepo "golang-server/cmd/product/makanplace/repositories/auth"
	outletrepo "golang-server/cmd/product/makanplace/repositories/outlet"

	goauthmux "golang-server/cmd/product/makanplace/controller/mux/oauth_google"
	outletmux "golang-server/cmd/product/makanplace/controller/mux/outlets"
	"golang-server/cmd/product/makanplace/controller/mux/ping"
	mksessionmux "golang-server/cmd/product/makanplace/controller/mux/session"
	
	"golang-server/cmd/product/makanplace/service/mk_user_session"
	goauthservice "golang-server/cmd/product/makanplace/service/oauth/google"
	"golang-server/cmd/product/makanplace/service/mk_outlet"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/cors"
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

	// services
	goauthService := goauthservice.NewService(Config.GoogleAuthConfig)
	mux := http.NewServeMux()
	mkAuthRepository := authrepo.New(DbConnPool)
	mkUserSessionService := mk_user_session.New(mkAuthRepository, Config.AdminConfig)
	makanTokenCookieKey := makanTokenCookieKey()

	// controller
	goauthloginurl := "/auth/google/login"

	defaultMiddlewares := middlewares.MiddewareStack{}

	//
	authMiddleware := defaultMiddlewares.Wrap(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("%s %s middleware:auth %s \n", r.Method, r.URL.Path, mklog.HttpRequestPrefix(r.Context()))
			cookie, _ := r.Cookie(makanTokenCookieKey)
			session := mkUserSessionService.GetSession(cookie.Value, true)
			if session == nil {
				log.Printf("%s %s middleware:auth %s session not found", r.Method, r.URL.Path, mklog.HttpRequestPrefix(r.Context()))
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			log.Printf("session found %#v\n", session)
			handler.ServeHTTP(w, r)
		})
	})

	goauthmux.Register(mux, makanTokenCookieKey, &goauthService, mkUserSessionService, goauthloginurl)
	ping.Register(mux, makanTokenCookieKey, mkUserSessionService, goauthloginurl, defaultMiddlewares)
	mksessionmux.Register(mux, makanTokenCookieKey, mkUserSessionService, goauthloginurl, defaultMiddlewares)

	outletRepo := outletrepo.New(DbConnPool)
	outletService := mk_outlet_service.NewOutletService(outletRepo)

	outletmux.Register(mux, makanTokenCookieKey, mkUserSessionService, authMiddleware, outletService)

	go func() {
		log.Println("Listening on " + Config.ServerConfig.Port)
		http.ListenAndServe(Config.ServerConfig.Port, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "METHOD", r.Method)
			ctx = context.WithValue(ctx, "URL", r.URL.String())
			ctx = context.WithValue(ctx, "ORIGIN", r.Header.Get("Origin"))
			r = r.WithContext(ctx)
			log.Printf("%s [middleware 0]\n", mklog.HttpRequestPrefix(r.Context()))

			c := cors.New(cors.Options{
				AllowedOrigins:   []string{"http://localhost:5173"},
				AllowCredentials: true,
				AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
				AllowedHeaders:   []string{makanTokenCookieKey, "Accept", "Authorization", "Content-Type", "X-Requested-With"},
				// Enable Debugging for testing, consider disabling in production
				Debug: true,
			}).Handler(mux)

			c.ServeHTTP(w, r)
		}))
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

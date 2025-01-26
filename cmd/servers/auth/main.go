package main

import (
	"context"
	"database/sql"
	"fmt"
	"golang-server/src/domain/email"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang-server/cmd/servers/auth/mux"
	"golang-server/src/config"
	gctx "golang-server/src/context"
	"golang-server/src/domain/auth"
	"golang-server/src/log"

	"github.com/redis/go-redis/v9"

	_ "github.com/lib/pq"
)

func main() {
	// Setup
	err := config.Init()
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}
	ctx := context.Background()
	ctx = gctx.AppendCallStack(ctx, "main")

	err = log.Set(os.Getenv(config.CONFIG_LOGGER_TYPE))
	if err != nil {
		slog.Error(err.Error())
	}

	l := log.Logger.With(gctx.AsAttributes(ctx)...)
	l.LogAttrs(ctx, log.LevelSystem, "started")

	ctx, _ = context.WithTimeoutCause(ctx, 100*time.Second, fmt.Errorf("timedout_main"))
	interruptSignal := make(chan os.Signal, 1)
	signal.Notify(interruptSignal, syscall.SIGINT /*keyboard input*/, syscall.SIGTERM /*process kill*/)
	// HTTP Server
	{
		pgConfig, _ := config.GetPostGresConfig()       // allow server to run without db conn.
		otpEmailConfig, _ := config.GetOtpEmailConfig() // allow server to run without db conn.

		db, err := sql.Open("postgres", pgConfig.CONNECTION_STRING)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		redisClient := redis.NewClient(&redis.Options{})
		authService, err := auth.NewService(
			&auth.Repository{}, redisClient, "issuerXX", "notsosecret",
		)

		if err != nil {
			l.LogAttrs(ctx, log.LevelSystem, "failed to initialize auth service")
			panic(err)
		}

		mailService, err := email.NewSimpleClientService(otpEmailConfig.CONFIG_OTP_EMAIL, otpEmailConfig.CONFIG_OTP_PASSWORD)
		if err != nil {
			l.LogAttrs(ctx, log.LevelSystem, "failed to initialize mail service")
			panic(err)
		}
		mux := mux.NewMux(&mux.MuxOpts{
			AuthMuxOpts: &mux.AuthMuxOpts{
				Auth: authService,
				Mail: mailService,
			},
		})
		http.ListenAndServe("", mux)
	}

	// defer ns.Shutdown() // not called, the library will handle.

	// Teardown
	select {
	case <-interruptSignal:
		l.LogAttrs(ctx, log.LevelSystem, "interrupt or terminated")
	case <-ctx.Done():
		l.LogAttrs(ctx, log.LevelSystem, "ctx.Done() received", slog.String("error", context.Cause(ctx).Error()))
	}

	l.LogAttrs(ctx, log.LevelSystem, "exited")
}

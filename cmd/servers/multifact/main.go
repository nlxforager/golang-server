package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang-server/cmd/servers/multifact/mux"
	"golang-server/src/config"
	gctx "golang-server/src/context"
	"golang-server/src/domain/auth"
	"golang-server/src/domain/email"
	"golang-server/src/log"

	"github.com/redis/go-redis/v9"

	// DO NOT EDIT ORDER
	_ "github.com/lib/pq"
	// DO NOT EDIT
	"database/sql"
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

	interruptSignal := make(chan os.Signal, 1)
	signal.Notify(interruptSignal, syscall.SIGINT /*keyboard input*/, syscall.SIGTERM /*process kill*/)
	// HTTP Server
	{
		pgConfig, _ := config.GetPostGresConfig() // allow server to run without db conn.
		db, err := sql.Open("postgres", pgConfig.CONNECTION_STRING)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		redisClient := redis.NewClient(&redis.Options{})

		otpEmailConfig, _ := config.GetOtpEmailConfig() // allow server to run without email conn.
		mailService, err := email.NewSimpleClientService(otpEmailConfig.CONFIG_OTP_EMAIL, otpEmailConfig.CONFIG_OTP_PASSWORD)
		if err != nil {
			l.LogAttrs(ctx, log.LevelSystem, "failed to initialize mail service")
			panic(err)
		}

		authService, err := auth.NewService(
			&auth.Repository{}, redisClient, "issuerXX", "notsosecret",
		)
		if err != nil {
			l.LogAttrs(ctx, log.LevelSystem, "failed to initialize auth service")
			panic(err)
		}

		mux := mux.NewMux(&mux.MuxOpts{
			AuthMuxOpts: &mux.AuthMuxOpts{
				Auth: authService,
				Mail: mailService,
			},
		})
		go http.ListenAndServe("", mux)
	}

	// defer ns.Shutdown() // not called, the library will handle.

	// Teardown
teardown:
	for {
		select {
		case <-time.NewTicker(5 * time.Second).C:
			l.LogAttrs(ctx, log.LevelSystem, "tick")
		case <-interruptSignal:
			l.LogAttrs(ctx, log.LevelSystem, "interrupt or terminated")
			break teardown
		case <-ctx.Done():
			l.LogAttrs(ctx, log.LevelSystem, "ctx.Done() received", slog.String("error", context.Cause(ctx).Error()))
			break teardown
		}
	}

	l.LogAttrs(ctx, log.LevelSystem, "exited")
}

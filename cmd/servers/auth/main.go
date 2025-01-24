package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang-server/src/config"
	gctx "golang-server/src/context"
	"golang-server/src/log"
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
		mux := NewMux()
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

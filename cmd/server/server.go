package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang-server/src/config"
	gctx "golang-server/src/context"
	"golang-server/src/log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error(err.Error())
	}
	ctx := context.Background()
	ctx = gctx.AppendCallStack(ctx, "main")

	err = log.Set(os.Getenv(config.CONFIG_LOGGER_TYPE))
	if err != nil {
		slog.Error(err.Error())
	}

	l := log.Logger.With(gctx.AsAttributes(ctx)...)
	l.LogAttrs(ctx, log.LevelSystem, "started")

	ctx, _ = context.WithTimeoutCause(ctx, 10*time.Second, err)

	interruptSignal := make(chan os.Signal, 1)
	signal.Notify(interruptSignal, syscall.SIGINT /*keyboard input*/, syscall.SIGTERM /*process kill*/)
	select {
	//case <-interruptSignal:
	//	l.LogAttrs(ctx, log.LevelSystem, "interrupt signal received")
	case <-interruptSignal:
		l.LogAttrs(ctx, log.LevelSystem, "interrupt or terminated")
	case <-ctx.Done():
		l.LogAttrs(ctx, log.LevelSystem, "ctx.Done() received")
	}

	l.LogAttrs(ctx, log.LevelSystem, "exited")
}

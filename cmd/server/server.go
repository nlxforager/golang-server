package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang-server/src/config"
	gctx "golang-server/src/context"
	"golang-server/src/log"
	"golang-server/src/natsinfra"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
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

	{ // Nats
		natsConfig, err := config.GetNatsConfig()
		if err != nil {
			panic(err)
		}
		{
			if natsConfig.Embedded {
				l.LogAttrs(ctx, log.LevelSystem, "NATS initializing embedded server")
				ns, err := server.NewServer(&server.Options{})
				if err != nil {
					panic(err)
				}
				go ns.Start()
				if !ns.ReadyForConnections(4 * time.Second) {
					panic("NATS embedded server not ready for connection")
				} else {
					l.LogAttrs(ctx, log.LevelSystem, "NATS embedded server ready")
				}
				natsConfig.Url = ns.ClientURL()
			}
			l.LogAttrs(ctx, log.LevelSystem, "NATS trying to client", slog.String("client url", natsConfig.Url))
			nc, err := nats.Connect(natsConfig.Url)
			if err != nil {
				panic(err)
			} else {
				l.LogAttrs(ctx, log.LevelSystem, "ok", slog.String("client url", natsConfig.Url))
			}

			nc.JetStream()
			natsinfra.Smoke(nc)
		}
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

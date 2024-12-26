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

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

var Server *server.Server
var JetStream jetstream.JetStream

func init() {
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
	l.LogAttrs(ctx, log.LevelSystem, "[nats_demo] init")

}

func main() {

	ctx := context.Background()
	ctx = gctx.AppendCallStack(ctx, "main")

	l := log.Logger.With(gctx.AsAttributes(ctx)...)
	l.LogAttrs(ctx, log.LevelSystem, "[nats_demo] started")

	var natsUrl string
	natsUrl = nats.DefaultURL
	l.LogAttrs(ctx, log.LevelSystem, "NATS trying to client", slog.String("client url", natsUrl))
	nc, err := nats.Connect(natsUrl)
	if err != nil {
		panic(err)
	} else {
		l.LogAttrs(ctx, slog.LevelInfo, "ok", slog.String("client url", natsUrl))
	}

	JetStream, err = jetstream.New(nc)
	if err != nil {
		panic(err)
	} else {
		l.LogAttrs(ctx, slog.LevelInfo, "ok JetStream", slog.String("client url", natsUrl))
	}
	ctx, _ = context.WithTimeoutCause(ctx, 299*time.Second, fmt.Errorf("timedout_main"))

	stream1Name := "stream-1"
	if err != nil {
		panic(err)
	}
	subject1 := "subject_a.1"

	stream1, err := JetStream.CreateStream(ctx, jetstream.StreamConfig{
		Name:     stream1Name,
		Subjects: []string{"ORDERS.*", subject1},
	})
	if err != nil {
		panic(err)
	}

	_, err = JetStream.PublishAsync(subject1, []byte("1_before_any_subscription11"))
	if err != nil {
		panic(err)
	}

	l.LogAttrs(ctx, slog.LevelInfo, "init consumers...")
	for i := range []int64{1, 2, 3, 4} {
		go func() {
			consumerName := fmt.Sprintf("consumer-%d", i)
			l.LogAttrs(ctx, slog.LevelInfo, "init consumer...", slog.String("consumerName", consumerName))
			//ctx, _ := context.WithDeadline(ctx, time.Now().Add(1*time.Second))
			_, err = JetStream.CreateConsumer(ctx, stream1Name, jetstream.ConsumerConfig{Name: consumerName})
			if err != nil {
				panic(err)
			}
			stream1.CachedInfo()

			l.LogAttrs(ctx, slog.LevelInfo, "getting consumer...", slog.String("consumerName", consumerName))
			// get consumer handle
			cons, err := JetStream.Consumer(ctx, stream1Name, consumerName)
			if err != nil {
				panic(err)
			}

			msgs, err := cons.Fetch(10)
			if err != nil {
				panic(err)
			}
			for msg := range msgs.Messages() {
				l.LogAttrs(ctx, slog.LevelInfo, "received message", slog.String("consumerName", consumerName), slog.String("data", string(msg.Data())))
				msg.Ack()
			}
		}()
	}
	_, err = JetStream.PublishAsync(subject1, []byte("2_after init consumers"))
	if err != nil {
		panic(err)
	}

	interruptSignal := make(chan os.Signal, 1)
	signal.Notify(interruptSignal, syscall.SIGINT /*keyboard input*/, syscall.SIGTERM /*process kill*/)
	if err != nil {
		panic(err)
	}
	select {
	//case <-interruptSignal:
	//	l.LogAttrs(ctx, log.LevelSystem, "interrupt signal received")
	case <-interruptSignal:
		l.LogAttrs(ctx, log.LevelSystem, "interrupt or terminated")
	case <-ctx.Done():
		l.LogAttrs(ctx, log.LevelSystem, "ctx.Done() received", slog.String("error", context.Cause(ctx).Error()))
	}

	l.LogAttrs(ctx, log.LevelSystem, "exited")
}

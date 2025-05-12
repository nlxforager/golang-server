package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"golang-server/src/config"
	gctx "golang-server/src/context"
	"golang-server/src/log"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

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

type JsonPayload struct {
	Some string    `json:"some"`
	Id   int64     `json:"id"`
	Time time.Time `json:"time"`
}

// main
// Prerequisite:
// Run nats server
// $ docker run -p 4222:4222 -ti nats:latest -js -m 8222
func main() {
	ctx := context.Background()
	ctx = gctx.AppendCallStack(ctx, "main")

	l := log.Logger.With(gctx.AsAttributes(ctx)...)
	l.LogAttrs(ctx, log.LevelSystem, "[nats_demo] started")

	natsUrl := nats.DefaultURL
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

	stream1Name := "stream-1"
	stream2Name := "stream-2"
	if err != nil {
		panic(err)
	}
	subject1 := "subject_a.1"
	subject1_Json := "subject_a_json.1"
	subject_QueueGroup := "subject_qgrp.1"

	l.LogAttrs(ctx, slog.LevelInfo, "handling kv store...")
	kv, err := JetStream.CreateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket:      "bucket-1",
		Description: "",
	})
	if err != nil {
		panic(err)
	}

	kv.Create(ctx, "key1", []byte("value1"))
	value, _ := kv.Get(ctx, "key1")

	l.LogAttrs(ctx, slog.LevelInfo, "kv.Get -> ", slog.Any("value.Value()", string(value.Value())))

	l.LogAttrs(ctx, slog.LevelInfo, "create stream...")
	_, err = JetStream.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     stream1Name,
		Subjects: []string{subject1, subject1_Json},
	})

	if err != nil {
		panic(err)
	}

	l.LogAttrs(ctx, slog.LevelInfo, "publishing some data...")
	_, err = JetStream.PublishAsync(subject1, []byte("1_before_any_subscription11"))
	if err != nil {
		panic(err)
	}
	b, _ := json.Marshal(JsonPayload{
		Some: "SomeValue",
	})

	_, err = JetStream.PublishAsync(subject1_Json, b)
	if err != nil {
		panic(err)
	}

	//go ConsumerWithIndividualStreams(l, ctx, err, stream1Name, subject1, subject1_Json)

	{
		l.LogAttrs(ctx, slog.LevelInfo, "[deliver_once] init stream...")
		//JetStream.DeleteStream(ctx, stream2Name)
		_, err = JetStream.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
			Name:      stream2Name,
			Subjects:  []string{subject_QueueGroup},
			Retention: jetstream.LimitsPolicy,
		})
		if err != nil {
			panic(err)
		}

		l.LogAttrs(ctx, slog.LevelInfo, "[deliver_once] init consumer copies in a deliver_once...")

		var globalCount = atomic.Int64{}
		for i, v := range []int64{1, 1} {
			i := i
			consumerName := fmt.Sprintf("qg_consumer-%d", v)
			l.LogAttrs(ctx, slog.LevelInfo, "[deliver_once] init consumers...", slog.String("consumerName", consumerName), slog.Int("pos", i))
			cons, err := JetStream.CreateOrUpdateConsumer(ctx, stream2Name, jetstream.ConsumerConfig{
				Name:           consumerName,
				Durable:        consumerName,
				MaxAckPending:  500,
				FilterSubjects: []string{subject_QueueGroup},
			})
			if err != nil {
				panic(err)
			}
			// Display relevant sequence information
			info, _ := cons.Info(ctx)
			l.LogAttrs(ctx, slog.LevelInfo, "[deliver_once] init consumers...", slog.String("consumerName", consumerName), slog.Int("pos", i), slog.Uint64("Last Delivered Sequence", info.Delivered.Stream),
				slog.Uint64("Last Acknowledged Sequence: %d", info.AckFloor.Stream),
				slog.Uint64("Pending Messages: %d", info.NumPending))

			go func() {
				<-time.Tick(5 * time.Second)
				var count atomic.Int64
				var maxPayload JsonPayload

				cons.Consume(func(msg jetstream.Msg) {
					defer msg.Ack()

					count.Add(1)
					globalCount.Add(1)

					var payload JsonPayload
					json.Unmarshal(msg.Data(), &payload)
					if maxPayload.Id < payload.Id {
						maxPayload = payload
					}
					meta, _ := msg.Metadata()

					l.LogAttrs(ctx, slog.LevelInfo, "[deliver_once] received messages (.Consume)", slog.String("consumerName", consumerName), slog.Int("pos", i), slog.Any("maxPayload", maxPayload),
						slog.Int64("count", count.Load()), slog.Uint64("meta.NumDelivered", meta.NumDelivered), slog.Uint64("meta.NumPending", meta.NumPending), slog.Int64("globalCount", globalCount.Load()))
				})

			}()

			go func() {
				<-time.Tick(5 * time.Second)
				var count atomic.Int64
				for range time.Tick(time.Second * 1) {
					batch, err := cons.FetchNoWait(200)
					if err != nil {
						panic(err)
					}

					var maxPayload JsonPayload
					for msg := range batch.Messages() {
						count.Add(1)
						globalCount.Add(1)

						var payload JsonPayload
						json.Unmarshal(msg.Data(), &payload)
						if maxPayload.Id < payload.Id {
							maxPayload = payload
						}
						//l.LogAttrs(ctx, slog.LevelInfo, "[deliver_once] received message", slog.String("consumerName", consumerName), slog.Int("pos", i), slog.String("subject", subject), slog.Any("data", data))
						msg.Ack()
					}
					l.LogAttrs(ctx, slog.LevelInfo, "[deliver_once] received messages (.FetchNoWait)", slog.String("consumerName", consumerName), slog.Int("pos", i), slog.Any("maxPayload", maxPayload),
						slog.Int64("count", count.Load()), slog.Int64("globalCount", globalCount.Load()))
				}
			}()

			go func() {
				for range time.Tick(1 * time.Second) {
					info, _ := cons.Info(ctx)
					l.LogAttrs(ctx, slog.LevelInfo, "[deliver_once] consumer info", slog.String("consumerName", consumerName), slog.Int("pos", i), slog.Uint64("Last Delivered Sequence", info.Delivered.Stream),
						slog.Uint64("Last Acknowledged Sequence: %d", info.AckFloor.Stream),
						slog.Uint64("Pending Messages: %d", info.NumPending))
					if info.NumAckPending >= info.Config.MaxAckPending {
						l.LogAttrs(ctx, slog.LevelWarn, "[deliver_once] MaxAckPending reached limit", slog.String("consumerName", consumerName), slog.Int("pos", i), slog.Uint64("Last Delivered Sequence", info.Delivered.Stream),
							slog.Uint64("Last Acknowledged Sequence: %d", info.AckFloor.Stream),
							slog.Uint64("Pending Messages: %d", info.NumPending))
						// Take action, e.g., send an error or log
					}
				}
			}()
		}
	}
	l.LogAttrs(ctx, slog.LevelInfo, "pulse")
	_, err = JetStream.PublishAsync(subject_QueueGroup, []byte("2_after init consumers"))
	if err != nil {
		panic(err)
	}

	go func() {
		var c atomic.Int64
		t := time.Now()
		for range time.NewTicker(time.Second).C {
			l.LogAttrs(ctx, slog.LevelInfo, "[deliver_once] pulsing messages", slog.Duration("since", time.Since(t)), slog.Int64("count", c.Load()))

			for i := 0; i < 1000; i++ {
				c.Add(1)
				payload := JsonPayload{
					Some: "queue??",
					Id:   c.Load(),
					Time: time.Now(),
				}
				payloadB, _ := json.Marshal(payload)
				ack, _ := JetStream.PublishAsync(subject_QueueGroup, payloadB)
				go func() {
					return
					pubAck := <-ack.Ok()

					//<-ack.Ok()
					l.LogAttrs(ctx, slog.LevelInfo, "[deliver_once] pulsing messages::ack", slog.Uint64("seq", pubAck.Sequence), slog.Int64("payload.Id", payload.Id))
				}()
			}
		}
	}()

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

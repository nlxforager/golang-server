package natsinfra

import (
	"fmt"
	
	"github.com/nats-io/nats.go"
)

func Smoke(nc *nats.Conn) {
	subject := "my-subject"

	nc.Subscribe(subject, func(msg *nats.Msg) {
		data := string(msg.Data)
		fmt.Println(data)
	})

	nc.Publish(subject, []byte("Smoke(nc *nats.Conn)"))
}

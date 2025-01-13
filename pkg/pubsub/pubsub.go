package pubsub

import (
	"fmt"
	"log"

	"github.com/bagaswibowo25/golang-search/pkg/config"
	"github.com/nats-io/nats.go"
)

type JetStream struct {
	StreamCtx nats.JetStreamContext
	Conn      *nats.Conn
	Conf      config.NatsConfig
}

func PublishMessage(message string, jServer *JetStream) {
	_, err := jServer.StreamCtx.Publish(jServer.Conf.Subject, []byte(message))
	if err != nil {
		log.Printf("Error publishing message: %v", err)
	} else {
		fmt.Printf("[%s] Published message: %s\n", jServer.Conf.ConsumerId, message)
	}
}

func SubscribeMessage(cb nats.MsgHandler, jServer *JetStream) {
	_, err := jServer.StreamCtx.Subscribe(jServer.Conf.Subject, cb, nats.Durable(jServer.Conf.ConsumerId), nats.ManualAck(), nats.SkipConsumerLookup())

	if err != nil {
		log.Fatalf("Error subscribing to subject: %v", err)
		return
	}
}

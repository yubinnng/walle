package event

import (
	"log"

	"github.com/nats-io/nats.go"
)

var client *nats.Conn

func ConnectNATS() {
	// Connect to a server
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal("cannot connect to NATS")
	}
	client = nc
	log.Println("connected to NATS")
}

func Publish(topic string, data []byte) {
	client.Publish(topic, data)
}

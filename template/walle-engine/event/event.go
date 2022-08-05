package event

import (
	"log"

	"github.com/nats-io/nats.go"
)

var client *nats.Conn

const NATS_URL = "nats://nats.walle:4222"

func ConnectNATS() {
	// Connect to a server
	nc, err := nats.Connect(NATS_URL)
	if err != nil {
		log.Fatal("Failed to connect to NATS, URL=" + NATS_URL)
	}
	client = nc
	log.Println("Connected to NATS, URL=" + NATS_URL)
}

func Publish(topic string, data []byte) {
	client.Publish(topic, data)
}

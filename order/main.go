package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {

	var (
		nc  *nats.Conn
		err error
	)

	for retries := 1; retries <= 5; retries++ {
		nc, err = nats.Connect("nats://nats:4222", nats.MaxReconnects(-1), nats.ReconnectWait(time.Second*2))
		if err == nil {
			break
		}
		log.Printf("Trying to connect to NATS, attemp %d...\n", retries)
		time.Sleep(time.Second * time.Duration(retries))
	}

	if err != nil {
		log.Fatalf("Cannot connect to NATS server: %v", err)
	}

	defer nc.Close()

	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, err = js.AddStream(&nats.StreamConfig{
		Name:              "ORDERS",
		Subjects:          []string{"order.created"},
		Storage:           nats.FileStorage,
		MaxMsgsPerSubject: 100_000_000,
		MaxMsgSize:        8 << 20,
		NoAck:             false,
		MaxAge:            7 * 24 * time.Hour,
	}, nats.Context(ctx))

	if err != nil {
		log.Fatalf("Error on creating JetStream stream: %v", err)
	}

	ticker := time.NewTicker(1 * time.Second)

	for range ticker.C {
		orderId := fmt.Sprintf("order-%d", time.Now().Unix())
		fmt.Printf("Sending shipping with JetStream: %s\n", orderId)
		_, err := js.Publish("order.created", []byte(orderId))
		if err != nil {
			log.Printf("Error al enviar pedido: %v", err)

		}
	}
}

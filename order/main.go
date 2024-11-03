package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	ticker := time.NewTicker(5 * time.Second)

	for range ticker.C {
		orderId := fmt.Sprintf("order-%d", time.Now().Unix())
		fmt.Printf("Enviando pedido: %s\n", orderId)
		err := nc.Publish("order.created", []byte(orderId))
		if err != nil {
			log.Printf("Error al enviar pedido: %v", err)

		}
	}
}

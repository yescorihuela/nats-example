package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect("nats://nats:4222", nats.MaxReconnects(-1), nats.ReconnectWait(time.Second*2))

	if err != nil {
		log.Fatal(err)
	}

	defer nc.Close()

	_, err = nc.Subscribe("inventory.checked", func(m *nats.Msg) {
		orderId := string(m.Data)
		fmt.Printf("Enviando notificación para el pedido: %s\n", orderId)
	})

	if err != nil {
		log.Fatalf("error on subscribing inventory.checked: %v", err)
	}

	select {}

}

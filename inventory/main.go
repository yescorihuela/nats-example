package main

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

func main() {

	nc, err := nats.Connect("nats://nats:4222")
	if err != nil {
		log.Fatal(err)
	}

	defer nc.Close()

	_, err = nc.Subscribe("order.created", func(m *nats.Msg) {
		orderId := string(m.Data)
		fmt.Printf("Procesado pedido en inventario: %s\n", orderId)
		nc.Publish("inventory.checked", []byte(orderId))

	})

	if err != nil {
		log.Fatalf("Error al suscribirse a order.created: %v", err)
	}

	select {}

}

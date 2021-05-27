package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"os"
)

var rabbit_host = os.Getenv("RABBIT_HOST")
var rabbit_port = os.Getenv("RABBIT_PORT")
var rabbit_user = os.Getenv("RABBIT_USERNAME")
var rabbit_password = os.Getenv("RABBIT_PASSWORD")

func main() {
	consume()
}

func consume() {

	conn, err := amqp.Dial("amqp://" + rabbit_user + ":" + rabbit_password + "@" + rabbit_host + ":" + rabbit_port + "/" )

	if err != nil {
		log.Fatalf("%s: %s\n", "Failed to connect to RabbitMQ\n", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("%s: %s\n", "Failed to open a channel\n", err)
	}

	q, err := ch.QueueDeclare(
		"publisher",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatalf("%s:%s\n", "Failed to declear a queue", err)
	}

	fmt.Println("Channel and the queue has been established")

	defer conn.Close()
	defer ch.Close()

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
		)



	if err != nil {
		log.Fatalf("%s:%s\n", "Failed to register consumer", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message %s\n", d.Body)
			d.Ack(false)
		}
	}()

	fmt.Println("Running...")

	<-forever
}
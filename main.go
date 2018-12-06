package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	go server()
	go client()

	var a string
	fmt.Scanln(&a)
}

func server() {
	conn, ch, q := getQueue()

	defer conn.Close()
	defer ch.Close()

	msg := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte("Hello world"),
	}
	for {
		ch.Publish("", q.Name, false, false, msg)
	}

}

func getQueue() (*amqp.Connection, *amqp.Channel, *amqp.Queue) {
	conn, err := amqp.Dial("amqp://guest@localhost:5672")
	failOnError(err, "Failed to connect to RabbitMQ")
	ch, err := conn.Channel()
	failOnError(err, "Failed to establish a channel")
	q, err := ch.QueueDeclare("hello",
		false,
		false,
		false,
		false,
		nil)
	failOnError(err, "Failed to declare a queue")

	return conn, ch, &q
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func client() {
	conn, ch, q := getQueue()
	defer conn.Close()
	defer ch.Close()

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil)

	failOnError(err, "Failed to establish client connection to the queue")
	for msg := range msgs {
		log.Printf("Received message with: %s", msg.Body)
	}
}

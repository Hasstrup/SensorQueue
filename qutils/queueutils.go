package qutils

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

const SensorListQueue = "SensorList"
const SensorDiscoveryExchange = "SensorDiscovery"
const PersistReadingsQueue = "PersistReadingsQueue"

func GetChannel(url string) (*amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial(url)
	failOnError(err, "Failed to establish connection")
	ch, err := conn.Channel()
	failOnError(err, "Something went wrong creating the channel")
	return conn, ch
}

func GetQueue(name string, ch *amqp.Channel, autoDelete bool) *amqp.Queue {
	q, err := ch.QueueDeclare(
		name,
		false,
		autoDelete,
		false,
		false,
		nil)
	failOnError(err, "Something went wrong fetching the queue")
	return &q
}
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s:", msg, err)
		panic(fmt.Sprintf("%s: %s:", msg, err))
	}
}

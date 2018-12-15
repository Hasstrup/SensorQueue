package main

import (
	"bytes"
	"encoding/gob"
	"log"

	"github.com/hasstrup/queue/datamanager"
	"github.com/hasstrup/queue/dto"
	"github.com/hasstrup/queue/qutils"
)

const url = "amqp://guest:guest@localhost:5672"

func main() {
	conn, ch := qutils.GetChannel(url)
	defer conn.Close()
	defer ch.Close()

	msgs, err := ch.Consume(
		qutils.PersistReadingsQueue,
		"",
		false,
		true,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatalf("Failed to get messages")
	}

	for msg := range msgs {
		buf := bytes.NewReader(msg.Body)
		dec := gob.NewDecoder(buf)
		sd := &dto.SensorMessage{}
		dec.Decode(sd)

		err := datamanager.SaveReading(sd)

		if err != nil {
			log.Printf("Failed to save reading from sensor %v", sd.Name)
		} else {
			// inform the queue to remove this message
			msg.Ack(false)
		}
	}
}

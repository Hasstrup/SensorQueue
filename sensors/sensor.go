package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/hasstrup/queue/dto"
	"github.com/hasstrup/queue/qutils"
	"github.com/streadway/amqp"
)

var url = "amqp:guest@localhost:5672"

var name = flag.String("name", "sensor", "name of sensor")
var freq = flag.Uint("freq", 5, "update frequency in cycles")
var max = flag.Float64("max", 5., "maximum value for the generated sensor")
var min = flag.Float64("min", 1., "minimum value for the sensor")
var stepSize = flag.Float64("step", 0.1, "maximum allowable change per measurement")

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

var value = r.Float64()*(*max-*min) + *min
var nom = (*min-*max)/2 + *min

func main() {
	flag.Parse()

	conn, ch := qutils.GetChannel(url)

	defer conn.Close()
	defer ch.Close()

	dataQueue := qutils.GetQueue(*name, ch)
	// the queue being created here holds tha value for
	// all the newly created queues
	msg := amqp.Publishing{Body: []byte(*name)}
	ch.Publish("amq.fanout",
		"", //doing this creates a fanout exchange, so that it sends the message to all connected queues
		false,
		false,
		msg)

	dur, _ := time.ParseDuration(strconv.Itoa(1000/int(*freq)) + "ms")
	signal := time.Tick(dur)
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)

	for range signal {
		calcValue()

		reading := dto.SensorMessage{
			Name:      *name,
			Value:     value,
			Timestamp: time.Now(),
		}

		buf.Reset()
		enc = gob.NewEncoder(buf)
		enc.Encode(reading)

		msg := amqp.Publishing{
			Body: buf.Bytes(),
		}
		ch.Publish(
			"",
			dataQueue.Name,
			false,
			false,
			msg)
		log.Printf("Sent :%v\n", value)
	}
}

func calcValue() {
	var minStep, maxStep float64

	if value < nom {
		maxStep = *stepSize
		minStep = -1 * *stepSize * (value - *min) / (nom - *min)
	} else {
		maxStep = *stepSize * (*max - value) / (*max - nom)
		minStep = -1 * *stepSize
	}

	value += r.Float64()*(maxStep-minStep) + minStep
}

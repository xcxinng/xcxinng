package main

import (
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/streadway/amqp"
)

func onError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func fibonacciRPC(n int) (res int, err error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	onError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	onError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	onError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	onError(err, "Failed to register a consumer")

	corrId := randomString(32)

	err = ch.Publish(
		"",          // exchange
		"rpc_queue", // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          []byte(strconv.Itoa(n)),
		})
	onError(err, "Failed to publish a message")

	for d := range msgs {
		if corrId == d.CorrelationId {
			res, err = strconv.Atoi(string(d.Body))
			onError(err, "Failed to convert body to integer")
			break
		}
	}

	return
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	n := bodyFrom(os.Args)

	log.Printf(" [x] Requesting fib(%d)", n)
	res, err := fibonacciRPC(n)
	onError(err, "Failed to handle RPC request")

	log.Printf(" [.] Got %d", res)
}

func bodyFrom(args []string) int {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "30"
	} else {
		s = strings.Join(args[1:], " ")
	}
	n, err := strconv.Atoi(s)
	onError(err, "Failed to convert arg to integer")
	return n
}

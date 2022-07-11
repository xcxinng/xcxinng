package main

import (
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	myExchange = "my-exchange"
	paramError = "param invalid, this program accept at least 2 params:\n" +
		"args1 >>>> queue name that you are targeting to produce\n" +
		"args2 >>>> a list of routing key\n" +
		"for example:\n" +
		"./produce my_queue key1 key2 key3"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func newTask(body []byte, key string) {
	// establish an ampq connection
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// open a channel which is multiplexing over ampq connection
	// one amqp connection can carry multiple channels.
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(myExchange, "direct", true, false, false, false, nil)
	if err != nil {
		failOnError(err, "error declare exchange")
	}

	err = ch.Publish(
		myExchange, // exchange
		key,        // routing key
		false,      // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         body,
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s", body)
}

// produce.go receive only 2 args
// args1 message body
// args2 a routing key of this message
//
// example:
// publish one message to queue "test1" and with routing key "test1"
// # go build -o produce .
// # ./produce test1 test1
func main() {
	if len(os.Args) < 3 {
		log.Fatal(paramError)
	}
	newTask([]byte(os.Args[1]), os.Args[2])
}

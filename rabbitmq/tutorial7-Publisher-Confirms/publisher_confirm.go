package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// Summary
// 1. Dial to establish connection
// 2. conn.Channel to open a channel over that connection
// 3. channel.Confirm(false) to tell client enter publish confirm mode (otherwise it will ignore ack from server)
// 4. channel.NotifyPublish(chan confirmation) to receive confirm messages from server
// 5. channel.QueueDeclare() get your queue
// 6. channel.Publish()
func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	confirms := make(chan amqp.Confirmation)
	ch.NotifyPublish(confirms)
	go func() {
		for confirm := range confirms {
			if confirm.Ack {
				// code when messages are confirmed
				log.Printf("Confirmed")
			} else {
				// code when messages are nack-ed
				log.Printf("Nacked")
			}
		}
	}()

	err = ch.Confirm(false)
	failOnError(err, "Failed to confirm")

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	consume(ch, q.Name)
	publish(ch, q.Name, "hello")

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	var forever chan struct{}
	<-forever
}

func consume(ch *amqp.Channel, qName string) {
	msgs, err := ch.Consume(
		qName, // queue
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
}

func publish(ch *amqp.Channel, qName, text string) {
	err := ch.Publish("", qName, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(text),
	})
	failOnError(err, "Failed to publish a message")
}

package main

import (
	"bytes"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func onError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

const (
	exchange            = "my-exchange"
	ErrorConsumingParam = "param invalid, this program only accept 2 params:\n" +
		"args1 >>> body of message\n" +
		"args2 >>> a routing key of this message\n" +
		"For example:\n" +
		"./consume my_queue routing_key_of_my_queue"
)

func consume(queueName string, routingKey ...string) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	onError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	onError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(exchange, "direct", true, false, false, false, nil)
	if err != nil {
		onError(err, "error declare exchange")
	}

	q, err := ch.QueueDeclare(queueName, false, true, false, false, nil)
	if err != nil {
		onError(err, "error declare")
	}

	for _, key := range routingKey {
		err = ch.QueueBind(q.Name, key, exchange, false, nil)
		if err != nil {
			onError(err, "error bind")
		}
	}

	//err = ch.Qos(
	//	// achieve fair dispatch not RR dispatch blindly
	//	1,     // prefetch count
	//	0,     // prefetch size
	//	false, // global
	//)

	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer

		//true, // auto-ack
		//set autoAck to "false" if you want to ack msg manually
		true,

		false,
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	onError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("consumed a message: %s", d.Body)
			dotCount := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dotCount)
			log.Printf("sleep %d seconds\n", dotCount)
			time.Sleep(t * time.Second)
			log.Printf("Done\n")
			// when set autoAck to false, remember to send ack manually
			//_ = d.Ack(false)
			//
			// As official docs said, it's a common easy error,but the consequences are serious.
			//
			// For troubleshooting :
			// [localhost:~]# rabbitmqctl list_queues name messages_ready messages_unacknowledged
			// timeout: 60.0 seconds ...
			// Listing queues for vhost / ...
			// name	messages_ready	messages_unacknowledged
			// work_queue	0	0
			// hello_world_queue	0	0
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

// consume.go receive at least 2 args
// args1 queue name that you are targeting to consume
// args2 a list of routing key, keys will be attached to message
//
// example1:
// consume messages to queue "test2" and with routing key "test2"
// # go build -o consume .
// # ./consume test2 test2
//
// example2:
// consume messages to queue "test3" and with routing key "test1" and "test2"
// # go build -o consume .
// # ./consume test3 test1 test2
//
// It can be represented graphically as ./direct_exchange.png
// binding table see: ./direct_exchange_rules.png
func runConsumer() {
	if len(os.Args) < 3 {
		log.Fatal(ErrorConsumingParam)
	}
	consume(os.Args[1], os.Args[2:]...)
}

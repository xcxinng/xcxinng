package rabbitmq

import (
    "bytes"
    "github.com/streadway/amqp"
    "log"
    "time"
)

func worker() {
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    failOnError(err, "Failed to connect to RabbitMQ")
    defer conn.Close()

    ch, err := conn.Channel()
    failOnError(err, "Failed to open a channel")
    defer ch.Close()

    q, err := ch.QueueDeclare(
        "task_queue", // name

        true, // durable
        // you can set durable to false that keep data only in memory.
        // false

        false, // delete when unused
        false, // exclusive
        false, // no-wait
        nil,   // arguments
    )
    failOnError(err, "Failed to declare a queue")

    err = ch.Qos(
        // achieve fair dispatch not RR dispatch blindly
        1,     // prefetch count
        0,     // prefetch size
        false, // global
    )

    msgs, err := ch.Consume(
        q.Name, // queue
        "",     // consumer

        //true, // auto-ack
        //set autoAck to "false" if you want to ack msg manually
        false,

        false,
        false, // no-local
        false, // no-wait
        nil,   // args
    )
    failOnError(err, "Failed to register a consumer")

    forever := make(chan bool)

    go func() {
        for d := range msgs {
            log.Printf("Received a message: %s", d.Body)
            dotCount := bytes.Count(d.Body, []byte("."))
            t := time.Duration(dotCount)
            time.Sleep(t * time.Second)
            log.Printf("Done")
            // when set autoAck to false, remember to send ack manually
            _ = d.Ack(false)
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

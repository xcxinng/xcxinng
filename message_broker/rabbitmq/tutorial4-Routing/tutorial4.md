#Routing
In the previous tutorial we built a simple logging system. We were able to broadcast log messages to many receivers.   

In this tutorial we're going to add a feature to it - we're going to make it possible to subscribe only to a subset of 
the messages. For example, we will be able to direct only critical error messages to the log file (to save disk space), 
while still being able to print all the log messages on the console.   

#Binding
In previous examples we were already creating bindings. You may recall code like:   
```go
// QueueBind bind an exchange with a queue through routing key.
// Or tell exchange what messages the queue is interested in.
err = ch.QueueBind(
  q.Name, // queue name
  "",     // routing key
  "logs", // exchange
  false,
  nil)
```
A binding is a relationship between an exchange and a queue. This can be simply read as: the queue is interested in 
messages from this exchange.   
   
Bindings can take an extra <mark>routing_key</mark> parameter. To avoid the confusion with a Channel.Publish parameter we're going
to call it a binding key. This is how we could create a binding with a key:   
```go
err = ch.QueueBind(
  q.Name,    // queue name
  "black",   // routing key
  "logs",    // exchange
  false,
  nil)
```
The meaning of a binding key depends on the exchange type. The <mark>fanout</mark> exchanges, which we used previously, simply **ignored** its value.

#Direct exchange
Our logging system from the previous tutorial broadcasts all messages to all consumers. We want to extend that to allow 
filtering messages based on their severity. For example we may want the script which is writing log messages to the disk
to only receive critical errors, and not waste disk space on warning or info log messages.   

We were using a fanout exchange, which doesn't give us much flexibility - it's only capable of mindless broadcasting.   
   
We will use a <mark>direct</mark> exchange instead. The routing algorithm behind a <mark>direct</mark> exchange is 
simple - a message goes to the queues whose <mark>binding key</mark> exactly matches the <mark>routing key</mark> of 
the message.   
To illustrate that, consider the following setup:   
![direct exchange](images/direct-exchange.png)   
In this setup, we can see the <mark>direct</mark> exchange <mark>X</mark> with two queues bound to it. The first queue
is bound with binding key <mark>orange</mark>, and the second has two bindings, one with binding key <mark>black</mark> 
and the other one with <mark>green</mark>.   
   
In such a setup a message published to the exchange with a routing key orange will be routed to queue Q1. Messages with
a routing key of black or green will go to Q2. All other messages will be discarded.   
   
#Multiple bindings
![multiple_bindings](images/direct-exchange-multiple.png)   

It is perfectly legal to bind multiple queues with the same binding key. In our example we could add a binding between
<mark>X</mark> and <mark>Q1</mark> with binding key <mark>black</mark>. In that case, the <mark>direct</mark> exchange 
will behave like fanout and will broadcast the message to all the matching queues. A message with routing key 
<mark>black</mark> will be delivered to both <mark>Q1</mark> and <mark>Q2</mark>.   
   
#Emitting logs
We'll use this model for our logging system. Instead of fanout we'll send messages to a direct exchange. We will supply
the log severity as a routing key. That way the receiving script will be able to select the severity it wants to receive.
Let's focus on emitting logs first.

We'll use this model for our logging system. Instead of fanout we'll send messages to a direct exchange. We will supply
the log severity as a routing key. That way the receiving script will be able to select the severity it wants to receive. 
Let's focus on emitting logs first.   

As always, we need to create an exchange first:   
```go
err = ch.ExchangeDeclare(
  "logs_direct", // name
  "direct",      // type
  true,          // durable
  false,         // auto-deleted
  false,         // internal
  false,         // no-wait
  nil,           // arguments
)
```
And we're ready to send a message:   
```go
err = ch.ExchangeDeclare(
  "logs_direct", // name
  "direct",      // type
  true,          // durable
  false,         // auto-deleted
  false,         // internal
  false,         // no-wait
  nil,           // arguments
)
failOnError(err, "Failed to declare an exchange")

body := bodyFrom(os.Args)
err = ch.Publish(
  "logs_direct",         // exchange
  severityFrom(os.Args), // routing key
  false, // mandatory
  false, // immediate
  amqp.Publishing{
    ContentType: "text/plain",
    Body:        []byte(body),
})
```
To simplify things we will assume that 'severity' can be one of **'info'**, **'warning'**,**'error'**.   

#Subscribing
Receiving messages will work just like in the previous tutorial, with one exception - we're going to create a new 
binding for each severity we're interested in.   
```go
q, err := ch.QueueDeclare(
  "",    // name
  false, // durable
  false, // delete when unused
  true,  // exclusive
  false, // no-wait
  nil,   // arguments
)
failOnError(err, "Failed to declare a queue")

if len(os.Args) < 2 {
  log.Printf("Usage: %s [info] [warning] [error]", os.Args[0])
  os.Exit(0)
}
for _, s := range os.Args[1:] {
  log.Printf("Binding queue %s to exchange %s with routing key %s",
     q.Name, "logs_direct", s)
  err = ch.QueueBind(
    q.Name,        // queue name
    s,             // routing key
    "logs_direct", // exchange
    false,
    nil)
  failOnError(err, "Failed to bind a queue")
}
```
#Putting it all together
![tutorial4](images/python-four.png)   
producer code in **emit_log_direct.go**   
consumer code in **receive_logs_direct.go**

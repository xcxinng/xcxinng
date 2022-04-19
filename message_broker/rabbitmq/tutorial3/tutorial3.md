#Publish/Subscribe
##Preface
When referring to the message model, 2 kinds of them occur to our mind:
>1. Producer/Consumer (one message one worker/consumer)
>2. Publish/Subscribe (one message multiple workers/consumers)   

In tutorial2, we created a work queue. The assumption behind a work queue is that each task is delivered to exactly one   
worker, such a pattern is known as "Producer/Consumer".

##title2
In this part we'll do something completely different: we'll deliver a message to multiple consumers. This one is   
known as "Publish/Subscribe".  

To illustrate the pattern, we're going to build a simple logging system. It will consist of two programs -- the first
will emit log messages and the second will receive and print them.

In our logging system every running copy of the receiver program will get the messages. That way we'll be able to run 
one receiver and direct the logs to disk; and at the same time we'll be able to run another receiver and see the logs
on the screen.

Essentially, published log messages are going to be broadcast to all the receivers.

##Exchanges
In previous parts of the tutorial we sent and received messages to and from a queue. Now it's time to introduce the 
full messaging model in Rabbit.   
Let's quickly go over what we covered in the previous tutorials:   
* A producer is a user application that sends messages.
* A queue is a buffer that stores messages.
* A consumer is a user application that receives messages.   

The core idea in the messaging model in RabbitMQ is that the producer never sends any messages directly to a queue. 
Actually, quite often the producer doesn't even know if a message will be delivered to any queue at all.  
   
Instead, the producer can only send messages to an exchange. An exchange is a very simple thing. On one side it 
receives messages from producers and the other side it pushes them to queues. The exchange must know exactly what 
to do with a message it receives. Should it be appended to a particular queue? Should it be appended to many queues? 
Or should it get discarded. The rules for that are defined by the exchange type.   
![exchange_model](images/exchanges.png)   

There are a few exchange types available, by default it has 7 types in rabbit:   
![exchange_types](images/exchange_type.png)   
We focus on fanout type, and let's create an exchange of this type, and call it logs:
```go
err = ch.ExchangeDeclare(
  "logs",   // name
  "fanout", // type
  true,     // durable
  false,    // auto-deleted
  false,    // internal
  false,    // no-wait
  nil,      // arguments
)
```
The fanout exchange is very simple. As you can probably guess from the name, it just broadcasts all the messages it 
receives to all the queues it knows. And that's exactly what we need for our logger.   
> ###Listing exchanges
>
> To list the exchanges on the server you can run the ever useful <mark>rabbitmqctl</mark>:   
> ```shell
>  sudo rabbitmqctl list_exchanges
> ```
> In this list there will be some amq.* exchanges and the default (unnamed) exchange. These are created by default, 
> but it is unlikely you'll need to use them at the moment.
> ###The default exchange
> In previous parts of the tutorial we knew nothing about exchanges, but still were able to send messages to queues. 
> That was possible because we were using a default exchange, which is identified by the empty string <mark>("")</mark>.
>    
> Recall how we published a message before:
> ```go
> err = ch.Publish(
>"",     // exchange
>q.Name, // routing key
>false,  // mandatory
>false,  // immediate
>amqp.Publishing{
>ContentType: "text/plain",
>Body:        []byte(body),
>})
>```
> Here we use the default or nameless exchange: messages are routed to the queue with the name specified by routing_key 
> parameter, if it exists. The default exchange <mark>""</mark> is bound to all queues, so it could send messages to all
> queues.

Now, we can publish to our named exchange instead:
```go
err = ch.ExchangeDeclare(
  "logs",   // name
  "fanout", // type
  true,     // durable
  false,    // auto-deleted
  false,    // internal
  false,    // no-wait
  nil,      // arguments
)
failOnError(err, "Failed to declare an exchange")

body := bodyFrom(os.Args)
err = ch.Publish(
  "logs", // exchange
  "",     // routing key
  false,  // mandatory
  false,  // immediate
  amqp.Publishing{
          ContentType: "text/plain",
          Body:        []byte(body),
  })
```
##Temporary queues
As you may remember previously we were using queues that had specific names (remember hello and task_queue?). Being able to name a queue was crucial for us -- we needed to point the workers to the same queue. Giving a queue a name is important when you want to share the queue between producers and consumers.

But that's not the case for our logger. We want to hear about all log messages, not just a subset of them. We're also interested only in currently flowing messages not in the old ones. To solve that we need two things.

Firstly, whenever we connect to Rabbit we need a fresh, empty queue. To do this we could create a queue with a random name, or, even better - let the server choose a random queue name for us.

Secondly, once we disconnect the consumer the queue should be automatically deleted.

In the amqp client, when we supply queue name as an empty string, we create a non-durable queue with a generated name:
```go
q, err := ch.QueueDeclare(
  "",    // name
  false, // durable
  false, // delete when unused
  true,  // exclusive
  false, // no-wait
  nil,   // arguments
)
```
When the method returns, the queue instance contains a random queue name generated by RabbitMQ. For example it may look like amq.gen-JzTY20BRgKO-HjmUJj0wLg.

When the connection that declared it closes, the queue will be deleted because it is declared as exclusive.

You can learn more about the exclusive flag and other queue properties in the [guide on queues](https://www.rabbitmq.com/queues.html).

##Bindings

![bindings](images/bindings.png)

We've already created a fanout exchange and a queue. Now we need to tell the exchange to send messages to our queue. 
That relationship between exchange and a queue is called a binding.
```go
err = ch.QueueBind(
  q.Name, // queue name
  "",     // routing key
  "logs", // exchange
  false,
  nil,
)
```
From now on the <mark>logs</mark> exchange will append messages to our queue.
>###Listing bindings
>
>You can list existing bindings using, you guessed it,
>
> ```shell
> rabbitmqctl list_bindings
> ```
##Putting it all together
![overall](images/python-three-overall.png)   

The producer program, which emits log messages, doesn't look much different from the previous tutorial. The most 
important change is that we now want to publish messages to our <mark>logs</mark> exchange instead of the nameless one. We need to
supply a <mark>routingKey</mark> when sending, but its value is ignored for <mark>fanout</mark> exchanges.    
The example of publish code see <mark>emit_log.go</mark>

As you see, after establishing the connection we declared the exchange. This step is necessary as publishing to a non-existing exchange is forbidden.

The messages will be lost if no queue is bound to the exchange yet, but that's okay for us; if no consumer is listening 
yet we can safely discard the message.

And the example code of subscribe see <mark>receive_logs.go</mark>

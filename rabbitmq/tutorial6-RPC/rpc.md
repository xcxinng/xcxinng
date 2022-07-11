# Remote procedure call (RPC)

In the [second tutorial](../tutorial2-Work-Queues) we learned how to use Work Queues to distribute time-consuming tasks
among multiple workers.

But what if we need to run a function on a remote computer and wait for the result? Well, that's a different story. This
pattern is commonly known as Remote Procedure Call or RPC.

In this tutorial we're going to use RabbitMQ to build an RPC system: a client and a scalable RPC server. As we don't
have any time-consuming tasks that are worth distributing, we're going to create a dummy RPC service that returns
Fibonacci numbers.

> ### A note on RPC
>
>Although RPC is a pretty common pattern in computing, it's often criticised. The problems arise when a programmer is
> not aware whether a function call is local or if it's a slow RPC. Confusions like that result in an unpredictable
> system and adds unnecessary complexity to debugging. Instead of simplifying software, misused RPC can result in
> unmaintainable spaghetti code.
>
>Bearing that in mind, consider the following advice:
>
>* Make sure it's obvious which function call is local and which is remote.
>* Document your system. Make the dependencies between components clear.
>* Handle error cases. How should the client react when the RPC server is down for a long time?
>
>When in doubt avoid RPC. If you can, you should use an asynchronous pipeline - instead of RPC-like blocking, results
> are asynchronously pushed to a next computation stage.

# Callback queue

In general doing RPC over RabbitMQ is easy. A client sends a request message and a server replies with a response
message. In order to receive a response we need to send a 'callback' queue address with the request. We can use the
default queue. Let's try it:

```go
q, err := ch.QueueDeclare(
"", // name
false, // durable
false, // delete when unused
true,  // exclusive
false, // noWait
nil,   // arguments
)

err = ch.Publish(
"", // exchange
"rpc_queue", // routing key
false,       // mandatory
false,       // immediate
amqp.Publishing{
ContentType:   "text/plain",
CorrelationId: corrId, // for the client to recognize which request/response
ReplyTo:       q.Name, // tell mq to which queue distribute
Body:          []byte(strconv.Itoa(n)),
})
```

> Message properties
>
>The AMQP 0-9-1 protocol predefines a set of 14 properties that go with a message. Most of the properties are rarely used,
> except the following:
>
>- <mark>persistent</mark>: Marks a message as persistent (with a value of true) or transient (false). You may remember
   > this property from the second tutorial.
>- <mark>content_type</mark>: Used to describe the mime-type of the encoding. For example for the often used JSON
   > encoding it is a good practice to set this property to: application/json.
>- <mark>reply_to</mark>: Commonly used to name a callback queue.
>- <mark>correlation_id</mark>: Useful to correlate RPC responses with requests.

# Correlation ID

In the method presented above we suggest creating a callback queue for every RPC request. That's pretty inefficient, but
fortunately there is a better way - let's create a single callback queue per client.

That raises a new issue, having received a response in that queue it's not clear to which request the response belongs.
That's when the correlation_id property is used. We're going to set it to a unique value for every request. Later, when
we receive a message in the callback queue we'll look at this property, and based on that we'll be able to match a
response with a request. If we see an unknown correlation_id value, we may safely discard the message - it doesn't
belong to our requests.

You may ask why should we ignore unknown messages in the callback queue, rather than failing with an error? It's due to
a possibility of a race condition on the server side. Although unlikely, it is possible that the RPC server will die
just after sending us the answer, but before sending an acknowledgment message for the request. If that happens, the
restarted RPC server will process the request again. That's why on the client we must handle the duplicate responses
gracefully, and the RPC should ideally be idempotent.

# Summary

![rpc](rpc-tutorial.png)

Our RPC will work like this:

- When the Client starts up, it creates an anonymous exclusive callback queue.
- For an RPC request, the Client sends a message with two properties:  reply_to, which is set to the callback queue and
  correlation_id, which is set to a unique value for every request.
- The request is sent to an rpc_queue queue.
- The RPC worker (aka: server) is waiting for requests on that queue. When a request appears, it does the job and sends
  a message with the result back to the Client, using the queue from the reply_to field.
- The client waits for data on the callback queue. When a message appears, it checks the correlation_id property. If it
  matches the value from the request it returns the response to the application.

#Put it all together
see code in **rpc_client.go** and **rpc_server.go**

#Advantage of RPC service based on rabbit
It may seem more complicated than RPC service without rabbit, therefor, why should we build RPC service based on rabbit?

After this tutorial, I've done some research on "RPC service based on rabbit", here are some advantages
I've found on internet:
1. decouple the RPC client and server
2. reduce the stress on the server
3. make it easier to expand horizontally
4. rabbit has supported RPC friendly

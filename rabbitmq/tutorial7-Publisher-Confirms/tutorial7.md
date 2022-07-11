# Publisher Confirms
Publisher confirms are a RabbitMQ extension to implement reliable publishing. When publisher confirms are enabled on a
channel, messages the client publishes are confirmed asynchronously by the broker, meaning they have been taken care of
on the server side.

# Overview
In this tutorial we're going to use publisher confirms to make sure published messages have safely reached the broker.
We will cover several strategies to using publisher confirms and explain their pros and cons.

# Enabling Publisher Confirms on a Channel
Publisher confirms are a RabbitMQ extension to the AMQP 0.9.1 protocol, so they are not enabled by default. Publisher
confirms are enabled at the channel level with the confirmSelect method:

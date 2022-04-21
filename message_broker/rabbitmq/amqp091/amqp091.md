# AMQP 0-9-1 Model Explained
## Overview
This guide provides an overview of the AMQP 0-9-1 protocol, one of the protocols supported by RabbitMQ
> 概述了rabbitMQ支持的一种高级消息队列协议：AMQP 0-9-1

## What is AMQP 0-9-1?
AMQP 0-9-1 (Advanced Message Queuing Protocol) is a messaging protocol that enables conforming client applications to
communicate with conforming messaging middleware brokers.
> AMQP 0-9-1（高级消息队列协议）是一种消息协议，使符合要求的客户端应用能够与符合要求的消息中间件代理服务器通信。

## Brokers and Their Role
Messaging brokers receive messages from publishers (applications that publish them, also known as producers) and
route them to consumers (applications that process them). Since it is a network protocol, the publishers, consumers
and the broker can all reside on different machines.
>消息代理服务器接收来自发布者(发布它们的应用程序，也称为生产者)的消息，并将其路由给消费者(处理它们的应用程序)。由于它是一种网络协议，
> 发布者、消费者和消息代理都可以部署在不同的机器上。

## AMQP 0-9-1 Model in Brief
The AMQP 0-9-1 Model has the following view of the world: messages are published to exchanges, which are often
compared to post offices or mailboxes. Exchanges then distribute message copies to queues using rules called bindings.
Then the broker either deliver messages to consumers subscribed to queues, or consumers fetch/pull messages from
queues on demand.
>AMQP 0-9-1模型具有以下世界观：消息发布到exchanges，exchange通常与邮局或邮箱进行比较。然后，exchange使用称为binding的规则将消息副本分
> 发给队列。然后，broker要么向订阅队列的消费者发送消息，要么消费者按需从队列中获取/拉取消息。

![routing](hello-world-example-routing.png)
When publishing a message, publishers may specify various message attributes (message meta-data). Some of this
meta-data may be used by the broker, however, the rest of it is completely opaque to the broker and is only used by
applications that receive the message.
>发布消息时，发布者可以指定各种消息属性（消息元数据）。其中一些元数据可能被broker使用，但除了那部分外，其余的消息元数据对broker来说完全不透明，
> 仅供接收消息的应用程序使用。

Networks are unreliable and applications may fail to process messages therefore the AMQP 0-9-1 model has a notion of
message acknowledgements: when a message is delivered to a consumer the consumer notifies the broker, either
automatically or as soon as the application developer chooses to do so. When message acknowledgements are in use,
a broker will only completely remove a message from a queue when it receives a notification for that message
(or group of messages).
> 网络不可靠时，应用程序可能无法处理消息，因此AMQP 0-9-1模型具有消息确认的概念：当消息发送给消费者时，消费者要么使用自动确认机制，要么由
> 应用程序开发人员选择自行通知broker。如果启用了消息确认，broker只会在收到该消息（或一组消息）的通知时才会从队列中完全删除消息。

In certain situations, for example, when a message cannot be routed, messages may be returned to publishers, dropped,
or, if the broker implements an extension, placed into a so-called "dead letter queue". Publishers choose how to
handle situations like this by publishing messages using certain parameters.
> 在某些情况下，例如，当消息无法路由时，消息可能会被返回给发布者，删除，或者，如果broker实现了某种插件，则放置在所谓的“dead letter queue”队列中。
> 发布者可以通过使用某些消息参数，在发布消息时决定如何处理这种情况。

Queues, exchanges and bindings are collectively referred to as AMQP entities.
> Queues, exchanges and bindings 统称为AMQP实体.

## AMQP 0-9-1 is a Programmable Protocol
AMQP 0-9-1 is a programmable protocol in the sense that AMQP 0-9-1 entities and routing schemes are primarily
defined by applications themselves, not a broker administrator. Accordingly, provision is made for protocol
operations that declare queues and exchanges, define bindings between them, subscribe to queues and so on.
> AMQP 0-9-1是一种可编程协议，即AMQP 0-9-1实体和路由方案主要由应用程序自己定义，而不是broker服务器管理员。
> 因此，AMQP协议为声明queue和exchange、定义它们之间的binding、订阅队列等操作做出了明确规定。

This gives application developers a lot of freedom but also requires them to be aware of potential definition
conflicts. In practice, definition conflicts are rare and often indicate a misconfiguration.
> 这给了应用程序开发人员很多自由，但也要求他们意识到潜在的定义冲突。在实践中，定义冲突很少见，通常(定义冲突)表示为配置错误。

Applications declare the AMQP 0-9-1 entities that they need, define necessary routing schemes and may choose
to delete AMQP 0-9-1 entities when they are no longer used.
> 应用程序声明他们需要的AMQP 0-9-1实体，定义必要的路由方案，也可以选择在不再使用这些AMQP 0-9-1实体时删除它们。

## Exchanges and Exchange Types
Exchanges are AMQP 0-9-1 entities where messages are sent to. Exchanges take a message and route it into zero or
more queues. The routing algorithm used depends on the exchange type and rules called bindings. AMQP 0-9-1
brokers provide four exchange types:
> exchange是消息发送到的AMQP 0-9-1实体。exchange会接收消息并将其路由到零个或多个队列中。使用的路由算法取决于交换类型和称为绑定的规则。
> AMQP 0-9-1经纪人提供四种交易所类型：

| Exchange Type|Default pre-declared names|
|---|---|
|Direct exchange|(Empty string) and amq.direct|
|Fanout exchange|amq.fanout|
|Topic exchange	|amq.topic|
|Headers exchange|	amq.match (and amq.headers in RabbitMQ)|

Besides the exchange type, exchanges are declared with a number of attributes, the most important of which are:
- Name
- Durability (exchanges survive broker restart)
- Auto-delete (exchange is deleted when last queue is unbound from it)
- Arguments (optional, used by plugins and broker-specific features)
> 除了类型外，exchange还有很多其他属性，其中最重要的是：
> * 名字
> * 耐久性(持久性)
> * 自动删除
> * 参数（可选的，由插件和broker其他特定功能使用）

Exhanges can be durable or transient. Durable exchanges survive broker restart whereas transient exchanges do not
(they have to be redeclared when broker comes back online). Not all scenarios and use cases require exchanges to
be durable.
> exchange可以是持久的（一直存在），也可以是短暂的（临时存在）。持久性的exchange实体可以在broker重启后继续存在，而临时的则会消失
> （当broker重新上线时，必须重新定义）。并非所有场景和用例都需要持久性的exchange。

### Default Exchange
The default exchange is a direct exchange with no name (empty string) pre-declared by the broker. It has one special
property that makes it very useful for simple applications: every queue that is created is automatically bound to it
with a routing key which is the same as the queue name.
> 默认的exchange是Direct Exchange，它是broker预定义好的，没有具体名称（空字符串）。它有一个特殊的属性，对简单的应用程序非常有用：
> **创建的每个队列都会自动绑定到它，routing key 与 queue name 相同**。

For example, when you declare a queue with the name of "search-indexing-online", the AMQP 0-9-1 broker will bind it
to the default exchange using "search-indexing-online" as the routing key (in this context sometimes referred to as
the binding key). Therefore, a message published to the default exchange with the routing key "search-indexing-online"
will be routed to the queue "search-indexing-online". In other words, the default exchange makes it seem like it is
possible to deliver messages directly to queues, even though that is not technically what is happening.
> 例如，当您声明名为“search-indexing-online”的队列时，AMQP 0-9-1代理将使用“search-indexing-online”作为routing key,
> 将其绑定到默认exchange（在这种情况下，有时被称为"binding key"）。因此，使用routing-key: “search-indexing-online”发布到
> 默认exchange的消息将被路由到队列“search-indexing-online”。
> 换句话说，默认exchange似乎可以将消息直接发送到队列，尽管从技术上来讲，事实并非所描述的那样。

### Direct Exchange
A direct exchange delivers messages to queues based on the message routing key. A direct exchange is ideal for the
unicast routing of messages (although they can be used for multicast routing as well). Here is how it works:
* A queue binds to the exchange with a routing key K
* When a new message with routing key R arrives at the direct exchange, the exchange routes it to the queue if K = R
> direct exchange 基于routing key属性把消息发送给队列。direct exchange 是单一路由消息（即只有一个目的队列的消息）的理想型。
> 当然，direct exchange 也可以用于多播。以下是它的工作原理：
> * 某个队列通过routing key K 绑定到exchange
> * 当一个带着routing key属性为R的消息到达direct exchange时，exchange会判断K是否等于R，如果相等则把消息入队


Direct exchanges are often used to distribute tasks between multiple workers (instances of the same application)
in a round robin manner. When doing so, it is important to understand that, in AMQP 0-9-1, messages are load
balanced between consumers and not between queues.
A direct exchange can be represented graphically as follows:
> direct exchange通常用于以轮询方式在多个worker（同一应用程序的不同实例）之间分配任务。这样做时，重要的是要了解，在AMQP 0-9-1中，
> **消息的负载均衡是在消费者之间而不是队列之间**。direct exchange图解如下：
![exchange direct](exchange-direct.png)


### Fanout Exchange
A fanout exchange routes messages to all of the queues that are bound to it and the routing key is ignored. If N queues
are bound to a fanout exchange, when a new message is published to that exchange a copy of the message is delivered to
all N queues. Fanout exchanges are ideal for the broadcast routing of messages.
> Fanout exchange将消息路由到绑定到它的**所有**队列，routing key会被直接忽略。
> 如果N个队列绑定到fanout exchange，则当新消息发布到该exchange时，消息的副本将发送到N个队列中。
> Fanout exchange是消息广播路由的理想选择。

Because a fanout exchange delivers a copy of a message to every queue bound to it, its use cases are quite similar:
* Massively multi-player online (MMO) games can use it for leaderboard updates or other global events
* Sport news sites can use fanout exchanges for distributing score updates to mobile clients in near real-time
* Distributed systems can broadcast various state and configuration updates
* Group chats can distribute messages between participants using a fanout exchange (although AMQP does not have a
built-in concept of presence, so XMPP may be a better choice).

A fanout exchange can be represented graphically as follows:
> 由于fanout exchange会向绑定到它的每个队列发送消息的副本，因此其用例非常相似：
> * 大型多人在线（MMO）游戏可用于排行榜更新或其他全球活动
> * 体育新闻网站可以使用fanout exchange来近实时地向移动客户端推送分数的更新
> * 分布式系统可以广播各种状态和配置更新
> * 群聊可以使用fanout exchange在参与者之间分发消息（尽管AMQP没有内置的presence概念，因此XMPP可能是一个更好的选择）

![fanout exchange](exchange-fanout.png)

### Topic Exchange
Topic exchanges route messages to one or many queues based on matching between a message routing key and the pattern that was used to bind a queue to an exchange. The topic exchange type is often used to implement various publish/subscribe pattern variations. Topic exchanges are commonly used for the multicast routing of messages.

Topic exchanges have a very broad set of use cases. Whenever a problem involves multiple consumers/applications that selectively choose which type of messages they want to receive, the use of topic exchanges should be considered.

Example uses:
* Distributing data relevant to specific geographic location, for example, points of sale
* Background task processing done by multiple workers, each capable of handling specific set of tasks
* Stocks price updates (and updates on other kinds of financial data)
* News updates that involve categorization or tagging (for example, only for a particular sport or team)
* Orchestration of services of different kinds in the cloud
* Distributed architecture/OS-specific software builds or packaging where each builder can handle only one architecture or OS
> topic exchange根据消息routing key和用于将queue绑定到exchange的模式之间的匹配，将消息路由到一个或多个队列。
> topic exchange通常用于实现各种发布/订阅模式变化。也通常用于消息的多播路由。
> topic exchange 有一套非常广泛的用例。每当问题涉及多个消费者/应用程序有选择地选择他们想要接收哪种类型的消息时，应考虑使用topic exchange。
>示例用途：
>
>* 分发与特定地理位置相关的数据，例如销售点
>* 由多个worker完成的后台任务处理，每个工人都能够处理特定的任务集
>* 股票价格更新（以及其他类型的财务数据更新）
>* 涉及分类或标记的新闻更新（例如，仅适用于特定运动或团队）
>* 在云端编排不同类型的服务
>* 分布式体系结构/特定于操作系统的软件构建或打包，每个构建器只能处理一个架构或操作系统

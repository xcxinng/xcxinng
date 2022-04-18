#Reference   
origin official page: https://www.rabbitmq.com/tutorials/tutorial-two-go.html   

#Communicating Model   

![image info](./one_producer_multiple_consumers.png)

#Read with the following questions
- How rabbitMQ dispatch tasks ?
- How to avoid losing messages when one of(or part of) consumers die ?
- How to avoid losing messages when rabbitMQ server dies(restart) ? (Data Persistence)
- Is there possible to dispatch tasks(messages) fairly ?

#Summary
1. By default, rabbitMQ dispatch tasks in Round-robin way(aka,RR)
2. Using a message acknowledgment mechanism to ensure messages do not get lost
3. rabbitMQ support to persist data in memory into a disk, enable it by setting <mark>durable=true</mark> on QueueDeclare
4. Using a QoS to achieve fair dispatch

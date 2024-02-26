# Service Discovery With etcd

## Before Start

A complete service discovery system include below three key functions:

1. Service registration: A service must register itself to some common place so that others can discover it
2. Health check: An instance must report its health status to the common place so that others can still talk to and discover it
3. Service discovery: A service initiates the service discovery process by initializing client with specified target service name to find

服务注册： 服务正式运行前，需要将自身地址注册到一个公共地方，以让其他服务知道它的存在以及如何找到它
健康检查： 服务注册后，还要定期更新状态，让其他服务知道它的健康运行
服务发现： 一个服务可以通过固定的服务名称找到对应服务的联系方式（通常是IP地址+端口之类的信息）

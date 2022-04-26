# Clustering Guide
An official guide can be found here: [cluster guide](https://www.rabbitmq.com/clustering.html)


## 1. Cluster Formation

### 1.1 What is rabbitMQ cluster?
A RabbitMQ cluster is a logical grouping of one or several nodes, each sharing users, virtual hosts, queues, exchanges,
bindings, runtime parameters and other distributed state.

### 1.2 Ways of Forming a Cluster
A RabbitMQ cluster can be formed in a number of ways:
* Declaratively by listing cluster nodes in config file
* Declaratively using DNS-based discovery
* Declaratively using AWS (EC2) instance discovery (via a plugin)
* Declaratively using Kubernetes discovery (via a plugin)
* Declaratively using Consul-based discovery (via a plugin)
* Declaratively using etcd-based discovery (via a plugin)
* Manually with rabbitmqctl

The composition of a cluster can be altered dynamically. All RabbitMQ brokers start out as running on a single node.
These nodes can be joined into clusters, and subsequently(随后) turned back into individual brokers again.

### 1.3 Node Names(Identifier)
Nodes are identified by node names. A node name has two parts: prefix(usually is "rabbit") and hostname,for example:
**rabbit@node1.msg.svc** (prefix **'rabbit'** and hostname **'node1.msg.svc'**)

Node names in a cluster must be unique. If more than one node is running on a given host (this is usually the case in
development and QA environments), they must use different prefixes, e.g. <mark>rabbit1@hostname</mark> and
<mark>rabbit2@hostname</mark>.

In a cluster, nodes identify and contact each other using node names. This means that the hostname part of every node
name must resolve. CLI tools also identify and address nodes using node names.

When a node starts up, it checks whether it has been assigned a node name. This is done via the <mark>RABBITMQ_NODENAME</mark>
environment variable. If no value was explicitly configured, the node resolves its hostname and prepends rabbit to it
to compute its node name.

If a system uses fully qualified domain names (FQDNs) for hostnames, RabbitMQ nodes and CLI tools must be configured to
use so-called long node names. For server nodes this is done by setting the RABBITMQ_USE_LONGNAME environment variable to true.

For CLI tools, either <mark>RABBITMQ_USE_LONGNAME</mark> must be set or the <mark>--longnames</mark> option must be specified.

### 1.4 Port Access
RabbitMQ has to bind to multiple ports(TCP Sockets) in order to accept client/CLI tool/inter-nodes connections.To ensure
rabbit work, make sure ports below can be used, otherwise, the node will fail to start.
>* 4369: epmd, a helper discovery daemon used by RabbitMQ nodes and CLI tools
>* 6000-6500: used by RabbitMQ Stream replication
>* 25672: used for inter-node and CLI tools communication
>* 35672-35682: used by CLI tools

Make sure that the firewall or SELinux allows external traffic to connect to the above ports!

### 1.5 Cookie File Locations
On UNIX systems, the cookie will be typically located in </mark>/var/lib/rabbitmq/.erlang.cookie</mark> (used by the
server) and <mark>$HOME/.erlang.cookie</mark> (used by CLI tools). Note that since the value of <mark>$HOME</mark>
varies from user to user, it's necessary to place a copy of the cookie file for each user that will be using the
CLI tools. This applies to both non-privileged users and root.

RabbitMQ nodes will log its effective user's home directory location early on boot.



## 2. Setup Cluster With <mark>rabbitmqctl</mark>
### 2.1 Prerequisite
1. setup hostname resolution by editing "/etc/hosts"
2. sync the **.erlang.cookie** file in both **"/var/lib/rabbitmq/.erlang.cookie"** and **"$HOME/.erlang.cookie"** on each node
3. don't forget to restart rabbitmq-server to take effect!

### 2.2 Join Cluster
![cluster_nodes](./cluster_nodes.png)

Let's make host147 and host148 join host146 cluster:
```shell
# on host147
rabbitmqctl stop_app
# => Stopping node rabbit@host147 ...done.

rabbitmqctl reset
# => Resetting node rabbit@host147 ...

rabbitmqctl join_cluster rabbit@host146
# => Clustering node rabbit@host147 with [rabbit@host146] ...done.

rabbitmqctl start_app
# => Starting node rabbit@host147 ...done.
```

```shell
# on host148
rabbitmqctl stop_app
# => Stopping node rabbit@host148 ...done.

# on host148
rabbitmqctl reset
# => Resetting node rabbit@host148 ...

rabbitmqctl join_cluster rabbit@host147
# => Clustering node rabbit@host148 with rabbit@host147 ...done.

rabbitmqctl start_app
# => Starting node rabbit@host148 ...done.
```
```shell
# on host146
rabbitmqctl cluster_status
# => Cluster status of node rabbit@host146 ...
# => [{nodes,[{disc,[rabbit@host146,rabbit@host147,rabbit@host148]}]},
# =>  {running_nodes,[rabbit@host148,rabbit@host147,rabbit@host146]}]
# => ...done.

# on host147
rabbitmqctl cluster_status
# => Cluster status of node rabbit@host147 ...
# => [{nodes,[{disc,[rabbit@host146,rabbit@host147,rabbit@host148]}]},
# =>  {running_nodes,[rabbit@host148,rabbit@host146,rabbit@host147]}]
# => ...done.

# on host148
rabbitmqctl cluster_status
# => Cluster status of node rabbit@host148 ...
# => [{nodes,[{disc,[rabbit@host148,rabbit@host147,rabbit@host146]}]},
# =>  {running_nodes,[rabbit@host147,rabbit@host146,rabbit@host148]}]
# => ...done.
```
### 2.3 Restarting Cluster Nodes
Nodes that have been joined to a cluster can be stopped at any time. They can also fail or be terminated by the OS.

In general, if the majority of nodes is still online after a node is stopped, this does not affect the rest of the
cluster, although client connection distribution, queue replica placement, and load distribution of the cluster will
change.

### 2.4 Schema Syncing from Online Peers
A restarted node will sync the schema and other information from its peers on boot. Before this process completes,
**the node won't be fully started and functional**.

It is therefore important to understand the process node go through when they are stopped and restarted.

A stopping node picks an online cluster member (**only disk nodes will be considered**) to sync with after restart.
Upon restart the node will try to contact that peer 10 times by default, with 30 second response timeouts.

In case the peer becomes available in that time interval, the node successfully starts, syncs what it needs from the
peer and keeps going.

If the peer does not become available, the restarted node will give up and voluntarily stop. Such condition can be
identified by the timeout (timeout_waiting_for_tables) warning messages in the logs that eventually lead to node
startup failure:

### 2.5 Change Cluster Node Type
```shell
rabbitmqctl stop_app
rabbitmqctl change_cluster_node_type {ram|disk}
rabbitmqctl start_app
rabbitmqctl cluster_status
```

### 2.6 Suspend And Resume Connections
```shell
# to stop new client connection on current node:
rabbitmqctl suspend_listeners
# resume new connection:
rabbitmqctl resume_listeners


# to stop host148 accepting new client connection, execute this command on an arbitrary node:
rabbitmqctl suspend_listeners -n rabbit@host148
# resume
rabbitmqctl resume_listeners -n rabbit@host148

```

## 3. Clustering and Observability
Client connections, channels and queues will be distributed across cluster nodes. Operators need to be able to inspect
and monitor such resources across all cluster nodes.

RabbitMQ <mark>CLI</mark> tools such as <mark>rabbitmq-diagnostics</mark> and <mark>rabbitmqctl</mark> provide commands
that inspect resources and cluster-wide state.

Commands focus on the state of a single node, for example:
>- rabbitmqctl-diagnostics environment
>- rabbitmqctl-diagnostics status

Commands focus on the state of cluster-wide, for example:
>- rabbitmqctl list_connections
>- rabbitmqctl list_mqtt_connections
>- rabbitmqctl list_stomp_connections
>- rabbitmqctl list_users
>- rabbitmqctl list_vhosts

Such "cluster-wide" commands will often contact one node first, discover cluster members and contact them all to
retrieve and combine their respective state.

Assuming a non-changing state of the cluster (e.g. no connections are closed or opened), two CLI commands executed
against two **different nodes** one after another will **produce identical or semantically identical** results.


<mark>Management UI</mark> works similarly: a node that has to respond to an HTTP API request will <mark>fan out</mark>
to other cluster members and aggregate their responses. In a cluster with multiple nodes that have management plugin
enabled, the operator can use any node to access management UI. The same goes for monitoring tools that use the HTTP
API to collect data about the state of the cluster. There is no need to issue a request to every cluster node in turn.

### 3.1 Health Checks
```shell
# will exit with an error for the nodes that are currently waiting for
# a peer to sync schema tables from
rabbitmq-diagnostics check_running
```
```shell
# a very basic check that will succeed for the nodes that are currently waiting for
# a peer to sync schema from
rabbitmq-diagnostics ping
```

## Node Failure Handling
### Broker Dimension
RabbitMQ brokers tolerate the failure of individual nodes. Nodes can be started and stopped at will, as long as they
can contact a cluster member node known at the time of shutdown.
>MQ集群能容忍单个节点的故障，节点可以随时开始运行或停止，只要重启时能与关闭服务时已知的成员列表中的一个进行通信即可，严格来说，应该是Disk类型的成员节点

### Queue Dimension
<mark>Quorum queue</mark> allows queue contents to be replicated across multiple cluster nodes with parallel replication and a
predictable <mark>leader election</mark> and data safety behavior as long as a majority of replicas are online.

Non-replicated classic queues can also be used in clusters. Non-mirrored queue behaviour in case of node failure depends on queue durability.

RabbitMQ clustering has several modes of dealing with [network partitions](https://www.rabbitmq.com/partitions.html),
primarily consistency oriented. Clustering is meant to be used across LAN. It is not recommended to run clusters that
span WAN. The <mark>Shovel</mark> or <mark>Federation</mark> plugins are better solutions for connecting brokers across
a WAN. Note that Shovel and Federation are not equivalent to clustering.

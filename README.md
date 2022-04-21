# awesomeGolang
For personal practice/learn purpose

# Commits Specification
All the commit messages of this project should comply with [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/).
And here are some special commit types in this codebase:
- <mark>tutorial</mark>: coding follow the official tutorial
- <mark>leetcode</mark>: leetcode algorithm stuffs
- <mark>demo</mark>: for personal learning habit, whenever after finishing some tutorials, I'll try to apply them to some
practical scenes.

# Naming Specification


# Coding Specification


## Project Layout
```
├── Makefile
├── README.md
├── README_zh.md
├── build
│   └── Dockerfile
├── concurrency
│   ├── pipeline.go
│   ├── rob_pike.go
│   └── xcx.go
├── doc
├── go.mod
├── go.sum
├── leetcode
│   └── leetcode.go
├── main.go
├── message_broker
│   └── rabbitmq
│       ├── tutorial2-Work-Queues
│       │   ├── new_task.go
│       │   ├── rabbit.png
│       │   ├── tutorial2.md
│       │   └── worker.go
│       ├── tutorial3-Publish-Subscribe
│       │   ├── emit_log.go
│       │   ├── images
│       │   │   ├── bindings.png
│       │   │   ├── exchange_type.png
│       │   │   ├── exchanges.png
│       │   │   ├── logs_exchange_create.png
│       │   │   └── python-three-overall.png
│       │   ├── receive_logs.go
│       │   └── tutorial3.md
│       ├── tutorial4-Routing
│       │   ├── emit_log_direct.go
│       │   ├── images
│       │   │   ├── direct-exchange-multiple.png
│       │   │   ├── direct-exchange.png
│       │   │   └── python-four.png
│       │   ├── receive_logs_direct.go
│       │   └── tutorial4.md
│       ├── tutorial5-Topics
│       │   ├── emit_log_topic.go
│       │   ├── images
│       │   │   └── python-five.png
│       │   ├── receive_logs_topic.go
│       │   └── topics.md
│       ├── tutorial6-RPC
│       │   ├── rpc-tutorial.png
│       │   ├── rpc.md
│       │   ├── rpc_client.go
│       │   └── rpc_server.go
│       └── tutorial7-Publisher-Confirms
│           ├── publisher_confirm.go
│           └── tutorial7.md
├── non-classified
│   ├── learn_generics.go
│   ├── learn_hash.go
│   ├── learn_linkname.go
│   └── learn_namespace.go
├── third_party
│   └── third_party.go
└── util
    ├── util.go
    └── whatever.s

```

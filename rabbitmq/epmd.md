## What is EPMD and How is it used?
EPMD(Erlang Port Mapping Daemon) a daemon to discover what ports a rabbitMQ node is listened on for
inter-node communication and CLI tools.

When a node or a CLI tool needs to talk to rabbit@hostname1, it will do the following:
* Resolve hostname1 to an IPv4 or IPv6 address using DNS way or some other ways
* Contact empd daemon running on hostname1 node using the above address
* Ask empd for the port used by node rabbit daemon on it
* Contact to the node using the above address and port
* Proceed the communication

## EPMD Interface

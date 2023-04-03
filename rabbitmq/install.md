### setup single node env
```shell
yum update
yum install -y epel-release

# get erlang-solutions
wget http://packages.erlang-solutions.com/erlang-solutions-1.0-1.noarch.rpm
rpm -Uvh erlang-solutions-1.0-1.noarch.rpm

# install erlang
yum install -y erlang

# if success, ought to see these:
#[root@host148 ~]# erl
#Erlang/OTP 24 [erts-12.2.1] [source] [64-bit] [smp:1:1] [ds:1:1:10] [async-threads:1]
#
#Eshell V12.2.1  (abort with ^G)
#1>
#(double type: "Ctrl + C" to quit)


# install rabbitmq-server
# download rpm package on https://www.rabbitmq.com/install-rpm.html#downloads
yum install -y rabbitmq-server-3.9.15-1.el8.noarch.rpm
systemctl enable rabbitmq-server
systemctl start rabbitmq-server

# clean up
rm -f erlang-solutions-1.0-1.noarch.rpm  rabbitmq-server-3.9.15-1.el8.noarch.rpm

# add firewall policy
# If confused with the ports below, see https://www.rabbitmq.com/networking.html#ports "Port Access" section for more detail
firewall-cmd --zone=public --permanent --add-port=4369/tcp --add-port=25672/tcp --add-port=5671-5672/tcp --add-port=15672/tcp  --add-port=61613-61614/tcp --add-port=1883/tcp --add-port=8883/tcp
firewall-cmd --reload

#delete guest and create admin
rabbitmqctl delete_user guest
rabbitmqctl add_user admin
rabbitmqctl set_user_tags admin administrator
rabbitmqctl set_permissions -p / admin ".*" ".*" ".*"


# optional to enable management plugin
rabbitmq-plugins enable rabbitmq_management

```

### Setup node name
```shell
hostname rabbit1.svc.local
# append to /etc/hosts:
127.0.0.1 rabbit1.svc.local
192.168.1.1 rabbit1.svc.local
192.168.1.2 rabbit2.svc.local
192.168.1.3 rabbit3.svc.local
```

package intelligentnetwork

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

var (
	serverName = flag.String("name", "", "")
)

func getEtcdClient() *clientv3.Client {
	// Etcd 服务器地址
	endpoints := []string{"127.0.0.1:2379"}
	clientConfig := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 2 * time.Second,
	}
	cli, err := clientv3.New(clientConfig)
	if err != nil {
		panic(err)
	}
	return cli
}

func main() {
	flag.Parse()

	client := getEtcdClient()
	// got a lease session
	leaseSession, err := concurrency.NewSession(client, concurrency.WithTTL(5))
	if err != nil {
		panic(err)
	}
	fmt.Println("session lessId is ", leaseSession.Lease())

	// go-client has already integrated the election api, owesome!
	election := concurrency.NewElection(leaseSession, "my-election")

	// [keep campaigning]
	go func() {
		// Campaign 做的事情很简单，带条件
		if err := election.Campaign(context.Background(), *serverName); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("campaign completed")
		}
	}()

	// [observe and get master name]
	var masterName string
	go func() {
		ctx, cancel := context.WithCancel(context.TODO())
		defer cancel()
		timer := time.NewTicker(time.Second)
		for range timer.C {
			timer.Reset(time.Second)
			resp := <-election.Observe(ctx)
			if len(resp.Kvs) > 0 {
				// 查看当前谁是 master
				masterName = string(resp.Kvs[0].Value)
				fmt.Println("get master with:", masterName)
			}
		}

	}()

	// [perform master job]
	go func() {
		timer := time.NewTicker(10 * time.Second)
		for range timer.C {
			if masterName == *serverName {
				fmt.Println("oh, i'm master, doing my job")
			} else {
				fmt.Println("I'll do nothing cause I'm not the leader")
			}
		}
	}()

	c := make(chan os.Signal, 1)
	// 接收 Ctrl C 中断
	signal.Notify(c, syscall.SIGHUP, syscall.SIGTERM)

	s := <-c
	fmt.Println("Got signal:", s)
	election.Resign(context.TODO())
}

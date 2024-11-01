package main

import (
	"context"
	"flag"
	"fmt"
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
	// 获取一个租期为60秒的租约，会话后台会自动续约，所以可以认为只要没有主动关闭会话，租约永远不会过期
	leaseSession, err := concurrency.NewSession(client, concurrency.WithTTL(5))
	if err != nil {
		panic(err)
	}
	defer leaseSession.Close()
	fmt.Println("session lessId is ", leaseSession.Lease())

	// go-client has already integrated the election api, owesome!
	election := concurrency.NewElection(leaseSession, "task")

	// [keep campaigning]
	// go func() {
	// 	// 会一致阻塞，直至连接出错或成为了leader
	// 	err := election.Campaign(context.Background(), *serverName)
	// 	if err != nil {
	// 		fmt.Println(err) // 出错了
	// 	} else {
	// 		fmt.Println("campaign completed") // 成为了leader
	// 	}
	// }()

	// 保持对当前leader的观察，这里单纯为了练习使用Observe API而用
	// 生产上的选举可以不用该方法，当然如果想在非leader的实例上定时打印
	// 当前leader是谁，可以这么做
	// var masterName string
	// go func() {
	// 	ctx, cancel := context.WithCancel(context.TODO())
	// 	defer cancel()
	// 	timer := time.NewTicker(time.Second)
	// 	for range timer.C {
	// 		timer.Reset(time.Second)
	// 		resp := <-election.Observe(ctx)
	// 		if len(resp.Kvs) > 0 {
	// 			// 查看当前谁是 master
	// 			masterName = string(resp.Kvs[0].Value)
	// 			fmt.Println("get master with:", masterName)
	// 		}
	// 	}

	// }()

	// [perform master job]
	// go func() {
	// 	timer := time.NewTicker(10 * time.Second)
	// 	for range timer.C {
	// 		if masterName == *serverName {
	// 			fmt.Println("oh, i'm master, doing my job")
	// 		} else {
	// 			fmt.Println("I'll do nothing cause I'm not the leader")
	// 		}
	// 	}
	// }()

	err = election.Campaign(context.Background(), *serverName)
	if err != nil {
		panic(err) // 出错了
	}

	// fmt.Println("租约ID: ", leaseSession.Lease())
	fmt.Println("我是领导人: ", *serverName)
	// c := make(chan os.Signal, 1)
	// // 接收 Ctrl C 中断
	// signal.Notify(c, syscall.SIGHUP, syscall.SIGTERM)

	// s := <-c
	// fmt.Println("Got signal:", s)
	time.Sleep(time.Second * 300)
	election.Resign(context.TODO())
}

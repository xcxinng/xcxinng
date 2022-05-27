package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-zookeeper/zk"
)

var (
	conn         *zk.Conn
	path         = flag.String("path", "/fw_lock", "lock path")
	clientID     = flag.Int("client_id", -1, "client id")
	jobTimeSpend = flag.Int("time_spend", 15, "job time sleep in second")
)

func main() {
	flag.Parse()
	if *clientID == -1 {
		panic("client id not found")
	}
	if !strings.HasPrefix(*path, "/fw_lock") {
		panic("path should has prefix '/fw_lock' and must be a child node of it")
	}

	logger := log.New(os.Stdout, fmt.Sprintf("client[%d]", *clientID), log.Ldate|log.Ltime)
	var err error
	conn, _, err = zk.Connect([]string{"127.0.0.1"}, time.Second) //*10)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	lock := zk.NewLock(conn, *path, zk.WorldACL(zk.PermAll))
	if err = lock.Lock(); err != nil {
		logger.Println("lock error:", err)
		return
	}
	logger.Println("lock success")
	time.Sleep(time.Second * (time.Duration(*jobTimeSpend)))
	if err = lock.Unlock(); err != nil {
		logger.Println("unlock error:", err)
		return
	}
	logger.Println("unlock success")
}

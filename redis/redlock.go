package main

import (
	"fmt"
	"time"

	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

var (
	rs *redsync.Redsync
)

func init() {
	// Create a pool with go-redis (or redigo) which is the pool redisync will
	// use while communicating with Redis. This can also be any pool that
	// implements the `redis.Pool` interface.
	pool := goredis.NewPool(goredislib.NewClient(&goredislib.Options{
		Addr: "host148:6379",
	}))

	pool2 := goredis.NewPool(goredislib.NewClient(&goredislib.Options{
		Addr: "host146:6379",
	}))

	pool3 := goredis.NewPool(goredislib.NewClient(&goredislib.Options{
		Addr: "host147:6379",
	}))

	// Create an instance of redisync to be used to obtain a mutual exclusion
	// lock.
	rs = redsync.New(pool, pool2, pool3)
}

func main() {
	// Obtain a new mutex by using the same name for all instances wanting the
	// same lock.
	mutexname := "my-global-mutex"
	mutex := rs.NewMutex(mutexname)

	// Obtain a lock for our given mutex. After this is successful, no one else
	// can obtain the same lock (the same mutex name) until we unlock it.
	if err := mutex.Lock(); err != nil {
		panic(err)
	}

	// Do your work that requires the lock.
	time.Sleep(time.Second * 2)
	fmt.Println("success")

	// Release the lock so other processes or threads can obtain a lock.
	if ok, err := mutex.Unlock(); !ok || err != nil {
		panic("unlock failed")
	}
}

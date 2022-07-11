// Consumer consumes data from kafka partition.
//
// # Setup Postgresql
//
// Type sql command below to create the table:
//
//     create table device_abnormal_event_record (
//     id serial primary key,
//     sysname varchar(32) not null,
//     ip varchar(64),
//     metric varchar(64) not null,
//     tag varchar(64) not null,
//     start_at timestamp not null,
//     recover_at timestamp,
//     );
//
// # Setup Redis
//
//  1. Append to redis config file: `notify-keyspace-events "Ex"`.
//  2. Watch key expired notification.
//
// The complete data flow is described in [data flow].
//
// [data flow]: ./self-monitor.drawio
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"xorm.io/xorm"
)

var (
	redisClient = redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "123456"})
	ctx         = context.TODO()
	pgClient    *xorm.Engine
)

func init() {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"localhost", "5432", "postgres", "123456", "postgres")
	client, err := xorm.NewEngine("postgres", connStr)
	if err != nil {
		logrus.Fatal("get xorm engine error:%v", err)
	}
	pgClient = client
}

// abnormal event metric name lists.
const (
	unreachable = "network.event.ping.unreachable"
	oidUnknown  = "network.event.oid.timeout"
	oidTimeout  = "network.event.oid.unknown" // also mean "no such object"
)

const (
	abnormalEventPrefix = "abnormalEvent"
)

// DeviceAbnormalEventRecord implements names.TableName
type DeviceAbnormalEventRecord struct {
	Id        int64  `xorm:"id pk autoincr serial"`
	Sysname   string `xorm:"sysname varchar(32) notnull"`
	Ip        string `xorm:"ip varchar(64)"`
	Metric    string `xorm:"metric  varchar(64) notnull"`
	Tag       string `xorm:"tag notnull"`
	StartAt   int64  `xorm:"start_at notnull"` // in int64-timestamp format
	RecoverAt int64  `xorm:"recover_at "`      // in int64-timestamp format
}

func (d DeviceAbnormalEventRecord) TableName() string {
	return "device_abnormal_event_record" // device_abnormal_event_record
}

// ConsumerMetric is exactly the same as Metric, ensure consumer can compile individually.
type ConsumerMetric struct {
	Endpoint    string      `json:"endpoint,omitempty"`
	Metric      string      `json:"metric,omitempty"`
	Value       interface{} `json:"value,omitempty"`
	Step        int         `json:"step,omitempty"`
	CounterType string      `json:"counter_type,omitempty"`
	Tags        string      `json:"tags,omitempty"`
	Timestamp   int64       `json:"timestamp,omitempty"`
}

func handleMessage(message *sarama.ConsumerMessage) {
	logrus.Info("begin to handle message")
	metricData := &ConsumerMetric{}
	err := json.Unmarshal(message.Value, metricData)
	if err != nil {
		logrus.Errorf("unmarshal message error:%v", err)
		return
	}

	// to get a redis lock
	var tmp = []string{"", metricData.Endpoint, metricData.Metric, metricData.Tags}
	tmp[0] = "lock"
	lockKey := strings.Join(tmp, "|")
	logrus.Debugf("lock key:%s", lockKey)
	if !getLock(lockKey) {
		logrus.Warning("lock failed,exited")
		redisClient.Del(ctx, lockKey)
		return
	} else {
		logrus.Infof("lock %s success", lockKey)
	}

	// get lock success
	defer redisClient.Del(ctx, lockKey)
	tmp[0] = abnormalEventPrefix
	eventKey := strings.Join(tmp, "|")
	result, _ := redisClient.Exists(ctx, eventKey).Result()
	if result == 1 { // such an abnormal event had been happened, reset event's expire time
		cmd := redisClient.Expire(ctx, eventKey, time.Minute*10)
		if cmd.Err() != nil {
			logrus.Errorf("set key %s expire time error:%v", eventKey, cmd.Err())
		}
		redisClient.Del(ctx, lockKey)
		return
	}

	// new abnormal event had occurred, write into both PG and redis.
	record := DeviceAbnormalEventRecord{
		Sysname: metricData.Endpoint,
		Metric:  metricData.Metric,
		Tag:     metricData.Tags,
		StartAt: metricData.Timestamp,
	}
	_, err = pgClient.Insert(&record)
	if err != nil {
		logrus.Errorf("insert pg failed,error:%v,data:%+v", err, record)
	} else {
		logrus.Infof("insert data:%v", record)
	}

	data, err := json.Marshal(record)
	if err != nil {
		logrus.Errorf("marshal error:%v", err)
	}
	err = insertRedis(eventKey, data)
	if err != nil {
		logrus.Errorf("insert redis failed,error:%v", err)
	} else {
		logrus.Infof("insert redis success")
	}
}

func getLock(key string) bool {
	lock, err := redisClient.SetNX(ctx, key, 1, time.Second*10).Result()
	if lock && err == nil {
		return true
	}
	if err != nil {
		logrus.Errorf("key %s set lock error:%v", key, err)
	}
	return false
}

func unLock(key string) error {
	count, err := redisClient.Del(ctx, key).Result()
	if count == 1 && err == nil {
		return nil
	}
	return fmt.Errorf("key %s set lock error:%v", key, err)
}

func insertRedis(key string, value []byte) error {
	if result := redisClient.Set(ctx, key, value, time.Minute*10); result.Err() != nil {
		//logrus.Errorf("set abnormal redis key %s error:%v",key,result.Err())
		return result.Err()
	}
	return nil
}

func handleExpireKeyEvent() {
	channels := redisClient.Subscribe(ctx, "__keyevent@0__:expired").Channel()
	for channel := range channels {
		recoverTime := time.Now().Unix()
		eventKey := channel.Payload
		logrus.Infof("receive expired key:%s", eventKey)
		if !strings.HasPrefix(eventKey, abnormalEventPrefix) {
			continue
		}

		data := strings.Split(eventKey, "|")
		if len(data) != 4 {
			continue
		}
		lockKey := "recover|" + eventKey
		if !getLock(lockKey) {
			continue
		}

		condition := DeviceAbnormalEventRecord{
			Sysname: data[1],
			Metric:  data[2],
			Tag:     data[3],
		}
		record := DeviceAbnormalEventRecord{RecoverAt: recoverTime}
		_, err := pgClient.Update(&record, &condition)
		if err != nil {
			logrus.Errorf("update device record error:%v", err)
		} else {
			logrus.Infof("update device %s success", data[0])
		}
		if err = unLock(lockKey); err != nil {
			logrus.Errorf("unlock key: %s error:%v", lockKey, err)
		}
	}
}

func main() {
	if redisClient == nil {
		logrus.Fatal("redis client was nil")
	}

	if pgClient == nil {
		logrus.Fatal("pg client was nil")
	}

	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, sarama.NewConfig())
	if err != nil {
		logrus.Fatal(err.Error())
	}
	defer consumer.Close()

	go handleExpireKeyEvent()

	partitionConsumer, err := consumer.ConsumePartition("self_monitoring_metric", 0, sarama.OffsetNewest)
	if err != nil {
		logrus.Fatal(err)
	}
	defer partitionConsumer.Close()

	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	consumed := 0
ConsumerLoop:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			handleMessage(msg)
			logrus.Printf("Consumed message offset %d\n", msg.Offset)
			consumed++
		case <-signals:
			break ConsumerLoop
		}
	}
}

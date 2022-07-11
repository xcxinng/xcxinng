// producer emulate to produce abnormal event message to kafka partition.
package main

import (
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
	"unsafe"

	"github.com/Shopify/sarama"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"

	"awesomeGolang/action_abnormal_event/common"
)

func cronJob() {
	brokers := []string{"localhost:9092"}
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		logrus.Fatal("get producer error:", err)
	}
	defer producer.Close()
	metrics := common.OpenFalconMetrics()
	msgs := make([]*sarama.ProducerMessage, len(metrics))
	for i := 0; i < len(metrics); i++ {
		value, err := json.Marshal(metrics[i])
		if err != nil {
			logrus.Println("marshal error:", err)
			continue
		}

		msgs[i] = &sarama.ProducerMessage{
			Topic: "self_monitoring_metric",
			Key:   sarama.StringEncoder(metrics[i].Endpoint),
			Value: sarama.StringEncoder(*(*string)(unsafe.Pointer(&value))),
		}
	}

	if errs := producer.SendMessages(msgs); err != nil {
		for _, err := range errs.(sarama.ProducerErrors) {
			logrus.Println("Write to kafka failed: ", err)
		}
		logrus.Println("write message error:", err)
	} else {
		logrus.Printf("produce %d messages to kafka success\n", len(msgs))
	}
}

func main() {
	logrus.Println("cron job starting...")
	c := cron.New()
	id, err := c.AddFunc("@every 180s", cronJob)
	if err != nil {
		logrus.Fatal(err)
	}
	c.Start()
	logrus.Printf("next cron job runs at:%v,please wait...\n", c.Entry(id).Next)

	chs := make(chan os.Signal)
	signal.Notify(chs, syscall.SIGKILL, syscall.SIGTERM)
	for true {
		select {
		case sig := <-chs:
			logrus.Println("receive signal:", sig, ", program exiting...")
			c.Stop()
			return
		}
	}
}

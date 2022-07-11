// producer emulate to produce abnormal event message to kafka partition.
package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
	"unsafe"

	"github.com/Shopify/sarama"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

func producerRandomInt() int {
	rand.Seed(time.Now().UnixNano())
	ipHostField := rand.Intn(254) + 1
	return ipHostField
}

func producerRandomCollectorIP() string {
	return fmt.Sprintf("100.1.1.%d", producerRandomInt())
}

func producerRandomDeviceIP() string {
	return fmt.Sprintf("30.3.3.%d", producerRandomInt())
}

var falconMetrics = []string{"cpu_utilization_ratio", "mem_utilization_ratio", "power_status", "fan_status", "in_bytes_total"}

func producerRandomMetric() string {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(falconMetrics))
	return falconMetrics[n]
}

// Metric is a metric defined by open falcon.
type Metric struct {
	Endpoint    string      `json:"endpoint,omitempty"`
	Metric      string      `json:"metric,omitempty"`
	Value       interface{} `json:"value,omitempty"`
	Step        int         `json:"step,omitempty"`
	CounterType string      `json:"counter_type,omitempty"`
	Tags        string      `json:"tags,omitempty"`
	Timestamp   int64       `json:"timestamp,omitempty"`
}

func genOpenFalconMetrics() []Metric {
	agentIP := producerRandomCollectorIP()
	metrics := make([]Metric, 0, 3)
	deviceIP := producerRandomDeviceIP()
	ts := time.Now().Unix()
	metrics = append(metrics, Metric{
		Endpoint:    deviceIP,
		Metric:      "network.event.ping.unreachable",
		Value:       1,
		Step:        60,
		CounterType: "GAUGE",
		Tags:        fmt.Sprintf("agent_ip=%s,device_ip=%s", agentIP, deviceIP),
		Timestamp:   ts,
	}, Metric{
		Endpoint:    deviceIP,
		Metric:      "network.event.oid.timeout",
		Value:       1,
		Step:        60,
		CounterType: "GAUGE",
		Tags:        fmt.Sprintf("agent_ip=%s,device_ip=%s,metric=%s", agentIP, deviceIP, producerRandomMetric()),
		Timestamp:   ts,
	}, Metric{
		Endpoint:    deviceIP,
		Metric:      "network.event.oid.unknown",
		Value:       1,
		Step:        60,
		CounterType: "GAUGE",
		Tags:        fmt.Sprintf("agent_ip=%s,device_ip=%s,metric=%s", agentIP, deviceIP, producerRandomMetric()),
		Timestamp:   ts,
	})
	return metrics
}

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
	metrics := genOpenFalconMetrics()
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

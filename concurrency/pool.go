package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

var logger = log.Default()

type Metric struct {
	Name      string
	Tag       string
	Value     interface{}
	Timestamp time.Time
}

type Transfer interface {
	Put([]*Metric)
	Send([]*Metric) error
}

type MetricTransfer struct {
	mutex sync.Mutex
	data  chan *Metric
}

var defaultMetricTransfer = MetricTransfer{
	mutex: sync.Mutex{},
	data:  make(chan *Metric),
}

func (m *MetricTransfer) Put(metric *Metric) {
	m.data <- metric
}

func (m *MetricTransfer) Send() error {
	tick := time.NewTicker(time.Second)
	data := make([]*Metric, 0, 10000)

	go func() {
		for {
			select {
			case item := <-m.data:
				data = append(data, item)
				if len(data) > 500 {
					_data := data
					consume(_data)
					data = make([]*Metric, 0, 500)
				}
			case <-tick.C:
				_data := data
				consume(_data)
				data = make([]*Metric, 500)
				tick.Reset(time.Second)
			}
		}
	}()
	return nil
}

func consume(metrics []*Metric) {
	fd, err := os.OpenFile("/dev/null", os.O_WRONLY, 0)
	if err != nil {
		logger.Println("open file failed,", err)
		return
	}
	_, err = fmt.Fprint(fd, metrics)
	if err != nil {
		logger.Println("consume metrics failed,", err)
		return
	}
	logger.Println("consume metrics count:", len(metrics))
}

func generateMetric(length int) []*Metric {
	data := make([]*Metric, length)
	for i := 0; i < length; i++ {
		data[i] = &Metric{
			Name:      "",
			Tag:       "tag=tag",
			Value:     time.Now().Unix(),
			Timestamp: time.Now(),
		}
	}
	return data
}

func ProduceMetric(totalCount int) {
	for i := 0; i < totalCount; i++ {
		for _, metric := range generateMetric(10000) {
			defaultMetricTransfer.Put(metric)
			time.Sleep(time.Second)
		}
	}
}

func test(t int) {
	go ProduceMetric(30)
	go defaultMetricTransfer.Send()
	time.Sleep(time.Duration(t) * time.Second)
}

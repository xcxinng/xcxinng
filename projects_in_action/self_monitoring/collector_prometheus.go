// collector_prometheus emulate an application to expose metric data
// in prometheus format.
//
// To ensure that each go file could be compiled successfully,
// each of them has individual functions or variables, even though
// they are exactly the same.
package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func getRandomInt() int {
	rand.Seed(time.Now().UnixNano())
	ipHostField := rand.Intn(254) + 1
	return ipHostField
}

func randomCollectorIP() string {
	return fmt.Sprintf("100.1.1.%d", getRandomInt())
}

func randomDeviceIP() string {
	return fmt.Sprintf("30.3.3.%d", getRandomInt())
}

func genRandomMetric() string {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(metrics))
	return metrics[n]
}

var metrics = []string{"cpu_utilization_ratio", "mem_utilization_ratio", "power_status", "fan_status", "in_bytes_total"}

// AbnormalEvent implements prometheus.Collector.
type AbnormalEvent struct {
	oidTimeout  *prometheus.Desc
	pingTimeout *prometheus.Desc
	oidUnknown  *prometheus.Desc
}

func (a AbnormalEvent) Describe(desc chan<- *prometheus.Desc) {
	desc <- a.oidUnknown
	desc <- a.pingTimeout
	desc <- a.oidTimeout
}

func (a AbnormalEvent) Collect(metrics chan<- prometheus.Metric) {
	collector := randomCollectorIP()
	metrics <- prometheus.MustNewConstMetric(a.oidUnknown, prometheus.GaugeValue, 1, randomDeviceIP(), collector, genRandomMetric())
	metrics <- prometheus.MustNewConstMetric(a.pingTimeout, prometheus.GaugeValue, 1, randomDeviceIP())
	metrics <- prometheus.MustNewConstMetric(a.oidTimeout, prometheus.GaugeValue, 1, randomDeviceIP(), collector, genRandomMetric())
}

// NewCollector return a prometheus.Collector implemented by AbnormalEvent.
func NewCollector() prometheus.Collector {
	tags := []string{"device_ip", "collector_ip", "metric"}
	namespace := "network_event_"
	return AbnormalEvent{
		oidTimeout:  prometheus.NewDesc(namespace+"oid_timeout", "oid timeout event", tags, nil),
		pingTimeout: prometheus.NewDesc(namespace+"device_unreachable", "ping device timeout", []string{"device_ip"}, nil),
		oidUnknown:  prometheus.NewDesc(namespace+"oid_unknown", "oid unknown", tags, nil),
	}
}

func runServer() {
	// set up url and corresponding handler
	http.Handle("/metrics", promhttp.Handler())

	// create http server and listen on the specific tcp port
	log.Fatal(http.ListenAndServe(":2112", nil))
}

func main() {
	prometheus.MustRegister(NewCollector())
	// unregister the default collector in order to keep output cleaner
	prometheus.Unregister(collectors.NewGoCollector())
	prometheus.Unregister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	go runServer()

	chs := make(chan os.Signal)
	signal.Notify(chs, syscall.SIGKILL, syscall.SIGTERM)
	for true {
		select {
		case sig := <-chs:
			log.Println("receive signal:", sig, ", program exiting...")
			return
		}
	}
}

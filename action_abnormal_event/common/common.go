package common

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func randomInt(max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max)
}

func DefaultRandInt() int {
	return randomInt(254) + 1
}

func RandAgentIP() string {
	return fmt.Sprintf("100.1.1.%d", DefaultRandInt())
}

func RandDeviceIP() string {
	return fmt.Sprintf("30.3.3.%d", DefaultRandInt())
}

var (
	falconMetrics = []string{"cpu_utilization_ratio", "mem_utilization_ratio", "power_status",
		"fan_status", "in_bytes_total"}
)

var mutex = sync.Mutex{}

func AppendMetric(s string) {
	mutex.Lock()
	defer mutex.Unlock()
	falconMetrics = append(falconMetrics, s)
}

func RandMetric() string {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(falconMetrics))
	return falconMetrics[n]
}

// Metric represent an open-falcon metric.
type Metric struct {
	Endpoint    string      `json:"endpoint,omitempty"`
	Metric      string      `json:"metric,omitempty"`
	Value       interface{} `json:"value,omitempty"`
	Step        int         `json:"step,omitempty"`
	CounterType string      `json:"counter_type,omitempty"`
	Tags        string      `json:"tags,omitempty"`
	Timestamp   int64       `json:"timestamp,omitempty"`
}

// abnormal event metric name lists.
const (
	Unreachable = "network.event.ping.unreachable"
	OidUnknown  = "network.event.oid.timeout"
	OidTimeout  = "network.event.oid.unknown" // also mean "no such object"
)

func OpenFalconMetrics() []Metric {
	agentIP := RandAgentIP()
	metrics := make([]Metric, 0, 3)
	deviceIP := RandAgentIP()
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
		Tags:        fmt.Sprintf("agent_ip=%s,device_ip=%s,metric=%s", agentIP, deviceIP, RandMetric()),
		Timestamp:   ts,
	}, Metric{
		Endpoint:    deviceIP,
		Metric:      "network.event.oid.unknown",
		Value:       1,
		Step:        60,
		CounterType: "GAUGE",
		Tags:        fmt.Sprintf("agent_ip=%s,device_ip=%s,metric=%s", agentIP, deviceIP, RandMetric()),
		Timestamp:   ts,
	})
	return metrics
}

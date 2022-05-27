package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Processor interface {
	GetStats() (Stats, error)
}

type process struct {
	name string
	host string
}

func NewProcessor(name, host string) *process {
	return &process{
		name: name,
		host: host,
	}
}

type Stats struct {
	data         map[string]float64
	instanceName string
	hostIP       string
}

func (h process) randomValue() float64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Float64() * 1000
}

func (h process) GetStats() (Stats, error) {
	stats := make(map[string]float64)
	stats["in_bytes"] = h.randomValue()
	stats["in_bytes"] = h.randomValue()
	stats["in_bytes"] = h.randomValue()
	stats["in_bytes"] = h.randomValue()
	stats["connected"] = 1
	info := Stats{
		data:         stats,
		instanceName: h.name,
		hostIP:       h.host,
	}
	return info, nil
}

type pseudoCollector struct {
	vpn       Processor
	namespace string
	inBytes   *prometheus.Desc
	outBytes  *prometheus.Desc
	inPkts    *prometheus.Desc
	outPkts   *prometheus.Desc
	connected *prometheus.Desc
}

func (p *pseudoCollector) Describe(descs chan<- *prometheus.Desc) {
	descs <- p.inBytes
	descs <- p.outBytes
	descs <- p.inPkts
	descs <- p.outPkts
	descs <- p.connected
}

func (p *pseudoCollector) Collect(metrics chan<- prometheus.Metric) {
	res, err := p.vpn.GetStats()
	if err != nil {
		log.Fatal(err)
	}
	stats := res.data
	metrics <- prometheus.MustNewConstMetric(p.inBytes, prometheus.CounterValue, stats["in_bytes"], res.instanceName, res.hostIP)
	metrics <- prometheus.MustNewConstMetric(p.outBytes, prometheus.CounterValue, stats["out_bytes"], res.instanceName, res.hostIP)
	metrics <- prometheus.MustNewConstMetric(p.inPkts, prometheus.CounterValue, stats["in_pkts"], res.instanceName, res.hostIP)
	metrics <- prometheus.MustNewConstMetric(p.outPkts, prometheus.CounterValue, stats["out_pkts"], res.instanceName, res.hostIP)
	metrics <- prometheus.MustNewConstMetric(p.connected, prometheus.GaugeValue, stats["connected"], res.instanceName, res.hostIP)
}

func NewPseudoCollector(namespace string, vpn Processor) prometheus.Collector {
	ns := ""
	if len(namespace) > 0 {
		ns = namespace + "_"
	}

	if vpn == nil {
		vpn = new(process)
	}
	tags := []string{"ip", "nfv_host_ip"}
	return &pseudoCollector{
		vpn:       vpn,
		inBytes:   prometheus.NewDesc(ns+"in_bytes", "Total in bytes", tags, nil),
		outBytes:  prometheus.NewDesc(ns+"out_bytes", "Total out bytes", tags, nil),
		inPkts:    prometheus.NewDesc(ns+"in_pkts", "Total in packets", tags, nil),
		outPkts:   prometheus.NewDesc(ns+"out_pkts", "Total out packets", tags, nil),
		connected: prometheus.NewDesc(ns+"is_connected", "Whether is be in connected", tags, nil),
	}
}

func main() {
	// unregister the default collector in order to keep output cleaner
	prometheus.Unregister(collectors.NewGoCollector())
	prometheus.Unregister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	// register custom collector
	prometheus.MustRegister(NewPseudoCollector("vpn_gateway", NewProcessor("1.1.1.1", "100.0.0.1")))
	prometheus.MustRegister(NewPseudoCollector("nat_instance", NewProcessor("2.2.2.2", "100.113.1.3")))

	// set up url and corresponding handler
	http.Handle("/metrics", promhttp.Handler())

	// create http server and listen on the specific tcp port
	log.Fatal(http.ListenAndServe(":2112", nil))
}

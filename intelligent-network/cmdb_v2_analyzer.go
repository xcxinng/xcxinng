package intelligentnetwork

import (
	"encoding/json"
	"fmt"
)

type RoutingTable struct {
	Destination string `json:"destination"`
	Nexthop     string `json:"nexthop"`
	Interface   string `json:"interface"`
	Protocol    string `json:"protocol"`
	Cost        int    `json:"cost"`
}

// There two slices are the example data, if nothing goes wrong, the result should be as below:
//   - record of 1.1.1.1 should be deleted
//   - record of 2.2.2.2 obviously has changed its nexthop to 10.0.0.3
//     (1) the old data in PG should be deleted
//     (2) the new data from Kafka should be added
//   - record of 3.3.3.3 nothing's changed, so nothing happens
//   - record of 4.4.4.4 does not exist in PG, so it should be added
//
// In a word, four events would be generated, and two of them are deletion events
// and the others are creation events.
//
// That's to say PG would delete 2 old data and insert 2 new data.
var RoutingTableInPG = []RoutingTable{
	{Destination: "1.1.1.1/32", Nexthop: "10.0.0.2", Interface: "eth1", Protocol: "RIP", Cost: 10},
	{Destination: "2.2.2.2/32", Nexthop: "10.0.0.1", Interface: "eth2", Protocol: "RIP", Cost: 10},
	{Destination: "3.3.3.3/32", Nexthop: "10.0.0.3", Interface: "eth3", Protocol: "RIP", Cost: 10},
}
var RoutingTableFromKafka = []RoutingTable{
	{Destination: "4.4.4.4/32", Nexthop: "10.0.0.10", Interface: "eth2", Protocol: "RIP", Cost: 10},
	{Destination: "2.2.2.2/32", Nexthop: "10.0.0.3", Interface: "eth3", Protocol: "RIP", Cost: 10},
	{Destination: "3.3.3.3/32", Nexthop: "10.0.0.3", Interface: "eth3", Protocol: "RIP", Cost: 10},
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type Event struct {
	// Action indicates the event type, which should be one of ["delete", "add","notice"].
	//
	// Note: If an event is a type of "notice" event, it usually represents that all of the
	// data are exact the same, the only thing the upstream need to do is to update the updateTime.
	Action string `json:"action"`

	// EntityType represents what resource this event is related to and provides
	// the necessary information to map the entity within the MongoDB.
	EntityType string `json:"event_type"`

	// RawData only take effects when action equals to "add", representing the data
	// that need to be inserted.
	RawData interface{} `json:"data"`
}

// Time Complexity = O(n)
// Space Complexity = O(n)
func DoAnalyze() {
	var pgSumMapping = map[RoutingTable]struct{}{}
	var kafkaSumMapping = map[RoutingTable]struct{}{}
	for i := 0; i < min(len(RoutingTableInPG), len(RoutingTableFromKafka)); i++ {
		pgSumMapping[RoutingTableInPG[i]] = struct{}{}
		kafkaSumMapping[RoutingTableFromKafka[i]] = struct{}{}
	}

	if len(RoutingTableInPG) > len(RoutingTableFromKafka) {
		for i := len(RoutingTableFromKafka); i < len(RoutingTableInPG); i++ {
			pgSumMapping[RoutingTableInPG[i]] = struct{}{}
		}
	} else {
		for i := len(RoutingTableInPG); i < len(RoutingTableFromKafka); i++ {
			kafkaSumMapping[RoutingTableFromKafka[i]] = struct{}{}
		}
	}

	var events []Event
	for data := range pgSumMapping {
		// If true, it means that data already exists in PG, do nothing
		_, exist := kafkaSumMapping[data]

		if !exist { // If false, it means that data should be deleted

			// If data from kafka got updated, its hash got changed too,
			// meaning that the old one in PG should be deleted as well.
			events = append(events, Event{Action: "delete", RawData: data, EntityType: "RoutingTable"})
		}
	}

	for data := range kafkaSumMapping {
		if _, exist := pgSumMapping[data]; !exist { // added
			events = append(events, Event{Action: "add", RawData: data, EntityType: "RoutingTable"})
		}
	}
	fmt.Println(PrettyStruct(events))
}

func PrettyStruct(data interface{}) string {
	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return ""
	}
	return string(val)
}

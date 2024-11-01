package intelligentnetwork

import "fmt"

// RouteEntry represents a single route entry in the switch route table
type RouteEntry struct {
	VRF         string
	Destination string
	NextHop     string
	Interface   string
	Metric      int
	Priority    int
	Protocol    string
}

// GenerateRouteTable generates switch route table entries based on the provided length
func GenerateRouteTable(length int) []RouteEntry {
	routeTable := make([]RouteEntry, 0, length)

	for i := 0; i < length; i++ {
		entry := RouteEntry{
			VRF:         fmt.Sprintf("VRF%d", i%3),
			Destination: fmt.Sprintf("192.168.%d.%d/24", i/256, i%256),
			NextHop:     fmt.Sprintf("10.0.%d.%d", i/256, i%256),
			Interface:   fmt.Sprintf("eth%d", i%4),
			Metric:      i + 1,
			Priority:    i % 16,
			Protocol:    "static",
		}
		routeTable = append(routeTable, entry)
	}

	return routeTable
}

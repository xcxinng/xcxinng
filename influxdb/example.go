package influxdb

import (
	"context"
	"fmt"
	"log"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

func getClient() influxdb2.Client {
	token := "UJdYMGP_n5GMDwDpSnrXJ1xr0uawlkL29PZnplgxGZjSIrxsWbZQIlcHALOzwhgR5uIbTozgtJXSjcDcHbJMNA=="
	url := "http://localhost:8086"
	client := influxdb2.NewClient(url, token)
	return client
}

var (
	org    = "xianchaoxing.com"
	bucket = "test"
)

// WriteData is an example of writing data into influxdb.
func WriteData() {
	client := getClient()

	writeAPI := client.WriteAPIBlocking(org, bucket)
	for value := 0; value < 5; value++ {
		tags := map[string]string{
			"tagname1": "tagvalue1",
		}
		fields := map[string]interface{}{
			"field1": value,
		}
		point := write.NewPoint("measurement1", tags, fields, time.Now())
		time.Sleep(1 * time.Second) // separate points by 1 second

		if err := writeAPI.WritePoint(context.Background(), point); err != nil {
			log.Fatal(err)
		}
	}
}

// Query demonstrates how to perform a query using Flux statement.
func Query() {
	client := getClient()
	queryAPI := client.QueryAPI(org)
	query := `from(bucket: "test")
                |> range(start: -10m)
                |> filter(fn: (r) => r._measurement == "measurement1")`
	results, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}
	for results.Next() {
		fmt.Println(results.Record())
	}
	if err := results.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println()
	fmt.Println()

	// perform a mean aggregation
	query = `from(bucket: "test")
                |> range(start: -10m)
                |> filter(fn: (r) => r._measurement == "measurement1")
                |> mean()`
	results, err = queryAPI.Query(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}
	for results.Next() {
		fmt.Println(results.Record())
	}
	if err := results.Err(); err != nil {
		log.Fatal(err)
	}
}

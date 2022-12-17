package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/influxdata/influxdb-client-go/v2"
)

var Client MQTT.Client
var token string

func main() {
	flag.StringVar(&token, "token", "token", "The token to operate the influxdb")
	flag.Parse()

	MqttConnect()

	for true {
		read()
		time.Sleep(5 * time.Second)
	}
}

func read() {
	bucket := "stats-cpu"
	org := "mapper-cpu"
	// token := "7UTkMfW4uuFdIJHDBapqggz8mSS_UWbqm6pmDuk3HGpT3zV1Gaf2dWnFIzveOLDO2Ec86jbFmSpoZVPNtXvMzA=="
	// Store the URL of your InfluxDB instance
	url := "http://127.0.0.1:8086"
	// Create client
	client := influxdb2.NewClient(url, token)
	// Get query client
	queryAPI := client.QueryAPI(org)
	// Get QueryTableResult
	result, err := queryAPI.Query(context.Background(), fmt.Sprintf(`from(bucket:"%s")|> range(start: -1h) |> filter(fn: (r) => r._measurement == "cpu")`, bucket))
	if err != nil {
		panic(err)
	}

	var time time.Time
	var value interface{}
	// Iterate over query response
	for result.Next() {
		//// Notice when group key has changed
		//if result.TableChanged() {
		//	fmt.Printf("table: %s\n", result.TableMetadata().String())
		//}
		// Access data
		//fmt.Printf("time: %v\n", result.Record().Time())
		//fmt.Printf("value: %v\n", result.Record().Value())
		time = result.Record().Time()
		value = result.Record().Value()
		//fmt.Printf("field: %v\n", result.Record().Field())
	}
	fmt.Printf("time: %v\n", time)
	fmt.Printf("value: %v\n", value)

	PublishMqtt("default/test", fmt.Sprintf("time: %v, cpu_stats: %v", time, value))

	// Check for an error
	if result.Err() != nil {
		fmt.Printf("query parsing error: %s\n", result.Err().Error())
	}

	// Ensures background processes finishes
	client.Close()
}

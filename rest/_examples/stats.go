package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	api "gopkg.in/ns1/ns1-go.v2/rest"
)

var client *api.Client

const (
	zoneName   = "example.com"
	recordName = "test.example.com"
	recordType = "A"
)

// Helper that initializes rest api client from environment variable.
func init() {
	k := os.Getenv("NS1_APIKEY")
	if k == "" {
		log.Fatal("NS1_APIKEY environment variable is not set, giving up!")
	}

	httpClient := &http.Client{Timeout: time.Second * 10}
	// Adds logging to each http request.
	doer := api.Decorate(httpClient, api.Logging(log.New(os.Stdout, "", log.LstdFlags)))
	client = api.NewClient(doer, api.SetAPIKey(k))
}

func main() {
	qps, _, err := client.Stats.GetQPS()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Account qps: %f\n", qps)

	qps, _, err = client.Stats.GetZoneQPS(zoneName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Zone qps: %f\n", qps)

	qps, _, err = client.Stats.GetRecordQPS(zoneName, recordName, recordType)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Record qps: %f\n", qps)
}

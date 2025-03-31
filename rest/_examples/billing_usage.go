package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	api "gopkg.in/ns1/ns1-go.v2/rest"
)

var client *api.Client

// Helper that initializes rest api client from environment variable.
func init() {
	k := os.Getenv("NS1_APIKEY")
	if k == "" {
		fmt.Println("NS1_APIKEY environment variable is not set, giving up")
	}

	httpClient := &http.Client{Timeout: time.Second * 30}
	// Adds logging to each http request.
	doer := api.Decorate(httpClient, api.Logging(log.New(os.Stdout, "", log.LstdFlags)))
	client = api.NewClient(doer, api.SetAPIKey(k))
}

func prettyPrint(header string, v interface{}) {
	fmt.Println(header)
	jsonBytes, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	fmt.Println(string(jsonBytes))
}

func main() {

	queries, _, err := client.BillingUsage.GetQueries(1738368000, 1740787199)
	if err != nil {
		log.Fatal(err)
	}
	prettyPrint("quiries", queries)

	decisions, _, err := client.BillingUsage.GetDecisions(1738368000, 1740787199)
	if err != nil {
		log.Fatal(err)
	}
	prettyPrint("decisions", decisions)

	limits, _, err := client.BillingUsage.GetLimits(1738368000, 1740787199)
	if err != nil {
		log.Fatal(err)
	}
	prettyPrint("limits", limits)

	monitors, _, err := client.BillingUsage.GetMonitors()
	if err != nil {
		log.Fatal(err)
	}
	prettyPrint("monitors", monitors)

	filterChains, _, err := client.BillingUsage.GetFilterChains()
	if err != nil {
		log.Fatal(err)
	}
	prettyPrint("filterChains", filterChains)

	records, _, err := client.BillingUsage.GetRecords()
	if err != nil {
		log.Fatal(err)
	}
	prettyPrint("records", records)

}

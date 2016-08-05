package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	ns1 "github.com/ns1/ns1-go"
	"github.com/ns1/ns1-go/rest"
)

func main() {
	k := os.Getenv("NS1_APIKEY")
	if k == "" {
		fmt.Println("NS1_APIKEY environment variable is not set, giving up")
	}

	httpClient := &http.Client{Timeout: time.Second * 10}
	client := rest.NewAPIClient(httpClient, rest.SetAPIKey(k))

	mjl, err := client.Monitors.List()
	if err != nil {
		log.Fatal(err)
	}
	for _, mj := range mjl {
		b, _ := json.MarshalIndent(mj, "", "  ")
		fmt.Println(string(b))
	}

	mj, err := client.Monitors.Get("52a99e559faa7fa6e546b9ca")
	if err != nil {
		log.Fatal(err)
	}
	b, _ := json.MarshalIndent(mj, "", "  ")
	fmt.Println(string(b))

	new_mj := &ns1.MonitoringJob{Config: map[string]interface{}{"host": "1.1.1.1"},
		Rules: []*ns1.MonitoringJobRule{
			&ns1.MonitoringJobRule{Key: "loss", Value: "5.0", Comparison: "<="},
		},
		JobType:      "ping",
		Regions:      []string{"master"},
		Active:       true,
		Frequency:    1,
		Policy:       "quorum",
		RegionScope:  "fixed",
		Name:         "Test monitor",
		NotifyRepeat: 0,
		NotifyDelay:  0,
	}

	err = client.Monitors.Create(new_mj)
	if err != nil {
		log.Fatal(err)
	}
	b, _ = json.MarshalIndent(new_mj, "", "  ")
	fmt.Println(string(b))

	new_mj.Frequency = 5
	err = client.Monitors.Update(new_mj)
	if err != nil {
		log.Fatal(err)
	}
	b, _ = json.MarshalIndent(new_mj, "", "  ")
	fmt.Println(string(b))

	err = client.Monitors.Delete(new_mj.ID)
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ns1/ns1-go/model/monitoring"
	"github.com/ns1/ns1-go/rest"
)

func main() {
	k := os.Getenv("NS1_APIKEY")
	if k == "" {
		fmt.Println("NS1_APIKEY environment variable is not set, giving up")
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 10,
	}

	// Adds logging to each http request.
	doer := rest.Decorate(
		httpClient, rest.Logging(log.New(os.Stdout, "", log.LstdFlags)))

	client := rest.NewClient(
		doer, rest.SetAPIKey(k), rest.SetEndpoint("https://api.dev.nsone.co/v1/"))

	mjl, err := client.Jobs.List()
	if err != nil {
		log.Fatal(err)
	}
	for _, mj := range mjl {
		b, _ := json.MarshalIndent(mj, "", "  ")
		fmt.Println(string(b))
	}

	mj, err := client.Jobs.Get("52a90e559faa7fa6e546b9ca")
	if err != nil {
		log.Fatal(err)
	}
	b, _ := json.MarshalIndent(mj, "", "  ")
	fmt.Println(string(b))

	new_mj := &monitoring.Job{Config: map[string]interface{}{"host": "1.1.1.1"},
		Rules: []*monitoring.JobRule{
			&monitoring.JobRule{Key: "loss", Value: "5.0", Comparison: "<="},
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

	err = client.Jobs.Create(new_mj)
	if err != nil {
		log.Fatal(err)
	}
	b, _ = json.MarshalIndent(new_mj, "", "  ")
	fmt.Println(string(b))

	new_mj.Frequency = 5
	err = client.Jobs.Update(new_mj)
	if err != nil {
		log.Fatal(err)
	}
	b, _ = json.MarshalIndent(new_mj, "", "  ")
	fmt.Println(string(b))

	err = client.Jobs.Delete(new_mj.ID)
	if err != nil {
		log.Fatal(err)
	}
}

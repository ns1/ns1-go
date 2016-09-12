package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ns1/ns1-go/model/dns"
	"github.com/ns1/ns1-go/model/dns/data"
	"github.com/ns1/ns1-go/model/dns/filter"
	api "github.com/ns1/ns1-go/rest"
)

// Helper that initializes rest api client from environment variable.
func initClient() *api.Client {
	k := os.Getenv("NS1_APIKEY")
	if k == "" {
		fmt.Println("NS1_APIKEY environment variable is not set, giving up")
	}

	httpClient := &http.Client{Timeout: time.Second * 10}
	// Adds logging to each http request.
	doer := api.Decorate(httpClient, api.Logging(log.New(os.Stdout, "", log.LstdFlags)))
	client := api.NewClient(doer, api.SetAPIKey(k))

	return client
}

// Helper that creates zone from given domain, and ignores error if zone already exists.
func setupZone(client *api.Client, domain string) error {
	z := dns.NewZone(domain)
	_, err := client.Zones.Create(z)
	if err != nil {
		// Ignore if zone already exists
		if err == api.ErrZoneExists {
			log.Printf("[Zone Create]Zone Exists %s: %s \n", z, err)
		} else {
			return err
		}
	}

	return nil
}

func prettyPrint(header string, v interface{}) {
	fmt.Println(header)
	fmt.Printf("%#v \n", v)
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}

// Define the type of update data we will send to the datasource.
// This isnt required, but meant for clarity.
type update map[string]data.Meta

func main() {
	// Initialize the rest api client.
	client := initClient()

	// Create the zone(if it doesnt already exist).
	domain := "myfailover.com"
	err := setupZone(client, domain)
	if err != nil {
		log.Fatal(err)
	}

	// Construct an NSONE API data source.
	source := data.NewSource("my api source", "nsone_v1")

	// Create the nsone_v1 api data source.
	// Note: this does not create the associated feeds.
	_, err = client.DataSources.Create(source)
	if err != nil {
		log.Fatal(err)
	}

	// Construct feeds which will drive the meta data for each answer.
	//  We'll use the id of these feeds when we connect the feeds to the
	//  answer meta below.
	feeds := map[string]*data.Feed{}

	// Construct the buffalo data feed.
	feeds["buf"] = data.NewFeed(
		"Buffalo Feed",
		data.Config{"label": "Buffalo-US"})

	// Construct the london data feed.
	feeds["lon"] = data.NewFeed(
		"London Feed",
		data.Config{"label": "London-UK"})

	// Create the buf/lon feeds through the rest api.
	for _, feed := range feeds {
		_, err = client.DataFeeds.Create(source.ID, feed)
		if err != nil {
			log.Fatal(err)
		}
	}

	record := dns.NewRecord(domain, "a", "A")                  // Construct the A record
	record.Regions["test"] = data.Region{data.Meta{Up: false}} // Add a region for examples sake.

	pAns := dns.NewAv4Answer("1.1.1.1")          // Construct the PRIMARY answer(with BUFFALO feed).
	pAns.Meta.Priority = 1                       // primary has higher priority
	pAns.Meta.Up = data.FeedPtr{feeds["buf"].ID} // Connect the primary answer to the Buffalo feed.

	sAns := dns.NewAv4Answer("2.2.2.2")          // Construct the SECONDARY answer(with LONDON feed).
	sAns.Meta.Priority = 2                       // secondary has lower priority.
	sAns.Meta.Up = data.FeedPtr{feeds["lon"].ID} // Connect the secondary answer to the London feed.

	record.AddAnswer(pAns) // Add primary answer to record
	record.AddAnswer(sAns) // Add secondary answer to record

	// Construct and add both filters to the record(ORDER MATTERS)
	record.AddFilter(filter.NewUp())
	record.AddFilter(filter.NewSelFirstN(1))

	// Helper to show record in json before sending PUT
	prettyPrint("record :", record)

	// Create the record
	_, err = client.Records.Create(record)
	if err != nil {
		log.Fatal(err)
	}

	// Create a linked record
	linkedRecord := dns.NewRecord(domain, "l", "A")
	linkedRecord.LinkTo(record.Domain)
	prettyPrint("linked record :", linkedRecord)

	// Create the linked record
	_, err = client.Records.Create(linkedRecord)
	if err != nil {
		log.Fatal(err)
	}

	// Flip which answer is 'Up' 5 times below.

	// Make an 'abort' goroutine for cancelling loop
	abort := make(chan struct{})
	go func() {
		os.Stdin.Read(make([]byte, 1)) // read a single byte
		abort <- struct{}{}
	}()

	// Every 5 sec, reverse which answer is 'up'.
	fmt.Println("Flipping answer every 5 sec. Press return to abort.")
	tick := time.Tick(5 * time.Second)
	var buffaloUp bool
	for countdown := 5; countdown > 0; countdown-- {
		fmt.Println(countdown)
		select {
		case <-tick:
			// Update the buffalo feed
			d := update{"Buffalo-US": data.Meta{Up: buffaloUp}}
			_, err = client.DataSources.Publish(source.ID, d)
			if err != nil {
				log.Fatal(err)
			}

			if buffaloUp {
				fmt.Printf("'dig %s' will respond with answer %s \n", record.Domain, pAns)
			} else {
				fmt.Printf("'dig %s' will respond with answer %s \n", record.Domain, sAns)
			}
			// Toggle status for next update.
			buffaloUp = !buffaloUp

		case <-abort:
			fmt.Println("Aborted")
			return
		}
	}
}

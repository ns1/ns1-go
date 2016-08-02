package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ns1/ns1-go"
)

func main() {
	k := os.Getenv("NSONE_APIKEY")
	if k == "" {
		fmt.Println("NSONE_APIKEY environment variable is not set, giving up")
	}
	// Initialize with only api key.
	api := nsone.New(k)

	api.Debug()
	api.RateLimitStrategySleep()

	// Create a zone to play in
	z := nsone.NewZone("testzone.net")
	err := api.CreateZone(z)
	if err != nil {
		log.Fatal(err)
	}

	// Construct/Create an NSONE API data source
	ds := nsone.NewDataSource("my api source", "nsone_v1")
	err = api.CreateDataSource(ds)
	if err != nil {
		log.Fatal(err)
	}

	// Construct feeds which will drive the meta data for each answer.
	// We'll use the id of these feeds when we connect the feeds to the
	// answer meta below.
	feed1 := nsone.NewDataFeed(ds.Id)
	feed1.Name = "feed to server1"
	feed1.Config = map[string]string{
		"label": "server1",
	}

	feed2 := nsone.NewDataFeed(ds.Id)
	feed2.Name = "feed to server2"
	feed2.Config = map[string]string{
		"label": "server2",
	}

	// Create the feeds through the rest api.
	err = api.CreateDataFeed(feed1)
	if err != nil {
		log.Fatal(err)
	}

	err = api.CreateDataFeed(feed2)
	if err != nil {
		log.Fatal(err)
	}

	// Create a record to connect to this source, with two answers
	// specify the up filter so we can send traffic to only those nodes
	// which are known to be up. We'll start with just the second answer up.
	// Each 'up' meta value is a feed pointer, pointing to the feeds we created above.
	r := nsone.NewRecord(z.Zone, "record", "A")
	r.Answers = []nsone.Answer{
		nsone.Answer{
			Answer: []string{"1.1.1.1"},
			Meta:   map[string]interface{}{"up": map[string]string{"feed": feed1.Id}},
		},
		nsone.Answer{
			Answer: []string{"9.9.9.9"},
			Meta:   map[string]interface{}{"up": map[string]string{"feed": feed2.Id}},
		},
	}
	r.Filters = []nsone.Filter{
		nsone.Filter{
			Filter: "up",
			Config: map[string]interface{}{},
		},
	}
	err = api.CreateRecord(r)
	if err != nil {
		log.Fatal(err)
	}

	// Now publish an update via feed to the records. here we push to both
	// feeds at once, but you can push to one or the other individually as well.
	update := make(map[string]interface{})
	update["server1"] = map[string]bool{"up": true}
	update["server2"] = map[string]bool{"up": false}
	err = api.PublishFeed(ds.Id, update)
	if err != nil {
		log.Fatal(err)
	}
}

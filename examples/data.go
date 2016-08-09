package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ns1/ns1-go/model/data"
	"github.com/ns1/ns1-go/model/dns"
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

	dataSources, _, err := client.DataSources.List()
	if err != nil {
		log.Fatal(err)
	}

	for _, ds := range dataSources {
		fmt.Println(ds)
	}

	// Create a zone to play in
	zone_name := "testdata.com"
	_, err = client.Zones.Delete(zone_name)
	if err != nil {
		log.Fatal(err)
	}
	z := dns.NewZone(zone_name)
	_, err = client.Zones.Create(z)
	if err != nil {
		log.Fatal(err)
	}

	// Construct/Create an NSONE API data source
	ds := data.NewSource("my api source", "nsone_v1")
	_, err = client.DataSources.Create(ds)
	if err != nil {
		log.Fatal(err)
	}

	// Construct feeds which will drive the meta data for each answer.
	// We'll use the id of these feeds when we connect the feeds to the
	// answer meta below.
	feed1 := data.NewFeed(ds.ID)
	feed1.Name = "feed to server1"
	feed1.Config = map[string]string{
		"label": "server1",
	}

	feed2 := data.NewFeed(ds.ID)
	feed2.Name = "feed to server2"
	feed2.Config = map[string]string{
		"label": "server2",
	}

	// Create the feeds through the rest api.
	_, err = client.DataFeeds.Create(feed1)
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.DataFeeds.Create(feed2)
	if err != nil {
		log.Fatal(err)
	}

	// Create a record to connect to this source, with two answers
	// specify the up filter so we can send traffic to only those nodes
	// which are known to be up. We'll start with just the second answer up.
	// Each 'up' meta value is a feed pointer, pointing to the feeds we created above.
	r := dns.NewRecord(z.Zone, "record", "A")
	r.Answers = []*dns.Answer{
		&dns.Answer{
			Answer: []string{"1.1.1.1"},
			Meta:   map[string]interface{}{"up": map[string]string{"feed": feed1.ID}},
		},
		&dns.Answer{
			Answer: []string{"9.9.9.9"},
			Meta:   map[string]interface{}{"up": map[string]string{"feed": feed2.ID}},
		},
	}
	r.Filters = []*dns.Filter{
		&dns.Filter{
			Filter: "up",
			Config: map[string]interface{}{},
		},
	}
	_, err = client.Records.Create(r)
	if err != nil {
		log.Fatal(err)
	}

	// Now publish an update via feed to the records. here we push to both
	// feeds at once, but you can push to one or the other individually as well.
	update := make(map[string]interface{})
	update["server1"] = map[string]bool{"up": true}
	update["server2"] = map[string]bool{"up": false}
	_, err = client.DataSources.PublishToFeed(ds.ID, update)
	if err != nil {
		log.Fatal(err)
	}
}

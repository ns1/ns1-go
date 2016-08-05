package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ns1/ns1-go"
	"github.com/ns1/ns1-go/rest"
)

func main() {
	k := os.Getenv("NS1_APIKEY")
	if k == "" {
		fmt.Println("NS1_APIKEY environment variable is not set, giving up")
	}

	httpClient := &http.Client{Timeout: time.Second * 10}
	client := rest.NewAPIClient(httpClient, rest.SetAPIKey(k))

	zones, err := client.Zones.List()
	if err != nil {
		log.Fatal(err)
	}

	for _, z := range zones {
		b, _ := json.MarshalIndent(z, "", "  ")
		fmt.Println(string(b))
	}

	// Construct/Create a zone.
	zone_name := "test.com"
	err = client.Zones.Delete(zone_name)
	if err != nil {
		log.Fatal(err)
	}

	z := ns1.NewZone(zone_name)
	z.NxTTL = 3600
	err = client.Zones.Create(z)
	if err != nil {
		log.Fatal(err)
	}

	// Update the zone.
	z.Retry = 5401
	err = client.Zones.Update(z)
	if err != nil {
		log.Fatal(err)
	}

	// Add an A record with a single static answer.
	orchidRec := ns1.NewRecord("test.com", "orchid", "A")
	orchidRec.Answers = []*ns1.Answer{&ns1.Answer{Answer: []string{"2.2.2.2"}}}
	err = client.Records.Create(orchidRec)
	if err != nil {
		log.Fatal(err)
	}

	// Add an A record with two static answers.
	honeyRec := ns1.NewRecord("test.com", "honey", "A")
	honeyRec.Answers = []*ns1.Answer{
		&ns1.Answer{Answer: []string{"1.2.3.4"}},
		&ns1.Answer{Answer: []string{"5.6.7.8"}},
	}
	err = client.Records.Create(honeyRec)
	if err != nil {
		log.Fatal(err)
	}

	// Add a cname
	potRec := ns1.NewRecord("test.com", "pot", "CNAME")
	potRec.Answers = []*ns1.Answer{&ns1.Answer{Answer: []string{"honey.test.com"}}}
	err = client.Records.Create(potRec)
	if err != nil {
		log.Fatal(err)
	}

	// Add a MX with two answers, priority 5 and 10
	mailRec := ns1.NewRecord("test.com", "mail", "MX")
	mailRec.Answers = []*ns1.Answer{
		&ns1.Answer{Answer: []string{"5", "mail1.test.com"}},
		&ns1.Answer{Answer: []string{"10", "mail2.test.com"}},
	}
	err = client.Records.Create(mailRec)
	if err != nil {
		log.Fatal(err)
	}

	// Add a AAAA, specify ttl of 300 seconds
	aaaaRec := ns1.NewRecord("test.com", "honey6", "AAAA")
	aaaaRec.TTL = 300
	aaaaRec.Answers = []*ns1.Answer{&ns1.Answer{Answer: []string{"2607:f8b0:4006:806::1010"}}}
	err = client.Records.Create(aaaaRec)
	if err != nil {
		log.Fatal(err)
	}

	// Add an A record using full answer format to specify 2 answers with meta data.
	// ensure edns-client-subnet is in use, and add two filters: geotarget_country,
	// and select_first_n, which has a filter config option N set to 1.
	bumbleRec := ns1.NewRecord("test.com", "bumble", "A")
	bumbleRec.Answers = []*ns1.Answer{
		&ns1.Answer{
			Answer: []string{"1.1.1.1"},
			Meta: map[string]interface{}{
				"up":      false,
				"country": []string{"US"},
			},
		},
		&ns1.Answer{
			Answer: []string{"9.9.9.9"},
			Meta: map[string]interface{}{
				"up":      true,
				"country": []string{"FR"},
			},
		},
	}
	bumbleRec.Filters = []*ns1.Filter{
		&ns1.Filter{
			Filter: "geotarget_country",
			Config: map[string]interface{}{},
		},
		&ns1.Filter{
			Filter: "select_first_n",
			Config: map[string]interface{}{"N": 1},
		},
	}
	err = client.Records.Create(bumbleRec)
	if err != nil {
		log.Fatal(err)
	}

}

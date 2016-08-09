package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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

	zones, _, err := client.Zones.List()
	if err != nil {
		log.Fatal(err)
	}

	for _, z := range zones {
		b, _ := json.MarshalIndent(z, "", "  ")
		fmt.Println(string(b))
	}

	// Construct/Create a zone.
	zone_name := "test.com"
	_, err = client.Zones.Delete(zone_name)
	if err != nil {
		log.Fatal(err)
	}

	z := dns.NewZone(zone_name)
	z.NxTTL = 3600
	_, err = client.Zones.Create(z)
	if err != nil {
		log.Fatal(err)
	}

	// Update the zone.
	z.Retry = 5401
	_, err = client.Zones.Update(z)
	if err != nil {
		log.Fatal(err)
	}

	// Add an A record with a single static answer.
	orchidRec := dns.NewRecord("test.com", "orchid", "A")
	orchidRec.Answers = []*dns.Answer{&dns.Answer{Answer: []string{"2.2.2.2"}}}
	_, err = client.Records.Create(orchidRec)
	if err != nil {
		log.Fatal(err)
	}

	// Add an A record with two static answers.
	honeyRec := dns.NewRecord("test.com", "honey", "A")
	honeyRec.Answers = []*dns.Answer{
		&dns.Answer{Answer: []string{"1.2.3.4"}},
		&dns.Answer{Answer: []string{"5.6.7.8"}},
	}
	_, err = client.Records.Create(honeyRec)
	if err != nil {
		log.Fatal(err)
	}

	// Add a cname
	potRec := dns.NewRecord("test.com", "pot", "CNAME")
	potRec.Answers = []*dns.Answer{&dns.Answer{Answer: []string{"honey.test.com"}}}
	_, err = client.Records.Create(potRec)
	if err != nil {
		log.Fatal(err)
	}

	// Add a MX with two answers, priority 5 and 10
	mailRec := dns.NewRecord("test.com", "mail", "MX")
	mailRec.Answers = []*dns.Answer{
		&dns.Answer{Answer: []string{"5", "mail1.test.com"}},
		&dns.Answer{Answer: []string{"10", "mail2.test.com"}},
	}
	_, err = client.Records.Create(mailRec)
	if err != nil {
		log.Fatal(err)
	}

	// Add a AAAA, specify ttl of 300 seconds
	aaaaRec := dns.NewRecord("test.com", "honey6", "AAAA")
	aaaaRec.TTL = 300
	aaaaRec.Answers = []*dns.Answer{&dns.Answer{Answer: []string{"2607:f8b0:4006:806::1010"}}}
	_, err = client.Records.Create(aaaaRec)
	if err != nil {
		log.Fatal(err)
	}

	// Add an A record using full answer format to specify 2 answers with meta data.
	// ensure edns-client-subnet is in use, and add two filters: geotarget_country,
	// and select_first_n, which has a filter config option N set to 1.
	bumbleRec := dns.NewRecord("test.com", "bumble", "A")
	bumbleRec.Answers = []*dns.Answer{
		&dns.Answer{
			Answer: []string{"1.1.1.1"},
			Meta: map[string]interface{}{
				"up":      false,
				"country": []string{"US"},
			},
		},
		&dns.Answer{
			Answer: []string{"9.9.9.9"},
			Meta: map[string]interface{}{
				"up":      true,
				"country": []string{"FR"},
			},
		},
	}
	bumbleRec.Filters = []*dns.Filter{
		&dns.Filter{
			Filter: "geotarget_country",
			Config: map[string]interface{}{},
		},
		&dns.Filter{
			Filter: "select_first_n",
			Config: map[string]interface{}{"N": 1},
		},
	}
	_, err = client.Records.Create(bumbleRec)
	if err != nil {
		log.Fatal(err)
	}

}

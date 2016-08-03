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

	// // Or initialize with http.Client
	// httpClient := &http.Client{
	// 	Timeout:   time.Second * 10,
	// }

	// // Adds logging to each http request.
	// doer := nsone.Decorate(httpClient,
	// 	nsone.Logging(log.New(os.Stdout, "", log.LstdFlags)))

	// api := nsone.NewAPIClient(doer, nsone.SetApiKey(k))

	// Construct/Create a zone.
	zone_name := "test.com"
	err := api.DeleteZone(zone_name)
	if err != nil {
		log.Fatal(err)
	}

	z := nsone.NewZone(zone_name)
	z.Nx_ttl = 3600
	err = api.CreateZone(z)
	if err != nil {
		log.Fatal(err)
	}

	// Update the zone.
	z.Retry = 5401
	err = api.UpdateZone(z)
	if err != nil {
		log.Fatal(err)
	}

	// Add an A record with a single static answer.
	orchidRec := nsone.NewRecord("test.com", "orchid", "A")
	orchidRec.Answers = []nsone.Answer{nsone.Answer{Answer: []string{"2.2.2.2"}}}
	err = api.CreateRecord(orchidRec)
	if err != nil {
		log.Fatal(err)
	}

	// Add an A record with two static answers.
	honeyRec := nsone.NewRecord("test.com", "honey", "A")
	honeyRec.Answers = []nsone.Answer{
		nsone.Answer{Answer: []string{"1.2.3.4"}},
		nsone.Answer{Answer: []string{"5.6.7.8"}},
	}
	err = api.CreateRecord(honeyRec)
	if err != nil {
		log.Fatal(err)
	}

	// Add a cname
	potRec := nsone.NewRecord("test.com", "pot", "CNAME")
	potRec.Answers = []nsone.Answer{nsone.Answer{Answer: []string{"honey.test.com"}}}
	err = api.CreateRecord(potRec)
	if err != nil {
		log.Fatal(err)
	}

	// Add a MX with two answers, priority 5 and 10
	mailRec := nsone.NewRecord("test.com", "mail", "MX")
	mailRec.Answers = []nsone.Answer{
		nsone.Answer{Answer: []string{"5", "mail1.test.com"}},
		nsone.Answer{Answer: []string{"10", "mail2.test.com"}},
	}
	err = api.CreateRecord(mailRec)
	if err != nil {
		log.Fatal(err)
	}

	// Add a AAAA, specify ttl of 300 seconds
	aaaaRec := nsone.NewRecord("test.com", "honey6", "AAAA")
	aaaaRec.Ttl = 300
	aaaaRec.Answers = []nsone.Answer{nsone.Answer{Answer: []string{"2607:f8b0:4006:806::1010"}}}
	err = api.CreateRecord(aaaaRec)
	if err != nil {
		log.Fatal(err)
	}

	// Add an A record using full answer format to specify 2 answers with meta data.
	// ensure edns-client-subnet is in use, and add two filters: geotarget_country,
	// and select_first_n, which has a filter config option N set to 1.
	bumbleRec := nsone.NewRecord("test.com", "bumble", "A")
	bumbleRec.Answers = []nsone.Answer{
		nsone.Answer{
			Answer: []string{"1.1.1.1"},
			Meta: map[string]interface{}{
				"up":      false,
				"country": []string{"US"},
			},
		},
		nsone.Answer{
			Answer: []string{"9.9.9.9"},
			Meta: map[string]interface{}{
				"up":      true,
				"country": []string{"FR"},
			},
		},
	}
	bumbleRec.Filters = []nsone.Filter{
		nsone.Filter{
			Filter: "geotarget_country",
			Config: map[string]interface{}{},
		},
		nsone.Filter{
			Filter: "select_first_n",
			Config: map[string]interface{}{"N": 1},
		},
	}
	err = api.CreateRecord(bumbleRec)
	if err != nil {
		log.Fatal(err)
	}

}

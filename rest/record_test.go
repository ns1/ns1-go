package rest

import (
	"fmt"
	"gopkg.in/ns1/ns1-go.v2/rest/model/dns"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

)

func TestCaaRecordIntegrationTest(t *testing.T){
	apiKey := os.Getenv("NS1_APIKEY")
	testDomain := os.Getenv("NS1_TESTDOMAIN")

	httpClient := &http.Client{Timeout: time.Second * 10}
	doer := Decorate(httpClient, Logging(log.New(os.Stdout, "", log.LstdFlags)))

	client := NewClient(doer, SetAPIKey(apiKey))

	records, _, err := client.Records.Get(testDomain, testDomain, "CAA")
	if err != nil {
		t.Error(err)
	}

	for _, a := range records.Answers {
		fmt.Println(a)
	}

	answer := dns.NewCAAAnswer(0, "issue", "globalchicken.com")
	records.AddAnswer(answer)

	response, err := client.Records.Update(records)

	if err != nil {
		t.Error(err)
	}

	println(response)



}

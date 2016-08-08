package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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

	teams, err := client.Teams.List()
	if err != nil {
		log.Fatal(err)
	}

	for _, t := range teams {
		b, _ := json.MarshalIndent(t, "", "  ")
		fmt.Println(string(b))
	}

	users, err := client.Users.List()
	if err != nil {
		log.Fatal(err)
	}

	for _, u := range users {
		b, _ := json.MarshalIndent(u, "", "  ")
		fmt.Println(string(b))
	}

	keys, err := client.APIKeys.List()
	if err != nil {
		log.Fatal(err)
	}

	for _, k := range keys {
		b, _ := json.MarshalIndent(k, "", "  ")
		fmt.Println(string(b))
	}
}

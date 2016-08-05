package main

import (
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

	httpClient := &http.Client{Timeout: time.Second * 10}
	client := rest.NewAPIClient(httpClient, rest.SetAPIKey(k))

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

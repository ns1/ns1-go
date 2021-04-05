package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	api "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/account"
)

var ddiClient *api.Client

// Helper that initializes rest api client from environment variable.
func init() {
	k := os.Getenv("NS1_APIKEY")
	if k == "" {
		log.Fatal("NS1_APIKEY environment variable is not set, giving up")
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Timeout: time.Second * 10}
	httpClient.Transport = tr
	// Adds logging to each http request.
	doer := api.Decorate(httpClient, api.Logging(log.New(os.Stdout, "", log.LstdFlags)))
	// In this example DDI is running on localhost
	ddiClient = api.NewClient(doer, api.SetAPIKey(k), api.SetDDIAPI(), api.SetEndpoint("https://127.0.0.1/v1/"))
}

func main() {
	key := account.APIKey{
		Name:              "foobar",
		TeamIDs:           []string{},
		IPWhitelist:       []string{"1.1.1.1", "2.2.2.2"},
		IPWhitelistStrict: true,
		Permissions: account.PermissionsMap{
			DNS: account.PermissionsDNS{
				ZonesAllow: []string{},
				ZonesDeny:  []string{},
			},
			DHCP: &account.PermissionsDHCP{
				TagsAllow: &[]account.AuthTag{},
				TagsDeny:  &[]account.AuthTag{},
			},
			IPAM: &account.PermissionsIPAM{
				TagsAllow: &[]account.AuthTag{},
				TagsDeny:  &[]account.AuthTag{},
			},
		},
	}

	_, err := ddiClient.APIKeys.Create(&key)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("New API key: %v", key)

	updatedIPAMPermissions := &account.PermissionsIPAM{
		TagsAllow: &[]account.AuthTag{
			{
				Name:  "auth:foo",
				Value: "",
			},
			{
				Name:  "auth:bar",
				Value: "baz",
			},
		},
		TagsDeny: &[]account.AuthTag{},
	}
	key.Permissions.IPAM = updatedIPAMPermissions

	_, err = ddiClient.APIKeys.Update(&key)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Updated API key: %v", key)
}

package nsone

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type APIClient struct {
	ApiKey string
}

type ZonePrimary struct {
	Enabled     bool
	Secondaries []string
}

type Zone struct {
	Id            string
	Ttl           int
	Nx_ttl        int
	Retry         int
	Zone          string
	Refresh       int
	Expiry        int
	Primary       ZonePrimary
	Dns_servers   []string
	Networks      []int
	Network_pools []string
	Hostmaster    string
	Pool          string
	Meta          map[string]string
}

func New(k string) *APIClient {
	return &APIClient{
		ApiKey: k,
	}
}

func (c APIClient) GetThing(uri string) []byte {
	client := &http.Client{}
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("X-NSONE-Key", c.ApiKey)
	fmt.Println(req)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	return body
}

func (c APIClient) GetZones() []Zone {
	var zl []Zone
	err := json.Unmarshal(c.GetThing("https://api.nsone.net/v1/zones"), &zl)
	if err != nil {
		panic(err)
	}
	return zl
}

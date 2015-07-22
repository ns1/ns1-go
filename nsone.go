package nsone

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type APIClient struct {
	ApiKey string
}

type ZonePrimary struct {
	Enabled     bool     `json:"enabled"`
	Secondaries []string `json:"secondaries,omitempty"`
}

type Zone struct {
	Id      string `json:"id,omitempty"`
	Ttl     int    `json:"ttl,omitempty"`
	Nx_ttl  int    `json:"nx_ttl,omitempty"`
	Retry   int    `json:"retry,omitempty"`
	Zone    string `json:"zone,omitempty"`
	Refresh int    `json:"refresh,omitempty"`
	Expiry  int    `json:"expiry,omitempty"`
	//	Primary       ZonePrimary       `json:"primary,omitempty"`
	Dns_servers   []string          `json:"dns_servers,omitempty"`
	Networks      []int             `json:"networks,omitempty"`
	Network_pools []string          `json:"network_pools,omitempty"`
	Hostmaster    string            `json:"hostmaster,omitempty"`
	Pool          string            `json:"pool,omitempty"`
	Meta          map[string]string `json:"meta,omitempty"`
}

func NewZone(zone string) *Zone {
	return &Zone{
		Zone: zone,
	}
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
	log.Printf("[DEBUG] Get: %s", uri)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	log.Println(resp)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(body))
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

func (c APIClient) GetZone(z *Zone) error {
	err := json.Unmarshal(c.GetThing(fmt.Sprintf("https://api.nsone.net/v1/zones/%s", z.Zone)), z)
	if err != nil {
		panic(err)
	}
	return err
}

func (c APIClient) DeleteZone(z *Zone) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("https://api.nsone.net/v1/zones/%s", z.Zone), nil)
	req.Header.Add("X-NSONE-Key", c.ApiKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	log.Println(resp)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(body))
	return err
}

func (c APIClient) CreateZone(z *Zone) error {
	body, err := json.Marshal(z)
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	req, err := http.NewRequest("PUT", fmt.Sprintf("https://api.nsone.net/v1/zones/%s", z.Zone), bytes.NewReader(body))
	req.Header.Add("X-NSONE-Key", c.ApiKey)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	log.Println(resp)
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)
	log.Println(string(body))
	err = json.Unmarshal(body, z)
	if err != nil {
		panic(err)
	}
	return nil
}

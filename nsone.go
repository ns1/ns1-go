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

type Record struct {
	Id     string `json:"id,omitempty"`
	Zone   string `json:"zone,omitempty"`
	Domain string `json:"domain,omitempty"`
	Type   string `json:"type,omitempty"`
}

func NewZone(zone string) *Zone {
	return &Zone{
		Zone: zone,
	}
}

func NewRecord(zone string, domain string, t string) *Record {
	return &Record{
		Zone:   zone,
		Domain: domain,
		Type:   t,
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

func (c APIClient) GetZone(zone string) (*Zone, error) {
	z := NewZone(zone)
	err := json.Unmarshal(c.GetThing(fmt.Sprintf("https://api.nsone.net/v1/zones/%s", z.Zone)), z)
	if err != nil {
		panic(err)
	}
	return z, err
}

func (c APIClient) GetRecord(zone string, domain string, t string) (*Record, error) {
	r := NewRecord(zone, domain, t)
	err := json.Unmarshal(c.GetThing(fmt.Sprintf("https://api.nsone.net/v1/zones/%s/%s/%s", r.Zone, r.Domain, r.Type)), r)
	return r, err
}

func (c APIClient) DeleteZone(zone string) error {
	return c.DeleteThing(fmt.Sprintf("https://api.nsone.net/v1/zones/%s", zone))
}

func (c APIClient) doRequest(t string, uri string) (*http.Response, error) {
	req, err := http.NewRequest(t, uri, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("X-NSONE-Key", c.ApiKey)
	client := &http.Client{}
	return client.Do(req)
}

func (c APIClient) DeleteThing(uri string) error {
	resp, err := c.doRequest("DELETE", uri)
	resp.Body.Close()
	if err != nil {
		panic(err)
	}
	// FIXME return status
	return err
}

func (c APIClient) DeleteRecord(zone string, domain string, t string) error {
	return c.DeleteThing(fmt.Sprintf("https://api.nsone.net/v1/zones/%s/%/%s", zone, domain, t))
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

func (c APIClient) CreateRecord(r *Record) error {
	return nil
}

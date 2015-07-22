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

func (c APIClient) GetZones() []Zone {
	var zl []Zone
	err := json.Unmarshal(c.doHTTP("GET", "https://api.nsone.net/v1/zones", nil), &zl)
	if err != nil {
		panic(err)
	}
	return zl
}

func (c APIClient) GetZone(zone string) (*Zone, error) {
	z := NewZone(zone)
	err := json.Unmarshal(c.doHTTP("GET", fmt.Sprintf("https://api.nsone.net/v1/zones/%s", z.Zone), nil), z)
	if err != nil {
		panic(err)
	}
	return z, err
}

func (c APIClient) GetRecord(zone string, domain string, t string) (*Record, error) {
	r := NewRecord(zone, domain, t)
	err := json.Unmarshal(c.doHTTP("GET", fmt.Sprintf("https://api.nsone.net/v1/zones/%s/%s/%s", r.Zone, r.Domain, r.Type), nil), r)
	return r, err
}

func (c APIClient) DeleteZone(zone string) error {
	return c.DeleteThing(fmt.Sprintf("https://api.nsone.net/v1/zones/%s", zone))
}

func (c APIClient) doHTTP(method string, uri string, rbody []byte) []byte {
	r := bytes.NewReader(rbody)
	log.Printf("[DEBUG] %s: %s", method, uri)
	req, err := http.NewRequest(method, uri, r)
	if err != nil {
		panic(err)
	}
	req.Header.Add("X-NSONE-Key", c.ApiKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	log.Println(resp)
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	log.Println(string(body))
	return body

}

func (c APIClient) DeleteThing(uri string) error {
	_ = c.doHTTP("DELETE", uri, nil)
	// FIXME return status
	return nil
}

func (c APIClient) DeleteRecord(zone string, domain string, t string) error {
	return c.DeleteThing(fmt.Sprintf("https://api.nsone.net/v1/zones/%s/%/%s", zone, domain, t))
}

func (c APIClient) CreateZone(z *Zone) error {
	rbody, err := json.Marshal(z)
	if err != nil {
		panic(err)
	}
	body := c.doHTTP("PUT", fmt.Sprintf("https://api.nsone.net/v1/zones/%s", z.Zone), rbody)
	log.Println("MOO: Response body")
	log.Println(string(body))

	err = json.Unmarshal(body, z)
	if err != nil {
		panic(err)
	}
	return nil
}

func (c APIClient) UpdateZone(z *Zone) error {
	rbody, err := json.Marshal(z)
	if err != nil {
		panic(err)
	}
	body := c.doHTTP("POST", fmt.Sprintf("https://api.nsone.net/v1/zones/%s", z.Zone), rbody)
	log.Println("MOO: Response body")
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

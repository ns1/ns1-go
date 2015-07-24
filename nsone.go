package nsone

import (
	"bytes"
	"encoding/json"
	"errors"
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

type ZoneSecondary struct {
	Status       string `json:"status,omitempty"`
	Last_xfr     int    `json:"last_xfr,omitempty"`
	Primary_ip   string `json:"primary_ip,omitempty"`
	Primary_port int    `json:"primary_port,omitempty"`
	Enabled      bool   `json:"enabled"`
	Expired      bool   `json:"expired,omitempty"`
}

type Zone struct {
	Id            string            `json:"id,omitempty"`
	Ttl           int               `json:"ttl,omitempty"`
	Nx_ttl        int               `json:"nx_ttl,omitempty"`
	Retry         int               `json:"retry,omitempty"`
	Zone          string            `json:"zone,omitempty"`
	Refresh       int               `json:"refresh,omitempty"`
	Expiry        int               `json:"expiry,omitempty"`
	Primary       ZonePrimary       `json:"primary,omitempty"`
	Dns_servers   []string          `json:"dns_servers,omitempty"`
	Networks      []int             `json:"networks,omitempty"`
	Network_pools []string          `json:"network_pools,omitempty"`
	Hostmaster    string            `json:"hostmaster,omitempty"`
	Pool          string            `json:"pool,omitempty"`
	Meta          map[string]string `json:"meta,omitempty"`
	Secondary     ZoneSecondary     `json:"secondary,omitempty"`
}

type Answer struct {
	Answer []string               `json:"answer,omitempty"`
	Meta   map[string]interface{} `json:"meta,omitempty"`
}

type Record struct {
	Id      string            `json:"id,omitempty"`
	Zone    string            `json:"zone,omitempty"`
	Domain  string            `json:"domain,omitempty"`
	Type    string            `json:"type,omitempty"`
	Link    string            `json:"link,omitempty"`
	Meta    map[string]string `json:"meta,omitempty"`
	Answers []Answer          `json:"answers,omitempty"`
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

func (c APIClient) GetZones() ([]Zone, error) {
	var zl []Zone
	body, err := c.doHTTP("GET", "https://api.nsone.net/v1/zones", nil)
	if err != nil {
		return zl, err
	}
	err = json.Unmarshal(body, &zl)
	return zl, err
}

func (c APIClient) GetZone(zone string) (*Zone, error) {
	z := NewZone(zone)
	body, err := c.doHTTP("GET", fmt.Sprintf("https://api.nsone.net/v1/zones/%s", z.Zone), nil)
	if err != nil {
		return z, err
	}
	err = json.Unmarshal(body, z)
	if err != nil {
		return z, err
	}
	return z, err
}

func (c APIClient) GetRecord(zone string, domain string, t string) (*Record, error) {
	r := NewRecord(zone, domain, t)
	body, err := c.doHTTP("GET", fmt.Sprintf("https://api.nsone.net/v1/zones/%s/%s/%s", r.Zone, r.Domain, r.Type), nil)
	if err != nil {
		return r, err
	}
	err = json.Unmarshal(body, r)
	return r, err
}

func (c APIClient) DeleteZone(zone string) error {
	return c.DeleteThing(fmt.Sprintf("https://api.nsone.net/v1/zones/%s", zone))
}

func (c APIClient) doHTTP(method string, uri string, rbody []byte) ([]byte, error) {
	var body []byte
	r := bytes.NewReader(rbody)
	log.Printf("[DEBUG] %s: %s (%s)", method, uri, string(rbody))
	req, err := http.NewRequest(method, uri, r)
	if err != nil {
		return body, err
	}
	req.Header.Add("X-NSONE-Key", c.ApiKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return body, err
	}
	log.Println(resp)
	body, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != 200 {
		return body, errors.New(fmt.Sprintf("%s: %s", resp.Status, string(body)))
	}
	log.Println(string(body))
	return body, nil

}

func (c APIClient) DeleteThing(uri string) error {
	_, err := c.doHTTP("DELETE", uri, nil)
	// FIXME return status
	return err
}

func (c APIClient) DeleteRecord(zone string, domain string, t string) error {
	return c.DeleteThing(fmt.Sprintf("https://api.nsone.net/v1/zones/%s/%s/%s", zone, domain, t))
}

func (c APIClient) CreateZone(z *Zone) error {
	rbody, err := json.Marshal(z)
	if err != nil {
		return err
	}
	body, err := c.doHTTP("PUT", fmt.Sprintf("https://api.nsone.net/v1/zones/%s", z.Zone), rbody)
	if err != nil {
		return err
	}
	log.Println("MOO: Response body")
	log.Println(string(body))

	err = json.Unmarshal(body, z)
	if err != nil {
		return err
	}
	return nil
}

func (c APIClient) UpdateZone(z *Zone) error {
	rbody, err := json.Marshal(z)
	if err != nil {
		return err
	}
	body, err := c.doHTTP("POST", fmt.Sprintf("https://api.nsone.net/v1/zones/%s", z.Zone), rbody)
	log.Println("MOO: Response body")
	log.Println(string(body))

	err = json.Unmarshal(body, z)
	if err != nil {
		return err
	}
	return nil
}

func (c APIClient) CreateRecord(r *Record) error {
	rbody, err := json.Marshal(r)
	if err != nil {
		return err
	}
	body, err := c.doHTTP("PUT", fmt.Sprintf("https://api.nsone.net/v1/zones/%s/%s/%s", r.Zone, r.Domain, r.Type), rbody)
	if err != nil {
		return err
	}
	log.Println("MOO: Response body")
	log.Println(string(body))

	err = json.Unmarshal(body, r)
	if err != nil {
		return err
	}
	return nil
}

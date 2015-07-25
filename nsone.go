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

func New(k string) *APIClient {
	return &APIClient{
		ApiKey: k,
	}
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
	log.Println(fmt.Sprintf("Response body: %s", string(body)))
	return body, nil

}

func (c APIClient) doHTTPUnmarshal(method string, uri string, rbody []byte, unpack_into interface{}) error {
	body, err := c.doHTTP(method, uri, rbody)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, unpack_into)
}

func (c APIClient) DeleteThing(uri string) error {
	_, err := c.doHTTP("DELETE", uri, nil)
	return err
}

func (c APIClient) GetZones() ([]Zone, error) {
	var zl []Zone
	err := c.doHTTPUnmarshal("GET", "https://api.nsone.net/v1/zones", nil, &zl)
	return zl, err
}

func (c APIClient) GetZone(zone string) (*Zone, error) {
	z := NewZone(zone)
	err := c.doHTTPUnmarshal("GET", fmt.Sprintf("https://api.nsone.net/v1/zones/%s", z.Zone), nil, z)
	return z, err
}

func (c APIClient) DeleteZone(zone string) error {
	return c.DeleteThing(fmt.Sprintf("https://api.nsone.net/v1/zones/%s", zone))
}

func (c APIClient) CreateZone(z *Zone) error {
	rbody, err := json.Marshal(z)
	if err != nil {
		return err
	}
	return c.doHTTPUnmarshal("PUT", fmt.Sprintf("https://api.nsone.net/v1/zones/%s", z.Zone), rbody, z)
}

func (c APIClient) UpdateZone(z *Zone) error {
	rbody, err := json.Marshal(z)
	if err != nil {
		return err
	}
	return c.doHTTPUnmarshal("POST", fmt.Sprintf("https://api.nsone.net/v1/zones/%s", z.Zone), rbody, z)
}

func (c APIClient) CreateRecord(r *Record) error {
	rbody, err := json.Marshal(r)
	if err != nil {
		return err
	}
	return c.doHTTPUnmarshal("PUT", fmt.Sprintf("https://api.nsone.net/v1/zones/%s/%s/%s", r.Zone, r.Domain, r.Type), rbody, r)
}

func (c APIClient) GetRecord(zone string, domain string, t string) (*Record, error) {
	r := NewRecord(zone, domain, t)
	err := c.doHTTPUnmarshal("GET", fmt.Sprintf("https://api.nsone.net/v1/zones/%s/%s/%s", r.Zone, r.Domain, r.Type), nil, r)
	return r, err
}

func (c APIClient) DeleteRecord(zone string, domain string, t string) error {
	return c.DeleteThing(fmt.Sprintf("https://api.nsone.net/v1/zones/%s/%s/%s", zone, domain, t))
}

func (c APIClient) UpdateRecord(r *Record) error {
	return errors.New("UpdateRecord not implemented")
}

func (c APIClient) CreateDataSource(r *DataSource) error {
	return errors.New("CreateDataSource not implemented")
}

func (c APIClient) GetDataSource(id string) (*DataSource, error) {
	return nil, errors.New("GetDataSource not implemented")
}

func (c APIClient) DeleteDataSource(id string) error {
	return errors.New("DeleteDataSource not implemented")
}

func (c APIClient) UpdateDataSource(d *DataSource) error {
	return errors.New("UpdateDataSource not implemented")
}

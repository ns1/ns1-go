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

func (c APIClient) doHTTPBoth(method string, uri string, s interface{}) error {
	rbody, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return c.doHTTPUnmarshal(method, uri, rbody, s)
}

func (c APIClient) doHTTPDelete(uri string) error {
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
	return c.doHTTPDelete(fmt.Sprintf("https://api.nsone.net/v1/zones/%s", zone))
}

func (c APIClient) CreateZone(z *Zone) error {
	return c.doHTTPBoth("PUT", fmt.Sprintf("https://api.nsone.net/v1/zones/%s", z.Zone), z)
}

func (c APIClient) UpdateZone(z *Zone) error {
	return c.doHTTPBoth("POST", fmt.Sprintf("https://api.nsone.net/v1/zones/%s", z.Zone), z)
}

func (c APIClient) CreateRecord(r *Record) error {
	return c.doHTTPBoth("PUT", fmt.Sprintf("https://api.nsone.net/v1/zones/%s/%s/%s", r.Zone, r.Domain, r.Type), r)
}

func (c APIClient) GetRecord(zone string, domain string, t string) (*Record, error) {
	r := NewRecord(zone, domain, t)
	err := c.doHTTPUnmarshal("GET", fmt.Sprintf("https://api.nsone.net/v1/zones/%s/%s/%s", r.Zone, r.Domain, r.Type), nil, r)
	return r, err
}

func (c APIClient) DeleteRecord(zone string, domain string, t string) error {
	return c.doHTTPDelete(fmt.Sprintf("https://api.nsone.net/v1/zones/%s/%s/%s", zone, domain, t))
}

func (c APIClient) UpdateRecord(r *Record) error {
	return c.doHTTPBoth("POST", fmt.Sprintf("https://api.nsone.net/v1/ones/%s/%s/%s", r.Zone, r.Domain, r.Type), r)
}

func (c APIClient) CreateDataSource(ds *DataSource) error {
	return c.doHTTPBoth("PUT", "https://api.nsone.net/v1/data/sources", ds)
}

func (c APIClient) GetDataSource(id string) (*DataSource, error) {
	ds := DataSource{}
	err := c.doHTTPUnmarshal("GET", fmt.Sprintf("https://api.nsone.net/v1/data/sources/%s", id), nil, &ds)
	return &ds, err
}

func (c APIClient) DeleteDataSource(id string) error {
	return c.doHTTPDelete(fmt.Sprintf("https://api.nsone.net/v1/data/sources/%s", id))
}

func (c APIClient) UpdateDataSource(ds *DataSource) error {
	return c.doHTTPBoth("POST", fmt.Sprintf("https://api.nsone.net/v1/data/sources/%s", ds.Id), ds)
}
func (c APIClient) CreateDataFeed(df *DataFeed) error {
	return c.doHTTPBoth("PUT", fmt.Sprintf("https://api.nsone.net/v1/data/feeds/%s", df.SourceId), df)
}

func (c APIClient) GetDataFeed(ds_id string, df_id string) (*DataFeed, error) {
	df := NewDataFeed(ds_id)
	err := c.doHTTPUnmarshal("GET", fmt.Sprintf("https://api.nsone.net/v1/data/sources/%s/%s", ds_id, df_id), nil, &df)
	return &df, err
}

func (c APIClient) DeleteDataFeed(ds_id string, df_id string) error {
	return c.doHTTPDelete(fmt.Sprintf("https://api.nsone.net/v1/data/sources/%s/%s", ds_id, df_id))
}

func (c APIClient) UpdateDataFeed(df *DataFeed) error {
	return c.doHTTPBoth("POST", fmt.Sprintf("https://api.nsone.net/v1/data/sources/%s", df.SourceId, df.Id), df)
}

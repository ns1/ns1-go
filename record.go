package nsone

import "fmt"

// Answer wraps the values of a Record's "filters" attribute
type Answer struct {
	Region string                 `json:"region,omitempty"`
	Answer []string               `json:"answer,omitempty"`
	Meta   map[string]interface{} `json:"meta,omitempty"`
}

// Filter wraps the values of a Record's "filters" attribute
type Filter struct {
	Filter   string                 `json:"filter"`
	Disabled bool                   `json:"disabled,omitempty"`
	Config   map[string]interface{} `json:"config"`
}

// Record wraps an NS1 /zone/{zone}/{domain}/{type} resource
type Record struct {
	Id              string            `json:"id,omitempty"`
	Zone            string            `json:"zone,omitempty"`
	Domain          string            `json:"domain,omitempty"`
	Type            string            `json:"type,omitempty"`
	Link            string            `json:"link,omitempty"`
	Meta            map[string]string `json:"meta,omitempty"`
	Answers         []Answer          `json:"answers"`
	Filters         []Filter          `json:"filters,omitempty"`
	Ttl             int               `json:"ttl,omitempty"`
	UseClientSubnet bool              `json:"use_client_subnet"`
	Regions         map[string]Region `json:"regions,omitempty"`
}

// Region wraps the values of a Record's "regions" attribute
type Region struct {
	Meta RegionMeta `json:"meta"`
}

// RegionMeta wraps the values of a Record's "regions.*.meta" attribute
type RegionMeta struct {
	GeoRegion []string `json:"georegion,omitempty"`
	Country   []string `json:"country,omitempty"`
	USState   []string `json:"us_state,omitempty"`
	Up        bool     `json:"up,omitempty"`
}

// MetaFeed wraps an Answer.Metadata element which points to a feed
type MetaFeed struct {
	Feed string `json:"feed"`
}

// MetaStatic wraps an Answer.Metadata element which just wraps a string
type MetaStatic string

// NewRecord takes a zone, domain and record type t and creates a *Record with UseClientSubnet: true & empty Answers
func NewRecord(zone string, domain string, t string) *Record {
	return &Record{
		Zone:            zone,
		Domain:          domain,
		Type:            t,
		UseClientSubnet: true,
		Answers:         make([]Answer, 0),
	}
}

// NewAnswer creates an empty Answer
func NewAnswer() Answer {
	return Answer{
		Meta: make(map[string]interface{}),
	}
}

// NewMetaFeed takes a feed_id and creates a MetaFeed
func NewMetaFeed(feedID string) MetaFeed {
	return MetaFeed{
		Feed: feedID,
	}
}

// LinkTo sets a Record Link to an FQDN and ensures no other record configuration is specified
func (r *Record) LinkTo(to string) {
	r.Meta = nil
	r.Answers = make([]Answer, 0)
	r.Link = to
}

// CreateRecord takes a *Record and creates a new DNS record in the specified zone, for the specified domain, of the given record type
func (c APIClient) CreateRecord(r *Record) error {
	return c.doHTTPBoth("PUT", fmt.Sprintf("https://api.nsone.net/v1/zones/%s/%s/%s", r.Zone, r.Domain, r.Type), r)
}

// GetRecord takes a zone, domain and record type t and returns full configuration for a DNS record
func (c APIClient) GetRecord(zone string, domain string, t string) (*Record, error) {
	r := NewRecord(zone, domain, t)
	status, err := c.doHTTPUnmarshal("GET", fmt.Sprintf("https://api.nsone.net/v1/zones/%s/%s/%s", r.Zone, r.Domain, r.Type), nil, r)
	if status == 404 {
		r.Id = ""
		r.Zone = ""
		r.Domain = ""
		r.Type = ""
		return r, nil
	}
	return r, err
}

// DeleteRecord takes a zone, domain and record type t and removes an existing record and all associated answers and configuration details
func (c APIClient) DeleteRecord(zone string, domain string, t string) error {
	return c.doHTTPDelete(fmt.Sprintf("https://api.nsone.net/v1/zones/%s/%s/%s", zone, domain, t))
}

// UpdateRecord takes a *Record and modifies configuration details for an existing DNS record
func (c APIClient) UpdateRecord(r *Record) error {
	return c.doHTTPBoth("POST", fmt.Sprintf("https://api.nsone.net/v1/zones/%s/%s/%s", r.Zone, r.Domain, r.Type), r)
}

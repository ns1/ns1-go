package ns1

import (
	"fmt"
	"strings"
)

// Record wraps an NS1 /zone/{zone}/{domain}/{type} resource
type Record struct {
	ID              string            `json:"id,omitempty"`
	Zone            string            `json:"zone,omitempty"`
	Domain          string            `json:"domain,omitempty"`
	Type            string            `json:"type,omitempty"`
	Link            string            `json:"link,omitempty"`
	Meta            map[string]string `json:"meta,omitempty"`
	Answers         []*Answer         `json:"answers"`
	Filters         []*Filter         `json:"filters,omitempty"`
	TTL             int               `json:"ttl,omitempty"`
	UseClientSubnet bool              `json:"use_client_subnet"`
	Regions         map[string]Region `json:"regions,omitempty"`
}

// Implementation of Stringer interface.
func (r Record) String() string {
	return fmt.Sprintf("%s %s", r.Domain, r.Type)
}

// NewRecord takes a zone, domain and record type t and creates a *Record with UseClientSubnet: true & empty Answers
func NewRecord(zone string, domain string, t string) *Record {
	if !strings.HasSuffix(domain, zone) {
		domain = fmt.Sprintf("%s.%s", domain, zone)
	}
	return &Record{
		Zone:            zone,
		Domain:          domain,
		Type:            t,
		UseClientSubnet: true,
	}
}

// MetaFeed wraps an Answer.Metadata element which points to a feed
type MetaFeed struct {
	Feed string `json:"feed"`
}

// NewMetaFeed takes a feed_id and creates a MetaFeed
func NewMetaFeed(feedID string) MetaFeed {
	return MetaFeed{
		Feed: feedID,
	}
}

// MetaStatic wraps an Answer.Metadata element which just wraps a string
type MetaStatic string

// LinkTo sets a Record Link to an FQDN and ensures no other record configuration is specified
func (r *Record) LinkTo(to string) {
	r.Meta = nil
	r.Link = to
}

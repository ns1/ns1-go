package dns

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

// Implementation of Stringer interface. Concatenates domain and record type.
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

// LinkTo sets a Record Link to an FQDN and ensures no other record configuration is specified
func (r *Record) LinkTo(to string) {
	r.Meta = nil
	r.Link = to
}

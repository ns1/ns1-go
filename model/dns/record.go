package dns

import (
	"fmt"
	"strconv"
	"strings"
)

// Record wraps an NS1 /zone/{zone}/{domain}/{type} resource
type Record struct {
	ID              string            `json:"id,omitempty"`
	Zone            string            `json:"zone"`
	Domain          string            `json:"domain"`
	Type            string            `json:"type"`
	Link            string            `json:"link,omitempty"`
	TTL             int               `json:"ttl,omitempty"`
	UseClientSubnet bool              `json:"use_client_subnet,omitempty"`
	Answers         []*Answer         `json:"answers,omitempty"`
	Filters         []*Filter         `json:"filters,omitempty"`
	Regions         map[string]Region `json:"regions,omitempty"`
	Meta            map[string]string `json:"meta,omitempty"`
}

// Implementation of Stringer interface. Concatenates domain and record type.
func (r Record) String() string {
	return fmt.Sprintf("%s %s", r.Domain, r.Type)
}

// NewRecord takes a zone, domain and record type t and creates a *Record with UseClientSubnet: true & empty Answers
//
// Answers are required on PUT, but not POST.
func NewRecord(zone string, domain string, t string) *Record {
	if !strings.HasSuffix(domain, zone) {
		domain = fmt.Sprintf("%s.%s", domain, zone)
	}
	return &Record{
		Zone:    zone,
		Domain:  domain,
		Type:    t,
		Answers: []*Answer{},
	}
}

// LinkTo sets a Record Link to an FQDN and ensures no other record configuration is specified
func (r *Record) LinkTo(to string) {
	r.Meta = nil
	r.Link = to
}

func (r *Record) AddAnswer(ans *Answer) {
	if r.Answers == nil {
		r.Answers = []*Answer{}
	}

	r.Answers = append(r.Answers, ans)
}

// Answer wraps the values of a Record's "filters" attribute
type Answer struct {
	Rdata  []string               `json:"answer"`
	Region string                 `json:"region,omitempty"`
	Meta   map[string]interface{} `json:"meta,omitempty"`
}

// NewAnswer creates a generic Answer with given rdata.
func NewAnswer(rdata []string) *Answer {
	return &Answer{Rdata: rdata}
}

func NewAv4Answer(host string) *Answer {
	return &Answer{Rdata: []string{host}}
}

func NewAv6Answer(host string) *Answer {
	return &Answer{Rdata: []string{host}}
}

func NewALIASAnswer(host string) *Answer {
	return &Answer{Rdata: []string{host}}
}

func NewCNAMEAnswer(name string) *Answer {
	return &Answer{Rdata: []string{name}}
}

func NewTXTAnswer(text string) *Answer {
	return &Answer{Rdata: []string{text}}
}

func NewMXAnswer(pri int, host string) *Answer {
	return &Answer{Rdata: []string{strconv.Itoa(pri), host}}
}

func NewSRVAnswer(priority, weight, port int, target string) *Answer {
	return &Answer{Rdata: []string{
		strconv.Itoa(priority),
		strconv.Itoa(weight),
		strconv.Itoa(port),
		target,
	}}
}

// Region wraps the values of a Record's "regions" attribute
type Region struct {
	Meta map[string]interface{} `json:"meta"`
}

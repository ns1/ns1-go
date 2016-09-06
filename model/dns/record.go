package dns

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ns1/ns1-go/model/filter"
	"github.com/ns1/ns1-go/model/meta"
)

// Record wraps an NS1 /zone/{zone}/{domain}/{type} resource
type Record struct {
	ID              string               `json:"id,omitempty"`
	Zone            string               `json:"zone"`
	Domain          string               `json:"domain"`
	Type            string               `json:"type"`
	Link            string               `json:"link,omitempty"`
	TTL             int                  `json:"ttl,omitempty"`
	UseClientSubnet bool                 `json:"use_client_subnet,omitempty"`
	Answers         []*Answer            `json:"answers,omitempty"`
	Filters         []*filter.Filter     `json:"filters,omitempty"`
	Regions         map[string]meta.Meta `json:"regions,omitempty"`
	Meta            *meta.Meta           `json:"meta,omitempty"`
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
		Regions: map[string]meta.Meta{},
		Meta:    meta.New(),
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

func (r *Record) AddFilter(fil *filter.Filter) {
	if r.Filters == nil {
		r.Filters = []*filter.Filter{}
	}

	r.Filters = append(r.Filters, fil)
}

// Answer wraps the values of a Record's "filters" attribute
type Answer struct {
	Rdata  []string   `json:"answer"`
	Region string     `json:"region,omitempty"`
	Meta   *meta.Meta `json:"meta,omitempty"`
}

func (a *Answer) SetRegion(region string) {
	a.Region = region
}

// NewAnswer creates a generic Answer with given rdata.
func NewAnswer(rdata []string) *Answer {
	return &Answer{
		Rdata: rdata,
		Meta:  meta.New(),
	}
}

func NewAv4Answer(host string) *Answer {
	return &Answer{
		Rdata: []string{host},
		Meta:  meta.New(),
	}
}

func NewAv6Answer(host string) *Answer {
	return &Answer{
		Rdata: []string{host},
		Meta:  meta.New(),
	}
}

func NewALIASAnswer(host string) *Answer {
	return &Answer{
		Rdata: []string{host},
		Meta:  meta.New(),
	}
}

func NewCNAMEAnswer(name string) *Answer {
	return &Answer{
		Rdata: []string{name},
		Meta:  meta.New(),
	}
}

func NewTXTAnswer(text string) *Answer {
	return &Answer{
		Rdata: []string{text},
		Meta:  meta.New(),
	}
}

func NewMXAnswer(pri int, host string) *Answer {
	return &Answer{
		Rdata: []string{strconv.Itoa(pri), host},
		Meta:  meta.New(),
	}
}

func NewSRVAnswer(priority, weight, port int, target string) *Answer {
	return &Answer{
		Rdata: []string{
			strconv.Itoa(priority),
			strconv.Itoa(weight),
			strconv.Itoa(port),
			target,
		},
		Meta: meta.New(),
	}
}

type Region struct {
	Meta meta.Meta `json:"meta,omitempty"`
}

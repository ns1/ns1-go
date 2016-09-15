package dns

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ns1/ns1-go/rest/model/data"
)

// Answer wraps the values of a Record's "filters" attribute
type Answer struct {
	Meta *data.Meta `json:"meta,omitempty"`

	// Answer response data. eg:
	// Av4: ["1.1.1.1"]
	// Av6: ["2001:db8:85a3::8a2e:370:7334"]
	// MX:  [10, "2.2.2.2"]
	Rdata []string `json:"answer"`

	// Region(grouping) that answer belongs to.
	RegionName string `json:"region,omitempty"`
}

// Implementation of Stringer interface. Simply displays response data.
func (a Answer) String() string {
	return strings.Trim(fmt.Sprint(a.Rdata), "[]")
}

func (a *Answer) SetRegion(name string) {
	a.RegionName = name
}

// NewAnswer creates a generic Answer with given rdata.
func NewAnswer(rdata []string) *Answer {
	return &Answer{
		Meta:  &data.Meta{},
		Rdata: rdata,
	}
}

func NewAv4Answer(host string) *Answer {
	return &Answer{
		Meta:  &data.Meta{},
		Rdata: []string{host},
	}
}

func NewAv6Answer(host string) *Answer {
	return &Answer{
		Meta:  &data.Meta{},
		Rdata: []string{host},
	}
}

func NewALIASAnswer(host string) *Answer {
	return &Answer{
		Meta:  &data.Meta{},
		Rdata: []string{host},
	}
}

func NewCNAMEAnswer(name string) *Answer {
	return &Answer{
		Meta:  &data.Meta{},
		Rdata: []string{name},
	}
}

func NewTXTAnswer(text string) *Answer {
	return &Answer{
		Meta:  &data.Meta{},
		Rdata: []string{text},
	}
}

func NewMXAnswer(pri int, host string) *Answer {
	return &Answer{
		Meta:  &data.Meta{},
		Rdata: []string{strconv.Itoa(pri), host},
	}
}

func NewSRVAnswer(priority, weight, port int, target string) *Answer {
	return &Answer{
		Meta: &data.Meta{},
		Rdata: []string{
			strconv.Itoa(priority),
			strconv.Itoa(weight),
			strconv.Itoa(port),
			target,
		},
	}
}

package dns

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/ns1/ns1-go.v2/rest/model/data"
)

// Answer wraps the values of a Record's "filters" attribute
type Answer struct {
	ID string `json:"id,omitempty"`

	Meta *data.Meta `json:"meta,omitempty"`

	// Answer response data. eg:
	// Av4: ["1.1.1.1"]
	// Av6: ["2001:db8:85a3::8a2e:370:7334"]
	// MX:  [10, "2.2.2.2"]
	Rdata []string `json:"answer"`

	// Region(grouping) that answer belongs to.
	RegionName string `json:"region,omitempty"`
}

// UnmarshalJSON parses responses to Answer and attempts to convert Rdata elements to string
func (a *Answer) UnmarshalJSON(b []byte) error {
	type TempAnswer struct {
		ID         string        `json:"id,omitempty"`
		Meta       *data.Meta    `json:"meta,omitempty"`
		Rdata      []interface{} `json:"answer"`
		RegionName string        `json:"region,omitempty"`
	}
	var temp TempAnswer

	if err := json.Unmarshal(b, &temp); err != nil {
		return err
	}

	a.ID = temp.ID
	a.Meta = temp.Meta
	a.RegionName = temp.RegionName

	var rd []string
	for _, record := range temp.Rdata {
		switch v := record.(type) {
		case string:
			rd = append(rd, v)
		case float64:
			rd = append(rd, strconv.Itoa(int(v)))
		default:
			return fmt.Errorf("Could not unmarshal Rdata value %v (type %T) as type string", v, v)
		}
	}
	a.Rdata = rd

	return nil
}

func (a Answer) String() string {
	return strings.Trim(fmt.Sprint(a.Rdata), "[]")
}

// SetRegion associates a region with this answer.
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

// NewAv4Answer creates an Answer for A record.
func NewAv4Answer(host string) *Answer {
	return &Answer{
		Meta:  &data.Meta{},
		Rdata: []string{host},
	}
}

// NewAv6Answer creates an Answer for AAAA record.
func NewAv6Answer(host string) *Answer {
	return &Answer{
		Meta:  &data.Meta{},
		Rdata: []string{host},
	}
}

// NewALIASAnswer creates an Answer for ALIAS record.
func NewALIASAnswer(host string) *Answer {
	return &Answer{
		Meta:  &data.Meta{},
		Rdata: []string{host},
	}
}

// NewCNAMEAnswer creates an Answer for CNAME record.
func NewCNAMEAnswer(name string) *Answer {
	return &Answer{
		Meta:  &data.Meta{},
		Rdata: []string{name},
	}
}

// NewTXTAnswer creates an Answer for TXT record.
func NewTXTAnswer(text string) *Answer {
	return &Answer{
		Meta:  &data.Meta{},
		Rdata: []string{text},
	}
}

// NewMXAnswer creates an Answer for MX record.
func NewMXAnswer(pri int, host string) *Answer {
	return &Answer{
		Meta:  &data.Meta{},
		Rdata: []string{strconv.Itoa(pri), host},
	}
}

// NewSRVAnswer creates an Answer for SRV record.
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

// NewCAAAnswer creates an Answer for a CAA record.
func NewCAAAnswer(flag int, tag, value string) *Answer {
	return &Answer{
		Meta:  &data.Meta{},
		Rdata: []string{strconv.Itoa(flag), tag, value},
	}
}

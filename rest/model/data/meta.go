package data

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// FeedPtr represents the dynamic metadata value in which a feed is providing the value.
type FeedPtr struct {
	FeedID string `json:"feed,omitempty"`
}

// Meta contains information on an entity's metadata table. Metadata key/value
// pairs are used by a record's filter pipeline during a dns query.
// All values can be a feed id as well, indicating real-time updates of these values.
// Structure/Precendence of metadata tables:
//  - Record
//    - Meta <- lowest precendence in filter
//    - Region(s)
//      - Meta <- middle precedence in filter chain
//      - ...
//    - Answer(s)
//      - Meta <- highest precedence in filter chain
//      - ...
//    - ...
type Meta struct {
	// STATUS

	// Indicates whether or not entity is considered 'up'
	// bool or FeedPtr.
	Up interface{} `json:"up,omitempty"`

	// Indicates the number of active connections.
	// Values must be positive.
	// int or FeedPtr.
	Connections interface{} `json:"connections,omitempty"`

	// Indicates the number of active requests (HTTP or otherwise).
	// Values must be positive.
	// int or FeedPtr.
	Requests interface{} `json:"requests,omitempty"`

	// Indicates the "load average".
	// Values must be positive, and will be rounded to the nearest tenth.
	// float64 or FeedPtr.
	LoadAvg interface{} `json:"loadavg,omitempty"`

	// The Job ID of a Pulsar telemetry gathering job and routing granularities
	// to associate with.
	// string or FeedPtr.
	Pulsar interface{} `json:"pulsar,omitempty"`

	// GEOGRAPHICAL

	// Must be between -180.0 and +180.0 where negative
	// indicates South and positive indicates North.
	// e.g., the longitude of the datacenter where a server resides.
	// float64 or FeedPtr.
	Latitude interface{} `json:"latitude,omitempty"`

	// Must be between -180.0 and +180.0 where negative
	// indicates West and positive indicates East.
	// e.g., the longitude of the datacenter where a server resides.
	// float64 or FeedPtr.
	Longitude interface{} `json:"longitude,omitempty"`

	// Valid geographic regions are: 'US-EAST', 'US-CENTRAL', 'US-WEST',
	// 'EUROPE', 'ASIAPAC', 'SOUTH-AMERICA', 'AFRICA'.
	// e.g., the rough geographic location of the Datacenter where a server resides.
	// []string or FeedPtr.
	Georegion interface{} `json:"georegion,omitempty"`

	// Countr(ies) must be specified as ISO3166 2-character country code(s).
	// []string or FeedPtr.
	Country interface{} `json:"country,omitempty"`

	// State(s) must be specified as standard 2-character state code(s).
	// []string or FeedPtr.
	USState interface{} `json:"us_state,omitempty"`

	// Canadian Province(s) must be specified as standard 2-character province
	// code(s).
	// []string or FeedPtr.
	CAProvince interface{} `json:"ca_province,omitempty"`

	// INFORMATIONAL

	// Notes to indicate any necessary details for operators.
	// Up to 256 characters in length.
	// string or FeedPtr.
	Note interface{} `json:"note,omitempty"`

	// NETWORK

	// IP (v4 and v6) prefixes in CIDR format ("a.b.c.d/mask").
	// May include up to 1000 prefixes.
	// e.g., "1.2.3.4/24"
	// []string or FeedPtr.
	IPPrefixes interface{} `json:"ip_prefixes,omitempty"`

	// Autonomous System (AS) number(s).
	// May include up to 1000 AS numbers.
	// []string or FeedPtr.
	ASN interface{} `json:"asn,omitempty"`

	// TRAFFIC

	// Indicates the "priority tier".
	// Lower values indicate higher priority.
	// Values must be positive.
	// int or FeedPtr.
	Priority interface{} `json:"priority,omitempty"`

	// Indicates a weight.
	// Filters that use weights normalize them.
	// Any positive values are allowed.
	// Values between 0 and 100 are recommended for simplicity's sake.
	// float64 or FeedPtr.
	Weight interface{} `json:"weight,omitempty"`

	// Indicates a "low watermark" to use for load shedding.
	// The value should depend on the metric used to determine
	// load (e.g., loadavg, connections, etc).
	// int or FeedPtr.
	LowWatermark interface{} `json:"low_watermark,omitempty"`

	// Indicates a "high watermark" to use for load shedding.
	// The value should depend on the metric used to determine
	// load (e.g., loadavg, connections, etc).
	// int or FeedPtr.
	HighWatermark interface{} `json:"high_watermark,omitempty"`
}

// StringMap returns a map[string]interface{} representation of metadata (for use with terraform in nested structures)
func (meta *Meta) StringMap() map[string]interface{} {
	m := make(map[string]interface{})
	v := reflect.Indirect(reflect.ValueOf(meta))
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fv := v.Field(i)
		if fv.IsNil() {
			continue
		}
		tag := f.Tag.Get("json")

		tag = strings.Split(tag, ",")[0]

		m[tag] = FormatInterface(fv.Interface())
	}
	return m
}

// FormatInterface takes an interface of types: string, bool, int, float64, []string, and FeedPtr, and returns a string representation of said interface
func FormatInterface(i interface{}) string {
	switch v := i.(type) {
	case string:
		return v
	case bool:
		if v {
			return "1"
		} else {
			return "0"
		}
	case int:
		return strconv.FormatInt(int64(v), 10)
	case float64:
		if isIntegral(v) {
			return strconv.FormatInt(int64(v), 10)
		}
		return strconv.FormatFloat(v, 'f', -1, 64)
	case []string:
		return strings.Join(v, ",")
	case FeedPtr:
		data, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}
		return string(data)
	default:
		panic(fmt.Sprintf("expected v to be convertible to a string, got: %+v, %T", v, v))
	}
}

// ParseType returns an interface containing a string, bool, int, float64, []string, or FeedPtr
// float64 values with no decimal may be returned as integers, but that should be ok because the api won't know the difference
// when it's json encoded
func ParseType(s string) interface{} {
	slc := strings.Split(s, ",")
	if len(slc) > 1 {
		return slc
	}

	feedptr := FeedPtr{}
	err := json.Unmarshal([]byte(s), &feedptr)
	if err == nil {
		return feedptr
	}

	f, err := strconv.ParseFloat(s, 64)
	if err == nil {
		if isIntegral(f) {
			return int(f)
		} else {
			return f
		}
	}

	return s
}

// isIntegral returns whether or not a float64 has a decimal place
func isIntegral(f float64) bool {
	return f == float64(int(f))
}

// MetaFromMap creates a *Meta and uses reflection to set fields from a map. This will panic if a value for a key is not a string.
// This it to ensure compatibility with terraform
func MetaFromMap(m map[string]interface{}) *Meta {
	meta := &Meta{}
	mv := reflect.Indirect(reflect.ValueOf(meta))
	mt := mv.Type()
	for k, v := range m {
		name := ToCamel(k)
		if _, ok := mt.FieldByName(name); ok {
			fv := mv.FieldByName(name)
			if name == "Up" {
				if v.(string) == "1" {
					fv.Set(reflect.ValueOf(true))
				} else {
					fv.Set(reflect.ValueOf(false))
				}
			} else {
				fv.Set(reflect.ValueOf(ParseType(v.(string))))
			}
		}
	}
	return meta
}

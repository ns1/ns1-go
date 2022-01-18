package data

import (
	"reflect"
	"testing"
)

func TestMeta_StringMap(t *testing.T) {
	meta := &Meta{}
	meta.Up = true
	meta.Latitude = 0.50
	meta.Note = "Hello, World!"
	meta.Longitude = FeedPtr{FeedID: "12345678"}
	meta.Georegion = []interface{}{"US-EAST"}
	meta.Priority = 10
	meta.Weight = 10.0
	meta.Cost = 25.0
	meta.Country = []string{"US", "UK"}
	meta.IPPrefixes = []string{"1.1.1.1/24", "2.2.2.2/24"}
	meta.ASN = []interface{}{float64(1), float64(2)}
	meta.Pulsar = []interface{}{map[string]interface{}{
		"job_id":     "abcdef",
		"bias":       "*0.55",
		"a5m_cutoff": 0.9,
	}}
	m := meta.StringMap()

	if m["up"].(string) != "1" {
		t.Fatal("up should be 1")
	}

	if m["latitude"].(string) != "0.5" {
		t.Fatal("latitude should be '0.5'")
	}

	if m["georegion"].(string) != "US-EAST" {
		t.Fatal("georegion should be 'US-EAST")
	}

	if m["note"].(string) != "Hello, World!" {
		t.Fatal("note should be 'Hello, World!'")
	}

	if m["priority"].(string) != "10" {
		t.Fatal("priority should be 10")
	}

	if m["weight"].(string) != "10" {
		t.Fatal("weight should be 10")
	}

	if m["cost"].(string) != "25" {
		t.Fatal("cost should be 25")
	}

	if m["country"].(string) != "US,UK" {
		t.Fatal("country should be 'US,UK'")
	}

	if m["ip_prefixes"].(string) != "1.1.1.1/24,2.2.2.2/24" {
		t.Fatal("IP prefixes should be '1.1.1.1/24,2.2.2.2/24")
	}

	if m["asn"].(string) != "1,2" {
		t.Fatal("ASN should be '1,2'")
	}

	expected := `[{"a5m_cutoff":0.9,"bias":"*0.55","job_id":"abcdef"}]`
	if m["pulsar"].(string) != expected {
		t.Fatal("pulsar should be", expected, "was", m["pulsar"].(string))
	}

	expected = `{"feed":"12345678"}`
	if m["longitude"].(string) != expected {
		t.Fatal("longitude should be", expected, "was", m["longitude"].(string))
	}

	meta.Up = false
	m = meta.StringMap()
	if m["up"].(string) != "0" {
		t.Fatal("up should be 0")
	}

	// Test mapping deserialized json feed ptr, for Terraform workaround
	meta.Up = map[string]interface{}{"feed": "12345678"}
	m = meta.StringMap()
	if m["up"].(string) != expected {
		t.Fatal("up should be", expected, "was", m["up"].(string))
	}

	meta.Subdivisions = map[string]interface{}{"BR": []string{"SP", "MG"}}
	stra := meta.StringMap()
	str_expected := "{\"BR\":[\"SP\",\"MG\"]}"
	if stra["subdivisions"] != str_expected {
		t.Fatal("expected:", str_expected, "got: ", stra["subdivisions"])
	}

	meta.Up = map[string]interface{}{"key": "12345678"}
	stra = meta.StringMap()
	str_expected = "{\"key\":\"12345678\"}"
	if stra["up"] != str_expected {
		t.Fatal("expected:", str_expected, "got: ", stra["up"])
	}
}

func TestParseType(t *testing.T) {
	fs := "3.14"

	v := ParseType(fs)

	if _, ok := v.(float64); !ok {
		t.Fatal("value should be float64, was", v)
	}

	cs := "hello,goodbye"

	v = ParseType(cs)

	if _, ok := v.([]string); !ok {
		t.Fatal("value should be []string, was", v)
	}

	is := "42"

	v = ParseType(is)

	if _, ok := v.(int); !ok {
		t.Fatal("value should be int, was", v)
	}

	s := "string value"

	v = ParseType(s)
	if _, ok := v.(string); !ok {
		t.Fatal("value should be string, was", v)
	}

	fp := `{"feed":"12345678"}`
	v = ParseType(fp)
	if _, ok := v.(FeedPtr); !ok {
		t.Fatal("value should be FeedPointer, was", v)
	}
}

func TestMetaFromMap(t *testing.T) {
	m := make(map[string]interface{})

	m["latitude"] = "0.50"
	m["up"] = "1"
	m["connections"] = "5"
	m["longitude"] = `{"feed":"12345678"}`
	m["ip_prefixes"] = "1.1.1.1/24,2.2.2.2/24"
	m["asn"] = "1"
	m["pulsar"] = `[{"job_id":"abcdef","bias":"*0.55","a5m_cutoff":0.9}]`
	meta := MetaFromMap(m)

	if meta.ASN.(string) != "1" {
		t.Fatal("meta.ASN should have been 1", meta.ASN)
	}

	if !meta.Up.(bool) {
		t.Fatal("meta.Up should be true")
	}

	if meta.Latitude.(float64) != 0.5 {
		t.Fatal("meta.Latitude should equal 0.5, was", meta.Latitude)
	}

	if meta.Connections.(int) != 5 {
		t.Fatal("meta.Connections should equal 5, was", meta.Connections)
	}

	if meta.Longitude.(FeedPtr).FeedID != "12345678" {
		t.Fatal("meta.Longitude should be a feed ptr with id 12345678, was", meta.Longitude)
	}

	expect := []map[string]interface{}{map[string]interface{}{
		"job_id":     "abcdef",
		"bias":       "*0.55",
		"a5m_cutoff": 0.9,
	}}
	if !reflect.DeepEqual(meta.Pulsar, expect) {
		t.Fatalf("meta.Pulsar should be %v, was %v", expect, meta.Pulsar)
	}

	expected := []string{"1.1.1.1/24", "2.2.2.2/24"}
	if !reflect.DeepEqual(meta.IPPrefixes.([]string), expected) {
		t.Fatal("meta.IPPrefixes should be a slice containing elements `1.1.1.1/24` and `2.2.2.2/24`")
	}

	m["up"] = "0"
	meta = MetaFromMap(m)

	if meta.Up.(bool) {
		t.Fatal("meta.Up should be false")
	}

	// Terraform 0.12 will no longer auto-convert bool values to 0 or 1, must support true and false
	m["up"] = "false"
	meta = MetaFromMap(m)

	if meta.Up.(bool) {
		t.Fatal("meta.Up should be false")
	}

	m["up"] = "true"
	meta = MetaFromMap(m)

	if !meta.Up.(bool) {
		t.Fatal("meta.Up should be true")
	}

	m["up"] = `{"feed":"12345678"}`
	meta = MetaFromMap(m)

	if meta.Up.(FeedPtr).FeedID != "12345678" {
		t.Fatal("meta.Up should be a feed ptr with id 12345678, was", meta.Up)
	}

	m["asn"] = "1,2,3"
	meta = MetaFromMap(m)
	expected = []string{"1", "2", "3"}
	if !reflect.DeepEqual(meta.ASN.([]string), expected) {
		t.Fatal("meta.ASN should be a slice containing elements `1`, `2`, and `3`")
	}

	sub := make(map[string]interface{})
	sub["BR"] = []interface{}{"SP", "MG"}
	sub["DZ"] = []interface{}{"01"}
	m["Subdivisions"] = sub
	meta = MetaFromMap(m)
	if !reflect.DeepEqual(meta.Subdivisions.(map[string]interface{}), sub) {
		t.Fatal("meta.Subdivisions should be a map[string]interface{} containing elements \"BR\":[\"SP\",\"MG\"],\"DZ\":[\"01\"] ")
	}

	subStr := "{\"BR\":[\"SP\",\"MG\"], \"DZ\": [\"01\"]}"
	m["Subdivisions"] = subStr
	meta = MetaFromMap(m)
	if !reflect.DeepEqual(meta.Subdivisions.(map[string]interface{}), sub) {
		t.Fatal(meta.Subdivisions.(map[string]interface{})["DZ"], sub)
	}
}

func TestGeokeyString(t *testing.T) {
	expected := "AFRICA,ASIAPAC,EUROPE,SOUTH-AMERICA,US-CENTRAL,US-EAST,US-WEST"
	got := geoKeyString()
	if expected != got {
		t.Fatalf("expected '%s', got '%s'", expected, got)
	}
}

func TestMeta_Validate(t *testing.T) {
	m := &Meta{}
	m.Up = true
	m.Georegion = "US-EAST"
	m.Connections = 5
	m.Longitude = 0.80
	m.Latitude = 0.80
	m.Country = "US"
	m.USState = "MA"
	m.CAProvince = "ON"
	m.Note = "Hello, just testing out this cool meta validation"
	m.LoadAvg = 3.14
	m.Weight = 40.0
	m.Cost = 60.0
	m.Requests = 10
	m.IPPrefixes = "10.0.0.1/24"
	m.Priority = 1
	m.Pulsar = []interface{}{map[string]interface{}{
		"job_id":     "abcd",
		"bias":       "*0.55",
		"a5m_cutoff": 0.9,
	}}
	errs := m.Validate()
	if len(errs) > 0 {
		t.Fatal("there should be 0 errors, but there were", len(errs), ":", errs)
	}

	m.IPPrefixes = []interface{}{"10.0.0.1/24", "10.0.0.2/24"}
	errs = m.Validate()
	if len(errs) > 0 {
		t.Fatal("there should be 0 errors, but there were", len(errs), ":", errs)
	}

	m.Georegion = "fantasy region"
	m.Up = "bad value"
	m.Connections = -5
	m.Longitude = 10000.0
	m.Latitude = -10000.0
	m.Country = "fantasy land"
	m.USState = "fantasy state"
	m.CAProvince = "quebec"
	m.Note = string(make([]rune, 257))
	m.LoadAvg = -3.14
	m.Weight = -40.0
	m.Cost = -70.0
	m.Requests = -1
	m.IPPrefixes = "1234567"
	m.Priority = -1
	m.Pulsar = []interface{}{map[string]interface{}{}}
	errs = m.Validate()
	if len(errs) != 16 {
		t.Fatal("expected 15 errors, but there were", len(errs), ":", errs)
	}

	m = &Meta{}
	m.Georegion = []interface{}{"US-EAST", "fantasy land"}
	m.Country = []interface{}{"US", "CANADA"}
	m.IPPrefixes = []interface{}{"1234567", "1.1.1.1/24"}
	m.Up = struct{}{}
	m.Pulsar = []interface{}{map[string]interface{}{
		"job_id": true,
	}}
	errs = m.Validate()
	if len(errs) != 5 {
		t.Fatal("expected 5 errors, but there were", len(errs), ":", errs)
	}

	// Test validation of []string and string values passed from Terraform
	m = &Meta{}
	m.Georegion = []string{"US-EAST", "US-CENTRAL"}
	m.Country = []string{"US", "CA"}
	m.IPPrefixes = []string{"1.1.1.1/24", "2.2.2.2/24"}
	m.Pulsar = `[{"job_id":"abcd","bias":"*0.55","a5m_cutoff":0.9}]`
	errs = m.Validate()
	if len(errs) > 0 {
		t.Fatal("there should be 0 errors, but there were", len(errs), ":", errs)
	}

	m = &Meta{}
	m.Georegion = []string{"US-EAST", "fantasy land"}
	m.Country = []string{"US", "CANADA"}
	m.IPPrefixes = []string{"1234567", "1.1.1.1/24"}
	m.Pulsar = "blah"
	errs = m.Validate()
	if len(errs) != 4 {
		t.Fatal("expected 4 errors, but there were", len(errs), ":", errs)
	}
}

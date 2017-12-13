package data

import "testing"

func TestMeta_StringMap(t *testing.T) {
	meta := &Meta{}
	meta.Up = true
	meta.Latitude = 0.50

	meta.Longitude = FeedPtr{FeedID: "12345678"}
	meta.Georegion = []interface{}{"US-EAST"}

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

	expected := `{"feed":"12345678"}`
	if m["longitude"].(string) != expected {
		t.Fatal("longitude should be", expected, "was", m["longitude"].(string))
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
}

func TestMetaFromMap(t *testing.T) {
	m := make(map[string]interface{})

	m["latitude"] = "0.50"
	m["up"] = "1"
	m["connections"] = "5"
	m["longitude"] = `{"feed":"12345678"}`
	meta := MetaFromMap(m)

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

}

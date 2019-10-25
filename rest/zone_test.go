package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZone_list_pagination(t *testing.T) {
	// It should follow Links and gather all Zones into the list
	testServer := mockAPI(t, responseMap{"/zones": handleZoneList})
	doer := testServer.Client()
	client := NewClient(doer, SetEndpoint(testServer.URL))

	zones, resp, err := client.Zones.List()
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	expected := []string{
		"example.com",
		"coolzone.cool",
		"example.org",
		"example.io",
	}
	assert.Equal(t, 4, len(zones))
	for idx := range zones {
		assert.Equal(t, expected[idx], zones[idx].Zone)
	}
}

func TestZone_get_pagination(t *testing.T) {
	// It should follow Links and gather all Records into the Zone
	testServer := mockAPI(t, responseMap{"/zones/coolzone.cool": handleZoneGet})
	doer := testServer.Client()
	client := NewClient(doer, SetEndpoint(testServer.URL))

	zone, resp, err := client.Zones.Get("coolzone.cool")
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	expected := []string{
		"a.coolzone.cool",
		"b.coolzone.cool",
		"c.coolzone.cool",
		"d.coolzone.cool",
	}
	assert.Equal(t, "coolzone.cool", zone.Zone)
	assert.Equal(t, 4, len(zone.Records))
	for idx := range zone.Records {
		assert.Equal(t, expected[idx],  zone.Records[idx].Domain)
	}

}

func TestZone_pagination_disabled(t *testing.T) {
	// It should not follow Links
	testServer := mockAPI(t, responseMap{"/zones": handleZoneList, "/zones/coolzone.cool": handleZoneGet})
	doer := testServer.Client()
	client := NewClient(doer, SetEndpoint(testServer.URL), SetFollowPagination(false))

	zones, resp, err := client.Zones.List()
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	expected := []string{
		"example.com",
		"coolzone.cool",
	}
	assert.Equal(t, 2, len(zones))
	for idx := range zones {
		assert.Equal(t, expected[idx], zones[idx].Zone)
	}

	zone, resp, err := client.Zones.Get("coolzone.cool")
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	expected = []string{
		"a.coolzone.cool",
		"b.coolzone.cool",
	}
	assert.Equal(t, "coolzone.cool", zone.Zone)
	assert.Equal(t, 2, len(zone.Records))
	for idx := range zone.Records {
		assert.Equal(t, expected[idx],  zone.Records[idx].Domain)
	}
}

func TestZone_http_errors(t *testing.T) {
	// It should have a non-nil resp with HTTP Errors, and
	// it should have a JSON error if the response body cannot be parsed as JSON
	testServer := mockAPI(t, responseMap{
		"/zones": func(res http.ResponseWriter, req *http.Request) error {
			res.WriteHeader(429)
			res.Write([]byte("{}"))
			return nil
		},
		"/zones/coolzone.cool": func(res http.ResponseWriter, req *http.Request) error {
			res.WriteHeader(500)
			res.Write([]byte("Internal Server Error"))
			return nil
		},
	})
	doer := testServer.Client()
	client := NewClient(doer, SetEndpoint(testServer.URL))

	zones, resp, err := client.Zones.List()
	assert.IsType(t, &Error{}, err)
	assert.Equal(t, 429, resp.StatusCode)
	assert.Nil(t, zones)

	zone, resp, err := client.Zones.Get("coolzone.cool")
	assert.IsType(t, &json.SyntaxError{}, err)
	assert.Equal(t, 500, resp.StatusCode)
	assert.Nil(t, zone)
}

func handleZoneList(res http.ResponseWriter, req *http.Request) error {
	const zoneListQuery = "?after=coolzone.cool&limit=2"
	var zoneListLinkHeader = fmt.Sprintf(`</zones?%s>; rel="next"`, zoneListQuery)
	const zoneListPageOne = `[
    {"zone": "example.com"},
    {"zone": "coolzone.cool"}
  ]`
	const zoneListPageTwo = `[
    {"zone": "example.org"},
    {"zone": "example.io"}
  ]`
	var body string
	if req.URL.RawQuery == "" {
		res.Header().Set("Link", zoneListLinkHeader)
		body = zoneListPageOne
	} else if req.URL.RawQuery == zoneListQuery {
		body = zoneListPageTwo
	} else {
		return fmt.Errorf("unhandled query: %s", req.URL.RawQuery)
	}
	res.WriteHeader(200)
	res.Write([]byte(body))
	return nil
}

func handleZoneGet(res http.ResponseWriter, req *http.Request) error {
	const zoneGetQuery = "after=b.coolzone.cool&limit=2"
	var zoneGetLinkHeader = fmt.Sprintf(
		`</zones/coolzone.cool?%s>; rel="next"`, zoneGetQuery,
	)
	const zoneGetPageOne = `{
    "zone": "coolzone.cool",
    "records": [{
      "domain":"a.coolzone.cool",
      "short_answers": ["1.1.1.1"],
      "type":"A"
    },
    {
      "domain":"b.coolzone.cool",
      "short_answers": ["2.2.2.2"],
      "type":"A"
    }]
  }`
	const zoneGetPageTwo = `{
    "zone": "coolzone.cool",
    "records": [{
      "domain":"c.coolzone.cool",
      "short_answers": ["3.3.3.3"],
      "type":"A"
    },
    {
      "domain":"d.coolzone.cool",
      "short_answers": ["4.4.4.4"],
      "type":"A"
    }]
  }`
	var body string
	if req.URL.RawQuery == "" {
		res.Header().Set("Link", zoneGetLinkHeader)
		body = zoneGetPageOne
	} else if req.URL.RawQuery == zoneGetQuery {
		body = zoneGetPageTwo
	} else {
		return fmt.Errorf("unhandled query: %s", req.URL.RawQuery)
	}
	res.WriteHeader(200)
	res.Write([]byte(body))
	return nil
}

type responseMap map[string]func(res http.ResponseWriter, req *http.Request) error

func mockAPI(t *testing.T, rm responseMap) *httptest.Server {
	handler := func(res http.ResponseWriter, req *http.Request) {
		if handler, ok := rm[req.URL.Path]; ok {
			if err := handler(res, req); err != nil {
				t.Fatalf(err.Error())
			}
		} else {
			t.Fatalf("unhandled path %s", req.URL.Path)
		}
	}
	return httptest.NewServer(http.HandlerFunc(handler))
}

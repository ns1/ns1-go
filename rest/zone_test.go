package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZone_List(t *testing.T) {
	// It should follow Links and gather all Zones into the list
	testServer := mockAPI(t, responseMap{"/zones": handleZoneList})
	doer := testServer.Client()
	client := NewClient(doer, SetEndpoint(testServer.URL))

	zones, resp, err := client.Zones.List()
	assert.NoError(t, err, "expected no error")
	assert.Equal(t, 200, resp.StatusCode, "status code does not match")

	expected := []string{
		"example.com",
		"coolzone.cool",
		"example.org",
		"example.io",
	}
	assert.Equal(t, 4, len(zones), "found wrong number of zones")
	for idx := range zones {
		assert.Equal(t, expected[idx], zones[idx].Zone, "zone.Zone does not match")
	}
}

func TestZone_ListWithPaginationDisabled(t *testing.T) {
	// It should not follow Links
	testServer := mockAPI(t, responseMap{"/zones": handleZoneList})
	doer := testServer.Client()
	client := NewClient(doer, SetEndpoint(testServer.URL), SetFollowPagination(false))

	zones, resp, err := client.Zones.List()
	assert.NoError(t, err, "expected no error")
	assert.Equal(t, 200, resp.StatusCode, "status code does not match")

	expected := []string{
		"example.com",
		"coolzone.cool",
	}
	assert.Equal(t, 2, len(zones))
	for idx := range zones {
		assert.Equal(t, expected[idx], zones[idx].Zone, "zone.Zone does not match")
	}
}

func TestZone_ListWithHTTPErrors(t *testing.T) {
	// It should have a non-nil resp, and it should have a JSON error if the
	// response body cannot be parsed as JSON, regardless of pagination.
	testServer := mockAPI(t, responseMap{
		"/zones": func(res http.ResponseWriter, req *http.Request) error {
			res.WriteHeader(429)
			res.Write([]byte("{}"))
			return nil
		},
	})
	doer := testServer.Client()
	client := NewClient(doer, SetEndpoint(testServer.URL))

	zones, resp, err := client.Zones.List()
	assert.IsType(t, &Error{}, err, "error type does not match")
	assert.Equal(t, 429, resp.StatusCode, "status code does not match")
	assert.Nil(t, zones, "zones is not nil")

	client = NewClient(doer, SetEndpoint(testServer.URL), SetFollowPagination(false))

	zones, resp, err = client.Zones.List()
	assert.IsType(t, &Error{}, err, "error type does not match")
	assert.Equal(t, 429, resp.StatusCode, "status code does not match")
	assert.Nil(t, zones, "zones is not nil")
}

func TestZone_ListWithNonHTTPErrors(t *testing.T) {
	// it should have a nil resp, regardless of pagination.
	client := NewClient(errorClient{}, SetEndpoint(""))

	zones, resp, err := client.Zones.List()
	assert.Nil(t, resp, "resp is not nil")
	assert.Error(t, err, "error was expected")
	assert.Nil(t, zones, "zones is not nil")

	client = NewClient(errorClient{}, SetEndpoint(""), SetFollowPagination(false))

	zones, resp, err = client.Zones.List()
	assert.Nil(t, resp, "resp is not nil")
	assert.Error(t, err, "error was expected")
	assert.Nil(t, zones, "zones is not nil")
}

func TestZone_Get(t *testing.T) {
	// It should follow Links and gather all Records into the Zone
	testServer := mockAPI(t, responseMap{"/zones/coolzone.cool": handleZoneGet})
	doer := testServer.Client()
	client := NewClient(doer, SetEndpoint(testServer.URL))

	zone, resp, err := client.Zones.Get("coolzone.cool")
	assert.NoError(t, err, "expected no error")
	assert.Equal(t, 200, resp.StatusCode, "status code does not match")

	expected := []string{
		"a.coolzone.cool",
		"b.coolzone.cool",
		"c.coolzone.cool",
		"d.coolzone.cool",
	}
	assert.Equal(t, "coolzone.cool", zone.Zone, "zone.Zone does not match")
	assert.Equal(t, 4, len(zone.Records), "found wrong number of records")
	for idx := range zone.Records {
		assert.Equal(t, expected[idx], zone.Records[idx].Domain)
	}

}

func TestZone_GetWithPaginationDisabled(t *testing.T) {
	// It should not follow Links
	testServer := mockAPI(t, responseMap{"/zones/coolzone.cool": handleZoneGet})
	doer := testServer.Client()
	client := NewClient(doer, SetEndpoint(testServer.URL), SetFollowPagination(false))

	zone, resp, err := client.Zones.Get("coolzone.cool")
	assert.NoError(t, err, "expected no error")
	assert.Equal(t, 200, resp.StatusCode, "status code does not match")

	expected := []string{
		"a.coolzone.cool",
		"b.coolzone.cool",
	}
	assert.Equal(t, "coolzone.cool", zone.Zone, "zone.Zone does not match")
	assert.Equal(t, 2, len(zone.Records), "found wrong number of records")
	for idx := range zone.Records {
		assert.Equal(
			t, expected[idx], zone.Records[idx].Domain, "record.Domain does not match",
		)
	}
}

func TestZone_GetWithHTTPErrors(t *testing.T) {
	// It should have a non-nil resp, and it should have a JSON error if the
	// response body cannot be parsed as JSON, regardless of pagination.
	testServer := mockAPI(t, responseMap{
		"/zones/coolzone.cool": func(res http.ResponseWriter, req *http.Request) error {
			res.WriteHeader(500)
			res.Write([]byte("Internal Server Error"))
			return nil
		},
	})
	doer := testServer.Client()
	client := NewClient(doer, SetEndpoint(testServer.URL))

	zone, resp, err := client.Zones.Get("coolzone.cool")
	assert.IsType(t, &json.SyntaxError{}, err, "error type does not match")
	assert.Equal(t, 500, resp.StatusCode, "status code does not match")
	assert.Nil(t, zone, "zone is not nil")

	client = NewClient(doer, SetEndpoint(testServer.URL), SetFollowPagination(false))

	zone, resp, err = client.Zones.Get("coolzone.cool")
	assert.IsType(t, &json.SyntaxError{}, err, "error type does not match")
	assert.Equal(t, 500, resp.StatusCode, "status code does not match")
	assert.Nil(t, zone, "zone is not nil")
}

func TestZone_GetWithNonHTTPErrors(t *testing.T) {
	// it should have a nil resp, regardless of pagination.
	client := NewClient(errorClient{}, SetEndpoint(""))

	zone, resp, err := client.Zones.Get("coolzone.cool")
	assert.Nil(t, resp, "resp is not nil")
	assert.Error(t, err, "error was expected")
	assert.Nil(t, zone, "zone is not nil")

	client = NewClient(errorClient{}, SetEndpoint(""), SetFollowPagination(false))

	zone, resp, err = client.Zones.Get("coolzone.cool")
	assert.Nil(t, resp, "resp is not nil")
	assert.Error(t, err, "error was expected")
	assert.Nil(t, zone, "zone is not nil")
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

type errorClient struct{}

func (c errorClient) Do(req *http.Request) (*http.Response, error) {
	return nil, errors.New("oops")
}

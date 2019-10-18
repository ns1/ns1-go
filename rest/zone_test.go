package rest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPagination_zoneList(t *testing.T) {
	// It should follow Links and gather all Zones into the list
	testServer := mockAPI(t)
	doer := testServer.Client()
	client := NewClient(doer, SetEndpoint(testServer.URL))

	zones, _, _ := client.Zones.List()

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

func TestPagination_zoneRecords(t *testing.T) {
	// It should follow Links and gather all Records into the Zone
	testServer := mockAPI(t)
	doer := testServer.Client()
	client := NewClient(doer, SetEndpoint(testServer.URL))

	zone, _, _ := client.Zones.Get("coolzone.cool")

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

func TestPagination_disabled(t *testing.T) {
	// It should not follow Links
	testServer := mockAPI(t)
	doer := testServer.Client()
	client := NewClient(doer, SetEndpoint(testServer.URL), SetFollowPagination(false))

	zones, _, _ := client.Zones.List()
	expected := []string{
		"example.com",
		"coolzone.cool",
	}
	assert.Equal(t, 2, len(zones))
	for idx := range zones {
		assert.Equal(t, expected[idx], zones[idx].Zone)
	}

	zone, _, _ := client.Zones.Get("coolzone.cool")
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

const zoneListPath = "/zones"
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

const zoneGetPath = "/zones/coolzone.cool"
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

func mockAPI(t *testing.T) *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			var body string
			switch path := req.URL.Path; path {

				case zoneListPath:
					if req.URL.RawQuery == "" {
						res.Header().Set("Link", zoneListLinkHeader)
						body = zoneListPageOne
					} else if req.URL.RawQuery == zoneListQuery {
						body = zoneListPageTwo
					} else {
						t.Fatalf("unexpected query: %s", req.URL.RawQuery)
					}

				case zoneGetPath:
					if req.URL.RawQuery == "" {
						res.Header().Set("Link", zoneGetLinkHeader)
						body = zoneGetPageOne
					} else if req.URL.RawQuery == zoneGetQuery {
						body = zoneGetPageTwo
					} else {
						t.Fatalf("unexpected query: %s", req.URL.RawQuery)
					}

				default:
					t.Fatalf("should not hit default case")
			}

			res.WriteHeader(200)
			res.Write([]byte(body))
		}),
	)
}

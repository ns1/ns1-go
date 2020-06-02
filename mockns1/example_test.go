package mockns1_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/ns1/ns1-go.v2/mockns1"
	api "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/dns"
)

func Example() {
	t := new(testing.T)

	// Setup the mock service
	mock, doer, err := mockns1.New(t)
	require.Nil(t, err)

	defer mock.Shutdown()

	// Create your NS1 client and configure it for the mock service
	ns1 := api.NewClient(doer, api.SetAPIKey("apikey"))
	ns1.Endpoint, _ = url.Parse("https://" + mock.Address + "/v1/")

	// Add your test case (zone list in this example)
	require.Nil(t, mock.AddTestCase(http.MethodGet, "zones", http.StatusOK, nil, nil, "",
		[]*dns.Zone{{Zone: "foo.bar"}}))

	// Perform your tests
	zones, _, err := ns1.Zones.List()
	require.Nil(t, err)
	require.Equal(t, 1, len(zones))
	require.Equal(t, "foo.bar", zones[0].Zone)

	// Output:
}

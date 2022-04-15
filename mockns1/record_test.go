package mockns1_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/ns1/ns1-go.v2/mockns1"
	api "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/dns"
)

var (
	zone   = "zone.test"
	domain = "domain.zone.test"
	rType  = "CNAME"
)

func TestRecord(t *testing.T) {
	mock, doer, err := mockns1.New(t)
	require.Nil(t, err)
	defer mock.Shutdown()

	client := api.NewClient(doer, api.SetEndpoint("https://"+mock.Address+"/v1/"))

	t.Run("AddRecordGetTestCase", func(t *testing.T) {
		response := &dns.Record{Zone: zone, Domain: domain, Type: rType}

		require.Nil(t,
			mock.AddRecordGetTestCase(zone, domain, rType, nil, nil, response),
		)

		record, _, err := client.Records.Get(zone, domain, rType)
		require.Nil(t, err)
		require.NotNil(t, record)
		require.Equal(t, record.Zone, zone)
		require.Equal(t, record.Domain, domain)
		require.Equal(t, record.Type, rType)
	})
}

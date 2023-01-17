package mockns1_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/ns1/ns1-go.v2/mockns1"
	api "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/dns"
)

func TestNetwork(t *testing.T) {
	mock, doer, err := mockns1.New(t)
	require.Nil(t, err)
	defer mock.Shutdown()

	client := api.NewClient(doer, api.SetEndpoint("https://"+mock.Address+"/v1/"))

	t.Run("NetworkGetTestCase", func(t *testing.T) {
		networks := []*dns.Network{
			{
				Name:      "a.name",
				NetworkID: 0,
				Label:     "a",
			},
			{
				Name:      "b.name",
				NetworkID: 1,
				Label:     "b",
			},
			{
				Name:      "c.name",
				NetworkID: 2,
				Label:     "c",
			},
		}

		require.Nil(t, mock.NetworkGetTestCase(nil, nil, networks))

		resp, _, err := client.Network.Get()
		require.Nil(t, err)
		require.NotNil(t, resp)
		require.Equal(t, len(networks), len(resp))

		for i := range networks {
			require.Equal(t, networks[i].Name, resp[i].Name, i)
			require.Equal(t, networks[i].NetworkID, resp[i].NetworkID, i)
			require.Equal(t, networks[i].Label, resp[i].Label, i)
		}
	})
}

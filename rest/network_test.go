package rest_test

import (
	"context"
	"net/http"
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

	t.Run("Get", func(t *testing.T) {

		t.Run("Standard", func(t *testing.T) {
			defer mock.ClearTestCases()

			client.FollowPagination = false
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

			respNetworks, _, err := client.Network.Get(context.Background())
			require.Nil(t, err)
			require.NotNil(t, respNetworks)
			require.Equal(t, len(networks), len(respNetworks))

			for i := range networks {
				require.Equal(t, networks[i].Name, respNetworks[i].Name, i)
				require.Equal(t, networks[i].NetworkID, respNetworks[i].NetworkID, i)
				require.Equal(t, networks[i].Label, respNetworks[i].Label, i)
			}
		})
		t.Run("Empty", func(t *testing.T) {
			defer mock.ClearTestCases()

			client.FollowPagination = false
			networks := []*dns.Network{}

			require.Nil(t, mock.NetworkGetTestCase(nil, nil, networks))

			respNetworks, _, err := client.Network.Get(context.Background())
			require.Nil(t, err)
			require.Equal(t, networks, respNetworks)

		})

		t.Run("Error", func(t *testing.T) {
			t.Run("HTTP", func(t *testing.T) {
				defer mock.ClearTestCases()

				require.Nil(t, mock.AddTestCase(
					http.MethodGet, "/networks", http.StatusNotFound,
					nil, nil, "", `{"message": "test error"}`,
				))

				networks, resp, err := client.Network.Get(context.Background())
				require.Nil(t, networks)
				require.NotNil(t, err)
				require.Contains(t, err.Error(), "test error")
				require.Equal(t, http.StatusNotFound, resp.StatusCode)
			})

			t.Run("Other", func(t *testing.T) {
				c := api.NewClient(errorClient{}, api.SetEndpoint(""))
				networks, resp, err := c.Network.Get(context.Background())
				require.Nil(t, resp)
				require.Error(t, err)
				require.Nil(t, networks)
			})
		})
	})

}

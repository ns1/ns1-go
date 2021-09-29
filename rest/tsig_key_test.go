package rest_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/ns1/ns1-go.v2/mockns1"

	api "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/dns"
)

func TestTsigKey(t *testing.T) {
	mock, doer, err := mockns1.New(t)

	require.Nil(t, err)
	defer mock.Shutdown()

	client := api.NewClient(doer, api.SetEndpoint("https://"+mock.Address+"/v1/"))

	// Tests for api.Client.TSIG.List()
	t.Run("List", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			tsigKeys := []*dns.Tsig_key{
				{
					Name:      "TsigKey1",
					Algorithm: "hmac-sha256",
					Secret:    "Ok1qR5IW1ajVka5cHPEJQIXfLyx5V3PSkFBROAzOn21JumDq6nIpoj6H8rfj5Uo+Ok55ZWQ0Wgrf302fDscHLA==",
				},
				{
					Name:      "TsigKey2",
					Algorithm: "hmac-sha256",
					Secret:    "1k1qR5IW1ajVka5cHPEJQIXfLyx5V3PSkFBROAzOn21JumDq6nIpoj6H8rfj5Uo+Ok55ZWQ0Wgrf302fDscHLA==",
				},
			}

			require.Nil(t, mock.AddTsigKeyListTestCase(nil, nil, tsigKeys))

			respTsigKeys, _, err := client.TSIG.List()
			require.Nil(t, err)
			require.NotNil(t, respTsigKeys)
			require.Equal(t, len(tsigKeys), len(respTsigKeys))

			for i := range tsigKeys {
				require.Equal(t, tsigKeys[i].Name, respTsigKeys[i].Name, i)
				require.Equal(t, tsigKeys[i].Algorithm, respTsigKeys[i].Algorithm, i)
				require.Equal(t, tsigKeys[i].Secret, respTsigKeys[i].Secret, i)
			}
		})

		t.Run("Error", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddTestCase(
				http.MethodGet, "tsig", http.StatusNotFound,
				nil, nil, "", `{"message": "test error"}`,
			))

			tsigKeys, resp, err := client.TSIG.List()
			require.Nil(t, tsigKeys)
			require.NotNil(t, err)
			require.Contains(t, err.Error(), "test error")
			require.Equal(t, http.StatusNotFound, resp.StatusCode)
		})
	})
}

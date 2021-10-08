package rest_test

import (
	"fmt"
	"net/http"
	"reflect"
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

			tsigKeys := []*dns.TSIGKey{
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

	// Tests for api.Client.TSIG.Get()
	t.Run("Get", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			tsigKey := &dns.TSIGKey{
				Name:      "TsigKey1",
				Algorithm: "hmac-sha256",
				Secret:    "Ok1qR5IW1ajVka5cHPEJQIXfLyx5V3PSkFBROAzOn21JumDq6nIpoj6H8rfj5Uo+Ok55ZWQ0Wgrf302fDscHLA==",
			}

			require.Nil(t, mock.AddTsigKeyGetTestCase("TsigKey1", nil, nil, tsigKey))

			respTsigKey, _, err := client.TSIG.Get("TsigKey1")

			require.Nil(t, err)
			require.True(t, reflect.DeepEqual(tsigKey, respTsigKey))
		})

		// Error TSIG key does not exist
		t.Run("TSIG key does not exist", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddTestCase(
				http.MethodGet, fmt.Sprintf("/tsig/%s", "TsigKey1"), http.StatusNotFound,
				nil, nil, "", `{"message": "TSIG key does not exist"}`,
			))
			tsigKey, resp, err := client.TSIG.Get("TsigKey1")
			require.Nil(t, tsigKey)
			require.NotNil(t, err)
			require.Contains(t, err.Error(), api.ErrTsigKeyMissing.Error())
			require.Equal(t, http.StatusNotFound, resp.StatusCode)
		})
	})

	t.Run("Create", func(t *testing.T) {
		tsigKey := &dns.TSIGKey{
			Name:      "TsigKey1",
			Algorithm: "hmac-sha256",
			Secret:    "Ok1qR5IW1ajVka5cHPEJQIXfLyx5V3PSkFBROAzOn21JumDq6nIpoj6H8rfj5Uo+Ok55ZWQ0Wgrf302fDscHLA==",
		}

		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddTsigKeyCreateTestCase(nil, nil, tsigKey, tsigKey))

			_, err := client.TSIG.Create(tsigKey)
			require.Nil(t, err)
		})

		t.Run("Error", func(t *testing.T) {
			// TSIG key already exists
			t.Run("TSIG key already exists", func(t *testing.T) {
				defer mock.ClearTestCases()

				require.Nil(t, mock.AddTestCase(
					http.MethodPut, fmt.Sprintf("tsig/%s", tsigKey.Name), http.StatusNotFound,
					nil, nil, tsigKey, `{"message": "TSIG key already exists"}`,
				))

				_, err = client.TSIG.Create(tsigKey)

				require.Contains(t, err.Error(), api.ErrTsigKeyExists.Error())
			})
		})
	})

	t.Run("Update", func(t *testing.T) {
		tsigKey := &dns.TSIGKey{
			Name:      "TsigKey1",
			Algorithm: "hmac-sha256",
			Secret:    "Ok1qR5IW1ajVka5cHPEJQIXfLyx5V3PSkFBROAzOn21JumDq6nIpoj6H8rfj5Uo+Ok55ZWQ0Wgrf302fDscHLA==",
		}

		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddTsigKeyUpdateTestCase(nil, nil, tsigKey, tsigKey))

			_, err := client.TSIG.Update(tsigKey)
			require.Nil(t, err)
		})

		t.Run("Error", func(t *testing.T) {
			// TSIG key does not exist
			t.Run("TSIG key does not exist", func(t *testing.T) {
				defer mock.ClearTestCases()

				require.Nil(t, mock.AddTestCase(
					http.MethodPost, fmt.Sprintf("/tsig/%s", tsigKey.Name), http.StatusNotFound,
					nil, nil, tsigKey, `{"message": "TSIG key does not exist"}`,
				))
				resp, err := client.TSIG.Update(tsigKey)

				require.NotNil(t, err)
				require.Contains(t, err.Error(), api.ErrTsigKeyMissing.Error())
				require.Equal(t, http.StatusNotFound, resp.StatusCode)
			})
		})
	})

	t.Run("Delete", func(t *testing.T) {
		tsigKey := &dns.TSIGKey{
			Name:      "TsigKey1",
			Algorithm: "hmac-sha256",
			Secret:    "Ok1qR5IW1ajVka5cHPEJQIXfLyx5V3PSkFBROAzOn21JumDq6nIpoj6H8rfj5Uo+Ok55ZWQ0Wgrf302fDscHLA==",
		}

		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddTsigKeyDeleteTestCase(nil, nil, tsigKey, nil))

			_, err := client.TSIG.Delete(tsigKey.Name)
			require.Nil(t, err)
		})

		t.Run("Error", func(t *testing.T) {
			// Error TSIG key does not exist
			t.Run("TSIG key does not exist", func(t *testing.T) {
				defer mock.ClearTestCases()

				require.Nil(t, mock.AddTestCase(
					http.MethodDelete, fmt.Sprintf("tsig/%s", tsigKey.Name), http.StatusNotFound,
					nil, nil, "", `{"message": "TSIG key does not exist"}`,
				))
				resp, err := client.TSIG.Delete(tsigKey.Name)

				require.NotNil(t, err)
				require.Contains(t, err.Error(), api.ErrTsigKeyMissing.Error())
				require.Equal(t, http.StatusNotFound, resp.StatusCode)
			})
		})
	})
}

package rest_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/ns1/ns1-go.v2/mockns1"
	api "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/account"
)

func TestGlobalIPWhitelist(t *testing.T) {
	mock, doer, err := mockns1.New(t)
	require.Nil(t, err)
	defer mock.Shutdown()
	client := api.NewClient(doer, api.SetEndpoint("https://"+mock.Address+"/v1/"))

	t.Run("List", func(t *testing.T) {
		whitelists := []*account.IPWhitelist{
			{
				ID:   "1",
				Name: "First Whitelist",
				Values: []string{
					"1.2.3.4",
				},
			},
			{
				ID:   "2",
				Name: "Second Whitelist",
				Values: []string{
					"1.2.3.4",
				},
			},
		}

		t.Run("list global IP whitelists", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddGlobalIPWhitelistListTestCase(nil, nil, whitelists))

			respWhitelists, _, err := client.GlobalIPWhitelist.List()
			require.Nil(t, err)
			require.NotNil(t, respWhitelists)
			require.Equal(t, len(whitelists), len(respWhitelists))

			for i := range whitelists {
				require.Equal(t, whitelists[i].Name, respWhitelists[i].Name, i)
				require.Equal(t, whitelists[i].ID, respWhitelists[i].ID, i)
				require.Equal(t, whitelists[i].Values, respWhitelists[i].Values, i)
			}
		})
	})

	t.Run("Get", func(t *testing.T) {
		whitelist := &account.IPWhitelist{
			ID:   "1",
			Name: "First Whitelist",
			Values: []string{
				"1.2.3.4",
			},
		}

		t.Run("get global IP whitelist", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddGlobalIPWhitelistGetTestCase("1", nil, nil, whitelist))

			respWhitelist, _, err := client.GlobalIPWhitelist.Get("1")
			require.Nil(t, err)
			require.NotNil(t, respWhitelist)

			require.Equal(t, whitelist.Name, respWhitelist.Name)
			require.Equal(t, whitelist.ID, respWhitelist.ID)
			require.Equal(t, whitelist.Values, respWhitelist.Values)

		})
	})

	t.Run("Create", func(t *testing.T) {
		whitelist := &account.IPWhitelist{
			ID:   "1",
			Name: "First Whitelist",
			Values: []string{
				"1.2.3.4",
			},
		}

		t.Run("create global IP whitelist", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddGlobalIPWhitelistCreateTestCase(nil, nil, whitelist, whitelist))

			_, err := client.GlobalIPWhitelist.Create(whitelist)
			require.Nil(t, err)
		})
	})

	t.Run("Update", func(t *testing.T) {
		defer mock.ClearTestCases()

		whitelist := &account.IPWhitelist{
			ID:   "1",
			Name: "First Whitelist",
			Values: []string{
				"1.2.3.4",
			},
		}

		t.Run("update global IP whitelist", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddGlobalIPWhitelistUpdateTestCase(nil, nil, whitelist, whitelist))

			_, err := client.GlobalIPWhitelist.Update(whitelist)
			require.Nil(t, err)

		})
	})

	t.Run("Delete", func(t *testing.T) {
		defer mock.ClearTestCases()

		whitelist := &account.IPWhitelist{
			ID:   "1",
			Name: "First Whitelist",
			Values: []string{
				"1.2.3.4",
			},
		}

		require.Nil(t, mock.AddGlobalIPWhitelistDeleteTestCase(whitelist.ID, nil, nil))

		_, err := client.GlobalIPWhitelist.Delete(whitelist.ID)
		require.Nil(t, err)
	})
}

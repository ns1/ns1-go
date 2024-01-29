package rest_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/ns1/ns1-go.v2/mockns1"
	api "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/redirect"
)

func TestRedirectService(t *testing.T) {
	mock, doer, err := mockns1.New(t)
	require.Nil(t, err)
	defer mock.Shutdown()

	client := api.NewClient(doer, api.SetEndpoint("https://"+mock.Address+"/v1/"))

	id := "87461586-c681-43b7-b283-f840be95e13c"

	t.Run("List", func(t *testing.T) {
		t.Run("Pagination", func(t *testing.T) {
			defer mock.ClearTestCases()

			client.FollowPagination = true
			cfgList := &redirect.ConfigurationList{
				Count: 4,
				Total: 4,
				Results: []*redirect.Configuration{
					{Domain: "a.com"},
					{Domain: "b.com"},
					{Domain: "c.com"},
					{Domain: "d.com"},
				}}
			require.Nil(t, mock.AddRedirectListTestCase(nil, nil, cfgList))

			respCfgs, _, err := client.Redirects.List()
			require.Nil(t, err)
			require.NotNil(t, respCfgs)
			require.Equal(t, len(cfgList.Results), len(respCfgs))

			for i := range cfgList.Results {
				require.Equal(t, cfgList.Results[i].Domain, respCfgs[i].Domain, i)
			}
		})

		t.Run("No Pagination", func(t *testing.T) {
			defer mock.ClearTestCases()

			client.FollowPagination = false
			cfgList := &redirect.ConfigurationList{
				Count: 2,
				Total: 4,
				Results: []*redirect.Configuration{
					{Domain: "a.com"},
					{Domain: "b.com"},
				}}

			header := http.Header{}
			header.Set("Link", `</redirect?after=b.com&limit=2>; rel="next"`)

			require.Nil(t, mock.AddRedirectListTestCase(nil, header, cfgList))

			respCfgs, resp, err := client.Redirects.List()
			require.Nil(t, err)
			require.NotNil(t, respCfgs)
			require.Equal(t, len(cfgList.Results), len(respCfgs))
			require.Contains(t, resp.Header.Get("Link"), "redirect?after=b.com")

			for i := range cfgList.Results {
				require.Equal(t, cfgList.Results[i].Domain, respCfgs[i].Domain, i)
			}
		})

		t.Run("Error", func(t *testing.T) {
			t.Run("HTTP", func(t *testing.T) {
				defer mock.ClearTestCases()

				require.Nil(t, mock.AddTestCase(
					http.MethodGet, "/redirect", http.StatusNotFound,
					nil, nil, "", `{"message": "test error"}`,
				))

				respCfgs, resp, err := client.Redirects.List()
				require.Nil(t, respCfgs)
				require.NotNil(t, err)
				require.Contains(t, err.Error(), "test error")
				require.Equal(t, http.StatusNotFound, resp.StatusCode)
			})

			t.Run("Other", func(t *testing.T) {
				c := api.NewClient(errorClient{}, api.SetEndpoint(""))
				respCfgs, resp, err := c.Redirects.List()
				require.Nil(t, resp)
				require.Error(t, err)
				require.Nil(t, respCfgs)
			})
		})
	})

	t.Run("Get", func(t *testing.T) {

		t.Run("Pagination", func(t *testing.T) {
			defer mock.ClearTestCases()

			client.FollowPagination = true
			cfg := &redirect.Configuration{
				Domain: "a.com",
			}
			require.Nil(t, mock.AddRedirectGetTestCase(id, nil, nil, cfg))

			respCfg, _, err := client.Redirects.Get(id)
			require.Nil(t, err)
			require.NotNil(t, respCfg)
			require.Equal(t, cfg.Domain, respCfg.Domain)
		})

		t.Run("Error", func(t *testing.T) {
			t.Run("HTTP", func(t *testing.T) {
				defer mock.ClearTestCases()

				require.Nil(t, mock.AddTestCase(
					http.MethodGet, "/redirect/"+id, http.StatusNotFound,
					nil, nil, "", `{"message": "test error"}`,
				))

				cfg, resp, err := client.Redirects.Get(id)
				require.Nil(t, cfg)
				require.NotNil(t, err)
				require.Contains(t, err.Error(), "test error")
				require.Equal(t, http.StatusNotFound, resp.StatusCode)
			})

			t.Run("Other", func(t *testing.T) {
				c := api.NewClient(errorClient{}, api.SetEndpoint(""))
				cfg, resp, err := c.Redirects.Get(id)
				require.Nil(t, resp)
				require.Error(t, err)
				require.Nil(t, cfg)
			})
		})
	})

	t.Run("Create", func(t *testing.T) {
		cfg := &redirect.Configuration{
			ID:     &id,
			Domain: "a.com",
		}

		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddRedirectCreateTestCase(nil, nil, cfg, cfg))

			_, _, err := client.Redirects.Create(cfg)
			require.Nil(t, err)
		})

		t.Run("Error", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddTestCase(
				http.MethodPut, "/redirect", http.StatusConflict,
				nil, nil, cfg, `{"message": "configuration already exists"}`,
			))

			_, _, err := client.Redirects.Create(cfg)
			require.Equal(t, api.ErrRedirectExists, err)
		})
	})

	t.Run("Update", func(t *testing.T) {
		cfg := &redirect.Configuration{
			ID:     &id,
			Domain: "a.com",
		}

		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddRedirectUpdateTestCase(nil, nil, cfg, cfg))

			_, _, err := client.Redirects.Update(cfg)
			require.Nil(t, err)
		})

		t.Run("Error", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddTestCase(
				http.MethodPost, "/redirect/"+id, http.StatusNotFound,
				nil, nil, cfg, `{"message": "configuration not found"}`,
			))

			_, _, err := client.Redirects.Update(cfg)
			require.Equal(t, api.ErrRedirectNotFound, err)
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddRedirectDeleteTestCase(id, nil, nil))

			_, err := client.Redirects.Delete(id)
			require.Nil(t, err)
		})

		t.Run("Error", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddTestCase(
				http.MethodDelete, "/redirect/"+id, http.StatusNotFound,
				nil, nil, "", `{"message": "configuration not found"}`,
			))

			_, err := client.Redirects.Delete(id)
			require.Equal(t, api.ErrRedirectNotFound, err)
		})
	})
}

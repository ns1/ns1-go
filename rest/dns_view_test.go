package rest_test

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/ns1/ns1-go.v2/mockns1"

	api "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/dns"
)

func TestDNSView(t *testing.T) {
	mock, doer, err := mockns1.New(t)

	require.Nil(t, err)
	defer mock.Shutdown()

	client := api.NewClient(doer, api.SetEndpoint("https://"+mock.Address+"/v1/"))

	// Tests for api.Client.View.List()
	t.Run("List", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			views := []*dns.DNSView{
				{
					Name:       "DNSView1",
					Preference: 1,
				},
				{
					Name:       "DNSView2",
					Preference: 2,
				},
			}

			require.Nil(t, mock.AddDNSViewListTestCase(nil, nil, views))

			respDNSViews, _, err := client.View.List(context.Background())
			require.Nil(t, err)
			require.NotNil(t, respDNSViews)
			require.Equal(t, len(views), len(respDNSViews))

			for i := range views {
				require.Equal(t, views[i].Name, respDNSViews[i].Name, i)
				require.Equal(t, views[i].Preference, respDNSViews[i].Preference, i)
			}
		})

		t.Run("Error", func(t *testing.T) {
			// Other errors
			t.Run("Other errors", func(t *testing.T) {
				defer mock.ClearTestCases()

				require.Nil(t, mock.AddTestCase(
					http.MethodGet, "views", http.StatusNotFound,
					nil, nil, "", `{"message": "test error"}`,
				))

				views, resp, err := client.View.List(context.Background())
				require.Nil(t, views)
				require.NotNil(t, err)
				require.Contains(t, err.Error(), "test error")
				require.Equal(t, http.StatusNotFound, resp.StatusCode)
			})
		})
	})

	// Tests for api.Client.View.Get()
	t.Run("Get", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			dnsView := myView

			require.Nil(t, mock.AddDNSViewGetTestCase(myView.Name, nil, nil, &dnsView))

			respDNSView, _, err := client.View.Get(context.Background(), myView.Name)

			require.Nil(t, err)
			require.True(t, reflect.DeepEqual(&dnsView, respDNSView))
		})

		t.Run("Error", func(t *testing.T) {
			// Error DNS View not found
			t.Run("DNS View not found", func(t *testing.T) {
				defer mock.ClearTestCases()

				require.Nil(t, mock.AddTestCase(
					http.MethodGet, "views/myView", http.StatusNotFound,
					nil, nil, "", `{"message": "Resource not found"}`,
				))
				dnsView, resp, err := client.View.Get(context.Background(), "myView")
				require.Nil(t, dnsView)
				require.NotNil(t, err)
				require.Equal(t, api.ErrViewMissing.Error(), err.Error())
				require.Equal(t, http.StatusNotFound, resp.StatusCode)
			})

			// Other errors
			t.Run("Other errors", func(t *testing.T) {
				defer mock.ClearTestCases()

				require.Nil(t, mock.AddTestCase(
					http.MethodGet, "views/myView", http.StatusBadRequest,
					nil, nil, "", `{"message": "test error"}`,
				))
				dnsView, resp, err := client.View.Get(context.Background(), "myView")
				require.Nil(t, dnsView)
				require.NotNil(t, err)
				require.Contains(t, err.Error(), "test error")
				require.Equal(t, http.StatusBadRequest, resp.StatusCode)
			})
		})
	})

	// Test for api.Client.View.Create()
	t.Run("Create", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddDNSViewCreateTestCase(nil, nil, &myView, &myView))

			_, err := client.View.Create(context.Background(), &myView)
			require.Nil(t, err)
		})

		t.Run("Error", func(t *testing.T) {
			// DNS View already exists
			t.Run("DNS View already exists", func(t *testing.T) {
				defer mock.ClearTestCases()

				dnsView := myView
				require.Nil(t, mock.AddTestCase(
					http.MethodPut, fmt.Sprintf("views/%s", myView.Name), http.StatusConflict,
					nil, nil, dnsView, `{"message": "conflicts with existing resource"}`,
				))

				_, err = client.View.Create(context.Background(), &dnsView)

				require.Equal(t, api.ErrViewExists.Error(), err.Error())
			})

			// Other errors
			t.Run("Other errors", func(t *testing.T) {
				defer mock.ClearTestCases()

				dnsView := myView
				require.Nil(t, mock.AddTestCase(
					http.MethodPut, fmt.Sprintf("views/%s", myView.Name), http.StatusBadGateway,
					nil, nil, dnsView, `{"message": "test error"}`,
				))

				_, err = client.View.Create(context.Background(), &dnsView)

				require.Contains(t, err.Error(), "test error")
			})
		})
	})

	// Test for api.Client.View.Update()
	t.Run("Update", func(t *testing.T) {
		dnsView := myView

		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddDNSViewUpdateTestCase(nil, nil, &dnsView, &dnsView))

			_, err := client.View.Update(context.Background(), &dnsView)
			require.Nil(t, err)
		})

		t.Run("Error", func(t *testing.T) {
			// Error DNS View not found
			t.Run("Resource not found", func(t *testing.T) {
				defer mock.ClearTestCases()

				require.Nil(t, mock.AddTestCase(
					http.MethodPost, fmt.Sprintf("views/%s", dnsView.Name), http.StatusNotFound,
					nil, nil, dnsView, `{"message": "Resource not found"}`,
				))
				resp, err := client.View.Update(context.Background(), &dnsView)
				require.NotNil(t, err)
				require.Equal(t, api.ErrViewMissing.Error(), err.Error())
				require.Equal(t, http.StatusNotFound, resp.StatusCode)
			})

			// Other errors
			t.Run("Other errors", func(t *testing.T) {
				defer mock.ClearTestCases()

				require.Nil(t, mock.AddTestCase(
					http.MethodPost, fmt.Sprintf("views/%s", dnsView.Name), http.StatusBadGateway,
					nil, nil, dnsView, `{"message": "test error"}`,
				))
				resp, err := client.View.Update(context.Background(), &dnsView)

				require.NotNil(t, err)
				require.Contains(t, err.Error(), "test error")
				require.Equal(t, http.StatusBadGateway, resp.StatusCode)
			})
		})
	})

	// Test for api.Client.View.GetPreferences()
	t.Run("GetPreferences", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()
			require.Nil(t, mock.AddDNSViewGetPreferencesTestCase(nil, nil, myMap))
			respMap, _, err := client.View.GetPreferences(context.Background())
			require.Nil(t, err)
			require.True(t, reflect.DeepEqual(myMap, respMap))
		})

		t.Run("Error", func(t *testing.T) {
			// Other errors
			t.Run("Other errors", func(t *testing.T) {
				defer mock.ClearTestCases()
				require.Nil(t, mock.AddTestCase(
					http.MethodGet, "config/views/preference", http.StatusBadGateway,
					nil, nil, "", `{"message": "test error"}`,
				))
				m, resp, err := client.View.GetPreferences(context.Background())
				require.Nil(t, m)
				require.NotNil(t, err)
				require.Contains(t, err.Error(), "test error")
				require.Equal(t, http.StatusBadGateway, resp.StatusCode)
			})
		})
	})

	// Test for api.Client.View.UpdatePreferences()
	t.Run("UpdatePreferences", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()
			require.Nil(t, mock.AddDNSViewUpdatePreferencesTestCase(nil, nil, myMap, myMap))
			respMap, _, err := client.View.UpdatePreferences(context.Background(), myMap)
			require.Nil(t, err)
			require.True(t, reflect.DeepEqual(myMap, respMap))
		})

		t.Run("Error", func(t *testing.T) {
			// DNS View not found error
			t.Run("Not found", func(t *testing.T) {
				defer mock.ClearTestCases()
				require.Nil(t, mock.AddTestCase(
					http.MethodPost, "config/views/preference", http.StatusNotFound,
					nil, nil, myMap, `{"message": "Resource not found"}`,
				))
				m, resp, err := client.View.UpdatePreferences(context.Background(), myMap)
				require.Nil(t, m)
				require.NotNil(t, err)
				require.Equal(t, api.ErrViewMissing.Error(), err.Error())
				require.Equal(t, http.StatusNotFound, resp.StatusCode)
			})

			// Other errors
			t.Run("Other errors", func(t *testing.T) {
				defer mock.ClearTestCases()
				require.Nil(t, mock.AddTestCase(
					http.MethodPost, "config/views/preference", http.StatusBadGateway,
					nil, nil, myMap, `{"message": "test error"}`,
				))
				m, resp, err := client.View.UpdatePreferences(context.Background(), myMap)
				require.Nil(t, m)
				require.NotNil(t, err)
				require.Contains(t, err.Error(), "test error")
				require.Equal(t, http.StatusBadGateway, resp.StatusCode)
			})
		})
	})
}

var (
	myView = dns.DNSView{
		Name:       "myView",
		Created_at: 123456789,
		Updated_at: 987654321,
		Read_acls: []string{
			"myACLS1",
			"myACLS2",
		},
		Update_acls: []string{
			"myACLS1",
			"myACLS2",
			"myACLS3",
		},
		Zones: []string{
			"myZones1",
			"myZones2",
			"myZones3",
		},
		Networks: []int{
			1, 2, 4,
		},
	}

	myMap = map[string]int{
		"view1": 5,
		"view2": 4,
		"view3": 3,
		"view4": 2,
		"view5": 1,
	}
)

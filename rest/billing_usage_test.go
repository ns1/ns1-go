package rest_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/ns1/ns1-go.v2/mockns1"
	api "gopkg.in/ns1/ns1-go.v2/rest"
	billingusage "gopkg.in/ns1/ns1-go.v2/rest/model/billingusage"
)

func TestBillingUsageService(t *testing.T) {

	mock, doer, err := mockns1.New(t)
	require.Nil(t, err)
	defer mock.Shutdown()

	from := 1738368000
	to := 1740787199

	client := api.NewClient(doer, api.SetEndpoint("https://"+mock.Address+"/v1/"))

	// Tests for api.Client.BillingUsage.GetQueries
	t.Run("Get Billing Usage Queries", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()
			queries := &billingusage.Queries{
				CleanQueries: 0,
				DdosQueries:  0,
				NxdResponses: 0,
				ByNetwork:    []billingusage.QueriesByNetwork{},
			}
			require.Nil(t, mock.AddBillingUsageQueriesGetTestCase(int32(from), int32(to), nil, nil, queries))

			respDt, _, err := client.BillingUsage.GetQueries(int32(from), int32(to))
			require.Nil(t, err)
			require.NotNil(t, respDt)
			require.Equal(t, queries.CleanQueries, respDt.CleanQueries)
		})

		t.Run("Error Get Billing Usage Queries", func(t *testing.T) {
			t.Run("HTTP", func(t *testing.T) {
				defer mock.ClearTestCases()

				require.Nil(t, mock.AddBillingUsageQueriesFailTestCase(
					http.MethodGet, int32(from), int32(to), http.StatusNotFound,
					nil, nil, `{"message": "not found"}`,
				))

				dt, resp, err := client.BillingUsage.GetQueries(int32(from), int32(to))
				require.Nil(t, dt)
				require.NotNil(t, err)
				require.Contains(t, err.Error(), "not found")
				require.Equal(t, http.StatusNotFound, resp.StatusCode)
			})

			t.Run("Other", func(t *testing.T) {
				c := api.NewClient(errorClient{}, api.SetEndpoint(""))
				dt, resp, err := c.BillingUsage.GetQueries(int32(from), int32(to))
				require.Nil(t, resp)
				require.Error(t, err)
				require.Nil(t, dt)
			})
		})
	})

	// Tests for api.Client.BillingUsage.GetLimits
	t.Run("Get Billing Usage Limits", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()
			limits := &billingusage.Limits{
				QueriesLimit:          0,
				ChinaQueriesLimit:     0,
				RecordsLimit:          0,
				FilterChainsLimit:     0,
				MonitorsLimit:         0,
				DecisionsLimit:        0,
				NxdProtectionEnabled:  false,
				DdosProtectionEnabled: false,
				IncludeDedicatedDnsNetworkInManagedDnsUsage: false,
			}
			require.Nil(t, mock.AddBillingUsageLimitsGetTestCase(int32(from), int32(to), nil, nil, limits))

			buLimits, _, err := client.BillingUsage.GetLimits(int32(from), int32(to))
			require.Nil(t, err)
			require.NotNil(t, buLimits)
			require.Equal(t, limits, buLimits)
		})

		t.Run("Error Get Billing Usage Limits", func(t *testing.T) {
			t.Run("HTTP", func(t *testing.T) {
				defer mock.ClearTestCases()

				require.Nil(t, mock.AddBillingUsageLimitsFailTestCase(
					http.MethodGet, int32(from), int32(to), http.StatusNotFound,
					nil, nil, `{"message": "not found"}`,
				))

				dt, resp, err := client.BillingUsage.GetLimits(int32(from), int32(to))
				require.Nil(t, dt)
				require.NotNil(t, err)
				require.Contains(t, err.Error(), "not found")
				require.Equal(t, http.StatusNotFound, resp.StatusCode)
			})

			t.Run("Other", func(t *testing.T) {
				c := api.NewClient(errorClient{}, api.SetEndpoint(""))
				dt, resp, err := c.BillingUsage.GetLimits(int32(from), int32(to))
				require.Nil(t, resp)
				require.Error(t, err)
				require.Nil(t, dt)
			})
		})
	})

	// Tests for api.Client.BillingUsage.GetDecisions
	t.Run("Get Billing Usage Decisions", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()
			decisions := &billingusage.TotalUsage{
				TotalUsage: 0,
			}
			require.Nil(t, mock.AddBillingUsageDecisionsGetTestCase(int32(from), int32(to), nil, nil, decisions))

			buDecisions, _, err := client.BillingUsage.GetDecisions(int32(from), int32(to))
			require.Nil(t, err)
			require.NotNil(t, buDecisions)
			require.Equal(t, decisions.TotalUsage, buDecisions.TotalUsage)
		})

		t.Run("Error Get Billing Usage Decisions", func(t *testing.T) {
			t.Run("HTTP", func(t *testing.T) {
				defer mock.ClearTestCases()

				require.Nil(t, mock.AddBillingUsageDecisionsFailTestCase(
					http.MethodGet, int32(from), int32(to), http.StatusNotFound,
					nil, nil, `{"message": "not found"}`,
				))

				dt, resp, err := client.BillingUsage.GetDecisions(int32(from), int32(to))
				require.Nil(t, dt)
				require.NotNil(t, err)
				require.Contains(t, err.Error(), "not found")
				require.Equal(t, http.StatusNotFound, resp.StatusCode)
			})

			t.Run("Other", func(t *testing.T) {
				c := api.NewClient(errorClient{}, api.SetEndpoint(""))
				dt, resp, err := c.BillingUsage.GetDecisions(int32(from), int32(to))
				require.Nil(t, resp)
				require.Error(t, err)
				require.Nil(t, dt)
			})
		})
	})

	// Tests for api.Client.BillingUsage.GetMonitors
	t.Run("Get Billing Usage Monitors", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()
			records := &billingusage.TotalUsage{
				TotalUsage: 0,
			}
			require.Nil(t, mock.AddBillingUsageMonitorsGetTestCase(nil, nil, records))

			buMonitors, _, err := client.BillingUsage.GetMonitors()
			require.Nil(t, err)
			require.NotNil(t, buMonitors)
			require.Equal(t, records.TotalUsage, buMonitors.TotalUsage)
		})

		t.Run("Error Get Billing Usage Monitors", func(t *testing.T) {
			t.Run("HTTP", func(t *testing.T) {
				defer mock.ClearTestCases()

				require.Nil(t, mock.AddBillingUsageMonitorsFailTestCase(
					http.MethodGet, http.StatusNotFound,
					nil, nil, `{"message": "not found"}`,
				))

				buMonitors, resp, err := client.BillingUsage.GetMonitors()
				require.Nil(t, buMonitors)
				require.NotNil(t, err)
				require.Contains(t, err.Error(), "not found")
				require.Equal(t, http.StatusNotFound, resp.StatusCode)
			})

			t.Run("Other", func(t *testing.T) {
				c := api.NewClient(errorClient{}, api.SetEndpoint(""))
				dt, resp, err := c.BillingUsage.GetMonitors()
				require.Nil(t, resp)
				require.Error(t, err)
				require.Nil(t, dt)
			})
		})
	})

	// Tests for api.Client.BillingUsage.GetFilterChains
	t.Run("Get Billing Usage FilterChains", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()
			filterChains := &billingusage.TotalUsage{
				TotalUsage: 0,
			}
			require.Nil(t, mock.AddBillingUsageFilterChainsGetTestCase(nil, nil, filterChains))

			buFilterChains, _, err := client.BillingUsage.GetFilterChains()
			require.Nil(t, err)
			require.NotNil(t, buFilterChains)
			require.Equal(t, filterChains.TotalUsage, buFilterChains.TotalUsage)
		})

		t.Run("Error Get Billing Usage FilterChains", func(t *testing.T) {
			t.Run("HTTP", func(t *testing.T) {
				defer mock.ClearTestCases()

				require.Nil(t, mock.AddBillingUsageFilterChainsFailTestCase(
					http.MethodGet, http.StatusNotFound,
					nil, nil, `{"message": "not found"}`,
				))

				buFilterChains, resp, err := client.BillingUsage.GetFilterChains()
				require.Nil(t, buFilterChains)
				require.NotNil(t, err)
				require.Contains(t, err.Error(), "not found")
				require.Equal(t, http.StatusNotFound, resp.StatusCode)
			})

			t.Run("Other", func(t *testing.T) {
				c := api.NewClient(errorClient{}, api.SetEndpoint(""))
				dt, resp, err := c.BillingUsage.GetFilterChains()
				require.Nil(t, resp)
				require.Error(t, err)
				require.Nil(t, dt)
			})
		})
	})

	// Tests for api.Client.BillingUsage.GetRecords
	t.Run("Get Billing Usage Records", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()
			records := &billingusage.TotalUsage{
				TotalUsage: 0,
			}
			require.Nil(t, mock.AddBillingUsageRecordsGetTestCase(nil, nil, records))

			buRecords, _, err := client.BillingUsage.GetRecords()
			require.Nil(t, err)
			require.NotNil(t, buRecords)
			require.Equal(t, records.TotalUsage, buRecords.TotalUsage)
		})

		t.Run("Error Get Billing Usage Records", func(t *testing.T) {
			t.Run("HTTP", func(t *testing.T) {
				defer mock.ClearTestCases()

				require.Nil(t, mock.AddBillingUsageRecordFailTestCase(
					http.MethodGet, http.StatusNotFound,
					nil, nil, `{"message": "not found"}`,
				))

				buRecords, resp, err := client.BillingUsage.GetRecords()
				require.Nil(t, buRecords)
				require.NotNil(t, err)
				require.Contains(t, err.Error(), "not found")
				require.Equal(t, http.StatusNotFound, resp.StatusCode)
			})

			t.Run("Other", func(t *testing.T) {
				c := api.NewClient(errorClient{}, api.SetEndpoint(""))
				dt, resp, err := c.BillingUsage.GetRecords()
				require.Nil(t, resp)
				require.Error(t, err)
				require.Nil(t, dt)
			})
		})
	})

}

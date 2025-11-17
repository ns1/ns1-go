package rest_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/ns1/ns1-go.v2/mockns1"

	api "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/pulsar"
)

func TestPulsarDecision(t *testing.T) {
	mock, doer, err := mockns1.New(t)

	require.Nil(t, err)
	defer mock.Shutdown()

	client := api.NewClient(doer, api.SetEndpoint("https://"+mock.Address+"/v1/"))

	// Tests for api.Client.PulsarDecisions.GetDecisions()
	t.Run("GetDecisions", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			decisionsResponse := &pulsar.DecisionsResponse{
				Graphs: []*pulsar.DecisionsGraph{
					{
						Count: 2,
						GraphData: []*pulsar.DecisionsGraphData{
							{Timestamp: 1234567890, Count: 1},
							{Timestamp: 1234567891, Count: 1},
						},
						Avg: 1,
						Tags: &pulsar.Tags{
							JobID: myJobID,
						},
					},
				},
				Total: 2,
			}

			require.Nil(t, mock.AddTestCase(
				http.MethodGet, "/pulsar/query/decisions", http.StatusOK,
				nil, nil, "", decisionsResponse,
			))

			resp, _, err := client.PulsarDecisions.GetDecisions(nil)
			require.Nil(t, err)
			require.NotNil(t, resp)
			require.Equal(t, len(decisionsResponse.Graphs), len(resp.Graphs))
			require.Equal(t, decisionsResponse.Total, resp.Total)
		})
	})

	// Tests for api.Client.PulsarDecisions.GetDecisionsGraphRegion()
	t.Run("GetDecisionsGraphRegion", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			unit := "requests"
			graphRegionResponse := &pulsar.DecisionsGraphRegionResponse{
				Data: []*pulsar.DecisionsGraphRegionData{
					{
						Region: "us-east",
						Counts: []*pulsar.JobIDCounts{
							{JobID: myJobID, Count: 1000},
						},
					},
				},
				Unit: &unit,
			}

			require.Nil(t, mock.AddTestCase(
				http.MethodGet, "/pulsar/query/decisions/graph/region", http.StatusOK,
				nil, nil, "", graphRegionResponse,
			))

			resp, _, err := client.PulsarDecisions.GetDecisionsGraphRegion(nil)
			require.Nil(t, err)
			require.NotNil(t, resp)
			require.Equal(t, 1, len(resp.Data))
			require.Equal(t, "us-east", resp.Data[0].Region)
		})
	})

	// Tests for api.Client.PulsarDecisions.GetDecisionsGraphTime()
	t.Run("GetDecisionsGraphTime", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			unit := "requests"
			graphTimeResponse := &pulsar.DecisionsGraphTimeResponse{
				Data: []*pulsar.DecisionsGraphTimeData{
					{
						Timestamp: 1234567890,
						Counts: []*pulsar.JobIDCounts{
							{JobID: myJobID, Count: 100},
						},
					},
					{
						Timestamp: 1234567900,
						Counts: []*pulsar.JobIDCounts{
							{JobID: myJobID, Count: 150},
						},
					},
				},
				Unit: &unit,
			}

			require.Nil(t, mock.AddTestCase(
				http.MethodGet, "/pulsar/query/decisions/graph/time", http.StatusOK,
				nil, nil, "", graphTimeResponse,
			))

			resp, _, err := client.PulsarDecisions.GetDecisionsGraphTime(nil)
			require.Nil(t, err)
			require.NotNil(t, resp)
			require.Equal(t, 2, len(resp.Data))
		})
	})

	// Tests for api.Client.PulsarDecisions.GetDecisionsArea()
	t.Run("GetDecisionsArea", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			areaResponse := &pulsar.DecisionsAreaResponse{
				Areas: []*pulsar.AreaData{
					{
						AreaName: "US",
						Parent:   stringPtr("GLOBAL"),
						Count:    1000,
						JobCounts: []*pulsar.JobIDCounts{
							{JobID: myJobID, Count: 1000},
						},
					},
				},
			}

			require.Nil(t, mock.AddTestCase(
				http.MethodGet, "/pulsar/query/decisions/area", http.StatusOK,
				nil, nil, "", areaResponse,
			))

			resp, _, err := client.PulsarDecisions.GetDecisionsArea(nil)
			require.Nil(t, err)
			require.NotNil(t, resp)
			require.Equal(t, 1, len(resp.Areas))
			require.Equal(t, "US", resp.Areas[0].AreaName)
		})
	})

	// Tests for api.Client.PulsarDecisions.GetDecisionsASN()
	t.Run("GetDecisionsASN", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			unit := "requests"
			asnResponse := &pulsar.DecisionsASNResponse{
				Data: []*pulsar.ASNData{
					{
						ASN:                 15169,
						Count:               500,
						TrafficDistribution: 0.5,
						PreviousDay:         450.0,
						PreviousWeek:        400.0,
					},
				},
				Unit: &unit,
			}

			require.Nil(t, mock.AddTestCase(
				http.MethodGet, "/pulsar/query/decisions/asn", http.StatusOK,
				nil, nil, "", asnResponse,
			))

			resp, _, err := client.PulsarDecisions.GetDecisionsASN(nil)
			require.Nil(t, err)
			require.NotNil(t, resp)
			require.Equal(t, 1, len(resp.Data))
			require.Equal(t, int64(15169), resp.Data[0].ASN)
		})
	})

	// Tests for api.Client.PulsarDecisions.GetDecisionsResultsTime()
	t.Run("GetDecisionsResultsTime", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			resultsTimeResponse := &pulsar.DecisionsResultsTimeResponse{
				Data: []*pulsar.ResultsTimeData{
					{
						Timestamp: 1234567890,
						Results: []*pulsar.ResultCount{
							{Result: "answer1.example.com", Count: 100},
							{Result: "answer2.example.com", Count: 50},
						},
					},
				},
			}

			require.Nil(t, mock.AddTestCase(
				http.MethodGet, "/pulsar/query/decisions/results/time", http.StatusOK,
				nil, nil, "", resultsTimeResponse,
			))

			resp, _, err := client.PulsarDecisions.GetDecisionsResultsTime(nil)
			require.Nil(t, err)
			require.NotNil(t, resp)
			require.Equal(t, 1, len(resp.Data))
			require.Equal(t, 2, len(resp.Data[0].Results))
		})
	})

	// Tests for api.Client.PulsarDecisions.GetDecisionsResultsArea()
	t.Run("GetDecisionsResultsArea", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			resultsAreaResponse := &pulsar.DecisionsResultsAreaResponse{
				Area: []*pulsar.DecisionsResultsArea{
					{
						Area:          "US",
						Parent:        "GLOBAL",
						DecisionCount: 1000,
						Results: []*pulsar.ResultCount{
							{Result: "answer1.example.com", Count: 600},
							{Result: "answer2.example.com", Count: 400},
						},
					},
				},
			}

			require.Nil(t, mock.AddTestCase(
				http.MethodGet, "/pulsar/query/decisions/results/area", http.StatusOK,
				nil, nil, "", resultsAreaResponse,
			))

			resp, _, err := client.PulsarDecisions.GetDecisionsResultsArea(nil)
			require.Nil(t, err)
			require.NotNil(t, resp)
			require.Equal(t, 1, len(resp.Area))
			require.Equal(t, "US", resp.Area[0].Area)
		})
	})

	// Tests for api.Client.PulsarDecisions.GetFiltersTime()
	t.Run("GetFiltersTime", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			filtersTimeResponse := &pulsar.FiltersTimeResponse{
				Filters: []*pulsar.FilterTimeData{
					{
						Timestamp: 1234567890,
						Filters:   map[string]int64{"filter1": 100, "filter2": 50},
					},
				},
			}

			require.Nil(t, mock.AddTestCase(
				http.MethodGet, "/pulsar/query/decisions/filters/time", http.StatusOK,
				nil, nil, "", filtersTimeResponse,
			))

			resp, _, err := client.PulsarDecisions.GetFiltersTime(nil)
			require.Nil(t, err)
			require.NotNil(t, resp)
			require.Equal(t, 1, len(resp.Filters))
		})
	})

	// Tests for api.Client.PulsarDecisions.GetDecisionCustomer()
	t.Run("GetDecisionCustomer", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			customerResponse := &pulsar.DecisionCustomerResponse{
				Data: []*pulsar.CustomerDecisionData{
					{
						Timestamp: 1234567890,
						Total:     1000,
						JobCounts: map[string]int64{myJobID: 1000},
					},
				},
			}

			require.Nil(t, mock.AddTestCase(
				http.MethodGet, "/pulsar/query/decision/customer/12345", http.StatusOK,
				nil, nil, "", customerResponse,
			))

			resp, _, err := client.PulsarDecisions.GetDecisionCustomer("12345", nil)
			require.Nil(t, err)
			require.NotNil(t, resp)
			require.Equal(t, 1, len(resp.Data))
		})
	})

	// Tests for api.Client.PulsarDecisions.GetDecisionCustomerUndetermined()
	t.Run("GetDecisionCustomerUndetermined", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			customerResponse := &pulsar.DecisionCustomerResponse{
				Data: []*pulsar.CustomerDecisionData{
					{
						Timestamp: 1234567890,
						Total:     50,
						JobCounts: map[string]int64{myJobID: 50},
					},
				},
			}

			require.Nil(t, mock.AddTestCase(
				http.MethodGet, "/pulsar/query/decision/customer/12345/undetermined", http.StatusOK,
				nil, nil, "", customerResponse,
			))

			resp, _, err := client.PulsarDecisions.GetDecisionCustomerUndetermined("12345", nil)
			require.Nil(t, err)
			require.NotNil(t, resp)
			require.Equal(t, 1, len(resp.Data))
		})
	})

	// Tests for api.Client.PulsarDecisions.GetDecisionRecord()
	t.Run("GetDecisionRecord", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			recordResponse := &pulsar.DecisionCustomerResponse{
				Data: []*pulsar.CustomerDecisionData{
					{
						Timestamp: 1234567890,
						Total:     500,
						JobCounts: map[string]int64{myJobID: 500},
					},
				},
			}

			require.Nil(t, mock.AddTestCase(
				http.MethodGet, "/pulsar/query/decision/customer/12345/record/example.com/A", http.StatusOK,
				nil, nil, "", recordResponse,
			))

			resp, _, err := client.PulsarDecisions.GetDecisionRecord("12345", "example.com", "A", nil)
			require.Nil(t, err)
			require.NotNil(t, resp)
			require.Equal(t, 1, len(resp.Data))
		})
	})

	// Tests for api.Client.PulsarDecisions.GetDecisionRecordUndetermined()
	t.Run("GetDecisionRecordUndetermined", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			recordResponse := &pulsar.DecisionCustomerResponse{
				Data: []*pulsar.CustomerDecisionData{
					{
						Timestamp: 1234567890,
						Total:     25,
						JobCounts: map[string]int64{myJobID: 25},
					},
				},
			}

			require.Nil(t, mock.AddTestCase(
				http.MethodGet, "/pulsar/query/decision/customer/12345/record/example.com/A/undetermined", http.StatusOK,
				nil, nil, "", recordResponse,
			))

			resp, _, err := client.PulsarDecisions.GetDecisionRecordUndetermined("12345", "example.com", "A", nil)
			require.Nil(t, err)
			require.NotNil(t, resp)
			require.Equal(t, 1, len(resp.Data))
		})
	})

	// Tests for api.Client.PulsarDecisions.GetDecisionTotal()
	t.Run("GetDecisionTotal", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			totalResponse := &pulsar.DecisionTotalResponse{
				Total: 10000,
			}

			require.Nil(t, mock.AddTestCase(
				http.MethodGet, "/pulsar/query/decision/customer/12345/total", http.StatusOK,
				nil, nil, "", totalResponse,
			))

			resp, _, err := client.PulsarDecisions.GetDecisionTotal("12345", nil)
			require.Nil(t, err)
			require.NotNil(t, resp)
			require.Equal(t, int64(10000), resp.Total)
		})
	})

	// Tests for api.Client.PulsarDecisions.GetPulsarDecisionsRecords()
	t.Run("GetPulsarDecisionsRecords", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			recordsResponse := &pulsar.DecisionsRecordsResponse{
				Total: 500,
				Records: map[string]*pulsar.Record{
					"example.com/A": {
						Count:             500,
						PercentageOfTotal: 100.0,
					},
				},
			}

			require.Nil(t, mock.AddTestCase(
				http.MethodGet, "/pulsar/query/decisions/records", http.StatusOK,
				nil, nil, "", recordsResponse,
			))

			resp, _, err := client.PulsarDecisions.GetPulsarDecisionsRecords(nil)
			require.Nil(t, err)
			require.NotNil(t, resp)
			require.Equal(t, int64(500), resp.Total)
			require.Equal(t, 1, len(resp.Records))
		})
	})

	// Tests for api.Client.PulsarDecisions.GetPulsarDecisionsResultsRecord()
	t.Run("GetPulsarDecisionsResultsRecord", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			resultsRecordResponse := &pulsar.DecisionsResultsRecordResponse{
				Record: map[string]*pulsar.Results{
					"example.com/A": {
						DecisionCount: 500,
						Results: map[string]int64{
							"192.0.2.1": 300,
							"192.0.2.2": 200,
						},
					},
				},
			}

			require.Nil(t, mock.AddTestCase(
				http.MethodGet, "/pulsar/query/decisions/results/record", http.StatusOK,
				nil, nil, "", resultsRecordResponse,
			))

			resp, _, err := client.PulsarDecisions.GetPulsarDecisionsResultsRecord(nil)
			require.Nil(t, err)
			require.NotNil(t, resp)
			require.Equal(t, 1, len(resp.Record))
			record := resp.Record["example.com/A"]
			require.NotNil(t, record)
			require.Equal(t, int64(500), record.DecisionCount)
			require.Equal(t, 2, len(record.Results))
		})
	})
}

// Helper functions for pointer types
func stringPtr(s string) *string {
	return &s
}

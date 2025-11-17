package pulsar

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecisionsQueryParams(t *testing.T) {
	params := &DecisionsQueryParams{
		Start:  1234567890,
		End:    1234567900,
		Period: "1h",
		Area:   "US",
		ASN:    "15169",
		Job:    "job123",
		Jobs:   []string{"job1", "job2"},
		Record: "example.com",
		Result: "192.0.2.1",
		Agg:    "sum",
		Geo:    "country",
		ZoneID: "zone123",
	}

	assert.Equal(t, int64(1234567890), params.Start)
	assert.Equal(t, int64(1234567900), params.End)
	assert.Equal(t, "1h", params.Period)
	assert.Equal(t, "US", params.Area)
	assert.Equal(t, "15169", params.ASN)
	assert.Equal(t, "job123", params.Job)
	assert.Equal(t, []string{"job1", "job2"}, params.Jobs)
	assert.Equal(t, "example.com", params.Record)
	assert.Equal(t, "192.0.2.1", params.Result)
	assert.Equal(t, "sum", params.Agg)
	assert.Equal(t, "country", params.Geo)
	assert.Equal(t, "zone123", params.ZoneID)
}

func TestDecisionsResponse(t *testing.T) {
	resp := &DecisionsResponse{
		Graphs: []*DecisionsGraph{
			{
				Count: 1,
				GraphData: []*DecisionsGraphData{
					{Timestamp: 1234567890, Count: 1},
				},
				Avg: 1,
				Tags: &Tags{
					JobID: "job123",
				},
			},
		},
		Total: 1,
	}

	assert.Equal(t, 1, len(resp.Graphs))
	assert.Equal(t, int64(1), resp.Graphs[0].Count)
	assert.Equal(t, int64(1), resp.Total)
	assert.Equal(t, "job123", resp.Graphs[0].Tags.JobID)
}

func TestDecisionsGraphRegionResponse(t *testing.T) {
	unit := "requests"
	resp := &DecisionsGraphRegionResponse{
		Data: []*DecisionsGraphRegionData{
			{
				Region: "us-east",
				Counts: []*JobIDCounts{
					{JobID: "job123", Count: 1000},
				},
			},
		},
		Unit: &unit,
	}

	assert.Equal(t, 1, len(resp.Data))
	assert.Equal(t, "us-east", resp.Data[0].Region)
	assert.Equal(t, 1, len(resp.Data[0].Counts))
	assert.Equal(t, int64(1000), resp.Data[0].Counts[0].Count)
	assert.Equal(t, "requests", *resp.Unit)
}

func TestDecisionsGraphTimeResponse(t *testing.T) {
	unit := "requests"
	resp := &DecisionsGraphTimeResponse{
		Data: []*DecisionsGraphTimeData{
			{
				Timestamp: 1234567890,
				Counts: []*JobIDCounts{
					{JobID: "job123", Count: 100},
				},
			},
		},
		Unit: &unit,
	}

	assert.Equal(t, 1, len(resp.Data))
	assert.Equal(t, int64(1234567890), resp.Data[0].Timestamp)
	assert.Equal(t, 1, len(resp.Data[0].Counts))
	assert.Equal(t, int64(100), resp.Data[0].Counts[0].Count)
	assert.Equal(t, "requests", *resp.Unit)
}

func TestDecisionsAreaResponse(t *testing.T) {
	parent := "GLOBAL"
	resp := &DecisionsAreaResponse{
		Areas: []*AreaData{
			{
				AreaName: "US",
				Parent:   &parent,
				Count:    1000,
				JobCounts: []*JobIDCounts{
					{JobID: "job123", Count: 1000},
				},
			},
		},
	}

	assert.Equal(t, 1, len(resp.Areas))
	assert.Equal(t, "US", resp.Areas[0].AreaName)
	assert.Equal(t, "GLOBAL", *resp.Areas[0].Parent)
	assert.Equal(t, int64(1000), resp.Areas[0].Count)
	assert.Equal(t, 1, len(resp.Areas[0].JobCounts))
	assert.Equal(t, "job123", resp.Areas[0].JobCounts[0].JobID)
}

func TestDecisionsASNResponse(t *testing.T) {
	unit := "requests"
	resp := &DecisionsASNResponse{
		Data: []*ASNData{
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

	assert.Equal(t, 1, len(resp.Data))
	assert.Equal(t, int64(15169), resp.Data[0].ASN)
	assert.Equal(t, int64(500), resp.Data[0].Count)
	assert.Equal(t, 0.5, resp.Data[0].TrafficDistribution)
	assert.Equal(t, 450.0, resp.Data[0].PreviousDay)
	assert.Equal(t, 400.0, resp.Data[0].PreviousWeek)
	assert.Equal(t, "requests", *resp.Unit)
}

func TestDecisionsResultsTimeResponse(t *testing.T) {
	resp := &DecisionsResultsTimeResponse{
		Data: []*ResultsTimeData{
			{
				Timestamp: 1234567890,
				Results: []*ResultCount{
					{Result: "answer1.example.com", Count: 100},
					{Result: "answer2.example.com", Count: 50},
				},
			},
		},
	}

	assert.Equal(t, 1, len(resp.Data))
	assert.Equal(t, int64(1234567890), resp.Data[0].Timestamp)
	assert.Equal(t, 2, len(resp.Data[0].Results))
	assert.Equal(t, "answer1.example.com", resp.Data[0].Results[0].Result)
	assert.Equal(t, int64(100), resp.Data[0].Results[0].Count)
}

func TestDecisionsResultsAreaResponse(t *testing.T) {
	resp := &DecisionsResultsAreaResponse{
		Area: []*DecisionsResultsArea{
			{
				Area:          "US",
				Parent:        "GLOBAL",
				DecisionCount: 1000,
				Results: []*ResultCount{
					{Result: "answer1.example.com", Count: 600},
				},
			},
		},
	}

	assert.Equal(t, 1, len(resp.Area))
	assert.Equal(t, "US", resp.Area[0].Area)
	assert.Equal(t, "GLOBAL", resp.Area[0].Parent)
	assert.Equal(t, int64(1000), resp.Area[0].DecisionCount)
	assert.Equal(t, 1, len(resp.Area[0].Results))
}

func TestFiltersTimeResponse(t *testing.T) {
	resp := &FiltersTimeResponse{
		Filters: []*FilterTimeData{
			{
				Timestamp: 1234567890,
				Filters:   map[string]int64{"filter1": 100, "filter2": 50},
			},
		},
	}

	assert.Equal(t, 1, len(resp.Filters))
	assert.Equal(t, int64(1234567890), resp.Filters[0].Timestamp)
	assert.Equal(t, int64(100), resp.Filters[0].Filters["filter1"])
	assert.Equal(t, int64(50), resp.Filters[0].Filters["filter2"])
}

func TestDecisionCustomerResponse(t *testing.T) {
	resp := &DecisionCustomerResponse{
		Data: []*CustomerDecisionData{
			{
				Timestamp: 1234567890,
				Total:     1000,
				JobCounts: map[string]int64{"job123": 1000},
			},
		},
	}

	assert.Equal(t, 1, len(resp.Data))
	assert.Equal(t, int64(1234567890), resp.Data[0].Timestamp)
	assert.Equal(t, int64(1000), resp.Data[0].Total)
	assert.Equal(t, int64(1000), resp.Data[0].JobCounts["job123"])
}

func TestDecisionTotalResponse(t *testing.T) {
	resp := &DecisionTotalResponse{
		Total: 10000,
	}

	assert.Equal(t, int64(10000), resp.Total)
}

func TestDecisionsRecordsResponse(t *testing.T) {
	resp := &DecisionsRecordsResponse{
		Total: 500,
		Records: map[string]*Record{
			"example.com/A": {
				Count:             500,
				PercentageOfTotal: 100.0,
			},
		},
	}

	assert.Equal(t, int64(500), resp.Total)
	assert.Equal(t, 1, len(resp.Records))
	record := resp.Records["example.com/A"]
	assert.NotNil(t, record)
	assert.Equal(t, int64(500), record.Count)
	assert.Equal(t, 100.0, record.PercentageOfTotal)
}

func TestDecisionsResultsRecordResponse(t *testing.T) {
	resp := &DecisionsResultsRecordResponse{
		Record: map[string]*Results{
			"example.com/A": {
				DecisionCount: 500,
				Results: map[string]int64{
					"192.0.2.1": 300,
					"192.0.2.2": 200,
				},
			},
		},
	}

	assert.Equal(t, 1, len(resp.Record))
	record := resp.Record["example.com/A"]
	assert.NotNil(t, record)
	assert.Equal(t, int64(500), record.DecisionCount)
	assert.Equal(t, 2, len(record.Results))
	assert.Equal(t, int64(300), record.Results["192.0.2.1"])
	assert.Equal(t, int64(200), record.Results["192.0.2.2"])
}

func TestDecisionsGraphJSONMarshaling(t *testing.T) {
	data := &DecisionsGraph{
		Count: 100,
		GraphData: []*DecisionsGraphData{
			{Timestamp: 1234567890, Count: 50},
			{Timestamp: 1234567900, Count: 50},
		},
		Avg: 50,
		Tags: &Tags{
			JobID: "job123",
		},
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(data)
	require.NoError(t, err)

	// Unmarshal back
	var unmarshaled DecisionsGraph
	err = json.Unmarshal(jsonData, &unmarshaled)
	require.NoError(t, err)

	// Verify
	assert.Equal(t, data.Count, unmarshaled.Count)
	assert.Equal(t, data.Avg, unmarshaled.Avg)
	assert.Equal(t, len(data.GraphData), len(unmarshaled.GraphData))
	assert.Equal(t, data.Tags.JobID, unmarshaled.Tags.JobID)
}

func TestDecisionsGraphTimeDataJSONMarshaling(t *testing.T) {
	data := &DecisionsGraphTimeData{
		Timestamp: 1234567890,
		Counts: []*JobIDCounts{
			{JobID: "job1", Count: 100},
			{JobID: "job2", Count: 50},
		},
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(data)
	require.NoError(t, err)

	// Unmarshal back
	var unmarshaled DecisionsGraphTimeData
	err = json.Unmarshal(jsonData, &unmarshaled)
	require.NoError(t, err)

	// Verify
	assert.Equal(t, data.Timestamp, unmarshaled.Timestamp)
	assert.Equal(t, len(data.Counts), len(unmarshaled.Counts))
	assert.Equal(t, data.Counts[0].JobID, unmarshaled.Counts[0].JobID)
	assert.Equal(t, data.Counts[0].Count, unmarshaled.Counts[0].Count)
}

func TestAreaDataWithNilParent(t *testing.T) {
	data := &AreaData{
		AreaName:  "US",
		Parent:    nil,
		Count:     1000,
		JobCounts: []*JobIDCounts{{JobID: "job123", Count: 1000}},
	}

	assert.Nil(t, data.Parent)
	assert.Equal(t, "US", data.AreaName)
	assert.Equal(t, int64(1000), data.Count)
}

func TestResultCountArray(t *testing.T) {
	results := []*ResultCount{
		{Result: "answer1.example.com", Count: 600},
		{Result: "answer2.example.com", Count: 400},
	}

	assert.Equal(t, 2, len(results))

	// Verify total count
	var total int64
	for _, r := range results {
		total += r.Count
	}
	assert.Equal(t, int64(1000), total)
}

func TestEmptyDecisionsResponse(t *testing.T) {
	resp := &DecisionsResponse{
		Graphs: []*DecisionsGraph{},
		Total:  0,
	}

	assert.NotNil(t, resp.Graphs)
	assert.Equal(t, 0, len(resp.Graphs))
	assert.Equal(t, int64(0), resp.Total)
}

func TestDecisionsQueryParamsEmpty(t *testing.T) {
	params := &DecisionsQueryParams{}

	assert.Equal(t, int64(0), params.Start)
	assert.Equal(t, int64(0), params.End)
	assert.Equal(t, "", params.Period)
	assert.Equal(t, "", params.Area)
	assert.Nil(t, params.Jobs)
}

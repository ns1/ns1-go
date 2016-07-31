package nsone

import "fmt"

// DataFeed wraps an NS1 /data/feeds resource
type DataFeed struct {
	SourceId string                 `json:"-"`
	Id       string                 `json:"id,omitempty"`
	Name     string                 `json:"name"`
	Config   map[string]string      `json:"config,omitempty"`
	Data     map[string]interface{} `json:"data,omitempty"`
}

// NewDataFeed takes a sourceID and creates a new *DataFeed
func NewDataFeed(sourceID string) *DataFeed {
	return &DataFeed{SourceId: sourceID}
}

// CreateDataFeed takes a *DataFeed and connects a new data feed to an existing data source
func (c APIClient) CreateDataFeed(df *DataFeed) error {
	return c.doHTTPBoth("PUT", fmt.Sprintf("https://api.nsone.net/v1/data/feeds/%s", df.SourceId), df)
}

// GetDataFeed takes a data source ID and a data feed ID and returns the details of a single data feed
func (c APIClient) GetDataFeed(dsID string, dfID string) (*DataFeed, error) {
	df := NewDataFeed(dsID)
	status, err := c.doHTTPUnmarshal("GET", fmt.Sprintf("https://api.nsone.net/v1/data/feeds/%s/%s", dsID, dfID), nil, df)
	if status == 404 {
		df.SourceId = ""
		df.Id = ""
		df.Name = ""
		return df, nil
	}
	return df, err
}

// DeleteDataFeed takes a data source ID and a data feed ID and disconnects the feed from the data source and all attached destination metadata tables
func (c APIClient) DeleteDataFeed(dsID string, dfID string) error {
	return c.doHTTPDelete(fmt.Sprintf("https://api.nsone.net/v1/data/feeds/%s/%s", dsID, dfID))
}

// UpdateDataFeed takes a *DataFeed and modifies and existing data feed
func (c APIClient) UpdateDataFeed(df *DataFeed) error {
	return c.doHTTPBoth("POST", fmt.Sprintf("https://api.nsone.net/v1/data/feeds/%s/%s", df.SourceId, df.Id), df)
}

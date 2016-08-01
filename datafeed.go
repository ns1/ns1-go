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
	path := fmt.Sprintf("data/feeds/%s", df.SourceId)

	req, err := c.NewRequest("PUT", path, &df)
	if err != nil {
		return err
	}

	// Update datafeed fields with data from api
	_, err = c.Do(req, &df)
	if err != nil {
		return err
	}

	return nil
}

// GetDataFeed takes a data source ID and a data feed ID and returns the details of a single data feed
func (c APIClient) GetDataFeed(dsID string, dfID string) (*DataFeed, error) {
	df := NewDataFeed(dsID)

	path := fmt.Sprintf("data/feeds/%s/%s", dsID, dfID)
	req, err := c.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	_, err = c.Do(req, df)
	if err != nil {
		return nil, err
	}

	return df, nil
}

// DeleteDataFeed takes a data source ID and a data feed ID and disconnects the feed from the data source and all attached destination metadata tables
func (c APIClient) DeleteDataFeed(dsID string, dfID string) error {
	path := fmt.Sprintf("data/feeds/%s/%s", dsID, dfID)
	req, err := c.NewRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	_, err = c.Do(req, nil)
	if err != nil {
		return err
	}

	return nil
}

// UpdateDataFeed takes a *DataFeed and modifies and existing data feed
func (c APIClient) UpdateDataFeed(df *DataFeed) error {
	path := fmt.Sprintf("data/feeds/%s/%s", df.SourceId, df.Id)

	req, err := c.NewRequest("POST", path, &df)
	if err != nil {
		return err
	}

	// Update datafeed fields with data from api(ensure consistent)
	_, err = c.Do(req, &df)
	if err != nil {
		return err
	}

	return nil
}

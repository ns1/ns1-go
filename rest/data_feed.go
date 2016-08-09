package rest

import (
	"fmt"

	"github.com/ns1/ns1-go/model/data"
)

// DataFeedsService handles 'data/feeds' endpoint.
type DataFeedsService service

// List returns all data feeds connected to a given data source.
//
// NS1 API docs: https://ns1.com/api/#feeds-get
func (s *DataFeedsService) List(dsID string) ([]*data.Feed, *Response, error) {
	path := fmt.Sprintf("data/feeds/%s", dsID)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	dfl := []*data.Feed{}
	resp, err = s.client.Do(req, &dfl)
	if err != nil {
		return nil, resp, err
	}

	return dfl, resp, nil
}

// Get takes a data source ID and a data feed ID and returns the details of a single data feed
//
// NS1 API docs: https://ns1.com/api/#feeds-feed-get
func (s *DataFeedsService) Get(dsID string, dfID string) (*data.Feed, *Response, error) {
	path := fmt.Sprintf("data/feeds/%s/%s", dsID, dfID)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var df data.Feed
	resp, err = s.client.Do(req, &df)
	if err != nil {
		return nil, err
	}

	return &df, resp, nil
}

// Create takes a *DataFeed and connects a new data feed to an existing data source.
//
// NS1 API docs: https://ns1.com/api/#feeds-put
func (s *DataFeedsService) Create(df *data.Feed) (*Response, error) {
	path := fmt.Sprintf("data/feeds/%s", df.SourceID)

	req, err := s.client.NewRequest("PUT", path, &df)
	if err != nil {
		return err
	}

	// Update datafeeds' fields with data from api(ensure consistent)
	resp, err = s.client.Do(req, &df)
	if err != nil {
		return err
	}

	return resp, nil
}

// Update takes a *Feed and modifies and existing data feed.
//
// NS1 API docs: https://ns1.com/api/#feeds-post
func (s *DataFeedsService) Update(df *data.Feed) (*Response, error) {
	path := fmt.Sprintf("data/feeds/%s/%s", df.SourceID, df.ID)

	req, err := s.client.NewRequest("POST", path, &df)
	if err != nil {
		return err
	}

	// Update df instance fields with data from api(ensure consistent)
	resp, err = s.client.Do(req, &df)
	if err != nil {
		return err
	}

	return resp, nil
}

// Delete takes a data source ID and a data feed ID and disconnects the feed from the data source and all attached destination metadata tables.
//
// NS1 API docs: https://ns1.com/api/#feeds-delete
func (s *DataFeedsService) Delete(dsID string, dfID string) (*Response, error) {
	path := fmt.Sprintf("data/feeds/%s/%s", dsID, dfID)

	req, err := s.client.NewRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	resp, err = s.client.Do(req, nil)
	if err != nil {
		return err
	}

	return resp, nil
}

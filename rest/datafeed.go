package rest

import (
	"fmt"

	ns1 "github.com/ns1/ns1-go"
)

const (
	dataFeedPath = "data/feeds"
)

// DataFeedsService handles 'data/feeds' endpoint.
type DataFeedsService service

// List returns all data feeds connected to a given data source.
//
// NS1 API docs: https://ns1.com/api/#feeds-get
func (s *DataFeedsService) List(dsID string) ([]*ns1.DataFeed, error) {
	path := fmt.Sprintf("%s/%s", dataFeedPath, dsID)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	dfl := []*ns1.DataFeed{}
	_, err = s.client.Do(req, &dfl)
	if err != nil {
		return nil, err
	}

	return dfl, nil
}

// Get takes a data source ID and a data feed ID and returns the details of a single data feed
//
// NS1 API docs: https://ns1.com/api/#feeds-feed-get
func (s *DataFeedsService) Get(dsID string, dfID string) (*ns1.DataFeed, error) {
	path := fmt.Sprintf("%s/%s/%s", dataFeedPath, dsID, dfID)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var df ns1.DataFeed
	_, err = s.client.Do(req, &df)
	if err != nil {
		return nil, err
	}

	return &df, nil
}

// Create takes a *DataFeed and connects a new data feed to an existing data source.
//
// NS1 API docs: https://ns1.com/api/#feeds-put
func (s *DataFeedsService) Create(df *ns1.DataFeed) error {
	path := fmt.Sprintf("%s/%s", dataFeedPath, df.SourceID)

	req, err := s.client.NewRequest("PUT", path, &df)
	if err != nil {
		return err
	}

	// Update datafeeds' fields with data from api(ensure consistent)
	_, err = s.client.Do(req, &df)
	if err != nil {
		return err
	}

	return nil
}

// Update takes a *DataFeed and modifies and existing data feed.
//
// NS1 API docs: https://ns1.com/api/#feeds-post
func (s *DataFeedsService) Update(df *ns1.DataFeed) error {
	path := fmt.Sprintf("%s/%s/%s", dataFeedPath, df.SourceID, df.ID)

	req, err := s.client.NewRequest("POST", path, &df)
	if err != nil {
		return err
	}

	// Update df instance fields with data from api(ensure consistent)
	_, err = s.client.Do(req, &df)
	if err != nil {
		return err
	}

	return nil
}

// Delete takes a data source ID and a data feed ID and disconnects the feed from the data source and all attached destination metadata tables.
//
// NS1 API docs: https://ns1.com/api/#feeds-delete
func (s *DataFeedsService) Delete(dsID string, dfID string) error {
	path := fmt.Sprintf("%s/%s/%s", dataFeedPath, dsID, dfID)

	req, err := s.client.NewRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(req, nil)
	if err != nil {
		return err
	}

	return nil
}

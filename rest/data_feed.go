package rest

import (
	"context"
	"fmt"
	"net/http"

	"gopkg.in/ns1/ns1-go.v2/rest/model/data"
)

// DataFeedsService handles 'data/feeds' endpoint.
type DataFeedsService service

// List returns all data feeds connected to a given data source.
//
// NS1 API docs: https://ns1.com/api/#feeds-get
func (s *DataFeedsService) List(sourceID string) ([]*data.Feed, *http.Response, error) {
	return s.ListWithContext(context.Background(), sourceID)
}

// ListWithContext is the same as List, but takes a context.
func (s *DataFeedsService) ListWithContext(ctx context.Context, sourceID string) ([]*data.Feed, *http.Response, error) {
	path := fmt.Sprintf("data/feeds/%s", sourceID)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}
	req = req.WithContext(ctx)

	dfl := []*data.Feed{}
	resp, err := s.client.Do(req, &dfl)
	if err != nil {
		return nil, resp, err
	}

	return dfl, resp, nil
}

// Get takes a data source ID and a data feed ID and returns the details of a single data feed
//
// NS1 API docs: https://ns1.com/api/#feeds-feed-get
func (s *DataFeedsService) Get(sourceID string, feedID string) (*data.Feed, *http.Response, error) {
	return s.GetWithContext(context.Background(), sourceID, feedID)
}

// GetWithContext is the same as Get, but takes a context.
func (s *DataFeedsService) GetWithContext(ctx context.Context, sourceID string, feedID string) (*data.Feed, *http.Response, error) {
	path := fmt.Sprintf("data/feeds/%s/%s", sourceID, feedID)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}
	req = req.WithContext(ctx)

	var df data.Feed
	resp, err := s.client.Do(req, &df)
	if err != nil {
		return nil, resp, err
	}

	return &df, resp, nil
}

// Create takes a *DataFeed and connects a new data feed to an existing data source.
//
// NS1 API docs: https://ns1.com/api/#feeds-put
func (s *DataFeedsService) Create(sourceID string, df *data.Feed) (*http.Response, error) {
	return s.CreateWithContext(context.Background(), sourceID, df)
}

// CreateWithContext is the same as Create, but takes a context.
func (s *DataFeedsService) CreateWithContext(ctx context.Context, sourceID string, df *data.Feed) (*http.Response, error) {
	path := fmt.Sprintf("data/feeds/%s", sourceID)

	req, err := s.client.NewRequest("PUT", path, &df)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	// Update datafeeds' fields with data from api(ensure consistent)
	resp, err := s.client.Do(req, &df)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Update takes a *Feed and modifies and existing data feed.
// Note:
//   - The 'data' portion of a feed does not actually
//     get updated during a POST. In order to update a feeds'
//     'data' attribute, one must use the Publish method.
//   - Both the 'destinations' and 'networks' attributes are
//     not updated during a POST.
//
// NS1 API docs: https://ns1.com/api/#feeds-post
func (s *DataFeedsService) Update(sourceID string, df *data.Feed) (*http.Response, error) {
	return s.UpdateWithContext(context.Background(), sourceID, df)
}

// UpdateWithContext is the same as Update, but takes a context.
func (s *DataFeedsService) UpdateWithContext(ctx context.Context, sourceID string, df *data.Feed) (*http.Response, error) {
	path := fmt.Sprintf("data/feeds/%s/%s", sourceID, df.ID)

	req, err := s.client.NewRequest("POST", path, &df)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	// Update df instance fields with data from api(ensure consistent)
	resp, err := s.client.Do(req, &df)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Delete takes a data source ID and a data feed ID and disconnects the feed from the data source and all attached destination metadata tables.
//
// NS1 API docs: https://ns1.com/api/#feeds-delete
func (s *DataFeedsService) Delete(sourceID string, feedID string) (*http.Response, error) {
	return s.DeleteWithContext(context.Background(), sourceID, feedID)
}

// DeleteWithContext is the same as Delete, but takes a context.
func (s *DataFeedsService) DeleteWithContext(ctx context.Context, sourceID string, feedID string) (*http.Response, error) {
	path := fmt.Sprintf("data/feeds/%s/%s", sourceID, feedID)

	req, err := s.client.NewRequest("DELETE", path, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

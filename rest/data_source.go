package rest

import (
	"fmt"

	"github.com/ns1/ns1-go/model/data"
)

// DataSourcesService handles 'data/sources' endpoint.
type DataSourcesService service

// List returns all connected data sources.
//
// NS1 API docs: https://ns1.com/api/#sources-get
func (s *DataSourcesService) List() ([]*data.Source, *Response, error) {
	req, err := s.client.NewRequest("GET", "data/sources", nil)
	if err != nil {
		return nil, nil, err
	}

	dsl := []*data.Source{}
	resp, err = s.client.Do(req, &dsl)
	if err != nil {
		return nil, resp, err
	}

	return dsl, resp, nil
}

// Get takes an ID returns the details for a single data source.
//
// NS1 API docs: https://ns1.com/api/#sources-source-get
func (s *DataSourcesService) Get(id string) (*data.Source, *Response, error) {
	path := fmt.Sprintf("data/sources/%s", id)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var ds data.Source
	resp, err = s.client.Do(req, &ds)
	if err != nil {
		return nil, resp, err
	}

	return &ds, resp, nil
}

// Create takes a *DataSource and creates a new data source.
//
// NS1 API docs: https://ns1.com/api/#sources-put
func (s *DataSourcesService) Create(ds *data.Source) (*Response, error) {
	req, err := s.client.NewRequest("PUT", "data/sources", &ds)
	if err != nil {
		return nil, err
	}

	// Update data sources' fields with data from api(ensure consistent)
	resp, err = s.client.Do(req, &ds)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Update takes a *DataSource modifies basic details of a data source.
// NOTE: This does not 'publish' data. See PublishToFeed
//
// NS1 API docs: https://ns1.com/api/#sources-post
func (s *DataSourcesService) Update(ds *data.Source) (*Response, error) {
	path := fmt.Sprintf("data/sources/%s", ds.ID)

	req, err := s.client.NewRequest("POST", path, &ds)
	if err != nil {
		return nil, err
	}

	// Update data sources' instance fields with data from api(ensure consistent)
	resp, err = s.client.Do(req, &ds)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Delete takes an ID and removes an existing data source and all connected feeds from the source.
//
// NS1 API docs: https://ns1.com/api/#sources-delete
func (s *DataSourcesService) Delete(id string) (*Response, error) {
	path := fmt.Sprintf("data/sources/%s", id)

	req, err := s.client.NewRequest("DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err = s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// PublishToFeed takes a datasources' id and data to publish to the feed.
//
// NS1 API docs: https://ns1.com/api/#feed-post
func (s *DataSourcesService) PublishToFeed(dataSourceId string, data interface{}) (*Response, error) {
	path := fmt.Sprintf("feed/%s", dataSourceId)

	req, err := s.client.NewRequest("POST", path, &data)
	if err != nil {
		return nil, err
	}

	resp, err = s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

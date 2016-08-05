package rest

import (
	"fmt"

	ns1 "github.com/ns1/ns1-go"
)

const (
	dataSourcePath = "data/sources"
)

// DataSourcesService handles 'data/sources' endpoint.
type DataSourcesService service

// List returns all connected data sources.
//
// NS1 API docs: https://ns1.com/api/#sources-get
func (s *DataSourcesService) List() ([]*ns1.DataSource, error) {
	req, err := s.client.NewRequest("GET", dataSourcePath, nil)
	if err != nil {
		return nil, err
	}

	dsl := []*ns1.DataSource{}
	_, err = s.client.Do(req, &dsl)
	if err != nil {
		return nil, err
	}

	return dsl, nil
}

// Get takes an ID returns the details for a single data source.
//
// NS1 API docs: https://ns1.com/api/#sources-source-get
func (s *DataSourcesService) Get(id string) (*ns1.DataSource, error) {
	path := fmt.Sprintf("%s/%s", dataSourcePath, id)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var ds ns1.DataSource
	_, err = s.client.Do(req, &ds)
	if err != nil {
		return nil, err
	}

	return &ds, nil
}

// Create takes a *DataSource and creates a new data source.
//
// NS1 API docs: https://ns1.com/api/#sources-put
func (s *DataSourcesService) Create(ds *ns1.DataSource) error {
	path := fmt.Sprintf("%s", dataSourcePath)

	req, err := s.client.NewRequest("PUT", path, &ds)
	if err != nil {
		return err
	}

	// Update data sources' fields with data from api(ensure consistent)
	_, err = s.client.Do(req, &ds)
	if err != nil {
		return err
	}

	return nil
}

// Update takes a *DataSource modifies basic details of a data source.
// NOTE: This does not 'publish' data. See PublishToFeed
//
// NS1 API docs: https://ns1.com/api/#sources-post
func (s *DataSourcesService) Update(ds *ns1.DataSource) error {
	path := fmt.Sprintf("%s/%s", dataSourcePath, ds.ID)

	req, err := s.client.NewRequest("POST", path, &ds)
	if err != nil {
		return err
	}

	// Update data sources' instance fields with data from api(ensure consistent)
	_, err = s.client.Do(req, &ds)
	if err != nil {
		return err
	}

	return nil
}

// Delete takes an ID and removes an existing data source and all connected feeds from the source.
//
// NS1 API docs: https://ns1.com/api/#sources-delete
func (s *DataSourcesService) Delete(id string) error {
	path := fmt.Sprintf("%s/%s", dataSourcePath, id)

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

// PublishToFeed takes a datasources' id and data to publish to the feed.
//
// NS1 API docs: https://ns1.com/api/#feed-post
func (s *DataSourcesService) PublishToFeed(dataSourceId string, data interface{}) error {
	path := fmt.Sprintf("feed/%s", dataSourceId)

	req, err := s.client.NewRequest("POST", path, &data)
	if err != nil {
		return err
	}

	_, err = s.client.Do(req, nil)
	if err != nil {
		return err
	}

	return nil
}

package rest

import (
	"context"
	"fmt"
	"net/http"

	"gopkg.in/ns1/ns1-go.v2/rest/model/data"
)

// DataSourcesService handles 'data/sources' endpoint.
type DataSourcesService service

// List returns all connected data sources.
//
// NS1 API docs: https://ns1.com/api/#sources-get
func (s *DataSourcesService) List() ([]*data.Source, *http.Response, error) {
	return s.ListWithContext(context.Background())
}

// ListWithContext is the same as List, but takes a context.
func (s *DataSourcesService) ListWithContext(ctx context.Context) ([]*data.Source, *http.Response, error) {
	req, err := s.client.NewRequest("GET", "data/sources", nil)
	if err != nil {
		return nil, nil, err
	}
	req = req.WithContext(ctx)

	dsl := []*data.Source{}
	resp, err := s.client.Do(req, &dsl)
	if err != nil {
		return nil, resp, err
	}

	return dsl, resp, nil
}

// Get takes an ID returns the details for a single data source.
//
// NS1 API docs: https://ns1.com/api/#sources-source-get
func (s *DataSourcesService) Get(id string) (*data.Source, *http.Response, error) {
	return s.GetWithContext(context.Background(), id)
}

// GetWithContext is the same as Get, but takes a context.
func (s *DataSourcesService) GetWithContext(ctx context.Context, id string) (*data.Source, *http.Response, error) {
	path := fmt.Sprintf("data/sources/%s", id)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}
	req = req.WithContext(ctx)

	var ds data.Source
	resp, err := s.client.Do(req, &ds)
	if err != nil {
		return nil, resp, err
	}

	return &ds, resp, nil
}

// Create takes a *DataSource and creates a new data source.
//
// NS1 API docs: https://ns1.com/api/#sources-put
func (s *DataSourcesService) Create(ds *data.Source) (*http.Response, error) {
	return s.CreateWithContext(context.Background(), ds)
}

// CreateWithContext is the same as Create, but takes a context.
func (s *DataSourcesService) CreateWithContext(ctx context.Context, ds *data.Source) (*http.Response, error) {
	req, err := s.client.NewRequest("PUT", "data/sources", &ds)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	// Update data sources' fields with data from api(ensure consistent)
	resp, err := s.client.Do(req, &ds)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Update takes a *DataSource modifies basic details of a data source.
// NOTE: This does not 'publish' data. See the Publish method.
//
// NS1 API docs: https://ns1.com/api/#sources-post
func (s *DataSourcesService) Update(ds *data.Source) (*http.Response, error) {
	return s.UpdateWithContext(context.Background(), ds)
}

// UpdateWithContext is the same as Update, but takes a context.
func (s *DataSourcesService) UpdateWithContext(ctx context.Context, ds *data.Source) (*http.Response, error) {
	path := fmt.Sprintf("data/sources/%s", ds.ID)

	req, err := s.client.NewRequest("POST", path, &ds)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	// Update data sources' instance fields with data from api(ensure consistent)
	resp, err := s.client.Do(req, &ds)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Delete takes an ID and removes an existing data source and all connected feeds from the source.
//
// NS1 API docs: https://ns1.com/api/#sources-delete
func (s *DataSourcesService) Delete(id string) (*http.Response, error) {
	return s.DeleteWithContext(context.Background(), id)
}

// DeleteWithContext is the same as Delete, but takes a context.
func (s *DataSourcesService) DeleteWithContext(ctx context.Context, id string) (*http.Response, error) {
	path := fmt.Sprintf("data/sources/%s", id)

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

// Publish takes a datasources' id and data to publish.
//
// NS1 API docs: https://ns1.com/api/#feed-post
func (s *DataSourcesService) Publish(dsID string, data interface{}) (*http.Response, error) {
	return s.PublishWithContext(context.Background(), dsID, data)
}

// PublishWithContext is the same as Publish, but takes a context.
func (s *DataSourcesService) PublishWithContext(ctx context.Context, dsID string, data interface{}) (*http.Response, error) {
	path := fmt.Sprintf("feed/%s", dsID)

	req, err := s.client.NewRequest("POST", path, &data)
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

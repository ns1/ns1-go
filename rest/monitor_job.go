package rest

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"gopkg.in/ns1/ns1-go.v2/rest/model/monitor"
)

// JobsService handles 'monitoring/jobs' endpoint.
type JobsService service

// List returns all monitoring jobs for the account.
//
// NS1 API docs: https://ns1.com/api/#jobs-get
func (s *JobsService) List() ([]*monitor.Job, *http.Response, error) {
	return s.ListWithContext(context.Background())
}

// ListWithContext is the same as List, but takes a context.
func (s *JobsService) ListWithContext(ctx context.Context) ([]*monitor.Job, *http.Response, error) {
	req, err := s.client.NewRequest("GET", "monitoring/jobs", nil)
	if err != nil {
		return nil, nil, err
	}
	req = req.WithContext(ctx)

	mjl := []*monitor.Job{}
	resp, err := s.client.Do(req, &mjl)
	if err != nil {
		return nil, resp, err
	}

	return mjl, resp, nil
}

// Get takes an ID and returns details for a specific monitoring job.
//
// NS1 API docs: https://ns1.com/api/#jobs-jobid-get
func (s *JobsService) Get(id string) (*monitor.Job, *http.Response, error) {
	return s.GetWithContext(context.Background(), id)
}

// GetWithContext is the same as Get, but takes a context.
func (s *JobsService) GetWithContext(ctx context.Context, id string) (*monitor.Job, *http.Response, error) {
	path := fmt.Sprintf("%s/%s", "monitoring/jobs", id)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}
	req = req.WithContext(ctx)

	var mj monitor.Job
	resp, err := s.client.Do(req, &mj)
	if err != nil {
		return nil, resp, err
	}

	return &mj, resp, nil
}

// Create takes a *MonitoringJob and creates a new monitoring job.
//
// NS1 API docs: https://ns1.com/api/#jobs-put
func (s *JobsService) Create(mj *monitor.Job) (*http.Response, error) {
	return s.CreateWithContext(context.Background(), mj)
}

// CreateWithContext is the same as Create, but takes a context.
func (s *JobsService) CreateWithContext(ctx context.Context, mj *monitor.Job) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", "monitoring/jobs", mj.ID)

	req, err := s.client.NewRequest("PUT", path, &mj)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	// Update mon jobs' fields with data from api(ensure consistent)
	resp, err := s.client.Do(req, &mj)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Update takes a *MonitoringJob and change the configuration details of an existing monitoring job.
//
// NS1 API docs: https://ns1.com/api/#jobs-jobid-post
func (s *JobsService) Update(mj *monitor.Job) (*http.Response, error) {
	return s.UpdateWithContext(context.Background(), mj)
}

// UpdateWithContext is the same as Update, but takes a context.
func (s *JobsService) UpdateWithContext(ctx context.Context, mj *monitor.Job) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", "monitoring/jobs", mj.ID)

	req, err := s.client.NewRequest("POST", path, &mj)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	// Update mon jobs' fields with data from api(ensure consistent)
	resp, err := s.client.Do(req, &mj)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Delete takes an ID and immediately terminates and deletes and existing monitoring job.
//
// NS1 API docs: https://ns1.com/api/#jobs-jobid-delete
func (s *JobsService) Delete(id string) (*http.Response, error) {
	return s.DeleteWithContext(context.Background(), id)
}

// DeleteWithContext is the same as Delete, but takes a context.
func (s *JobsService) DeleteWithContext(ctx context.Context, id string) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", "monitoring/jobs", id)

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

// History takes an ID and returns status log history for a specific monitoring job.
//
// NS1 API docs: https://ns1.com/api/#history-get
func (s *JobsService) History(id string, opts ...func(*url.Values)) ([]*monitor.StatusLog, *http.Response, error) {
	return s.HistoryWithContext(context.Background(), id, opts...)
}

// HistoryWithContext is the same as History, but takes a context.
func (s *JobsService) HistoryWithContext(ctx context.Context, id string, opts ...func(*url.Values)) ([]*monitor.StatusLog, *http.Response, error) {
	v := url.Values{}
	for _, opt := range opts {
		opt(&v)
	}

	path := fmt.Sprintf("%s/%s?%s", "monitoring/history", id, v.Encode())

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}
	req = req.WithContext(ctx)

	var slgs []*monitor.StatusLog
	resp, err := s.client.Do(req, &slgs)
	if err != nil {
		return nil, resp, err
	}

	return slgs, resp, nil
}

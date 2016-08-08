package rest

import (
	"fmt"

	"github.com/ns1/ns1-go/monitoring"
)

const (
	monitorPath = "monitoring/jobs"
)

// MonitorsService handles 'monitoring/jobs' endpoint.
type MonitorsService service

// List returns all monitoring jobs for the account.
//
// NS1 API docs: https://ns1.com/api/#jobs-get
func (s *MonitorsService) List() ([]*monitoring.Job, error) {
	req, err := s.client.NewRequest("GET", monitorPath, nil)
	if err != nil {
		return nil, err
	}

	mjl := []*monitoring.Job{}
	_, err = s.client.Do(req, &mjl)
	if err != nil {
		return nil, err
	}

	return mjl, nil
}

// Get takes an ID and returns details for a specific monitoring job.
//
// NS1 API docs: https://ns1.com/api/#jobs-jobid-get
func (s *MonitorsService) Get(id string) (*monitoring.Job, error) {
	path := fmt.Sprintf("%s/%s", monitorPath, id)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var mj monitoring.Job
	_, err = s.client.Do(req, &mj)
	if err != nil {
		return nil, err
	}

	return &mj, nil
}

// Create takes a *MonitoringJob and creates a new monitoring job.
//
// NS1 API docs: https://ns1.com/api/#jobs-put
func (s *MonitorsService) Create(mj *monitoring.Job) error {
	path := fmt.Sprintf("%s/%s", monitorPath, mj.ID)

	req, err := s.client.NewRequest("PUT", path, &mj)
	if err != nil {
		return err
	}

	// Update mon jobs' fields with data from api(ensure consistent)
	_, err = s.client.Do(req, &mj)
	if err != nil {
		return err
	}

	return nil
}

// Update takes a *MonitoringJob and change the configuration details of an existing monitoring job.
//
// NS1 API docs: https://ns1.com/api/#jobs-jobid-post
func (s *MonitorsService) Update(mj *monitoring.Job) error {
	path := fmt.Sprintf("%s/%s", monitorPath, mj.ID)

	req, err := s.client.NewRequest("POST", path, &mj)
	if err != nil {
		return err
	}

	// Update mon jobs' fields with data from api(ensure consistent)
	_, err = s.client.Do(req, &mj)
	if err != nil {
		return err
	}

	return nil
}

// Delete takes an ID and immediately terminates and deletes and existing monitoring job.
//
// NS1 API docs: https://ns1.com/api/#jobs-jobid-delete
func (s *MonitorsService) Delete(id string) error {
	path := fmt.Sprintf("%s/%s", monitorPath, id)

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

// ListTypes returns all available monitoring job types.
//
// NS1 API docs: https://ns1.com/api/#jobtypes-get
func (s *MonitorsService) ListTypes() ([]*monitoring.JobType, error) {
	path := "monitoring/jobtypes"

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	mjt := []*monitoring.JobType{}
	_, err = s.client.Do(req, &mjt)
	if err != nil {
		return nil, err
	}

	return mjt, nil
}

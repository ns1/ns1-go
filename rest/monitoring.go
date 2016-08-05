package rest

import (
	"fmt"

	ns1 "github.com/ns1/ns1-go"
)

const (
	monitorPath = "monitoring/jobs"
)

type MonitorsService service

// List returns all monitoring jobs for the account.
//
// NS1 API docs: https://ns1.com/api/#jobs-get
func (s *MonitorsService) List() ([]*ns1.MonitoringJob, error) {
	req, err := s.client.NewRequest("GET", monitorPath, nil)
	if err != nil {
		return nil, err
	}

	mj := []*ns1.MonitoringJob{}
	_, err = s.client.Do(req, &mj)
	if err != nil {
		return nil, err
	}

	return mj, nil
}

// Get takes an ID and returns details for a specific monitoring job.
//
// NS1 API docs: https://ns1.com/api/#jobs-jobid-get
func (s *MonitorsService) Get(id string) (*ns1.MonitoringJob, error) {
	path := fmt.Sprintf("%s/%s", monitorPath, id)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var mj ns1.MonitoringJob
	_, err = s.client.Do(req, &mj)
	if err != nil {
		return nil, err
	}

	return &mj, nil
}

// Create takes a *MonitoringJob and creates a new monitoring job.
//
// NS1 API docs: https://ns1.com/api/#jobs-put
func (s *MonitorsService) Create(mj *ns1.MonitoringJob) error {
	path := fmt.Sprintf("%s/%s", monitorPath, mj.Id)

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
func (s *MonitorsService) Update(mj *ns1.MonitoringJob) error {
	path := fmt.Sprintf("%s/%s", monitorPath, mj.Id)

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
func (s *MonitorsService) ListTypes() ([]*ns1.MonitoringJobType, error) {
	path := "monitoring/jobtypes"

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	mjt := []*ns1.MonitoringJobType{}
	_, err = s.client.Do(req, &mjt)
	if err != nil {
		return nil, err
	}

	return mjt, nil
}

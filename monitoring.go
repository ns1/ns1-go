package ns1

import "fmt"

const (
	monitorPath = "monitoring/jobs"
)

// MonitoringJob wraps an NS1 /monitoring/jobs resource
type MonitoringJob struct {
	Id             string                         `json:"id,omitempty"`
	Config         map[string]interface{}         `json:"config"`
	Status         map[string]MonitoringJobStatus `json:"status,omitempty"`
	Rules          []MonitoringJobRule            `json:"rules"`
	JobType        string                         `json:"job_type"`
	Regions        []string                       `json:"regions"`
	Active         bool                           `json:"active"`
	Frequency      int                            `json:"frequency"`
	Policy         string                         `json:"policy"`
	RegionScope    string                         `json:"region_scope"`
	Notes          string                         `json:"notes,omitempty"`
	Name           string                         `json:"name"`
	NotifyRepeat   int                            `json:"notify_repeat"`
	RapidRecheck   bool                           `json:"rapid_recheck"`
	NotifyDelay    int                            `json:"notify_delay"`
	NotifyList     string                         `json:"notify_list"`
	NotifyRegional bool                           `json:"notidy_regional"`
	NotifyFailback bool                           `json:"notify_failback"`
}

// MonitoringJobType wraps an element of MonitoringJobTypes
type MonitoringJobType struct {
	ShortDesc string                             `json:"shortdesc"`
	Config    map[string]interface{}             `json:"config"`
	Results   map[string]MonitoringJobTypeResult `json:"results"`
	Desc      string                             `json:"desc"`
}

// MonitoringJobTypeResult wraps an element of a MonitoringJobType's "results" attribute
type MonitoringJobTypeResult struct {
	Comparators []string `json:"comparators"`
	Metric      bool     `json:"metric"`
	Validator   string   `json:"validator"`
	ShortDesc   string   `json:"shortdesc"`
	Type        string   `json:"type"`
	Desc        string   `json:"desc"`
}

// MonitoringJobStatus wraps an value of a MonitoringJob's "status" attribute
type MonitoringJobStatus struct {
	Since  int    `json:"since"`
	Status string `json:"status"`
}

// MonitoringJobRule wraps an element of a MonitoringJob's "rules" attribute
type MonitoringJobRule struct {
	Key        string      `json:"key"`
	Value      interface{} `json:"value"`
	Comparison string      `json:"comparison"`
}

type MonitorsService service

// List returns all monitoring jobs for the account.
//
// NS1 API docs: https://ns1.com/api/#jobs-get
func (s *MonitorsService) List() ([]*MonitoringJob, error) {
	req, err := s.client.NewRequest("GET", monitorPath, nil)
	if err != nil {
		return nil, err
	}

	mj := []*MonitoringJob{}
	_, err = s.client.Do(req, &mj)
	if err != nil {
		return nil, err
	}

	return mj, nil
}

// Get takes an ID and returns details for a specific monitoring job.
//
// NS1 API docs: https://ns1.com/api/#jobs-jobid-get
func (s *MonitorsService) Get(id string) (*MonitoringJob, error) {
	path := fmt.Sprintf("%s/%s", monitorPath, id)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var mj MonitoringJob
	_, err = s.client.Do(req, &mj)
	if err != nil {
		return nil, err
	}

	return &mj, nil
}

// Create takes a *MonitoringJob and creates a new monitoring job.
//
// NS1 API docs: https://ns1.com/api/#jobs-put
func (s *MonitorsService) Create(mj *MonitoringJob) error {
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
func (s *MonitorsService) Update(mj *MonitoringJob) error {
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
func (s *MonitorsService) ListTypes() ([]*MonitoringJobType, error) {
	path := "monitoring/jobtypes"

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	mjt := []*MonitoringJobType{}
	_, err = s.client.Do(req, &mjt)
	if err != nil {
		return nil, err
	}

	return mjt, nil
}

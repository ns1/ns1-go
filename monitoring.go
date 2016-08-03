package ns1

import "fmt"

// MonitoringJobTypes wraps an NS1 /monitoring/jobtypes resource
type MonitoringJobTypes map[string]MonitoringJobType

// MonitoringJobType wraps an element of MonitoringJobTypes
type MonitoringJobType struct {
	ShortDesc string                   `json:"shortdesc"`
	Config    MonitoringJobTypeConfig  `json:"config"`
	Results   MonitoringJobTypeResults `json:"results"`
	Desc      string                   `json:"desc"`
}

// MonitoringJobTypeConfig wraps a MonitoringJobType's "config" attribute
type MonitoringJobTypeConfig map[string]interface{}

// MonitoringJobTypeResults wraps a MonitoringJobType's "results" attribute
type MonitoringJobTypeResults map[string]MonitoringJobTypeResult

// MonitoringJobTypeResult wraps an element of a MonitoringJobType's "results" attribute
type MonitoringJobTypeResult struct {
	Comparators []string `json:"comparators"`
	Metric      bool     `json:"metric"`
	Validator   string   `json:"validator"`
	ShortDesc   string   `json:"shortdesc"`
	Type        string   `json:"type"`
	Desc        string   `json:"desc"`
}

// MonitoringJobs is just a MonitoringJob array
type MonitoringJobs []MonitoringJob

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

// GetMonitoringJobTypes returns the list of all available monitoring job types
func (c APIClient) GetMonitoringJobTypes() (MonitoringJobTypes, error) {
	path := "monitoring/jobtypes"
	req, err := c.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var mjt MonitoringJobTypes
	_, err = c.Do(req, &mjt)
	if err != nil {
		return nil, err
	}

	return mjt, nil
}

// GetMonitoringJobs returns the list of all monitoring jobs for the account
func (c APIClient) GetMonitoringJobs() (MonitoringJobs, error) {
	path := "monitoring/jobs"
	req, err := c.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var mj MonitoringJobs
	_, err = c.Do(req, &mj)
	if err != nil {
		return nil, err
	}

	return mj, nil
}

// GetMonitoringJob takes an ID and returns details for a specific monitoring job
func (c APIClient) GetMonitoringJob(id string) (MonitoringJob, error) {
	var mj MonitoringJob
	path := fmt.Sprintf("monitoring/jobs/%s", id)
	req, err := c.NewRequest("GET", path, nil)
	if err != nil {
		return mj, err
	}

	_, err = c.Do(req, &mj)
	if err != nil {
		return mj, err
	}

	return mj, nil
}

// // CreateMonitoringJob takes a *MonitoringJob and creates a new monitoring job
// func (c APIClient) CreateMonitoringJob(mj *MonitoringJob) error {
// 	path := fmt.Sprintf("monitoring/jobs/%s", MonitoringJob.Id)
// 	req, err := c.NewRequest("PUT", path, &mj)
// 	if err != nil {
// 		return err
// 	}

// 	// Update mon jobs' fields with data from api(ensure consistent)
// 	_, err = c.Do(req, &mj)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// DeleteMonitoringJob takes an ID and immediately terminates and deletes and existing monitoring job
func (c APIClient) DeleteMonitoringJob(id string) error {
	path := fmt.Sprintf("monitoring/jobs/%s", id)
	req, err := c.NewRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	_, err = c.Do(req, nil)
	if err != nil {
		return err
	}

	return nil
}

// // UpdateMonitoringJob takes a *MonitoringJob and change the configuration details of an existing monitoring job
// func (c APIClient) UpdateMonitoringJob(mj *MonitoringJob) error {
// 	path := fmt.Sprintf("monitoring/jobs/%s", MonitoringJob.Id)
// 	req, err := c.NewRequest("POST", path, &mj)
// 	if err != nil {
// 		return err
// 	}

// 	// Update mon jobs' fields with data from api(ensure consistent)
// 	_, err = c.Do(req, &mj)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

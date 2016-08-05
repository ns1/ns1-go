package ns1

// MonitoringJob wraps an NS1 /monitoring/jobs resource
type MonitoringJob struct {
	ID             string                         `json:"id,omitempty"`
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

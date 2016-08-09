package monitoring

// Job wraps an NS1 /monitoring/jobs resource
type Job struct {
	ID string `json:"id,omitempty"`

	// Configuration dictionary(key/vals depend on the JobType.
	Config map[string]interface{} `json:"config"`

	// The current status of the monitor.
	Status map[string]JobStatus `json:"status,omitempty"`

	// Rules for determining failure conditions.
	Rules []*JobRule `json:"rules"`

	// Type of monitor to be run.
	JobType string `json:"job_type"`

	// List of regions in which to run the monitor.
	Regions []string `json:"regions"`

	// Indicates if the job is active or temporarily disabled.
	Active bool `json:"active"`

	// Frequency(in seconds), at which to run the monitor.
	Frequency int `json:"frequency"`

	// The policy for determining the monitor's global status based
	// on the status of the job in all regions.
	Policy string `json:"policy"`

	// Controls behavior of how the job is assigned to monitoring regions.
	// Currently this must be fixed — indicating monitoring regions are explicitly chosen.
	RegionScope string `json:"region_scope"`

	// Freeform notes to be included in any notifications about this job,
	// e.g., instructions for operators who will receive the notifications.
	Notes string `json:"notes,omitempty"`

	// A free-form display name for the monitoring job.
	Name string `json:"name"`

	// Time(in seconds) between repeat notifications of a failed job.
	// Set to 0 to disable repeating notifications.
	NotifyRepeat int `json:"notify_repeat"`

	// If true, on any apparent state change, the job is quickly re-run after
	// one second to confirm the state change before notification.
	RapidRecheck bool `json:"rapid_recheck"`

	// Time(in seconds) after a failure to wait before sending a notification.
	NotifyDelay int `json:"notify_delay"`

	// The id of the notification list to send notifications to.
	NotifyList string `json:"notify_list"`

	// If true, notifications are sent for any regional failure (and failback if desired),
	// in addition to global state notifications.
	NotifyRegional bool `json:"notidy_regional"`

	// If true, a notification is sent when a job returns to an "up" state.
	NotifyFailback bool `json:"notify_failback"`
}

// JobType wraps an element of JobTypes
type JobType struct {
	ShortDesc string                   `json:"shortdesc"`
	Config    map[string]interface{}   `json:"config"`
	Results   map[string]JobTypeResult `json:"results"`
	Desc      string                   `json:"desc"`
}

// JobTypeResult wraps an element of a JobType's "results" attribute
type JobTypeResult struct {
	Comparators []string `json:"comparators"`
	Metric      bool     `json:"metric"`
	Validator   string   `json:"validator"`
	ShortDesc   string   `json:"shortdesc"`
	Type        string   `json:"type"`
	Desc        string   `json:"desc"`
}

// JobStatus wraps an value of a Job's "status" attribute
type JobStatus struct {
	Since  int    `json:"since"`
	Status string `json:"status"`
}

// JobRule wraps an element of a Job's "rules" attribute
type JobRule struct {
	Key        string      `json:"key"`
	Value      interface{} `json:"value"`
	Comparison string      `json:"comparison"`
}

package nsone

type MonitoringJobs map[string]MonitoringJob
type MonitoringJob struct {
	ShortDesc string
	Config    MonitoringJobConfig
	Results   MonitoringJobResults
	Desc      string
}

type MonitoringJobConfig map[string]interface{}
type MonitoringJobResults map[string]MonitoringJobResult
type MonitoringJobResult struct {
	Comparators []string
	Metric      bool
	Validator   string
	ShortDesc   string
	Type        string
	Desc        string
}

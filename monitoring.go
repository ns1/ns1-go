package nsone

type MonitoringJobTypes map[string]MonitoringJobType
type MonitoringJobType struct {
	ShortDesc string
	Config    MonitoringJobTypeConfig
	Results   MonitoringJobTypeResults
	Desc      string
}

type MonitoringJobTypeConfig map[string]interface{}
type MonitoringJobTypeResults map[string]MonitoringJobTypeResult
type MonitoringJobTypeResult struct {
	Comparators []string
	Metric      bool
	Validator   string
	ShortDesc   string
	Type        string
	Desc        string
}

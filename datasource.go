package nsone

type FeedDestination struct {
	Destid   string `json"destif"`
	Desttype string `json"desttype"`
	Record   string `json"record"`
}

type DataSource struct {
	Id           string            `json:"id,omitempty"`
	Name         string            `json:"name"`
	SourceType   string            `json"source_type"`
	Config       map[string]string `json"config,omitempty"`
	Status       string            `json"status,omitempty"`
	Destinations []FeedDestination `json"destinations,omitempty"`
}

func NewDataSource(name string, source_type string) *DataSource {
	return &DataSource{
		Name:       name,
		SourceType: source_type,
	}
}

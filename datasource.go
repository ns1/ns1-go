package nsone

// FeedDestination wraps an element of a DataSource's "destinations" attribute
type FeedDestination struct {
	DestID   string `json:"destid"`
	DestType string `json:"desttype"`
	Record   string `json:"record"`
}

// DataSource wraps an NS1 /data/sources resource
type DataSource struct {
	ID           string            `json:"id,omitempty"`
	Name         string            `json:"name"`
	SourceType   string            `json:"sourcetype"`
	Config       map[string]string `json:"config,omitempty"`
	Status       string            `json:"status,omitempty"`
	Destinations []FeedDestination `json:"destinations,omitempty"`
}

// NewDataSource takes a name and sourceType and creates a new *DataSource
func NewDataSource(name string, sourceType string) *DataSource {
	cf := make(map[string]string, 0)
	return &DataSource{
		Name:       name,
		SourceType: sourceType,
		Config:     cf,
	}
}

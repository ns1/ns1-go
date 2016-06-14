package nsone

// DataFeed wraps an NS1 /data/feeds resource
type DataFeed struct {
	SourceID string                 `json:"-"`
	ID       string                 `json:"id,omitempty"`
	Name     string                 `json:"name"`
	Config   map[string]string      `json:"config,omitempty"`
	Data     map[string]interface{} `json:"data,omitempty"`
}

// NewDataFeed takes a sourceID and creates a new *DataFeed
func NewDataFeed(id string) *DataFeed {
	return &DataFeed{SourceID: id}
}

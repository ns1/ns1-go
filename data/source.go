package data

// FeedDestination wraps an element of a Source's "destinations" attribute
type FeedDestination struct {
	Destid   string `json:"destid"`
	Desttype string `json:"desttype"`
	Record   string `json:"record"`
}

// Source wraps an NS1 /data/sources resource
type Source struct {
	ID           string             `json:"id,omitempty"`
	Name         string             `json:"name"`
	SourceType   string             `json:"sourcetype"`
	Config       map[string]string  `json:"config,omitempty"`
	Status       string             `json:"status,omitempty"`
	Destinations []*FeedDestination `json:"destinations,omitempty"`
}

// NewSource takes a name and sourceType and creates a new *Source
func NewSource(name string, sourceType string) *Source {
	cf := make(map[string]string, 0)
	return &Source{
		Name:       name,
		SourceType: sourceType,
		Config:     cf,
	}
}

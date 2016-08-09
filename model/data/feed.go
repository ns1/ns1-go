package data

// Feed wraps an NS1 /data/feeds resource
type Feed struct {
	SourceID string                 `json:"-"`
	ID       string                 `json:"id,omitempty"`
	Name     string                 `json:"name"`
	Config   map[string]string      `json:"config,omitempty"`
	Data     map[string]interface{} `json:"data,omitempty"`
}

// NewFeed takes a sourceID and creates a new *Feed
func NewFeed(sourceID string) *Feed {
	return &Feed{SourceID: sourceID}
}

package data

// Destination is the target resource the receives data from a feed/source.
type Destination struct {
	ID string `json:"destid"`

	// All destinations must point to a record.
	RecordID string `json:"record"`

	// Type is the 'level' at which to apply the filters(on the targeted record).
	// Options:
	//   - answer (highest precedence)
	//   - region
	//   - record (lowest precendence)
	Type string `json:"desttype"`

	SourceID string `json:"-"`
}

func NewDestination() *Destination {
	return &Destination{}
}

// DataFeed wraps an NS1 /data/feeds resource
type Feed struct {
	ID     string `json:"id,omitempty"`
	Name   string `json:"name"`
	Config Config `json:"config,omitempty"`
	Data   Meta   `json:"data,omitempty"`

	SourceID string
}

func NewFeed(name string, cfg Config) *Feed {
	return &Feed{Name: name, Config: cfg}
}

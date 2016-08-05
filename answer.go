package ns1

// Answer wraps the values of a Record's "filters" attribute
type Answer struct {
	Region string                 `json:"region,omitempty"`
	Answer RecordAnswer           `json:"answer,omitempty"`
	Meta   map[string]interface{} `json:"meta,omitempty"`
}

type RecordAnswer []string

// NewAnswer creates an empty Answer
func NewAnswer() Answer {
	return Answer{
		Meta: make(map[string]interface{}),
	}
}

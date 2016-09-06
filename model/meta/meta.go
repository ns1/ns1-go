package meta

import (
	"encoding/json"
	"sync"
)

type Meta struct {
	mu   sync.Mutex // guards meta
	data map[string]interface{}
}

func (m *Meta) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.data)
}

func (m *Meta) UnmarshalJSON(b []byte) (err error) {
	d := map[string]interface{}{}
	if err = json.Unmarshal(b, &d); err == nil {
		m.data = d
		return
	}
	return
}

func New() *Meta {
	return &Meta{data: map[string]interface{}{}}
}

func (m *Meta) SetFeed(key string, feedID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = map[string]string{"feed": feedID}
}

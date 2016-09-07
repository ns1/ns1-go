package filter

// Filter wraps the values of a Record's "filters" attribute
type Filter struct {
	Filter   string                 `json:"filter"`
	Disabled bool                   `json:"disabled,omitempty"`
	Config   map[string]interface{} `json:"config"`
}

func (f *Filter) Enable() {
	f.Disabled = false
}

func (f *Filter) Disable() {
	f.Disabled = true
}

// NewSelFirstN returns a filter that eliminates all but the
// first N answers from the list.
func NewSelFirstN(n int) *Filter {
	config := map[string]interface{}{"N": n}
	return &Filter{
		Filter: "select_first_n",
		Config: config,
	}
}

// NewShuffle returns a filter that randomly sorts the answers.
func NewShuffle() *Filter {
	return &Filter{
		Filter: "shuffle",
		Config: map[string]interface{}{},
	}
}

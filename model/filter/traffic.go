package filter

// NewWeightedShuffle returns a filter that shuffles answers
// randomly based on their weight.
func NewWeightedShuffle() *Filter {
	return &Filter{
		Filter: "weighted_shuffle",
		Config: map[string]interface{}{},
	}
}

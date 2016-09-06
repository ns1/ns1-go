package filter

// NewUp returns a filter that eliminates all answers where
// the 'up' metadata field is not true.
func NewUp() *Filter {
	return &Filter{
		Filter: "up",
		Config: map[string]interface{}{},
	}
}

// NewPriority returns a filter that fails over according to
// prioritized answer tiers.
func NewPriority() *Filter {
	return &Filter{
		Filter: "priority",
		Config: map[string]interface{}{},
	}
}

// NewShedLoad returns a filter that "sheds" traffic to answers
// based on load, using one of several load metrics. You must set
// values for low_watermark, high_watermark, and the configured
// load metric, for each answer you intend to subject to load
// shedding.
func NewShedLoad(metric string) *Filter {
	config := map[string]interface{}{
		"metric": metric,
	}
	return &Filter{
		Filter: "shed_load",
		Config: config,
	}
}

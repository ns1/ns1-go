package meta

// This integer value indicates the "priority tier" of an
// answer or region. Lower values indicate higher priority. Values
// must be positive integers.
func (m *Meta) SetPriority(pri int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data["priority"] = pri
}

func (m Meta) Priority() interface{} {
	return m.data["priority"]
}

// This positive decimal value indicates a weight to assign
// to an answer or region. Filters that use weights normalize them,
// so you may use any positive values, but we recommend values
// between 0 and 100 for simplicity's sake.
func (m *Meta) SetWeight(weight float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data["weight"] = weight
}

func (m Meta) Weight() interface{} {
	return m.data["weight"]
}

// This positive number indicates a "low watermark" to use
// for load shedding. The value you should use depends on the metric
// you're using to determine load (e.g., loadavg, connections, etc).
func (m *Meta) SetLowWatermark(watermark int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data["low_matermark"] = watermark
}

func (m Meta) LowWatermark() interface{} {
	return m.data["low_watermark"]
}

// This positive number indicates a "high watermark" to use
// for load shedding. The value you should use depends on the metric
// you're using to determine load (e.g., loadavg, connections, etc).
func (m *Meta) SetHighWatermark(watermark int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data["high_matermark"] = watermark
}

func (m Meta) HighWatermark() interface{} {
	return m.data["high_watermark"]
}

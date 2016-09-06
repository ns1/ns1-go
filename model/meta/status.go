package meta

// This boolean value indicates "upness" of answers or
// regions. If true, the answer/region is "up".
// If false it is "down".
func (m *Meta) SetUp(val bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data["up"] = val
}

func (m Meta) Up() interface{} {
	return m.data["up"]
}

// This integer value indicates the number of active
// connections for an answer or region. Values must be positive
// integers.
func (m *Meta) SetConnections(numConns int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data["connections"] = numConns
}

func (m Meta) Connections() interface{} {
	return m.data["connections"]
}

// This integer value indicates the number of active
// requests (HTTP or otherwise) for an answer or region. Values must
// be positive integers.
func (m *Meta) SetRequests(numReqs int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data["requests"] = numReqs
}

func (m Meta) Requests() interface{} {
	return m.data["requests"]
}

// This decimal value indicates the "load average" for an
// answer or region. Values must be positive decimal numbers, and will be
// rounded to the nearest tenth.
func (m *Meta) SetLoadAvg(ldAvg float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data["loadavg"] = ldAvg
}

func (m Meta) LoadAvg() interface{} {
	return m.data["loadavg"]
}

// Pulsar telemetry gathering job and routing granularities to
// associate with the answer or region
func (m *Meta) SetPulsar(jobid string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data["pulsar"] = jobid
}

func (m Meta) Pulsar() interface{} {
	return m.data["pulsar"]
}

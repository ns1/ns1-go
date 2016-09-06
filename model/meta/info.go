package meta

// Use these freeform notes to indicate any necessary
// details about the answer or region for operators. Up to 256
// characters in length.
func (m *Meta) SetNote(notes string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data["note"] = notes
}
func (m Meta) Note() interface{} {
	return m.data["note"]
}

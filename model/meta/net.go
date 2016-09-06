package meta

// IP (v4 and v6) prefixes to associate with the answer or
// region, in CIDR format ("a.b.c.d/mask"), e.g., "1.2.3.4/24". You
// may include up to 1000 prefixes.
func (m *Meta) SetIPPrefixes(prefixes []string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data["ip_prefixes"] = prefixes
}

func (m Meta) IPPrefixes() interface{} {
	return m.data["ip_prefixes"]
}

// Autonomous System (AS) numbers to associate with the
// answer or region. You may include up to 1000 AS numbers.
func (m *Meta) SetASN(asns []string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data["asn"] = asns
}

func (m Meta) ASN() interface{} {
	return m.data["asn"]
}

package meta

// Geographic latitude associated with the answer or
// region, e.g., the latitude of the datacenter where your server
// resides. Must be a decimal number between -180 and +180 where
// negative indicates South and positive indicates North.
func (m *Meta) SetLatitude(lat float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data["latitude"] = lat
}

func (m Meta) Latitude() interface{} {
	return m.data["latitude"]
}

// Geographic longitude associated with the answer or
// region, e.g., the longitude of the datacenter where your server
// resides. Must be a decimal number between -180 and +180 where
// negative indicates West and positive indicates East.
func (m *Meta) SetLongitude(long float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data["longitude"] = long
}

func (m Meta) Longitude() interface{} {
	return m.data["longitude"]
}

// Geographic region(s) associated with this answer/region,
// e.g., the rough geographic location of the datacenter where your
// server resides. Valid geographic regions are 'US-EAST',
// 'US-CENTRAL', 'US-WEST', 'EUROPE', 'ASIAPAC', 'SOUTH-AMERICA',
// 'AFRICA'. You may assign multiple regions if desired.
func (m *Meta) SetGeographical(region []string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data["georegion"] = region
}
func (m Meta) Geographical() interface{} {
	return m.data["georegion"]
}

// Country or countries associated with this answer/region.
// Countries must be specified as ISO3166 2-character country codes.
// You may assign multiple countries if desired.
func (m *Meta) SetCountry(code []string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data["country"] = code
}
func (m Meta) Country() interface{} {
	return m.data["country"]
}

// US state(s) associated with this answer/region. States
// must be specified as standard 2-character state codes. You may
// assign multiple states if desired.
func (m *Meta) SetUSState(code []string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data["us_state"] = code
}
func (m Meta) USState() interface{} {
	return m.data["us_state"]
}

// Canadian province(s) associated with this answer/region.
// Provinces must be specified as standard 2-character province
// codes. You may assign multiple provinces if desired.
func (m *Meta) SetCAProvince(code []string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data["ca_province"] = code
}
func (m Meta) CAProvince() interface{} {
	return m.data["ca_province"]
}

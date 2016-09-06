package filter

// NewSelFirstRegion returns a filter that keeps only the answers
// that are in the same region as the first answer.
func NewSelFirstRegion() *Filter {
	return &Filter{
		Filter: "select_first_n",
		Config: map[string]interface{}{},
	}
}

// NewStickyRegion first sorts regions uniquely depending on the IP
// address of the requester, and then groups all answers together by
// region. The same requester always gets the same ordering of regions,
// but answers within each region may be in any order. byNetwork indicates
// whether to apply the 'stickyness' by subnet(not individual IP).
func NewStickyRegion(byNetwork bool) *Filter {
	config := map[string]interface{}{
		"sticky_by_network": byNetwork,
	}
	return &Filter{
		Filter: "sticky_region",
		Config: config,
	}
}

// NewGeofenceCountry returns a filter that fences using "country",
// "us_state", and "ca_province" metadata fields in answers. Only
// answers in the same country/state/province as the user (or
// answers with no specified location) are returned. rmNoLoc determines
// whether to remove answers without location on any match.
func NewGeofenceCounty(rmNoLoc bool) *Filter {
	config := map[string]interface{}{
		"remove_no_location": rmNoLoc,
	}
	return &Filter{
		Filter: "geofence_country",
		Config: config,
	}
}

// NewGeofenceRegional returns a filter that restricts to answers in
// same geographical region as requester. rmNoGeo determines whether
// to remove answers without georegion on any match.
func NewGeofenceRegional(rmNoGeo bool) *Filter {
	config := map[string]interface{}{
		"remove_no_georegion": rmNoGeo,
	}
	return &Filter{
		Filter: "geofence_regional",
		Config: config,
	}
}

// NewGeotargetCountry returns a filter that sorts answers by distance
// to requester by country, US state, and/or Canadian province.
func NewGeotargetCounty() *Filter {
	return &Filter{
		Filter: "geofence_country",
		Config: map[string]interface{}{},
	}
}

// NewGeotargetLatLong returns a filter that sorts answers by distance
// to user using lat/long.
func NewGeotargetLatLong() *Filter {
	return &Filter{
		Filter: "geotarget_latlong",
		Config: map[string]interface{}{},
	}
}

// NewGeotargetRegional returns a filter that sorts answers by distance
// to user by geographical region.
func NewGeotargetRegional() *Filter {
	return &Filter{
		Filter: "geotarget_regional",
		Config: map[string]interface{}{},
	}
}

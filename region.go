package ns1

// Region wraps the values of a Record's "regions" attribute
type Region struct {
	Meta RegionMeta `json:"meta"`
}

// RegionMeta wraps the values of a Record's "regions.*.meta" attribute
type RegionMeta struct {
	GeoRegion []string `json:"georegion,omitempty"`
	Country   []string `json:"country,omitempty"`
	USState   []string `json:"us_state,omitempty"`
	Up        bool     `json:"up,omitempty"`
}

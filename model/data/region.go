package data

type Region struct {
	Meta Meta `json:"meta,omitempty"`
}

type Regions map[string]Region

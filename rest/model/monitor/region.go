package monitor

type Region struct {
	// The region code to use when provisioning.
	Code string `json:"code"`

	// The display name for the region.
	Name string `json:"name"`

	// The list of all subnets from which monitoring traffic is generated
	// in that region.
	Subnets []string `json:"subnets"`
}

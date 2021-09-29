package dns

// Zone wraps an NS1 /zone resource
type Tsig_key struct {
	Name      string `json:"name,omitempty"`
	Algorithm string `json:"algorithm,omitempty"`
	Secret    string `json:"secret,omitempty"`
}

// NewTSIGKey takes a name, algorithm and secret and creates a new zone.
func NewTsigKey(name string, algorithm string, secret string) *Tsig_key {
	tsig_key := Tsig_key{
		Name:      name,
		Algorithm: algorithm,
		Secret:    secret,
	}
	return &tsig_key
}

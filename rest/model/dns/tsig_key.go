package dns

// TSIGKey wraps an NS1 /tsig resource
type TSIGKey struct {
	Name      string `json:"name,omitempty"`
	Algorithm string `json:"algorithm,omitempty"`
	Secret    string `json:"secret,omitempty"`
}

// NewTSIGKey takes a name, algorithm and secret and creates a new TSIG key.
func NewTsigKey(name string, algorithm string, secret string) *TSIGKey {
	tsig_key := TSIGKey{
		Name:      name,
		Algorithm: algorithm,
		Secret:    secret,
	}
	return &tsig_key
}

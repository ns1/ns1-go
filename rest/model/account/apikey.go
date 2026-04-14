package account

// APIKey wraps an NS1 /account/apikeys resource
type APIKey struct {
	// Read-only fields
	ID         string `json:"id,omitempty"`
	Key        string `json:"key,omitempty"`
	LastAccess int    `json:"last_access,omitempty"`

	Name              string         `json:"name"`
	TeamIDs           []string       `json:"teams"`
	Permissions       PermissionsMap `json:"permissions"`
	IPWhitelist       []string       `json:"ip_whitelist"`
	IPWhitelistStrict bool           `json:"ip_whitelist_strict"`

	// API Key Rotation fields
	ExpiryDuration string          `json:"expiry_duration,omitempty"`
	Secrets        []*APIKeySecret `json:"secrets,omitempty"`
}

// APIKeySecret represents a rotating secret for an API key with expiry duration
type APIKeySecret struct {
	ID         string `json:"secret_id,omitempty"`
	Key        string `json:"secret,omitempty"`
	ExpiresAt  string `json:"expires_at,omitempty"`
	LastAccess string `json:"last_access,omitempty"`
	Enabled    bool   `json:"enabled,omitempty"`
}

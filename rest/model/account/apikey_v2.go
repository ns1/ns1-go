package account

// APIKeyV2 wraps an NS1 /account/apikeys resource
type APIKeyV2 struct {
	// Read-only fields
	ID         string `json:"id,omitempty"`
	Key        string `json:"key,omitempty"`
	LastAccess int    `json:"last_access,omitempty"`

	Name              string            `json:"name"`
	TeamIDs           []string          `json:"teams"`
	Permissions       *PermissionsMapV2 `json:"permissions,omitempty"`
	IPWhitelist       []string          `json:"ip_whitelist"`
	IPWhitelistStrict bool              `json:"ip_whitelist_strict"`
}

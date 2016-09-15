package account

// APIKey wraps an NS1 /account/apikeys resource
type APIKey struct {
	ID          string         `json:"id,omitempty"`
	Name        string         `json:"name"`
	Key         string         `json:"key,omitempty"`
	LastAccess  int            `json:"last_access,omitempty"`
	Teams       []string       `json:"teams"`
	Permissions PermissionsMap `json:"permissions"`
}

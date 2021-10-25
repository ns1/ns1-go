package account

// Team wraps an NS1 /accounts/teams resource
type TeamV2 struct {
	ID          string            `json:"id,omitempty"`
	Name        string            `json:"name"`
	Permissions *PermissionsMapV2 `json:"permissions,omitempty"`
	IPWhitelist []IPWhitelist     `json:"ip_whitelist"`
}

// Same as v1
// // IPWhitelist wraps the IP whitelist for Teams.
// type IPWhitelist struct {
// 	ID     string   `json:"id,omitempty"`
// 	Name   string   `json:"name"`
// 	Values []string `json:"values"`
// }

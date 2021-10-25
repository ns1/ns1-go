package rest

import "gopkg.in/ns1/ns1-go.v2/rest/model/account"

// ddiTeamV2 wraps an NS1 /accounts/teams resource for DDI.
// Used for internally mapping between DDI permissions to maintain backwards compatibility.
type ddiTeamV2 struct {
	ID          string                `json:"id,omitempty"`
	Name        string                `json:"name"`
	Permissions *ddiPermissionsMapV2  `json:"permissions,omitempty"`
	IPWhitelist []account.IPWhitelist `json:"ip_whitelist"`
}

// ddiUserV2 wraps an NS1 /account/users resource for DDI.
// Used for internally mapping between DDI permissions to maintain backwards compatibility.
type ddiUserV2 struct {
	// Read-only fields
	LastAccess float64 `json:"last_access"`

	Name              string                       `json:"name"`
	Username          string                       `json:"username"`
	Email             string                       `json:"email"`
	TeamIDs           []string                     `json:"teams"`
	Notify            account.NotificationSettings `json:"notify"`
	IPWhitelist       []string                     `json:"ip_whitelist"`
	IPWhitelistStrict bool                         `json:"ip_whitelist_strict"`

	Permissions *ddiPermissionsMapV2 `json:"permissions,omitempty"`
}

// ddiAPIKeyV2 wraps an NS1 /account/apikeys resource for DDI specifically.
// Used for internally mapping between DDI permissions to maintain backwards compatibility.
type ddiAPIKeyV2 struct {
	// Read-only fields
	ID         string `json:"id,omitempty"`
	Key        string `json:"key,omitempty"`
	LastAccess int    `json:"last_access,omitempty"`

	Name              string   `json:"name"`
	TeamIDs           []string `json:"teams"`
	IPWhitelist       []string `json:"ip_whitelist"`
	IPWhitelistStrict bool     `json:"ip_whitelist_strict"`

	Permissions *ddiPermissionsMapV2 `json:"permissions,omitempty"`
}

// ddiPermissionsMapV2 wraps a User's "permissions" attribute for DDI.
// Used for internally mapping between DDI permissions to maintain backwards compatibility.
type ddiPermissionsMapV2 struct {
	DNS      *account.PermissionsDNSV2  `json:"dns,omitempty"`
	Data     *account.PermissionsDataV2 `json:"data,omitempty"`
	Account  *permissionsDDIAccountV2   `json:"account,omitempty"`
	Security *permissionsDDISecurityV2  `json:"security,omitempty"`
	DHCP     *account.PermissionsDHCPV2 `json:"dhcp,omitempty"`
	IPAM     *account.PermissionsIPAMV2 `json:"ipam,omitempty"`
}

// permissionsDDIAccountV2 wraps a User's "permissions.account" attribute for DDI.
// Used for internally mapping between DDI permissions to maintain backwards compatibility.
type permissionsDDIAccountV2 struct {
	ManageUsers           *bool `json:"manage_users"`
	ManageTeams           *bool `json:"manage_teams"`
	ManageApikeys         *bool `json:"manage_apikeys"`
	ManageAccountSettings *bool `json:"manage_account_settings"`
	ViewActivityLog       *bool `json:"view_activity_log"`
}

// permissionsDDISecurity wraps a User's "permissions.security" attribute for DDI.
// Used for internally mapping between DDI permissions to maintain backwards compatibility.
type permissionsDDISecurityV2 struct {
	ManageGlobal2FA       *bool `json:"manage_global_2fa"`
	ManageActiveDirectory *bool `json:"manage_active_directory"`
}

// Convert UserV2, TeamV2 and APIKeyV2 permissions to DDI permissions
func convertDDIPermissionsV2(m *account.PermissionsMapV2) *ddiPermissionsMapV2 {
	if m == nil {
		return nil
	}

	ddiMap := &ddiPermissionsMapV2{}
	if m.Account != nil {
		ddiMap.Account = &permissionsDDIAccountV2{
			ManageUsers:           m.Account.ManageUsers,
			ManageTeams:           m.Account.ManageTeams,
			ManageApikeys:         m.Account.ManageApikeys,
			ManageAccountSettings: m.Account.ManageAccountSettings,
			ViewActivityLog:       m.Account.ViewActivityLog,
		}
	}

	if m.Security != nil {
		ddiMap.Security = &permissionsDDISecurityV2{
			ManageGlobal2FA:       m.Security.ManageGlobal2FA,
			ManageActiveDirectory: m.Security.ManageActiveDirectory,
		}
	}

	if m.DHCP != nil {
		ddiMap.DHCP = m.DHCP
	}

	if m.IPAM != nil {
		ddiMap.IPAM = m.IPAM
	}

	if m.DNS != nil {
		ddiMap.DNS = m.DNS
	}

	if m.Data != nil {
		ddiMap.Data = m.Data
	}

	return ddiMap
}

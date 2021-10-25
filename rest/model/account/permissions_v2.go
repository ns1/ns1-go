package account

// PermissionsMapV2 wraps a User's "permissions" attribute
type PermissionsMapV2 struct {
	DNS        *PermissionsDNSV2        `json:"dns,omitempty"`
	Data       *PermissionsDataV2       `json:"data,omitempty"`
	Account    *PermissionsAccountV2    `json:"account,omitempty"`
	Monitoring *PermissionsMonitoringV2 `json:"monitoring,omitempty"`
	Security   *PermissionsSecurityV2   `json:"security,omitempty"`

	// DHCP and IPAM are only relevant for DDI and should not be provided in managed.
	DHCP *PermissionsDHCPV2 `json:"dhcp,omitempty"`
	IPAM *PermissionsIPAMV2 `json:"ipam,omitempty"`
}

// PermissionsDNSV2 wraps a User's "permissions.dns" attribute
type PermissionsDNSV2 struct {
	ViewZones           *bool               `json:"view_zones,omitempty"`
	ManageZones         *bool               `json:"manage_zones,omitempty"`
	ZonesAllowByDefault *bool               `json:"zones_allow_by_default,omitempty"`
	ZonesDeny           []string            `json:"zones_deny"`
	ZonesAllow          []string            `json:"zones_allow"`
	RecordsAllow        []PermissionsRecord `json:"records_allow"`
	RecordsDeny         []PermissionsRecord `json:"records_deny"`
}

// PermissionsDataV2 wraps a User's "permissions.data" attribute
type PermissionsDataV2 struct {
	PushToDatafeeds   *bool `json:"push_to_datafeeds,omitempty"`
	ManageDatasources *bool `json:"manage_datasources,omitempty"`
	ManageDatafeeds   *bool `json:"manage_datafeeds,omitempty"`
}

// PermissionsAccountV2 wraps a User's "permissions.account" attribute
type PermissionsAccountV2 struct {
	ManageUsers           *bool `json:"manage_users,omitempty"`
	ManagePaymentMethods  *bool `json:"manage_payment_methods,omitempty"`
	ManagePlan            *bool `json:"manage_plan,omitempty"`
	ManageTeams           *bool `json:"manage_teams,omitempty"`
	ManageApikeys         *bool `json:"manage_apikeys,omitempty"`
	ManageAccountSettings *bool `json:"manage_account_settings,omitempty"`
	ViewActivityLog       *bool `json:"view_activity_log,omitempty"`
	ViewInvoices          *bool `json:"view_invoices,omitempty"`
	ManageIPWhitelist     *bool `json:"manage_ip_whitelist,omitempty"`
}

// PermissionsSecurityV2 wraps a User's "permissions.security" attribute.
type PermissionsSecurityV2 struct {
	ManageGlobal2FA *bool `json:"manage_global_2fa,omitempty"`

	// This field is only relevant for DDI and should not be set to true for managed.
	ManageActiveDirectory *bool `json:"manage_active_directory,omitempty"`
}

// PermissionsMonitoringV2 wraps a User's "permissions.monitoring" attribute
// Only relevant for the managed product.
type PermissionsMonitoringV2 struct {
	ManageLists *bool `json:"manage_lists,omitempty"`
	ManageJobs  *bool `json:"manage_jobs,omitempty"`
	ViewJobs    *bool `json:"view_jobs,omitempty"`
}

// PermissionsDHCPV2 wraps a User's "permissions.dhcp" attribute for DDI.
type PermissionsDHCPV2 struct {
	ManageDHCP *bool `json:"manage_dhcp,omitempty"`
	ViewDHCP   *bool `json:"view_dhcp,omitempty"`
	// The fields below are only relevant in DDI v2.5+
	TagsAllow *[]AuthTag `json:"tags_allow"`
	TagsDeny  *[]AuthTag `json:"tags_deny"`
}

// PermissionsIPAMV2 wraps a User's "permissions.ipam" attribute for DDI.
type PermissionsIPAMV2 struct {
	ManageIPAM *bool `json:"manage_ipam,omitempty"`
	ViewIPAM   *bool `json:"view_ipam,omitempty"`
	// The fields below are only relevant in DDI v2.5+
	TagsAllow *[]AuthTag `json:"tags_allow"`
	TagsDeny  *[]AuthTag `json:"tags_deny"`
}

// Same as v1
// // AuthTag wraps the tags used in "tags_allow" and "tags_deny" in DDI and IPAM permissions in DDI.
// // Tag Names must start with prefix "auth:"
// type AuthTag struct {
// 	Name  string `json:"name"`
// 	Value string `json:"value"`
// }

// type PermissionsRecord struct {
// 	Domain     string `json:"domain"`
// 	Subdomains bool   `json:"include_subdomains"`
// 	Zone       string `json:"zone"`
// 	RecordType string `json:"type"`
// }

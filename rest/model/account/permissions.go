package account

// PermissionsMap wraps a User's "permissions" attribute
type PermissionsMap struct {
	DNS        *PermissionsDNS        `json:"dns,omitempty"`
	Data       *PermissionsData       `json:"data,omitempty"`
	Account    *PermissionsAccount    `json:"account,omitempty"`
	Monitoring *PermissionsMonitoring `json:"monitoring,omitempty"`
	Security   *PermissionsSecurity   `json:"security,omitempty"`

	// DHCP and IPAM are only relevant for DDI and should not be provided in managed.
	DHCP *PermissionsDHCP `json:"dhcp,omitempty"`
	IPAM *PermissionsIPAM `json:"ipam,omitempty"`
}

// PermissionsDNS wraps a User's "permissions.dns" attribute
type PermissionsDNS struct {
	ViewZones           *bool               `json:"view_zones,omitempty"`
	ManageZones         *bool               `json:"manage_zones,omitempty"`
	ZonesAllowByDefault *bool               `json:"zones_allow_by_default,omitempty"`
	ZonesDeny           []string            `json:"zones_deny"`
	ZonesAllow          []string            `json:"zones_allow"`
	RecordsAllow        []PermissionsRecord `json:"records_allow"`
	RecordsDeny         []PermissionsRecord `json:"records_deny"`
}

// PermissionsData wraps a User's "permissions.data" attribute
type PermissionsData struct {
	PushToDatafeeds   *bool `json:"push_to_datafeeds,omitempty"`
	ManageDatasources *bool `json:"manage_datasources,omitempty"`
	ManageDatafeeds   *bool `json:"manage_datafeeds,omitempty"`
}

// PermissionsAccount wraps a User's "permissions.account" attribute
type PermissionsAccount struct {
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

// PermissionsSecurity wraps a User's "permissions.security" attribute.
type PermissionsSecurity struct {
	ManageGlobal2FA *bool `json:"manage_global_2fa,omitempty"`

	// This field is only relevant for DDI and should not be set to true for managed.
	ManageActiveDirectory *bool `json:"manage_active_directory,omitempty"`
}

// PermissionsMonitoring wraps a User's "permissions.monitoring" attribute
// Only relevant for the managed product.
type PermissionsMonitoring struct {
	ManageLists *bool `json:"manage_lists,omitempty"`
	ManageJobs  *bool `json:"manage_jobs,omitempty"`
	ViewJobs    *bool `json:"view_jobs,omitempty"`
}

// PermissionsDHCP wraps a User's "permissions.dhcp" attribute for DDI.
type PermissionsDHCP struct {
	ManageDHCP *bool `json:"manage_dhcp,omitempty"`
	ViewDHCP   *bool `json:"view_dhcp,omitempty"`
	// The fields below are only relevant in DDI v2.5+
	TagsAllow *[]AuthTag `json:"tags_allow,omitempty"`
	TagsDeny  *[]AuthTag `json:"tags_deny,omitempty"`
}

// PermissionsIPAM wraps a User's "permissions.ipam" attribute for DDI.
type PermissionsIPAM struct {
	ManageIPAM *bool `json:"manage_ipam,omitempty"`
	ViewIPAM   *bool `json:"view_ipam,omitempty"`
	// The fields below are only relevant in DDI v2.5+
	TagsAllow *[]AuthTag `json:"tags_allow,omitempty"`
	TagsDeny  *[]AuthTag `json:"tags_deny,omitempty"`
}

// AuthTag wraps the tags used in "tags_allow" and "tags_deny" in DDI and IPAM permissions in DDI.
// Tag Names must start with prefix "auth:"
type AuthTag struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type PermissionsRecord struct {
	Domain     string `json:"domain"`
	Subdomains bool   `json:"include_subdomains"`
	Zone       string `json:"zone"`
	RecordType string `json:"type"`
}

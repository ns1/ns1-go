package billing_usage

import "errors"

type BillingUsage string

// constants values for billing-usage module implementations
const (
	BillingUsageQueries      = BillingUsage("queries")
	BillingUsageLimits       = BillingUsage("limits")
	BillingUsageDecisions    = BillingUsage("decisions")
	BillingUsageMonitors     = BillingUsage("monitors")
	BillingUsageFilterChains = BillingUsage("filter-chains")
	BillingUsageRecords      = BillingUsage("records")
)

var (
	// ErrBillingUsageNotFound bundles GET not found errors.
	ErrBillingUsageNotFound = errors.New("billing_usage not found")
	NotFound                = " not found"
)

// Queries wraps an NS1 /billing-usage-queries resource
type Queries struct {
	CleanQueries int32              `json:"clean_queries"`
	DdosQueries  int32              `json:"ddos_queries"`
	NxdResponses int32              `json:"nxd_responses"`
	ByNetwork    []QueriesByNetwork `json:"by_network"`
}

// Limits wraps an NS1 /billing-usage-limits resource
type Limits struct {
	QueriesLimit                                int32 `json:"queries_limit"`
	ChinaQueriesLimit                           int32 `json:"china_queries_limit"`
	RecordsLimit                                int32 `json:"records_limit"`
	FilterChainsLimit                           int32 `json:"tefilter_chains_limitams"`
	MonitorsLimit                               int32 `json:"monitors_limit"`
	DecisionsLimit                              int32 `json:"decisions_limit"`
	NxdProtectionEnabled                        bool  `json:"nxd_protection_enabled"`
	DdosProtectionEnabled                       bool  `json:"ddos_protection_enabled"`
	IncludeDedicatedDnsNetworkInManagedDnsUsage bool  `json:"include_dedicated_dns_network_in_managed_dns_usage"`
}

// QueriesByNetwork wraps an NS1 /billing-usage-queries.ByNetwork resource
type QueriesByNetwork struct {
	Network         int32                   `json:"network"`
	CleanQueries    int32                   `json:"clean_queries"`
	DdosQueries     int32                   `json:"ddos_queries"`
	NxdResponses    int32                   `json:"nxd_responses"`
	BillableQueries int32                   `json:"billable_queries"`
	Daily           []QueriesByNetworkDaily `json:"daily"`
}

// QueriesByNetworkDaily an NS1 /billing-usage-queries.ByNetwork.Daily resource
type QueriesByNetworkDaily struct {
	Timestamp    int32 `json:"timestamp"`
	CleanQueries int32 `json:"clean_queries"`
	DdosQueries  int32 `json:"ddos_queries"`
	NxdResponses int32 `json:"nxd_responses"`
}

// wraps an NS1 /billing-usage-monitors, billing-usage-decisions, billing-usage-records resource
type TotalUsage struct {
	TotalUsage int32 `json:"total_usage"`
}

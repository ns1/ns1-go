package mockns1

import (
	"fmt"
	"net/http"

	billingusage "gopkg.in/ns1/ns1-go.v2/rest/model/billingusage"
)

// Should be identical to rest.billingusage
const billingUsagePath = "../billing-usage/v1/"

// AddBillingUsageQueriesGetTestCase sets up a test case for the api.Client.BillingUsage.GetQueries(int32(from), int32(to))
// function
func (s *Service) AddBillingUsageQueriesGetTestCase(
	from int32,
	to int32,
	requestHeaders, responseHeaders http.Header,
	response *billingusage.Queries,
) error {

	path := fmt.Sprintf("%s/%s?from=%d&to=%d", billingUsagePath, billingusage.BillingUsageQueries, from, to)
	return s.AddTestCase(
		http.MethodGet, path, http.StatusOK, requestHeaders,
		responseHeaders, "", response,
	)
}

// AddBillingUsageQueriesFailTestCase sets up a test case for the api.Client.BillingUsage.GetQueries(int32(from), int32(to))
// functions that fails.
func (s *Service) AddBillingUsageQueriesFailTestCase(
	method string, from int32, to int32, returnStatus int,
	requestHeaders, responseHeaders http.Header,
	responseBody string,
) error {
	path := fmt.Sprintf("%s/%s?from=%d&to=%d", billingUsagePath, billingusage.BillingUsageQueries, from, to)
	return s.AddTestCase(
		method, path, returnStatus,
		nil, nil, "", responseBody)
}

// AddBillingUsageLimitsGetTestCase sets up a test case for the api.Client.BillingUsage.GetLimits(int32(from), int32(to))
// function
func (s *Service) AddBillingUsageLimitsGetTestCase(
	from int32,
	to int32,
	requestHeaders, responseHeaders http.Header,
	response *billingusage.Limits,
) error {

	path := fmt.Sprintf("%s/%s?from=%d&to=%d", billingUsagePath, billingusage.BillingUsageLimits, from, to)
	return s.AddTestCase(
		http.MethodGet, path, http.StatusOK, requestHeaders,
		responseHeaders, "", response,
	)
}

// AddBillingUsageLimitsFailTestCase sets up a test case for the api.Client.BillingUsage.GetLimits(int32(from), int32(to))
// functions that fails.
func (s *Service) AddBillingUsageLimitsFailTestCase(
	method string, from int32, to int32, returnStatus int,
	requestHeaders, responseHeaders http.Header,
	responseBody string,
) error {
	path := fmt.Sprintf("%s/%s?from=%d&to=%d", billingUsagePath, billingusage.BillingUsageLimits, from, to)
	return s.AddTestCase(
		method, path, returnStatus,
		nil, nil, "", responseBody)
}

// AddBillingUsageDecisionsGetTestCase sets up a test case for the api.Client.BillingUsage.GetDecisions(int32(from), int32(to))
// function
func (s *Service) AddBillingUsageDecisionsGetTestCase(
	from int32,
	to int32,
	requestHeaders, responseHeaders http.Header,
	response *billingusage.TotalUsage,
) error {

	path := fmt.Sprintf("%s/%s?from=%d&to=%d", billingUsagePath, billingusage.BillingUsageDecisions, from, to)
	return s.AddTestCase(
		http.MethodGet, path, http.StatusOK, requestHeaders,
		responseHeaders, "", response,
	)
}

// AddBillingUsageDecisionsFailTestCase sets up a test case for the api.Client.BillingUsage.GetDecisions(int32(from), int32(to))
// functions that fails.
func (s *Service) AddBillingUsageDecisionsFailTestCase(
	method string, from int32, to int32, returnStatus int,
	requestHeaders, responseHeaders http.Header,
	responseBody string,
) error {
	path := fmt.Sprintf("%s/%s?from=%d&to=%d", billingUsagePath, billingusage.BillingUsageDecisions, from, to)
	return s.AddTestCase(
		method, path, returnStatus,
		nil, nil, "", responseBody)
}

// AddBillingUsageMonitorsGetTestCase sets up a test case for the api.Client.BillingUsage.GetMonitors()
// function
func (s *Service) AddBillingUsageMonitorsGetTestCase(
	requestHeaders, responseHeaders http.Header,
	response *billingusage.TotalUsage,
) error {
	path := fmt.Sprintf("%s/%s", billingUsagePath, billingusage.BillingUsageMonitors)
	return s.AddTestCase(
		http.MethodGet, path, http.StatusOK, requestHeaders,
		responseHeaders, "", response,
	)
}

// AddBillingUsageMonitorsFailTestCase sets up a test case for the api.Client.BillingUsage.GetMonitors()
// functions that fails.
func (s *Service) AddBillingUsageMonitorsFailTestCase(
	method string, returnStatus int,
	requestHeaders, responseHeaders http.Header,
	responseBody string,
) error {
	path := fmt.Sprintf("%s/%s", billingUsagePath, billingusage.BillingUsageMonitors)
	return s.AddTestCase(
		method, path, returnStatus,
		nil, nil, "", responseBody)
}

// AddBillingUsageFilterChainsGetTestCase sets up a test case for the api.Client.BillingUsage.FilterChains()
// function
func (s *Service) AddBillingUsageFilterChainsGetTestCase(
	requestHeaders, responseHeaders http.Header,
	response *billingusage.TotalUsage,
) error {
	path := fmt.Sprintf("%s/%s", billingUsagePath, billingusage.BillingUsageFilterChains)
	return s.AddTestCase(
		http.MethodGet, path, http.StatusOK, requestHeaders,
		responseHeaders, "", response,
	)
}

// AddBillingUsageFilterChainsFailTestCase sets up a test case for the api.Client.BillingUsage.FilterChains()
// functions that fails.
func (s *Service) AddBillingUsageFilterChainsFailTestCase(
	method string, returnStatus int,
	requestHeaders, responseHeaders http.Header,
	responseBody string,
) error {
	path := fmt.Sprintf("%s/%s", billingUsagePath, billingusage.BillingUsageFilterChains)
	return s.AddTestCase(
		method, path, returnStatus,
		nil, nil, "", responseBody)
}

// AddBillingUsageRecordsGetTestCase sets up a test case for the api.Client.BillingUsage.GetRecords()
// function
func (s *Service) AddBillingUsageRecordsGetTestCase(
	requestHeaders, responseHeaders http.Header,
	response *billingusage.TotalUsage,
) error {
	path := fmt.Sprintf("%s/%s", billingUsagePath, billingusage.BillingUsageRecords)
	return s.AddTestCase(
		http.MethodGet, path, http.StatusOK, requestHeaders,
		responseHeaders, "", response,
	)
}

// AddBillingUsageRecordFailTestCase sets up a test case for the api.Client.BillingUsage.GetRecords()
// functions that fails.
func (s *Service) AddBillingUsageRecordFailTestCase(
	method string, returnStatus int,
	requestHeaders, responseHeaders http.Header,
	responseBody string,
) error {
	path := fmt.Sprintf("%s/%s", billingUsagePath, billingusage.BillingUsageRecords)
	return s.AddTestCase(
		method, path, returnStatus,
		nil, nil, "", responseBody)
}

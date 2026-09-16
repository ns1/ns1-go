package rest

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	billingusage "gopkg.in/ns1/ns1-go.v2/rest/model/billingusage"
)

// BillingUsageService handles 'billimg-usage/v1' endpoint.
type BillingUsageService service

// The base for the billing-usage api relative to /v1
// client.NewRequest will call ResolveReference and remove /v1/../
const billingUsageRelativeBase = "../billing-usage/v1"

// GetQueries takes the timeframe input "from" and "to", returns all its queries.
// NS1 API docs: https://developer.ibm.com/apis/catalog/ns1--ibm-ns1-connect-api/api/API--ns1--ibm-ns1-connect-api#getBillingUsageQueries
func (bu *BillingUsageService) GetQueries(from int32, to int32) (*billingusage.Queries, *http.Response, error) {
	path := fmt.Sprintf("%s/%s?from=%d&to=%d", billingUsageRelativeBase, billingusage.BillingUsageQueries, from, to)
	req, err := bu.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var queries billingusage.Queries

	resp, err := bu.client.Do(req, &queries)
	if err != nil {
		var clientErr *Error
		switch {
		case errors.As(err, &clientErr):
			if strings.HasSuffix(clientErr.Message, billingusage.NotFound) {
				return nil, resp, billingusage.ErrBillingUsageNotFound
			}
		}
		return nil, resp, err
	}

	return &queries, resp, nil
}

// GetDecisions takes the timeframe input "from" and "to", returns all its decisions.
// NS1 API docs: https://developer.ibm.com/apis/catalog/ns1--ibm-ns1-connect-api/api/API--ns1--ibm-ns1-connect-api#getBillingUsageDecisions
func (bu *BillingUsageService) GetDecisions(from int32, to int32) (*billingusage.TotalUsage, *http.Response, error) {
	path := fmt.Sprintf("%s/%s?from=%d&to=%d", billingUsageRelativeBase, billingusage.BillingUsageDecisions, from, to)

	req, err := bu.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var decisions billingusage.TotalUsage

	resp, err := bu.client.Do(req, &decisions)
	if err != nil {
		var clientErr *Error
		switch {
		case errors.As(err, &clientErr):
			if strings.HasSuffix(clientErr.Message, billingusage.NotFound) {
				return nil, resp, billingusage.ErrBillingUsageNotFound
			}
		}
		return nil, resp, err
	}

	return &decisions, resp, nil
}

// GetLimits takes the timeframe input "from" and "to", returns all its limits.
// NS1 API docs: https://developer.ibm.com/apis/catalog/ns1--ibm-ns1-connect-api/api/API--ns1--ibm-ns1-connect-api#getBillingUsageLimits
func (bu *BillingUsageService) GetLimits(from int32, to int32) (*billingusage.Limits, *http.Response, error) {
	path := fmt.Sprintf("%s/%s?from=%d&to=%d", billingUsageRelativeBase, billingusage.BillingUsageLimits, from, to)

	req, err := bu.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var limits billingusage.Limits

	resp, err := bu.client.Do(req, &limits)
	if err != nil {
		var clientErr *Error
		switch {
		case errors.As(err, &clientErr):
			if strings.HasSuffix(clientErr.Message, billingusage.NotFound) {
				return nil, resp, billingusage.ErrBillingUsageNotFound
			}
		}
		return nil, resp, err
	}

	return &limits, resp, nil
}

// GetMonitors returns total no. of monitors.
// NS1 API docs: https://developer.ibm.com/apis/catalog/ns1--ibm-ns1-connect-api/api/API--ns1--ibm-ns1-connect-api#getBillingUsageMonitors
func (bu *BillingUsageService) GetMonitors() (*billingusage.TotalUsage, *http.Response, error) {
	path := fmt.Sprintf("%s/%s", billingUsageRelativeBase, billingusage.BillingUsageMonitors)

	req, err := bu.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var monitors billingusage.TotalUsage

	resp, err := bu.client.Do(req, &monitors)
	if err != nil {
		var clientErr *Error
		switch {
		case errors.As(err, &clientErr):
			if strings.HasSuffix(clientErr.Message, billingusage.NotFound) {
				return nil, resp, billingusage.ErrBillingUsageNotFound
			}
		}
		return nil, resp, err
	}

	return &monitors, resp, nil
}

// GetFilterChains returns total no. of filter-chains.
// NS1 API docs: https://developer.ibm.com/apis/catalog/ns1--ibm-ns1-connect-api/api/API--ns1--ibm-ns1-connect-api#getBillingUsageFilterChains
func (bu *BillingUsageService) GetFilterChains() (*billingusage.TotalUsage, *http.Response, error) {
	path := fmt.Sprintf("%s/%s", billingUsageRelativeBase, billingusage.BillingUsageFilterChains)

	req, err := bu.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var filterChains billingusage.TotalUsage

	resp, err := bu.client.Do(req, &filterChains)
	if err != nil {
		var clientErr *Error
		switch {
		case errors.As(err, &clientErr):
			if strings.HasSuffix(clientErr.Message, billingusage.NotFound) {
				return nil, resp, billingusage.ErrBillingUsageNotFound
			}
		}
		return nil, resp, err
	}

	return &filterChains, resp, nil
}

// GetRecords returns total no. of records.
// NS1 API docs: https://developer.ibm.com/apis/catalog/ns1--ibm-ns1-connect-api/api/API--ns1--ibm-ns1-connect-api#getBillingUsageRecords
func (bu *BillingUsageService) GetRecords() (*billingusage.TotalUsage, *http.Response, error) {
	path := fmt.Sprintf("%s/%s", billingUsageRelativeBase, billingusage.BillingUsageRecords)

	req, err := bu.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var records billingusage.TotalUsage

	resp, err := bu.client.Do(req, &records)
	if err != nil {
		var clientErr *Error
		switch {
		case errors.As(err, &clientErr):
			if strings.HasSuffix(clientErr.Message, billingusage.NotFound) {
				return nil, resp, billingusage.ErrBillingUsageNotFound
			}
		}
		return nil, resp, err
	}

	return &records, resp, nil
}

// GetRedirects returns total no. of redirects.
// NS1 API docs: https://developer.ibm.com/apis/catalog/ns1--ibm-ns1-connect-api/api/API--ns1--ibm-ns1-connect-api#getBillingUsageRedirects
func (bu *BillingUsageService) GetRedirects() (*billingusage.TotalUsage, *http.Response, error) {
	path := fmt.Sprintf("%s/%s", billingUsageRelativeBase, billingusage.BillingUsageRedirects)

	req, err := bu.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var redirects billingusage.TotalUsage

	resp, err := bu.client.Do(req, &redirects)
	if err != nil {
		var clientErr *Error
		switch {
		case errors.As(err, &clientErr):
			if strings.HasSuffix(clientErr.Message, billingusage.NotFound) {
				return nil, resp, billingusage.ErrBillingUsageNotFound
			}
		}
		return nil, resp, err
	}

	return &redirects, resp, nil
}

package rest

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	billing_usage "gopkg.in/ns1/ns1-go.v2/rest/model/billing-usage"
)

// BillingUsageService handles 'billimg-usage/v1' endpoint.
type BillingUsageService service

// The base for the billing-usage api relative to /v1
// client.NewRequest will call ResolveReference and remove /v1/../
const billingUsageRelativeBase = "../billing-usage/v1"

// GetQueries takes the timeframe input "from" and "to", returns all its queries.
// NS1 API docs: https://ns1.com/api/#billing-usage-queries-get
func (bu *BillingUsageService) GetQueries(from int32, to int32) (*billing_usage.Queries, *http.Response, error) {
	path := fmt.Sprintf("%s/%s?from=%d&to=%d", billingUsageRelativeBase, billing_usage.BillingUsageQueries, from, to)
	req, err := bu.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var queries billing_usage.Queries

	resp, err := bu.client.Do(req, &queries)
	if err != nil {
		var clientErr *Error
		switch {
		case errors.As(err, &clientErr):
			if strings.HasSuffix(clientErr.Message, billing_usage.NotFound) {
				return nil, resp, billing_usage.ErrBillingUsageNotFound
			}
		}
		return nil, resp, err
	}

	return &queries, resp, nil
}

// GetDecisions takes the timeframe input "from" and "to", returns all its decisions.
// NS1 API docs: https://ns1.com/api/#billing-usage-decisions-get
func (bu *BillingUsageService) GetDecisions(from int32, to int32) (*billing_usage.TotalUsage, *http.Response, error) {
	path := fmt.Sprintf("%s/%s?from=%d&to=%d", billingUsageRelativeBase, billing_usage.BillingUsageDecisions, from, to)

	req, err := bu.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var decisions billing_usage.TotalUsage

	resp, err := bu.client.Do(req, &decisions)
	if err != nil {
		var clientErr *Error
		switch {
		case errors.As(err, &clientErr):
			if strings.HasSuffix(clientErr.Message, billing_usage.NotFound) {
				return nil, resp, billing_usage.ErrBillingUsageNotFound
			}
		}
		return nil, resp, err
	}

	return &decisions, resp, nil
}

// GetLimits takes the timeframe input "from" and "to", returns all its limits.
// NS1 API docs: https://ns1.com/api/#billing-usage-limits-get
func (bu *BillingUsageService) GetLimits(from int32, to int32) (*billing_usage.Limits, *http.Response, error) {
	path := fmt.Sprintf("%s/%s?from=%d&to=%d", billingUsageRelativeBase, billing_usage.BillingUsageLimits, from, to)

	req, err := bu.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var limits billing_usage.Limits

	resp, err := bu.client.Do(req, &limits)
	if err != nil {
		var clientErr *Error
		switch {
		case errors.As(err, &clientErr):
			if strings.HasSuffix(clientErr.Message, billing_usage.NotFound) {
				return nil, resp, billing_usage.ErrBillingUsageNotFound
			}
		}
		return nil, resp, err
	}

	return &limits, resp, nil
}

// GetMonitors returns total no. of monitors.
// NS1 API docs: https://ns1.com/api/#billing-usage-monitors-get
func (bu *BillingUsageService) GetMonitors() (*billing_usage.TotalUsage, *http.Response, error) {
	path := fmt.Sprintf("%s/%s", billingUsageRelativeBase, billing_usage.BillingUsageMonitors)

	req, err := bu.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var monitors billing_usage.TotalUsage

	resp, err := bu.client.Do(req, &monitors)
	if err != nil {
		var clientErr *Error
		switch {
		case errors.As(err, &clientErr):
			if strings.HasSuffix(clientErr.Message, billing_usage.NotFound) {
				return nil, resp, billing_usage.ErrBillingUsageNotFound
			}
		}
		return nil, resp, err
	}

	return &monitors, resp, nil
}

// GetFilterChains returns total no. of filter-chains.
// NS1 API docs: https://ns1.com/api/#billing-usage-filter-chains-get
func (bu *BillingUsageService) GetFilterChains() (*billing_usage.TotalUsage, *http.Response, error) {
	path := fmt.Sprintf("%s/%s", billingUsageRelativeBase, billing_usage.BillingUsageFilterChains)

	req, err := bu.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var filterChains billing_usage.TotalUsage

	resp, err := bu.client.Do(req, &filterChains)
	if err != nil {
		var clientErr *Error
		switch {
		case errors.As(err, &clientErr):
			if strings.HasSuffix(clientErr.Message, billing_usage.NotFound) {
				return nil, resp, billing_usage.ErrBillingUsageNotFound
			}
		}
		return nil, resp, err
	}

	return &filterChains, resp, nil
}

// GetRecords returns total no. of records.
// NS1 API docs: https://ns1.com/api/#billing-usage-records-get
func (bu *BillingUsageService) GetRecords() (*billing_usage.TotalUsage, *http.Response, error) {
	path := fmt.Sprintf("%s/%s", billingUsageRelativeBase, billing_usage.BillingUsageRecords)

	req, err := bu.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var records billing_usage.TotalUsage

	resp, err := bu.client.Do(req, &records)
	if err != nil {
		var clientErr *Error
		switch {
		case errors.As(err, &clientErr):
			if strings.HasSuffix(clientErr.Message, billing_usage.NotFound) {
				return nil, resp, billing_usage.ErrBillingUsageNotFound
			}
		}
		return nil, resp, err
	}

	return &records, resp, nil
}

package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"gopkg.in/ns1/ns1-go.v2/rest/model/alerting"
)

// UsageAlertListParams defines the parameters for listing usage alerts
type UsageAlertListParams struct {
	Limit           int
	Next            string
	OrderDescending bool
}

// CreateUsageAlert creates a new usage alert
//
// NS1 API docs: https://ns1.com/api/#alert-post
func (s *AlertsService) CreateUsageAlert(name, subtype string, alertAtPercent int, notifierListIds []string) (*alerting.Alert, *http.Response, error) {
	alert := alerting.NewUsageAlert(name, subtype, alertAtPercent, notifierListIds)
	
	if err := alerting.ValidateUsageAlert(alert); err != nil {
		return nil, nil, err
	}
	
	resp, err := s.Create(alert)
	if err != nil {
		return nil, resp, err
	}
	
	return alert, resp, nil
}

// GetUsageAlert returns a specific usage alert by ID
//
// NS1 API docs: https://ns1.com/api/#alert-alertid-get
func (s *AlertsService) GetUsageAlert(alertID string) (*alerting.Alert, *http.Response, error) {
	alert, resp, err := s.Get(alertID)
	if err != nil {
		return nil, resp, err
	}
	
	// Verify it's a usage alert
	if alert.Type == nil || *alert.Type != alerting.AlertTypeAccount {
		return nil, resp, fmt.Errorf("alert is not a usage alert (type is not 'account')")
	}
	
	if alert.Subtype == nil {
		return nil, resp, fmt.Errorf("alert is missing subtype")
	}
	
	if _, ok := alerting.AllowedUsageSubtypes[*alert.Subtype]; !ok {
		return nil, resp, fmt.Errorf("alert is not a usage alert (invalid subtype)")
	}
	
	return alert, resp, nil
}

// UsageAlertPatch defines the fields that can be updated on a usage alert
// Intentionally excluding Type and Subtype which cannot be changed via PATCH
type UsageAlertPatch struct {
	Name            *string             `json:"name,omitempty"`
	Data            json.RawMessage     `json:"data,omitempty"`
	NotifierListIds *[]string           `json:"notifier_list_ids,omitempty"`
	ZoneNames       *[]string           `json:"zone_names,omitempty"`
}

// PatchUsageAlert updates a usage alert
//
// NS1 API docs: https://ns1.com/api/#alert-alertid-patch
func (s *AlertsService) PatchUsageAlert(alertID string, patch *UsageAlertPatch) (*alerting.Alert, *http.Response, error) {
	// First get the alert to ensure it's a usage alert
	_, resp, err := s.GetUsageAlert(alertID)
	if err != nil {
		return nil, resp, err
	}
	
	// Create alert object with only fields that can be patched
	alert := &alerting.Alert{
		ID: &alertID,
	}
	
	// Only set fields that are provided
	if patch.Name != nil {
		alert.Name = patch.Name
	}
	
	if patch.Data != nil {
		alert.Data = patch.Data
	}
	
	if patch.NotifierListIds != nil {
		alert.NotifierListIds = *patch.NotifierListIds
	}
	
	if patch.ZoneNames != nil {
		alert.ZoneNames = *patch.ZoneNames
	}
	
	// IMPORTANT: Do not set Type or Subtype in PATCH operations
	// The server will ignore these fields, but we explicitly exclude them
	// for clarity and to prevent any potential issues
	
	resp, err = s.Update(alert)
	if err != nil {
		return nil, resp, err
	}
	
	// Re-fetch the alert to get the updated values
	return s.GetUsageAlert(alertID)
}

// PatchUsageAlertSimple provides a simplified interface to patch a usage alert
// with basic fields that are commonly updated
//
// NS1 API docs: https://ns1.com/api/#alert-alertid-patch
func (s *AlertsService) PatchUsageAlertSimple(alertID string, name string, alertAtPercent int, notifierListIds []string) (*alerting.Alert, *http.Response, error) {
	patch := &UsageAlertPatch{}
	
	// Only set fields that are specified
	if name != "" {
		patch.Name = &name
	}
	
	if alertAtPercent > 0 {
		data := alerting.UsageAlertData{AlertAtPercent: alertAtPercent}
		patch.Data = alerting.MarshalUsageAlertData(data)
	}
	
	if notifierListIds != nil {
		patch.NotifierListIds = &notifierListIds
	}
	
	return s.PatchUsageAlert(alertID, patch)
}

// DeleteUsageAlert deletes a usage alert
//
// NS1 API docs: https://ns1.com/api/#alert-alertid-delete
func (s *AlertsService) DeleteUsageAlert(alertID string) (*http.Response, error) {
	// First get the alert to ensure it's a usage alert
	_, resp, err := s.GetUsageAlert(alertID)
	if err != nil {
		return resp, err
	}
	
	return s.Delete(alertID)
}

// ListUsageAlerts returns all usage alerts with pagination support
// This method filters the results to only include usage alerts (type="account" and valid subtypes)
//
// NS1 API docs: https://ns1.com/api/#alerts-get
func (s *AlertsService) ListUsageAlerts(params *UsageAlertListParams) ([]*alerting.Alert, *alertListResponse, *http.Response, error) {
	path := fmt.Sprintf("%s/%s", alertingRelativeBase, "alerts")
	
	// Add query parameters
	v := url.Values{}
	if params != nil {
		if params.Limit > 0 {
			v.Set("limit", strconv.Itoa(params.Limit))
		}
		if params.Next != "" {
			v.Set("next", params.Next)
		}
		if params.OrderDescending {
			v.Set("order_descending", "true")
		}
	}
	
	if queryStr := v.Encode(); queryStr != "" {
		path = fmt.Sprintf("%s?%s", path, queryStr)
	}
	
	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, nil, err
	}
	
	alertListResp := alertListResponse{}
	var resp *http.Response
	if s.client.FollowPagination {
		resp, err = s.client.DoWithPagination(req, &alertListResp, s.nextAlerts)
	} else {
		resp, err = s.client.Do(req, &alertListResp)
	}
	if err != nil {
		return nil, nil, resp, err
	}
	
	// Filter for usage alerts only
	usageAlerts := make([]*alerting.Alert, 0)
	for _, alert := range alertListResp.Results {
		if alert.Type != nil && *alert.Type == alerting.AlertTypeAccount && 
		   alert.Subtype != nil && alerting.AllowedUsageSubtypes[*alert.Subtype] {
			usageAlerts = append(usageAlerts, alert)
		}
	}
	
	return usageAlerts, &alertListResp, resp, nil
}

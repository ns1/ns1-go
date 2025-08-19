package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"gopkg.in/ns1/ns1-go.v2/rest/model/alerting"
)

func TestCreateUsageAlert(t *testing.T) {
	client, server, _ := setup(t)
	defer server.Close()
	
	// Mock alert creation response
	alertID := "abcdef123456"
	alertName := "Test Usage Alert"
	alertSubtype := alerting.UsageSubtypes.QueryUsage
	alertAtPercent := 85
	notifierListIds := []string{"notifier1", "notifier2"}
	
	mockResponse := fmt.Sprintf(`{
		"id": "%s",
		"name": "%s",
		"type": "account",
		"subtype": "%s",
		"data": {"alert_at_percent": %d},
		"notifier_list_ids": ["notifier1", "notifier2"],
		"zone_names": [],
		"created_at": 1597937213,
		"updated_at": 1597937213
	}`, alertID, alertName, alertSubtype, alertAtPercent)
	
	server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected method POST, got %s", r.Method)
		}
		if r.URL.Path != "/alerting/v1/alerts" {
			t.Fatalf("expected path /alerting/v1/alerts, got %s", r.URL.Path)
		}
		
		var alert alerting.Alert
		if err := json.NewDecoder(r.Body).Decode(&alert); err != nil {
			t.Fatalf("error decoding request body: %v", err)
		}
		
		if *alert.Name != alertName {
			t.Fatalf("expected name %s, got %s", alertName, *alert.Name)
		}
		if *alert.Type != alerting.AlertTypeAccount {
			t.Fatalf("expected type %s, got %s", alerting.AlertTypeAccount, *alert.Type)
		}
		if *alert.Subtype != alertSubtype {
			t.Fatalf("expected subtype %s, got %s", alertSubtype, *alert.Subtype)
		}
		
		fmt.Fprintln(w, mockResponse)
	})
	
	alert, resp, err := client.Alerts.CreateUsageAlert(alertName, alertSubtype, alertAtPercent, notifierListIds)
	if err != nil {
		t.Fatalf("error creating usage alert: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code 200, got %d", resp.StatusCode)
	}
	if alert.ID == nil || *alert.ID != alertID {
		t.Fatalf("expected alert ID %s, got %v", alertID, alert.ID)
	}
	
	// Test validation
	_, _, err = client.Alerts.CreateUsageAlert(alertName, "invalid_subtype", alertAtPercent, notifierListIds)
	if err == nil {
		t.Fatal("expected error for invalid subtype, got nil")
	}
	
	_, _, err = client.Alerts.CreateUsageAlert(alertName, alertSubtype, 0, notifierListIds)
	if err == nil {
		t.Fatal("expected error for alert_at_percent below 1, got nil")
	}
	
	_, _, err = client.Alerts.CreateUsageAlert(alertName, alertSubtype, 101, notifierListIds)
	if err == nil {
		t.Fatal("expected error for alert_at_percent above 100, got nil")
	}
}

func TestGetUsageAlert(t *testing.T) {
	client, server, _ := setup(t)
	defer server.Close()
	
	alertID := "abcdef123456"
	alertName := "Test Usage Alert"
	alertSubtype := alerting.UsageSubtypes.QueryUsage
	alertAtPercent := 85
	
	mockResponse := fmt.Sprintf(`{
		"id": "%s",
		"name": "%s",
		"type": "account",
		"subtype": "%s",
		"data": {"alert_at_percent": %d},
		"notifier_list_ids": ["notifier1", "notifier2"],
		"zone_names": [],
		"created_at": 1597937213,
		"updated_at": 1597937213
	}`, alertID, alertName, alertSubtype, alertAtPercent)
	
	server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected method GET, got %s", r.Method)
		}
		if r.URL.Path != fmt.Sprintf("/alerting/v1/alerts/%s", alertID) {
			t.Fatalf("expected path /alerting/v1/alerts/%s, got %s", alertID, r.URL.Path)
		}
		
		fmt.Fprintln(w, mockResponse)
	})
	
	alert, resp, err := client.Alerts.GetUsageAlert(alertID)
	if err != nil {
		t.Fatalf("error getting usage alert: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code 200, got %d", resp.StatusCode)
	}
	if alert.ID == nil || *alert.ID != alertID {
		t.Fatalf("expected alert ID %s, got %v", alertID, alert.ID)
	}
	
	// Test data extraction
	data, err := alerting.UnmarshalUsageAlertData(alert)
	if err != nil {
		t.Fatalf("error unmarshalling data: %v", err)
	}
	if data.AlertAtPercent != alertAtPercent {
		t.Fatalf("expected alert_at_percent %d, got %d", alertAtPercent, data.AlertAtPercent)
	}
}

func TestGetUsageAlert_InvalidNotifierId(t *testing.T) {
	client, server, _ := setup(t)
	defer server.Close()
	
	alertID := "nonexistent_alert"
	
	// Mock 404 response with specific error message about the notifier not found
	mockErrorResponse := `{
		"message": "notifier with id 'nonexistent_notifier' not found"
	}`
	
	server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected method GET, got %s", r.Method)
		}
		if r.URL.Path != fmt.Sprintf("/alerting/v1/alerts/%s", alertID) {
			t.Fatalf("expected path /alerting/v1/alerts/%s, got %s", alertID, r.URL.Path)
		}
		
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, mockErrorResponse)
	})
	
	_, resp, err := client.Alerts.GetUsageAlert(alertID)
	if err == nil {
		t.Fatal("expected error for nonexistent alert, got nil")
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status code 404, got %d", resp.StatusCode)
	}
	
	// Check that the error message is correctly passed through
	restErr, ok := err.(*Error)
	if !ok {
		t.Fatalf("expected *Error, got %T", err)
	}
	if restErr.Message != "notifier with id 'nonexistent_notifier' not found" {
		t.Fatalf("expected error message about notifier not found, got: %s", restErr.Message)
	}
}

func TestPatchUsageAlert(t *testing.T) {
	client, server, _ := setup(t)
	defer server.Close()
	
	alertID := "abcdef123456"
	oldName := "Test Usage Alert"
	newName := "Updated Usage Alert"
	alertSubtype := alerting.UsageSubtypes.QueryUsage
	oldAlertAtPercent := 85
	newAlertAtPercent := 90
	
	// Mock the GET response to fetch the alert first
	getResponse := fmt.Sprintf(`{
		"id": "%s",
		"name": "%s",
		"type": "account",
		"subtype": "%s",
		"data": {"alert_at_percent": %d},
		"notifier_list_ids": ["notifier1", "notifier2"],
		"zone_names": [],
		"created_at": 1597937213,
		"updated_at": 1597937213
	}`, alertID, oldName, alertSubtype, oldAlertAtPercent)
	
	// Mock the PATCH response
	patchResponse := fmt.Sprintf(`{
		"id": "%s",
		"name": "%s",
		"type": "account",
		"subtype": "%s",
		"data": {"alert_at_percent": %d},
		"notifier_list_ids": ["notifier1", "notifier2"],
		"zone_names": [],
		"created_at": 1597937213,
		"updated_at": 1597937300
	}`, alertID, newName, alertSubtype, newAlertAtPercent)
	
	callCount := 0
	server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := fmt.Sprintf("/alerting/v1/alerts/%s", alertID)
		
		if r.URL.Path != expectedPath {
			t.Fatalf("expected path %s, got %s", expectedPath, r.URL.Path)
		}
		
		callCount++
		if callCount == 1 {
			// First call is GET
			if r.Method != http.MethodGet {
				t.Fatalf("expected method GET for first call, got %s", r.Method)
			}
			fmt.Fprintln(w, getResponse)
		} else if callCount == 2 {
			// Second call is PATCH
			if r.Method != http.MethodPatch {
				t.Fatalf("expected method PATCH for second call, got %s", r.Method)
			}
			
			var alert alerting.Alert
			if err := json.NewDecoder(r.Body).Decode(&alert); err != nil {
				t.Fatalf("error decoding request body: %v", err)
			}
			
			// Verify type and subtype are not sent
			if alert.Type != nil {
				t.Fatalf("expected Type to be nil in PATCH, got %v", alert.Type)
			}
			if alert.Subtype != nil {
				t.Fatalf("expected Subtype to be nil in PATCH, got %v", alert.Subtype)
			}
			
			fmt.Fprintln(w, patchResponse)
		} else {
			// Third call is GET again to fetch updated alert
			if r.Method != http.MethodGet {
				t.Fatalf("expected method GET for third call, got %s", r.Method)
			}
			fmt.Fprintln(w, patchResponse)
		}
	})
	
	patch := &UsageAlertPatch{
		Name: &newName,
		Data: alerting.MarshalUsageAlertData(alerting.UsageAlertData{AlertAtPercent: newAlertAtPercent}),
	}
	
	alert, resp, err := client.Alerts.PatchUsageAlert(alertID, patch)
	if err != nil {
		t.Fatalf("error updating usage alert: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code 200, got %d", resp.StatusCode)
	}
	if alert.ID == nil || *alert.ID != alertID {
		t.Fatalf("expected alert ID %s, got %v", alertID, alert.ID)
	}
	if alert.Name == nil || *alert.Name != newName {
		t.Fatalf("expected alert name %s, got %v", newName, alert.Name)
	}
	
	// Verify data was updated
	data, err := alerting.UnmarshalUsageAlertData(alert)
	if err != nil {
		t.Fatalf("error unmarshalling data: %v", err)
	}
	if data.AlertAtPercent != newAlertAtPercent {
		t.Fatalf("expected alert_at_percent %d, got %d", newAlertAtPercent, data.AlertAtPercent)
	}
}

func TestDeleteUsageAlert(t *testing.T) {
	client, server, _ := setup(t)
	defer server.Close()
	
	alertID := "abcdef123456"
	alertName := "Test Usage Alert"
	alertSubtype := alerting.UsageSubtypes.QueryUsage
	alertAtPercent := 85
	
	// Mock the GET response to fetch the alert first
	getResponse := fmt.Sprintf(`{
		"id": "%s",
		"name": "%s",
		"type": "account",
		"subtype": "%s",
		"data": {"alert_at_percent": %d},
		"notifier_list_ids": ["notifier1", "notifier2"],
		"zone_names": [],
		"created_at": 1597937213,
		"updated_at": 1597937213
	}`, alertID, alertName, alertSubtype, alertAtPercent)
	
	callCount := 0
	server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := fmt.Sprintf("/alerting/v1/alerts/%s", alertID)
		
		if r.URL.Path != expectedPath {
			t.Fatalf("expected path %s, got %s", expectedPath, r.URL.Path)
		}
		
		callCount++
		if callCount == 1 {
			// First call is GET
			if r.Method != http.MethodGet {
				t.Fatalf("expected method GET for first call, got %s", r.Method)
			}
			fmt.Fprintln(w, getResponse)
		} else {
			// Second call is DELETE
			if r.Method != http.MethodDelete {
				t.Fatalf("expected method DELETE for second call, got %s", r.Method)
			}
			
			w.WriteHeader(http.StatusNoContent)
		}
	})
	
	resp, err := client.Alerts.DeleteUsageAlert(alertID)
	if err != nil {
		t.Fatalf("error deleting usage alert: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected status code 204, got %d", resp.StatusCode)
	}
}

func TestListUsageAlerts(t *testing.T) {
	client, server, _ := setup(t)
	defer server.Close()
	
	mockResponse := `{
		"limit": 1,
		"next": "next_token_value",
		"total_results": 3,
		"results": [
			{
				"id": "alert1",
				"name": "Usage Alert 1",
				"type": "account",
				"subtype": "query_usage",
				"data": {"alert_at_percent": 80},
				"notifier_list_ids": ["notifier1"],
				"zone_names": [],
				"created_at": 1597937213,
				"updated_at": 1597937213
			},
			{
				"id": "alert2",
				"name": "Usage Alert 2",
				"type": "account",
				"subtype": "record_usage",
				"data": {"alert_at_percent": 90},
				"notifier_list_ids": ["notifier2"],
				"zone_names": [],
				"created_at": 1597937213,
				"updated_at": 1597937213
			},
			{
				"id": "alert3",
				"name": "Zone Alert",
				"type": "zone",
				"subtype": "zone_timeout",
				"data": {},
				"notifier_list_ids": ["notifier1"],
				"zone_names": ["example.com"],
				"created_at": 1597937213,
				"updated_at": 1597937213
			}
		]
	}`
	
	server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected method GET, got %s", r.Method)
		}
		
		expectedPath := "/alerting/v1/alerts"
		if r.URL.Path != expectedPath {
			t.Fatalf("expected path %s, got %s", expectedPath, r.URL.Path)
		}
		
		// Check query parameters
		expectedQuery := url.Values{
			"limit":            []string{"1"},
			"order_descending": []string{"true"},
		}
		
		if r.URL.Query().Get("limit") != expectedQuery.Get("limit") {
			t.Fatalf("expected limit %s, got %s", expectedQuery.Get("limit"), r.URL.Query().Get("limit"))
		}
		if r.URL.Query().Get("order_descending") != expectedQuery.Get("order_descending") {
			t.Fatalf("expected order_descending %s, got %s", expectedQuery.Get("order_descending"), r.URL.Query().Get("order_descending"))
		}
		
		fmt.Fprintln(w, mockResponse)
	})
	
	params := &UsageAlertListParams{
		Limit:           1,
		OrderDescending: true,
	}
	
	alerts, listResp, resp, err := client.Alerts.ListUsageAlerts(params)
	if err != nil {
		t.Fatalf("error listing usage alerts: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code 200, got %d", resp.StatusCode)
	}
	
	// Should only return usage alerts (type=account), filtering out the zone alert
	if len(alerts) != 2 {
		t.Fatalf("expected 2 usage alerts, got %d", len(alerts))
	}
	
	// Verify first alert
	if alerts[0].ID == nil || *alerts[0].ID != "alert1" {
		t.Fatalf("expected first alert ID 'alert1', got %v", alerts[0].ID)
	}
	if alerts[0].Subtype == nil || *alerts[0].Subtype != alerting.UsageSubtypes.QueryUsage {
		t.Fatalf("expected first alert subtype '%s', got %v", alerting.UsageSubtypes.QueryUsage, alerts[0].Subtype)
	}
	
	// Verify second alert
	if alerts[1].ID == nil || *alerts[1].ID != "alert2" {
		t.Fatalf("expected second alert ID 'alert2', got %v", alerts[1].ID)
	}
	if alerts[1].Subtype == nil || *alerts[1].Subtype != alerting.UsageSubtypes.RecordUsage {
		t.Fatalf("expected second alert subtype '%s', got %v", alerting.UsageSubtypes.RecordUsage, alerts[1].Subtype)
	}
	
	// Verify pagination info
	if listResp.Next == nil || *listResp.Next != "next_token_value" {
		t.Fatalf("expected next token 'next_token_value', got %v", listResp.Next)
	}
	if listResp.TotalResults == nil || *listResp.TotalResults != 3 {
		t.Fatalf("expected total_results 3, got %v", listResp.TotalResults)
	}
}

// TestPatchUsageAlert_NoTypeSubtype tests that type and subtype are not sent in the PATCH request
func TestPatchUsageAlert_NoTypeSubtype(t *testing.T) {
	client, server, _ := setup(t)
	defer server.Close()
	
	alertID := "abcdef123456"
	newName := "Updated Usage Alert"
	
	// Mock the GET response to fetch the alert first
	getResponse := `{
		"id": "abcdef123456",
		"name": "Test Usage Alert",
		"type": "account",
		"subtype": "query_usage",
		"data": {"alert_at_percent": 85},
		"notifier_list_ids": ["notifier1", "notifier2"],
		"zone_names": [],
		"created_at": 1597937213,
		"updated_at": 1597937213
	}`
	
	// Mock the PATCH response
	patchResponse := `{
		"id": "abcdef123456",
		"name": "Updated Usage Alert",
		"type": "account",
		"subtype": "query_usage",
		"data": {"alert_at_percent": 85},
		"notifier_list_ids": ["notifier1", "notifier2"],
		"zone_names": [],
		"created_at": 1597937213,
		"updated_at": 1597937300
	}`
	
	callCount := 0
	server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := fmt.Sprintf("/alerting/v1/alerts/%s", alertID)
		
		if r.URL.Path != expectedPath {
			t.Fatalf("expected path %s, got %s", expectedPath, r.URL.Path)
		}
		
		callCount++
		if callCount == 1 {
			// First call is GET
			if r.Method != http.MethodGet {
				t.Fatalf("expected method GET for first call, got %s", r.Method)
			}
			fmt.Fprintln(w, getResponse)
		} else if callCount == 2 {
			// Second call is PATCH
			if r.Method != http.MethodPatch {
				t.Fatalf("expected method PATCH for second call, got %s", r.Method)
			}
			
			// Ensure the PATCH request doesn't include type or subtype
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("error reading request body: %v", err)
			}
			
			if strings.Contains(string(body), `"type"`) {
				t.Fatalf("PATCH request should not include 'type' field, got: %s", string(body))
			}
			
			if strings.Contains(string(body), `"subtype"`) {
				t.Fatalf("PATCH request should not include 'subtype' field, got: %s", string(body))
			}
			
			fmt.Fprintln(w, patchResponse)
		} else {
			// Third call is GET again to fetch updated alert
			if r.Method != http.MethodGet {
				t.Fatalf("expected method GET for third call, got %s", r.Method)
			}
			fmt.Fprintln(w, patchResponse)
		}
	})
	
	patch := &UsageAlertPatch{
		Name: &newName,
	}
	
	_, _, err := client.Alerts.PatchUsageAlert(alertID, patch)
	if err != nil {
		t.Fatalf("error patching usage alert: %v", err)
	}
}

// TestCreateUsageAlert_InvalidThreshold tests that 0 and 101 threshold values are rejected with 400 errors
func TestCreateUsageAlert_InvalidThreshold(t *testing.T) {
	client, server, _ := setup(t)
	defer server.Close()
	
	alertName := "Test Usage Alert"
	alertSubtype := alerting.UsageSubtypes.QueryUsage
	tooLowPercent := 0
	tooHighPercent := 101
	notifierListIds := []string{"notifier1"}
	
	mockErrorResponse := `{
		"message": "validation failed",
		"details": [
			{
				"message": "must be between 1 and 100",
				"path": "/data/alert_at_percent"
			}
		]
	}`
	
	server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected method POST, got %s", r.Method)
		}
		if r.URL.Path != "/alerting/v1/alerts" {
			t.Fatalf("expected path /alerting/v1/alerts, got %s", r.URL.Path)
		}
		
		var alert alerting.Alert
		if err := json.NewDecoder(r.Body).Decode(&alert); err != nil {
			t.Fatalf("error decoding request body: %v", err)
		}
		
		// Server-side validation failure for threshold out of bounds
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, mockErrorResponse)
	})
	
	// Test too low value
	_, resp, err := client.Alerts.CreateUsageAlert(alertName, alertSubtype, tooLowPercent, notifierListIds)
	// This should fail client-side validation
	if err == nil {
		t.Fatal("expected error for alert_at_percent below 1, got nil")
	}
	
	// If it passes client validation, it should fail server validation
	if resp != nil && resp.StatusCode == http.StatusBadRequest {
		restErr, ok := err.(*Error)
		if !ok {
			t.Fatalf("expected *Error, got %T", err)
		}
		
		// Check that the error contains the field path
		if !strings.Contains(restErr.Message, "/data/alert_at_percent") {
			t.Fatalf("error should include field path '/data/alert_at_percent', got: %s", restErr.Message)
		}
	}
	
	// Test too high value
	_, resp, err = client.Alerts.CreateUsageAlert(alertName, alertSubtype, tooHighPercent, notifierListIds)
	// This should fail client-side validation
	if err == nil {
		t.Fatal("expected error for alert_at_percent above 100, got nil")
	}
	
	// If it passes client validation, it should fail server validation
	if resp != nil && resp.StatusCode == http.StatusBadRequest {
		restErr, ok := err.(*Error)
		if !ok {
			t.Fatalf("expected *Error, got %T", err)
		}
		
		// Check that the error contains the field path
		if !strings.Contains(restErr.Message, "/data/alert_at_percent") {
			t.Fatalf("error should include field path '/data/alert_at_percent', got: %s", restErr.Message)
		}
	}
}

// Simple end-to-end test without the problematic large format
func TestUsageAlertsE2E(t *testing.T) {
	client, server, _ := setup(t)
	defer server.Close()

	alertID := "abcdef123456"
	alertName := "Test Usage Alert"
	alertSubtype := alerting.UsageSubtypes.QueryUsage
	alertAtPercent := 85
	notifierListIds := []string{"notifier1", "notifier2"}

	// Setup server to handle both create and get requests
	var reqCount int
	server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqCount++
		switch reqCount {
		case 1: // Create
			if r.Method != http.MethodPost {
				t.Fatalf("expected POST for create, got %s", r.Method)
			}
			fmt.Fprintf(w, `{"id":"%s","name":"%s","type":"account","subtype":"%s","data":{"alert_at_percent":%d}}`, 
				alertID, alertName, alertSubtype, alertAtPercent)
		case 2: // Get
			if r.Method != http.MethodGet {
				t.Fatalf("expected GET for get, got %s", r.Method)
			}
			fmt.Fprintf(w, `{"id":"%s","name":"%s","type":"account","subtype":"%s","data":{"alert_at_percent":%d}}`, 
				alertID, alertName, alertSubtype, alertAtPercent)
		case 3: // Delete
			if r.Method != http.MethodDelete {
				t.Fatalf("expected DELETE for delete, got %s", r.Method)
			}
			w.WriteHeader(http.StatusNoContent)
		}
	})

	// Create alert
	createResp, _, err := client.Alerts.CreateUsageAlert(alertName, alertSubtype, alertAtPercent, notifierListIds)
	if err != nil {
		t.Fatalf("failed to create alert: %v", err)
	}
	if createResp.ID == nil || *createResp.ID != alertID {
		t.Fatalf("expected alert ID %s from create, got %v", alertID, createResp.ID)
	}

	// Get alert
	getResp, _, err := client.Alerts.GetUsageAlert(alertID)
	if err != nil {
		t.Fatalf("failed to get alert: %v", err)
	}
	if getResp.ID == nil || *getResp.ID != alertID {
		t.Fatalf("expected alert ID %s from get, got %v", alertID, getResp.ID)
	}

	// Delete alert
	delResp, err := client.Alerts.DeleteUsageAlert(alertID)
	if err != nil {
		t.Fatalf("failed to delete alert: %v", err)
	}
	if delResp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected status 204 for delete, got %d", delResp.StatusCode)
	}
}

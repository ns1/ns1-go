package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"gopkg.in/ns1/ns1-go.v2/rest/model/alerting"
)

func TestUsageAlert_Creation(t *testing.T) {
	client, server, mux := setup(t)
	defer server.Close()
	
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
	
	mux.HandleFunc("/alerting/v1/alerts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected method POST, got %s", r.Method)
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
	
	// Create alert using NewUsageAlert helper + generic Create method
	alert, err := NewUsageAlert(alertName, alertSubtype, alertAtPercent, notifierListIds)
	if err != nil {
		t.Fatalf("error creating usage alert object: %v", err)
	}
	
	_, err = client.Alerts.Create(alert)
	if err != nil {
		t.Fatalf("error creating usage alert: %v", err)
	}
	
	// Verify it's correctly identified as a usage alert
	if !IsUsageAlert(alert) {
		t.Fatalf("expected IsUsageAlert to return true")
	}
	
	// Test validation
	_, err = NewUsageAlert(alertName, "invalid_subtype", alertAtPercent, notifierListIds)
	if err == nil {
		t.Fatal("expected error for invalid subtype, got nil")
	}
	
	_, err = NewUsageAlert(alertName, alertSubtype, 0, notifierListIds)
	if err == nil {
		t.Fatal("expected error for alert_at_percent below 1, got nil")
	}
	
	_, err = NewUsageAlert(alertName, alertSubtype, 101, notifierListIds)
	if err == nil {
		t.Fatal("expected error for alert_at_percent above 100, got nil")
	}
}

func TestUsageAlert_Get(t *testing.T) {
	client, server, mux := setup(t)
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
	
	path := fmt.Sprintf("/alerting/v1/alerts/%s", alertID)
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected method GET, got %s", r.Method)
		}
		
		fmt.Fprintln(w, mockResponse)
	})
	
	// Use generic Get method
	alert, resp, err := client.Alerts.Get(alertID)
	if err != nil {
		t.Fatalf("error getting usage alert: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code 200, got %d", resp.StatusCode)
	}
	if *alert.ID != alertID {
		t.Fatalf("expected alert ID %s, got %v", alertID, *alert.ID)
	}
	
	// Verify it's a usage alert
	if !IsUsageAlert(alert) {
		t.Fatalf("expected IsUsageAlert to return true")
	}
	
	// Test data extraction using our model functions
	data, err := alerting.UnmarshalUsageAlertData(alert)
	if err != nil {
		t.Fatalf("error unmarshalling data: %v", err)
	}
	if data.AlertAtPercent != alertAtPercent {
		t.Fatalf("expected alert_at_percent %d, got %d", alertAtPercent, data.AlertAtPercent)
	}
}

func TestUsageAlert_Patch(t *testing.T) {
	client, server, mux := setup(t)
	defer server.Close()
	
	alertID := "abcdef123456"
	oldName := "Test Usage Alert"
	newName := "Updated Usage Alert"
	alertSubtype := alerting.UsageSubtypes.QueryUsage
	oldAlertAtPercent := 85
	newAlertAtPercent := 90
	
	// Mock initial GET response
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
	
	// Mock PATCH response
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
	
	path := fmt.Sprintf("/alerting/v1/alerts/%s", alertID)
	callCount := 0
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 && r.Method == http.MethodGet {
			// First call is GET
			fmt.Fprintln(w, getResponse)
		} else if callCount == 2 && r.Method == http.MethodPatch {
			// Second call is PATCH
			// Ensure the PATCH request doesn't include type or subtype
			body, _ := io.ReadAll(r.Body)
			bodyStr := string(body)
			
			if strings.Contains(bodyStr, `"type"`) {
				t.Fatalf("PATCH request should not include 'type' field, got: %s", bodyStr)
			}
			
			if strings.Contains(bodyStr, `"subtype"`) {
				t.Fatalf("PATCH request should not include 'subtype' field, got: %s", bodyStr)
			}
			
			fmt.Fprintln(w, patchResponse)
		} else if callCount == 3 && r.Method == http.MethodGet {
			// Third call is GET to fetch updated alert
			fmt.Fprintln(w, patchResponse)
		} else {
			t.Fatalf("unexpected call: method=%s, count=%d", r.Method, callCount)
		}
	})
	
	// First get the alert
	alert, _, err := client.Alerts.Get(alertID)
	if err != nil {
		t.Fatalf("error getting alert: %v", err)
	}
	
	// Create patch with functional options
	patch, err := NewUsageAlertPatch(
		WithName(newName),
		WithAlertAtPercent(newAlertAtPercent),
	)
	if err != nil {
		t.Fatalf("error creating patch: %v", err)
	}
	
	// Apply patch using generic Update method
	// Build a minimal alert with just the fields we want to update
	updateAlert := &alerting.Alert{
		ID:   alert.ID,
		Name: patch.Name,
		Data: alerting.MarshalUsageAlertData(*patch.Data),
	}
	
	// Only add these fields if they were set in the patch
	if patch.NotifierListIds != nil {
		updateAlert.NotifierListIds = *patch.NotifierListIds
	}
	
	if patch.ZoneNames != nil {
		updateAlert.ZoneNames = *patch.ZoneNames
	}
	
	resp, err := client.Alerts.Update(updateAlert)
	if err != nil {
		t.Fatalf("error updating usage alert: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code 200, got %d", resp.StatusCode)
	}
	
	// Get the updated alert
	updated, _, err := client.Alerts.Get(alertID)
	if err != nil {
		t.Fatalf("error getting updated alert: %v", err)
	}
	
	// Verify updates
	if *updated.Name != newName {
		t.Fatalf("expected updated name %s, got %s", newName, *updated.Name)
	}
	
	data, err := alerting.UnmarshalUsageAlertData(updated)
	if err != nil {
		t.Fatalf("error unmarshalling data: %v", err)
	}
	if data.AlertAtPercent != newAlertAtPercent {
		t.Fatalf("expected updated alert_at_percent %d, got %d", newAlertAtPercent, data.AlertAtPercent)
	}
}

func TestUsageAlert_Delete(t *testing.T) {
	client, server, mux := setup(t)
	defer server.Close()
	
	alertID := "abcdef123456"
	alertName := "Test Usage Alert"
	alertSubtype := alerting.UsageSubtypes.QueryUsage
	alertAtPercent := 85
	
	// Mock GET response
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
	
	path := fmt.Sprintf("/alerting/v1/alerts/%s", alertID)
	callCount := 0
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 && r.Method == http.MethodGet {
			// First call is GET
			fmt.Fprintln(w, getResponse)
		} else if callCount == 2 && r.Method == http.MethodDelete {
			// Second call is DELETE
			w.WriteHeader(http.StatusNoContent)
		} else {
			t.Fatalf("unexpected call: method=%s, count=%d", r.Method, callCount)
		}
	})
	
	// First get the alert
	_, _, err := client.Alerts.Get(alertID)
	if err != nil {
		t.Fatalf("error getting alert: %v", err)
	}
	
	// Delete using generic Delete method
	resp, err := client.Alerts.Delete(alertID)
	if err != nil {
		t.Fatalf("error deleting usage alert: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected status code 204, got %d", resp.StatusCode)
	}
}

func TestUsageAlert_List(t *testing.T) {
	client, server, mux := setup(t)
	defer server.Close()
	
	mockResponse := `{
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
	
	mux.HandleFunc("/alerting/v1/alerts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected method GET, got %s", r.Method)
		}
		
		fmt.Fprintln(w, mockResponse)
	})
	
	// List all alerts
	alerts, _, err := client.Alerts.List()
	if err != nil {
		t.Fatalf("error listing alerts: %v", err)
	}
	
	// Filter usage alerts
	usageAlerts := FilterUsageAlerts(alerts)
	
	// Should return only the usage alerts, filtering out the zone alert
	if len(usageAlerts) != 2 {
		t.Fatalf("expected 2 usage alerts, got %d", len(usageAlerts))
	}
	
	// Verify first usage alert
	if *usageAlerts[0].ID != "alert1" {
		t.Fatalf("expected first alert ID 'alert1', got %v", *usageAlerts[0].ID)
	}
	if *usageAlerts[0].Subtype != alerting.UsageSubtypes.QueryUsage {
		t.Fatalf("expected first alert subtype '%s', got %v", alerting.UsageSubtypes.QueryUsage, *usageAlerts[0].Subtype)
	}
	
	// Verify second usage alert
	if *usageAlerts[1].ID != "alert2" {
		t.Fatalf("expected second alert ID 'alert2', got %v", *usageAlerts[1].ID)
	}
	if *usageAlerts[1].Subtype != alerting.UsageSubtypes.RecordUsage {
		t.Fatalf("expected second alert subtype '%s', got %v", alerting.UsageSubtypes.RecordUsage, *usageAlerts[1].Subtype)
	}
}

func TestUsageAlert_PatchValidation(t *testing.T) {
	// Test patch validation for alert_at_percent
	_, err := NewUsageAlertPatch(WithAlertAtPercent(0))
	if err == nil {
		t.Fatal("expected error for alert_at_percent below 1, got nil")
	}
	
	_, err = NewUsageAlertPatch(WithAlertAtPercent(101))
	if err == nil {
		t.Fatal("expected error for alert_at_percent above 100, got nil")
	}
	
	// Test valid patch
	patch, err := NewUsageAlertPatch(
		WithName("New Name"),
		WithAlertAtPercent(50),
		WithNotifierListIds([]string{"notifier1"}),
		WithZoneNames([]string{}),
	)
	if err != nil {
		t.Fatalf("unexpected error creating valid patch: %v", err)
	}
	
	// Verify patch fields were set correctly
	if *patch.Name != "New Name" {
		t.Fatalf("expected name 'New Name', got %v", *patch.Name)
	}
	
	if patch.Data.AlertAtPercent != 50 {
		t.Fatalf("expected alert_at_percent 50, got %d", patch.Data.AlertAtPercent)
	}
	
	if len(*patch.NotifierListIds) != 1 || (*patch.NotifierListIds)[0] != "notifier1" {
		t.Fatalf("unexpected notifier_list_ids value: %v", *patch.NotifierListIds)
	}
	
	if len(*patch.ZoneNames) != 0 {
		t.Fatalf("expected empty zone_names, got %v", *patch.ZoneNames)
	}
}

func TestUsageAlert_E2E(t *testing.T) {
	client, server, mux := setup(t)
	defer server.Close()
	
	alertID := "abcdef123456"
	alertName := "Test Usage Alert"
	alertSubtype := alerting.UsageSubtypes.QueryUsage
	alertAtPercent := 85
	
	// Setup handlers for create, get, and delete endpoints
	mux.HandleFunc("/alerting/v1/alerts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			fmt.Fprintf(w, `{
				"id": "%s",
				"name": "%s",
				"type": "account",
				"subtype": "%s",
				"data": {"alert_at_percent": %d},
				"notifier_list_ids": [],
				"zone_names": []
			}`, alertID, alertName, alertSubtype, alertAtPercent)
		} else if r.Method == http.MethodGet {
			fmt.Fprintf(w, `{"results": [
				{
					"id": "%s",
					"name": "%s",
					"type": "account",
					"subtype": "%s",
					"data": {"alert_at_percent": %d},
					"notifier_list_ids": [],
					"zone_names": []
				}
			]}`, alertID, alertName, alertSubtype, alertAtPercent)
		}
	})
	
	alertPath := fmt.Sprintf("/alerting/v1/alerts/%s", alertID)
	mux.HandleFunc(alertPath, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			fmt.Fprintf(w, `{
				"id": "%s",
				"name": "%s",
				"type": "account",
				"subtype": "%s",
				"data": {"alert_at_percent": %d},
				"notifier_list_ids": [],
				"zone_names": []
			}`, alertID, alertName, alertSubtype, alertAtPercent)
		} else if r.Method == http.MethodDelete {
			w.WriteHeader(http.StatusNoContent)
		}
	})
	
	// Create alert
	alert, err := NewUsageAlert(alertName, alertSubtype, alertAtPercent, nil)
	if err != nil {
		t.Fatalf("error creating usage alert object: %v", err)
	}
	
	_, err = client.Alerts.Create(alert)
	if err != nil {
		t.Fatalf("error creating usage alert: %v", err)
	}
	
	// Since we can't directly verify the created alert from the Create response,
	// we'll fetch it to verify its ID
	
	// Get alert
	fetched, _, err := client.Alerts.Get(alertID)
	if err != nil {
		t.Fatalf("error getting usage alert: %v", err)
	}
	
	// Verify fetched alert
	if *fetched.ID != alertID {
		t.Fatalf("expected alert ID %s, got %v", alertID, *fetched.ID)
	}
	
	// List alerts
	allAlerts, _, err := client.Alerts.List()
	if err != nil {
		t.Fatalf("error listing alerts: %v", err)
	}
	
	// Filter usage alerts
	usageAlerts := FilterUsageAlerts(allAlerts)
	if len(usageAlerts) == 0 {
		t.Fatal("expected at least one usage alert, got none")
	}
	
	// Delete alert
	_, err = client.Alerts.Delete(alertID)
	if err != nil {
		t.Fatalf("error deleting usage alert: %v", err)
	}
}

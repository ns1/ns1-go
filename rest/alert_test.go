package rest_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/ns1/ns1-go.v2/mockns1"
	"gopkg.in/ns1/ns1-go.v2/rest/model/alerting"

	api "gopkg.in/ns1/ns1-go.v2/rest"
)

func strPtr(str string) *string {
	return &str
}

func int64Ptr(n int64) *int64 {
	return &n
}

func TestAlert(t *testing.T) {
	mock, doer, err := mockns1.New(t)

	require.Nil(t, err)
	defer mock.Shutdown()

	client := api.NewClient(doer, api.SetEndpoint("https://"+mock.Address+"/v1/"))

	// Tests for api.Client.View.List()
	t.Run("List", func(t *testing.T) {
		t.Run("List", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddAlertListTestCase("", nil, nil, alertList))

			respAlerts, _, err := client.Alerts.List()
			require.Nil(t, err)
			require.NotNil(t, respAlerts)
			compareAlertLists(t, alertList, respAlerts)
		})

		t.Run("List with pagination", func(t *testing.T) {
			defer mock.ClearTestCases()

			linkHeader := http.Header{}
			linkHeader.Set("Link", `</alerting/v1/alerts?next=`+*alertList[1].Name+`>; rel="next"`)
			require.Nil(t, mock.AddAlertListTestCase("", nil, linkHeader, alertList[0:2]))
			require.Nil(t, mock.AddAlertListTestCase("next="+*alertList[1].Name, nil, nil, alertList[2:5]))

			respAlerts, _, err := client.Alerts.List()
			require.Nil(t, err)
			compareAlertLists(t, alertList, respAlerts)
		})

		t.Run("List without pagination", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddAlertListTestCase("", nil, nil, alertList))

			respAlerts, _, err := client.Alerts.List()
			require.Nil(t, err)
			require.NotNil(t, respAlerts)

			compareAlertLists(t, alertList, respAlerts)
		})
	})

	t.Run("Get", func(t *testing.T) {

		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			client.FollowPagination = true
			alert := alertList[0]
			require.Nil(t, mock.AddAlertGetTestCase(*alert.ID, nil, nil, alert))

			respAlert, _, err := client.Alerts.Get(*alert.ID)
			require.Nil(t, err)
			require.NotNil(t, respAlert)
			compareAlerts(t, alert, respAlert)
		})
		t.Run("Not Found", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddAlertFailTestCase(
				http.MethodGet, "abcd-efgh-ijkl", http.StatusNotFound,
				nil, nil, `{"message": "test error"}`,
			))

			respAlert, resp, err := client.Alerts.Get("abcd-efgh-ijkl")
			require.Nil(t, respAlert)
			require.NotNil(t, err)
			require.Contains(t, err.Error(), "test error")
			require.Equal(t, http.StatusNotFound, resp.StatusCode)
		})

		t.Run("Other Error", func(t *testing.T) {
			c := api.NewClient(errorClient{}, api.SetEndpoint(""))
			respAlert, resp, err := c.Alerts.Get("abcd-efgh-ijkl")
			require.Nil(t, resp)
			require.Error(t, err)
			require.Nil(t, respAlert)
		})
	})

	t.Run("Create", func(t *testing.T) {

		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			alertToCreate := alerting.Alert{
				ID:   strPtr("f136f0cd-eca8-4b2b-9b2d-9a6631200d51"),
				Data: json.RawMessage(nil),
				Name: strPtr("first_alert"),
				NotifierListIds: []string{
					"66d07ca6e113eb00014fe257",
					"66d07caf8519c000011cdda6",
					"6707da567cd4f300012cd7e4",
				},
				RecordIds: []string{},
				Type:      strPtr("zone"),
				Subtype:   strPtr("transfer_failed"),
				ZoneNames: []string{
					"alerttest.com", "alerttest.net",
				},
			}

			alertResponse := alertList[0]

			// Pass by value here. alertToCreate will be modified by client.Alerts.Create
			require.Nil(t, mock.AddAlertCreateTestCase(nil, nil, alertToCreate, *alertResponse))

			_, err := client.Alerts.Create(&alertToCreate)
			require.Nil(t, err)
			// alerttoCreate will be modified by the call to Create and fully populated.
			compareAlerts(t, alertResponse, &alertToCreate)
		})

		t.Run("Error Duplicate Name", func(t *testing.T) {
			defer mock.ClearTestCases()

			alertToCreate := alerting.Alert{
				ID:   strPtr("f136f0cd-eca8-4b2b-9b2d-9a6631200d51"),
				Data: json.RawMessage(nil),
				Name: strPtr("duplicate_alert"),
				NotifierListIds: []string{
					"66d07ca6e113eb00014fe257",
					"66d07caf8519c000011cdda6",
				},
				RecordIds: []string{},
				Type:      strPtr("zone"),
				Subtype:   strPtr("transfer_failed"),
				ZoneNames: []string{
					"alerttest.com", "alerttest.net",
				},
			}

			require.Nil(t, mock.AddAlertFailTestCaseWithReqBody(
				http.MethodPost, "", http.StatusConflict,
				nil, nil, alertToCreate, `{"message": "test error"}`,
			))

			resp, err := client.Alerts.Create(&alertToCreate)
			require.NotNil(t, err)
			require.Equal(t, api.ErrAlertExists, err)
			require.Equal(t, http.StatusConflict, resp.StatusCode)
		})

		t.Run("Error Bad Alert Object", func(t *testing.T) {
			defer mock.ClearTestCases()

			alertToCreate := alerting.Alert{
				ID:   strPtr("f136f0cd-eca8-4b2b-9b2d-9a6631200d51"),
				Data: json.RawMessage(nil),
				Name: strPtr("duplicate_alert"),
				NotifierListIds: []string{
					"66d07ca6e113eb00014fe257",
					"66d07caf8519c000011cdda6",
				},
				RecordIds: []string{},
				Type:      strPtr("fakeType"), // Bad alert type
				Subtype:   strPtr("transfer_failed"),
				ZoneNames: []string{
					"alerttest.com", "alerttest.net",
				},
			}

			require.Nil(t, mock.AddAlertFailTestCaseWithReqBody(
				http.MethodPost, "", http.StatusBadRequest,
				nil, nil, alertToCreate, `{"message": "test error"}`,
			))

			resp, err := client.Alerts.Create(&alertToCreate)
			require.NotNil(t, err)
			require.Contains(t, err.Error(), "test error")
			require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("Success SSO Alert", func(t *testing.T) {
			defer mock.ClearTestCases()

			notifierList := []string{
				"66d07ca6e113eb00014fe257",
				"66d07caf8519c000011cdda6",
				"6707da567cd4f300012cd7e4",
			}
			alertToCreate := alerting.NewSSOAlert("sso_alert", notifierList)
			alertResponse := alertList[0]

			// Pass by value here. alertToCreate will be modified by client.Alerts.Create
			require.Nil(t, mock.AddAlertCreateTestCase(nil, nil, *alertToCreate, *alertResponse))

			_, err := client.Alerts.Create(alertToCreate)
			require.Nil(t, err)
			// alerttoCreate will be modified by the call to Create and fully populated.
			compareAlerts(t, alertResponse, alertToCreate)
		})

		t.Run("Success Redirect Alert", func(t *testing.T) {
			defer mock.ClearTestCases()

			notifierList := []string{
				"66d07ca6e113eb00014fe257",
				"66d07caf8519c000011cdda6",
				"6707da567cd4f300012cd7e4",
			}

			alertToCreate := alerting.NewRedirectAlert("redirect_alert", notifierList)
			alertResponse := alertList[0]

			// Pass by value here. alertToCreate will be modified by client.Alerts.Create
			require.Nil(t, mock.AddAlertCreateTestCase(nil, nil, *alertToCreate, *alertResponse))

			_, err := client.Alerts.Create(alertToCreate)
			require.Nil(t, err)
			// alerttoCreate will be modified by the call to Create and fully populated.
			compareAlerts(t, alertResponse, alertToCreate)
		})
	})

	t.Run("Update", func(t *testing.T) {

		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			alertUpdate := alerting.Alert{
				ID:   strPtr("f136f0cd-eca8-4b2b-9b2d-9a6631200d51"),
				Name: strPtr("renamed_alert"),
			}

			alertResponse := alerting.Alert{
				ID:        strPtr("f136f0cd-eca8-4b2b-9b2d-9a6631200d51"),
				CreatedAt: int64Ptr(1728637379),
				CreatedBy: strPtr("testapikey"),
				Data:      json.RawMessage(nil),
				Name:      strPtr("renamed_alert"),
				NotifierListIds: []string{
					"66d07ca6e113eb00014fe257",
					"66d07caf8519c000011cdda6",
					"6707da567cd4f300012cd7e4",
				},
				RecordIds: []string{},
				Type:      strPtr("zone"),
				Subtype:   strPtr("transfer_failed"),
				UpdatedAt: int64Ptr(1728637379),
				UpdatedBy: strPtr("testapikey"),
				ZoneNames: []string{
					"alerttest.com", "alerttest.net",
				},
			}

			// Pass by value here. alertUpdate will be modified by client.Alerts.Update
			require.Nil(t, mock.AddAlertUpdateTestCase(nil, nil, alertUpdate, alertResponse))

			_, err := client.Alerts.Update(&alertUpdate)
			require.Nil(t, err)
			// alertUpdate will be modified by the call to Update and fully populated.
			compareAlerts(t, &alertResponse, &alertUpdate)
		})

		t.Run("Error", func(t *testing.T) {
			defer mock.ClearTestCases()

			alertUpdate := alerting.Alert{
				ID:   strPtr("abcd-efgh-ijkl"),
				Name: strPtr("renamed_alert"),
			}

			require.Nil(t, mock.AddAlertFailTestCaseWithReqBody(
				http.MethodPatch, "abcd-efgh-ijkl", http.StatusNotFound,
				nil, nil, alertUpdate, `{"message": "test error"}`,
			))

			resp, err := client.Alerts.Update(&alertUpdate)
			require.NotNil(t, err)
			require.Contains(t, err.Error(), "test error")
			require.Equal(t, http.StatusNotFound, resp.StatusCode)
		})
	})

	t.Run("Replace", func(t *testing.T) {

		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			// With replace, the whole alert is passed with the update.
			updatedAlert := alerting.Alert{
				ID:        strPtr("f136f0cd-eca8-4b2b-9b2d-9a6631200d51"),
				CreatedAt: int64Ptr(1728637379),
				CreatedBy: strPtr("testapikey"),
				Data:      json.RawMessage(nil),
				Name:      strPtr("renamed_alert"), // This would be the change.
				NotifierListIds: []string{
					"66d07ca6e113eb00014fe257",
					"66d07caf8519c000011cdda6",
					"6707da567cd4f300012cd7e4",
				},
				RecordIds: []string{},
				Type:      strPtr("zone"),
				Subtype:   strPtr("transfer_failed"),
				UpdatedAt: int64Ptr(1728637379),
				UpdatedBy: strPtr("testapikey"),
				ZoneNames: []string{
					"alerttest.com", "alerttest.net",
				},
			}

			alertResponse := alerting.Alert{
				ID:        strPtr("f136f0cd-eca8-4b2b-9b2d-9a6631200d51"),
				CreatedAt: int64Ptr(1728637379),
				CreatedBy: strPtr("testapikey"),
				Data:      json.RawMessage(nil),
				Name:      strPtr("renamed_alert"),
				NotifierListIds: []string{
					"66d07ca6e113eb00014fe257",
					"66d07caf8519c000011cdda6",
					"6707da567cd4f300012cd7e4",
				},
				RecordIds: []string{},
				Type:      strPtr("zone"),
				Subtype:   strPtr("transfer_failed"),
				UpdatedAt: int64Ptr(1728640000),
				UpdatedBy: strPtr("anotherapikey"),
				ZoneNames: []string{
					"alerttest.com", "alerttest.net",
				},
			}

			// Pass by value here. updatedAlert will be modified by client.Alerts.Replace
			require.Nil(t, mock.AddAlertReplaceTestCase(nil, nil, updatedAlert, alertResponse))

			_, err := client.Alerts.Replace(&updatedAlert)
			require.Nil(t, err)
			// updatedAlert will be modified by the call to Update and fully populated.
			compareAlerts(t, &alertResponse, &updatedAlert)
		})

		t.Run("Error", func(t *testing.T) {
			defer mock.ClearTestCases()

			alertUpdate := alerting.Alert{
				ID:   strPtr("abcd-efgh-ijkl"),
				Name: strPtr("renamed_alert"),
			}

			require.Nil(t, mock.AddAlertFailTestCaseWithReqBody(
				http.MethodPut, "abcd-efgh-ijkl", http.StatusNotFound,
				nil, nil, alertUpdate, `{"message": "test error"}`,
			))

			resp, err := client.Alerts.Replace(&alertUpdate)
			require.NotNil(t, err)
			require.Contains(t, err.Error(), "test error")
			require.Equal(t, http.StatusNotFound, resp.StatusCode)
		})

		t.Run("Error Immutable Field", func(t *testing.T) {
			defer mock.ClearTestCases()

			alertUpdate := alerting.Alert{
				ID:        strPtr("abcd-efgh-ijkl"),
				CreatedBy: strPtr("testUser"),
			}

			require.Nil(t, mock.AddAlertFailTestCaseWithReqBody(
				http.MethodPut, "abcd-efgh-ijkl", http.StatusConflict,
				nil, nil, alertUpdate, `{"message": "test error"}`,
			))

			resp, err := client.Alerts.Replace(&alertUpdate)
			require.NotNil(t, err)
			require.Contains(t, err.Error(), "test error")
			require.Equal(t, http.StatusConflict, resp.StatusCode)
		})
	})

	t.Run("Delete", func(t *testing.T) {

		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			client.FollowPagination = true
			alert := alertList[0]
			require.Nil(t, mock.AddAlertDeleteTestCase(*alert.ID, nil, nil))

			resp, err := client.Alerts.Delete(*alert.ID)
			require.Nil(t, err)
			require.Equal(t, http.StatusNoContent, resp.StatusCode)
		})
		t.Run("Not Found", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddAlertFailTestCase(
				http.MethodDelete, "abcd-efgh-ijkl", http.StatusNotFound,
				nil, nil, `{"message": "test error"}`,
			))

			resp, err := client.Alerts.Delete("abcd-efgh-ijkl")
			require.NotNil(t, err)
			require.Contains(t, err.Error(), "test error")
			require.Equal(t, http.StatusNotFound, resp.StatusCode)
		})

		t.Run("Other Error", func(t *testing.T) {
			c := api.NewClient(errorClient{}, api.SetEndpoint(""))
			resp, err := c.Alerts.Delete("abcd-efgh-ijkl")
			require.Nil(t, resp)
			require.Error(t, err)
		})
	})

	t.Run("Test", func(t *testing.T) {

		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			client.FollowPagination = true
			alert := alertList[0]
			require.Nil(t, mock.AddAlertTestPostTestCase(*alert.ID, nil, nil))

			resp, err := client.Alerts.Test(*alert.ID)
			require.Nil(t, err)
			require.Equal(t, http.StatusNoContent, resp.StatusCode)
		})
		t.Run("Not Found", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddAlertFailTestCase(
				http.MethodPost, "abcd-efgh-ijkl/test", http.StatusNotFound,
				nil, nil, `{"message": "test error"}`,
			))

			resp, err := client.Alerts.Test("abcd-efgh-ijkl")
			require.NotNil(t, err)
			require.Contains(t, err.Error(), "test error")
			require.Equal(t, http.StatusNotFound, resp.StatusCode)
		})

		t.Run("Other Error", func(t *testing.T) {
			c := api.NewClient(errorClient{}, api.SetEndpoint(""))
			resp, err := c.Alerts.Test("abcd-efgh-ijkl")
			require.Nil(t, resp)
			require.Error(t, err)
		})
	})
}

func compareAlertLists(t *testing.T, expected []*alerting.Alert, actual []*alerting.Alert) {
	require.Equal(t, len(expected), len(actual))
	for i := range expected {
		compareAlerts(t, expected[i], actual[i])
	}
}
func compareAlerts(t *testing.T, expected *alerting.Alert, actual *alerting.Alert) {
	require.Equal(t, expected.ID, actual.ID)
	require.Equal(t, expected.CreatedAt, actual.CreatedAt)
	require.Equal(t, expected.CreatedBy, actual.CreatedBy)
	if expected.Data != nil && len(expected.Data) == 0 {
		require.Equal(t, json.RawMessage(nil), actual.Data)
	} else {
		require.Equal(t, expected.Data, actual.Data)
	}
	require.Equal(t, expected.Name, actual.Name)
	require.Equal(t, len(expected.NotifierListIds), len(actual.NotifierListIds))
	for i := range expected.NotifierListIds {
		require.Equal(t, expected.NotifierListIds[i], actual.NotifierListIds[i])
	}
	require.Equal(t, len(expected.RecordIds), len(actual.RecordIds))
	for i := range expected.RecordIds {
		require.Equal(t, expected.RecordIds[i], actual.RecordIds[i])
	}
	require.Equal(t, expected.Type, actual.Type)
	require.Equal(t, expected.Subtype, actual.Subtype)
	require.Equal(t, expected.UpdatedAt, actual.UpdatedAt)
	require.Equal(t, expected.UpdatedBy, actual.UpdatedBy)
	require.Equal(t, len(expected.ZoneNames), len(actual.ZoneNames))
	for i := range expected.ZoneNames {
		require.Equal(t, expected.ZoneNames[i], actual.ZoneNames[i])
	}

}

var alertList = []*alerting.Alert{
	{
		ID:        strPtr("f136f0cd-eca8-4b2b-9b2d-9a6631200d51"),
		CreatedAt: int64Ptr(1728637379),
		CreatedBy: strPtr("testapikey"),
		Data:      json.RawMessage(nil),
		Name:      strPtr("first_alert"),
		NotifierListIds: []string{
			"66d07ca6e113eb00014fe257",
			"66d07caf8519c000011cdda6",
			"6707da567cd4f300012cd7e4",
		},
		RecordIds: []string{},
		Type:      strPtr("zone"),
		Subtype:   strPtr("transfer_failed"),
		UpdatedAt: int64Ptr(1728637379),
		UpdatedBy: strPtr("testapikey"),
		ZoneNames: []string{
			"alerttest.com", "alerttest.net",
		},
	},
	{
		ID:        strPtr("3a81d9fa-6f03-4baf-83e4-be3f16411c4f"),
		CreatedAt: int64Ptr(1728637233),
		CreatedBy: strPtr("testapikey"),
		Data:      json.RawMessage(nil),
		Name:      strPtr("second_alert"),
		NotifierListIds: []string{
			"66d07ca6e113eb00014fe242",
			"66d07caf8519c000011cddb7",
			"6707da567cd4f300012cd7f9",
		},
		RecordIds: []string{},
		Type:      strPtr("zone"),
		Subtype:   strPtr("external_primary_failed"),
		UpdatedAt: int64Ptr(1728637233),
		UpdatedBy: strPtr("testapikey"),
		ZoneNames: []string{
			"alerttest2.com",
		},
	},
	{
		ID:        strPtr("3a81d9fa-6f03-4baf-83e4-be3f16411c4f"),
		CreatedAt: int64Ptr(1728637233),
		CreatedBy: strPtr("testapikey"),
		Data:      json.RawMessage(nil),
		Name:      strPtr("third_alert"),
		NotifierListIds: []string{
			"66d07ca6e113eb00014fe242",
			"6707da567cd4f300012cd7f9",
		},
		RecordIds: []string{},
		Type:      strPtr("zone"),
		Subtype:   strPtr("external_primary_failed"),
		UpdatedAt: int64Ptr(1728637233),
		UpdatedBy: strPtr("testapikey"),
		ZoneNames: []string{
			"alerttest1.com", "alerttest2.com",
		},
	},
	{
		ID:        strPtr("3a81d9fa-6f03-4baf-83e4-be3f16411c4f"),
		CreatedAt: int64Ptr(1728637833),
		CreatedBy: strPtr("testapikey"),
		Data:      json.RawMessage(nil),
		Name:      strPtr("fourth_alert"),
		NotifierListIds: []string{
			"66d07ca6e113eb00014fe257",
			"66d07caf8519c000011cdda6",
			"6707da567cd4f300012cd7e4",
			"66d07ca6e113eb00014fe242",
			"66d07caf8519c000011cddb7",
			"6707da567cd4f300012cd7f9",
		},
		RecordIds: []string{},
		Type:      strPtr("zone"),
		Subtype:   strPtr("transfer_failed"),
		UpdatedAt: int64Ptr(1728639233),
		UpdatedBy: strPtr("testapikey"),
		ZoneNames: []string{
			"alerttest2.com", "alerttest3.com",
		},
	},
	{
		ID:              strPtr("e2e64e2b-575b-4669-8a9b-72512e8bdc6f"),
		CreatedAt:       int64Ptr(1728637836),
		CreatedBy:       strPtr("testapikey"),
		Data:            json.RawMessage(nil),
		Name:            strPtr("fifth_alert - empty lists"),
		NotifierListIds: []string{},
		RecordIds:       []string{},
		Type:            strPtr("zone"),
		Subtype:         strPtr("transfer_failed"),
		UpdatedAt:       int64Ptr(1728639236),
		UpdatedBy:       strPtr("testapikey"),
		ZoneNames:       []string{},
	},
}

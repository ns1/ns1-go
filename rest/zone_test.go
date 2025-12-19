package rest_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/ns1/ns1-go.v2/mockns1"
	api "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/dns"
)

func TestZone(t *testing.T) {
	mock, doer, err := mockns1.New(t)
	require.Nil(t, err)
	defer mock.Shutdown()

	client := api.NewClient(doer, api.SetEndpoint("https://"+mock.Address+"/v1/"))

	t.Run("List", func(t *testing.T) {
		t.Run("Pagination", func(t *testing.T) {
			defer mock.ClearTestCases()

			client.FollowPagination = true
			zones := []*dns.Zone{
				{Zone: "a.list.zone"},
				{Zone: "b.list.zone"},
				{Zone: "c.list.zone"},
				{Zone: "d.list.zone"},
			}
			require.Nil(t, mock.AddZoneListTestCase(nil, nil, zones))

			respZones, _, err := client.Zones.List()
			require.Nil(t, err)
			require.NotNil(t, respZones)
			require.Equal(t, len(zones), len(respZones))

			for i := range zones {
				require.Equal(t, zones[i].Zone, respZones[i].Zone, i)
			}
		})

		t.Run("No Pagination", func(t *testing.T) {
			defer mock.ClearTestCases()

			client.FollowPagination = false
			zones := []*dns.Zone{
				{Zone: "a.list.zone"},
				{Zone: "b.list.zone"},
			}

			header := http.Header{}
			header.Set("Link", `</zones?b.list.zone&limit=2>; rel="next"`)

			require.Nil(t, mock.AddZoneListTestCase(nil, header, zones))

			respZones, resp, err := client.Zones.List()
			require.Nil(t, err)
			require.NotNil(t, respZones)
			require.Equal(t, len(zones), len(respZones))
			require.Contains(t, resp.Header.Get("Link"), "zones?b.list.zone")

			for i := range zones {
				require.Equal(t, zones[i].Zone, respZones[i].Zone, i)
			}
		})

		t.Run("Error", func(t *testing.T) {
			t.Run("HTTP", func(t *testing.T) {
				defer mock.ClearTestCases()

				require.Nil(t, mock.AddTestCase(
					http.MethodGet, "/zones", http.StatusNotFound,
					nil, nil, "", `{"message": "test error"}`,
				))

				zones, resp, err := client.Zones.List()
				require.Nil(t, zones)
				require.NotNil(t, err)
				require.Contains(t, err.Error(), "test error")
				require.Equal(t, http.StatusNotFound, resp.StatusCode)
			})

			t.Run("Other", func(t *testing.T) {
				c := api.NewClient(errorClient{}, api.SetEndpoint(""))
				zones, resp, err := c.Zones.List()
				require.Nil(t, resp)
				require.Error(t, err)
				require.Nil(t, zones)
			})
		})
	})

	t.Run("Get", func(t *testing.T) {
		zoneName := "a.get.zone"

		t.Run("Pagination", func(t *testing.T) {
			defer mock.ClearTestCases()

			client.FollowPagination = true
			zone := &dns.Zone{
				Zone: "a.get.zone",
				Records: []*dns.ZoneRecord{
					{Domain: "1.a.get.zone"},
					{Domain: "2.a.get.zone"},
					{Domain: "3.a.get.zone"},
					{Domain: "4.a.get.zone"},
				},
			}
			require.Nil(t, mock.AddZoneGetTestCase(zoneName, nil, nil, zone, true))

			respZone, _, err := client.Zones.Get(zoneName, true)
			require.Nil(t, err)
			require.NotNil(t, respZone)
			require.Equal(t, len(zone.Records), len(respZone.Records))

			for i := range zone.Records {
				require.Equal(t, zone.Records[i].Domain, respZone.Records[i].Domain, i)
			}
		})

		t.Run("No Pagination", func(t *testing.T) {
			defer mock.ClearTestCases()

			client.FollowPagination = false
			zone := &dns.Zone{
				Zone: "a.get.zone",
				Records: []*dns.ZoneRecord{
					{Domain: "1.a.get.zone"},
					{Domain: "2.a.get.zone"},
				},
			}

			link := `zones/` + zoneName + `?after=3.` + zoneName
			header := http.Header{}
			header.Set("Link", `</`+link+`&limit=2>; rel="next"`)

			require.Nil(t, mock.AddZoneGetTestCase(zoneName, nil, header, zone, true))

			respZone, resp, err := client.Zones.Get(zoneName, true)
			require.Nil(t, err)
			require.NotNil(t, respZone)
			require.Equal(t, len(zone.Records), len(respZone.Records))
			require.Contains(t, resp.Header.Get("Link"), link)

			for i := range zone.Records {
				require.Equal(t, zone.Records[i].Domain, respZone.Records[i].Domain, i)
			}
		})

		t.Run("Error", func(t *testing.T) {
			t.Run("HTTP", func(t *testing.T) {
				defer mock.ClearTestCases()

				require.Nil(t, mock.AddTestCase(
					http.MethodGet, "/zones/"+zoneName, http.StatusNotFound,
					nil, nil, "", `{"message": "test error"}`,
				))

				zone, resp, err := client.Zones.Get(zoneName, true)
				require.Nil(t, zone)
				require.NotNil(t, err)
				require.Contains(t, err.Error(), "test error")
				require.Equal(t, http.StatusNotFound, resp.StatusCode)
			})

			t.Run("Other", func(t *testing.T) {
				c := api.NewClient(errorClient{}, api.SetEndpoint(""))
				zones, resp, err := c.Zones.Get(zoneName, true)
				require.Nil(t, resp)
				require.Error(t, err)
				require.Nil(t, zones)
			})
		})

		t.Run("No Records", func(t *testing.T) {
			defer mock.ClearTestCases()

			client.FollowPagination = true
			zone := &dns.Zone{
				Zone:    "a.get.zone",
				Records: []*dns.ZoneRecord{},
			}
			require.Nil(t, mock.AddZoneGetTestCase(zoneName, nil, nil, zone, false))

			respZone, _, err := client.Zones.Get(zoneName, false)
			require.Nil(t, err)
			require.NotNil(t, respZone)
			require.Equal(t, 0, len(respZone.Records))

		})

	})

	t.Run("Create", func(t *testing.T) {
		zone := &dns.Zone{
			Zone: "create.zone",
			TTL:  42,
		}

		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddZoneCreateTestCase(nil, nil, zone, zone))

			_, err := client.Zones.Create(zone)
			require.Nil(t, err)
		})

		t.Run("Error - zone already exists", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddTestCase(
				http.MethodPut, "/zones/create.zone", http.StatusNotFound,
				nil, nil, zone, `{"message": "zone already exists"}`,
			))

			_, err := client.Zones.Create(zone)
			require.Equal(t, api.ErrZoneExists, err)
		})

		t.Run("Error - invalid: FQDN already exists", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddTestCase(
				http.MethodPut, "/zones/create.zone", http.StatusNotFound,
				nil, nil, zone, `{"message": "invalid: FQDN already exists"}`,
			))

			_, err := client.Zones.Create(zone)
			require.Equal(t, api.ErrZoneExists, err)
		})

		t.Run("Error - invalid: FQDN already exists in the view", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddTestCase(
				http.MethodPut, "/zones/create.zone", http.StatusNotFound,
				nil, nil, zone, `{"message": "invalid: FQDN already exists in the view"}`,
			))

			_, err := client.Zones.Create(zone)
			require.Equal(t, api.ErrZoneExists, err)
		})
	})

	t.Run("Update", func(t *testing.T) {
		zone := &dns.Zone{
			Zone: "update.zone",
			TTL:  42,
		}

		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddZoneUpdateTestCase(nil, nil, zone, zone))

			_, err := client.Zones.Update(zone)
			require.Nil(t, err)
		})

		t.Run("Error", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddTestCase(
				http.MethodPost, "/zones/update.zone", http.StatusNotFound,
				nil, nil, zone, `{"message": "zone not found"}`,
			))

			_, err := client.Zones.Update(zone)
			require.Equal(t, api.ErrZoneMissing, err)
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddZoneDeleteTestCase("delete.zone", nil, nil))

			_, err := client.Zones.Delete("delete.zone")
			require.Nil(t, err)
		})

		t.Run("Error", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddTestCase(
				http.MethodDelete, "/zones/delete.zone", http.StatusNotFound,
				nil, nil, "", `{"message": "zone not found"}`,
			))

			_, err := client.Zones.Delete("delete.zone")
			require.Equal(t, api.ErrZoneMissing, err)
		})

		t.Run("ExportZonefile", func(t *testing.T) {
			zoneName := "export.zone"

			t.Run("Success", func(t *testing.T) {
				defer mock.ClearTestCases()

				expectedStatus := &dns.ZoneFileExportStatus{
					Status:      "GENERATING",
					Message:     "export has been started",
					GeneratedAt: "2025-12-08T18:15:00Z",
				}

				require.Nil(t, mock.AddTestCase(
					http.MethodPut, "/zones/"+zoneName+"/export/zonefile", http.StatusOK,
					nil, nil, map[string]interface{}{}, expectedStatus,
				))

				status, resp, err := client.Zones.ExportZonefile(zoneName)
				require.Nil(t, err)
				require.NotNil(t, status)
				require.NotNil(t, resp)
				require.Equal(t, expectedStatus.Status, status.Status)
				require.Equal(t, expectedStatus.Message, status.Message)
				require.Equal(t, expectedStatus.GeneratedAt, status.GeneratedAt)
			})

			t.Run("Error - zone not found", func(t *testing.T) {
				defer mock.ClearTestCases()

				require.Nil(t, mock.AddTestCase(
					http.MethodPut, "/zones/"+zoneName+"/export/zonefile", http.StatusNotFound,
					nil, nil, map[string]interface{}{}, `{"message": "zone not found"}`,
				))

				status, resp, err := client.Zones.ExportZonefile(zoneName)
				require.Nil(t, status)
				require.NotNil(t, resp)
				require.Equal(t, api.ErrZoneMissing, err)
			})

			t.Run("Error - HTTP", func(t *testing.T) {
				defer mock.ClearTestCases()

				require.Nil(t, mock.AddTestCase(
					http.MethodPut, "/zones/"+zoneName+"/export/zonefile", http.StatusInternalServerError,
					nil, nil, map[string]interface{}{}, `{"message": "internal server error"}`,
				))

				status, resp, err := client.Zones.ExportZonefile(zoneName)
				require.Nil(t, status)
				require.NotNil(t, err)
				require.Contains(t, err.Error(), "internal server error")
				require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
			})

			t.Run("Error - Other", func(t *testing.T) {
				c := api.NewClient(errorClient{}, api.SetEndpoint(""))
				status, resp, err := c.Zones.ExportZonefile(zoneName)
				require.Nil(t, resp)
				require.Error(t, err)
				require.Nil(t, status)
			})
		})

		t.Run("GetExportZonefileStatus", func(t *testing.T) {
			zoneName := "export.zone"

			t.Run("Success - Completed", func(t *testing.T) {
				defer mock.ClearTestCases()

				expectedStatus := &dns.ZoneFileExportStatus{
					Status:      "completed",
					GeneratedAt: "2012-04-23T18:25:43.511Z",
				}

				require.Nil(t, mock.AddTestCase(
					http.MethodGet, "/export/zonefile/"+zoneName+"/status", http.StatusOK,
					nil, nil, "", expectedStatus,
				))

				status, resp, err := client.Zones.GetExportZonefileStatus(zoneName)
				require.Nil(t, err)
				require.NotNil(t, status)
				require.NotNil(t, resp)
				require.Equal(t, expectedStatus.Status, status.Status)
				require.Equal(t, expectedStatus.GeneratedAt, status.GeneratedAt)
			})

			t.Run("Success - Generating", func(t *testing.T) {
				defer mock.ClearTestCases()

				expectedStatus := &dns.ZoneFileExportStatus{
					Status:  "generating",
					Message: "Export is in progress.",
				}

				require.Nil(t, mock.AddTestCase(
					http.MethodGet, "/export/zonefile/"+zoneName+"/status", http.StatusOK,
					nil, nil, "", expectedStatus,
				))

				status, resp, err := client.Zones.GetExportZonefileStatus(zoneName)
				require.Nil(t, err)
				require.NotNil(t, status)
				require.NotNil(t, resp)
				require.Equal(t, expectedStatus.Status, status.Status)
				require.Equal(t, expectedStatus.Message, status.Message)
			})

			t.Run("Success - Failed", func(t *testing.T) {
				defer mock.ClearTestCases()

				expectedStatus := &dns.ZoneFileExportStatus{
					Status:  "failed",
					Message: "An internal error occurred during export.",
				}

				require.Nil(t, mock.AddTestCase(
					http.MethodGet, "/export/zonefile/"+zoneName+"/status", http.StatusOK,
					nil, nil, "", expectedStatus,
				))

				status, resp, err := client.Zones.GetExportZonefileStatus(zoneName)
				require.Nil(t, err)
				require.NotNil(t, status)
				require.NotNil(t, resp)
				require.Equal(t, expectedStatus.Status, status.Status)
				require.Equal(t, expectedStatus.Message, status.Message)
			})

			t.Run("Error - No export found", func(t *testing.T) {
				defer mock.ClearTestCases()

				require.Nil(t, mock.AddTestCase(
					http.MethodGet, "/export/zonefile/"+zoneName+"/status", http.StatusNotFound,
					nil, nil, "", `{"message": "No export found for the specified zone."}`,
				))

				status, resp, err := client.Zones.GetExportZonefileStatus(zoneName)
				require.Nil(t, status)
				require.NotNil(t, resp)
				require.Error(t, err)
			})

			t.Run("Error - Zone not found", func(t *testing.T) {
				defer mock.ClearTestCases()

				require.Nil(t, mock.AddTestCase(
					http.MethodGet, "/export/zonefile/"+zoneName+"/status", http.StatusNotFound,
					nil, nil, "", `{"message": "zone not found"}`,
				))

				status, resp, err := client.Zones.GetExportZonefileStatus(zoneName)
				require.Nil(t, status)
				require.NotNil(t, resp)
				require.Equal(t, api.ErrZoneMissing, err)
			})

			t.Run("Error - HTTP", func(t *testing.T) {
				defer mock.ClearTestCases()

				require.Nil(t, mock.AddTestCase(
					http.MethodGet, "/export/zonefile/"+zoneName+"/status", http.StatusInternalServerError,
					nil, nil, "", `{"message": "internal server error"}`,
				))

				status, resp, err := client.Zones.GetExportZonefileStatus(zoneName)
				require.Nil(t, status)
				require.NotNil(t, err)
				require.Contains(t, err.Error(), "internal server error")
				require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
			})

			t.Run("Error - Other", func(t *testing.T) {
				c := api.NewClient(errorClient{}, api.SetEndpoint(""))
				status, resp, err := c.Zones.GetExportZonefileStatus(zoneName)
				require.Nil(t, resp)
				require.Error(t, err)
				require.Nil(t, status)
			})
		})
	})

	t.Run("DownloadZonefile", func(t *testing.T) {
		zoneName := "export.zone"
		zonefileContent := "; Zone file for export.zone\n$ORIGIN export.zone.\n@ IN SOA ns1.export.zone. admin.export.zone. 1 3600 600 604800 86400\n"

		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddTestCase(
				http.MethodGet, "/export/zonefile/"+zoneName, http.StatusOK,
				nil, nil, "", zonefileContent,
			))

			buf, resp, err := client.Zones.DownloadZonefile(zoneName)
			require.Nil(t, err)
			require.NotNil(t, buf)
			require.NotNil(t, resp)
			require.Equal(t, zonefileContent, buf.String())
		})

		t.Run("Error - No export found", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddTestCase(
				http.MethodGet, "/export/zonefile/"+zoneName, http.StatusNotFound,
				nil, nil, "", `{"message": "No export found."}`,
			))

			buf, resp, err := client.Zones.DownloadZonefile(zoneName)
			require.Nil(t, buf)
			require.NotNil(t, resp)
			require.Error(t, err)
		})

		t.Run("Error - Zone not found", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddTestCase(
				http.MethodGet, "/export/zonefile/"+zoneName, http.StatusNotFound,
				nil, nil, "", `{"message": "zone not found"}`,
			))

			buf, resp, err := client.Zones.DownloadZonefile(zoneName)
			require.Nil(t, buf)
			require.NotNil(t, resp)
			require.Equal(t, api.ErrZoneMissing, err)
		})

		t.Run("Error - HTTP", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddTestCase(
				http.MethodGet, "/export/zonefile/"+zoneName, http.StatusInternalServerError,
				nil, nil, "", `{"message": "internal server error"}`,
			))

			buf, resp, err := client.Zones.DownloadZonefile(zoneName)
			require.Nil(t, buf)
			require.NotNil(t, err)
			require.Contains(t, err.Error(), "internal server error")
			require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		})

		t.Run("Error - Other", func(t *testing.T) {
			c := api.NewClient(errorClient{}, api.SetEndpoint(""))
			buf, resp, err := c.Zones.DownloadZonefile(zoneName)
			require.Nil(t, resp)
			require.Error(t, err)
			require.Nil(t, buf)
		})
	})
}

type errorClient struct{}

func (c errorClient) Do(req *http.Request) (*http.Response, error) {
	return nil, errors.New("oops")
}

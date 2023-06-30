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
	})
}

type errorClient struct{}

func (c errorClient) Do(req *http.Request) (*http.Response, error) {
	return nil, errors.New("oops")
}

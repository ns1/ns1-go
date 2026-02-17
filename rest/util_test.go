package rest_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	api "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/data"
	"gopkg.in/ns1/ns1-go.v2/rest/model/dns"
	"gopkg.in/ns1/ns1-go.v2/rest/model/filter"
	"gopkg.in/ns1/ns1-go.v2/rest/model/redirect"
)

func TestConfig(t *testing.T) {
	client := api.NewClient(
		&http.Client{Timeout: time.Second * 10},
		api.SetAPIKey("DoK9U6p6rddzUdL6KwOm"),
		api.SetEndpoint("https://api.nszero.com/v1/"),
	)

	t.Run("List", func(t *testing.T) {
		cfgList, resp, err := client.Redirects.List()
		require.NoError(t, err)
		require.Equal(t, 200, resp.StatusCode)
		require.GreaterOrEqual(t, len(cfgList), 1)
		t.Logf("+++ found %d configs", len(cfgList))
	})

	t.Run("Get", func(t *testing.T) {
		cfg, resp, err := client.Redirects.Get("cb1e64f2-ffb8-4806-8a25-77759f92c56a")
		require.NoError(t, err)
		require.Equal(t, 200, resp.StatusCode)
		require.NotNil(t, cfg)
		t.Logf("[%s]/[%s]/[%s]", cfg.Domain, cfg.Path, cfg.Target)
	})

	t.Run("Create/Get/Update/Get/Delete", func(t *testing.T) {
		cfg := redirect.NewConfigurationMinimal("test3.fformica.com", "/path", "http://localhost")
		cfg.Tags = []string{"test"}
		cfg, resp, err := client.Redirects.Create(cfg)
		require.NoError(t, err)
		require.Equal(t, 201, resp.StatusCode)
		require.NotNil(t, resp)
		require.NotNil(t, cfg)
		require.NotNil(t, cfg.ID)
		require.Equal(t, 1, len(cfg.Tags))
		t.Logf("Create [%s]=[%s]/[%s]/[%s]/[%s]", *cfg.ID, cfg.Domain, cfg.Path, cfg.Target, strings.Join(cfg.Tags, ","))
		cfg, resp, err = client.Redirects.Get(*cfg.ID)
		require.NoError(t, err)
		require.Equal(t, 200, resp.StatusCode)
		require.NotNil(t, resp)
		require.NotNil(t, cfg)
		require.NotNil(t, cfg.ID)
		require.Equal(t, 1, len(cfg.Tags))
		t.Logf("Get [%s]=[%s]/[%s]/[%s]/[%s]", *cfg.ID, cfg.Domain, cfg.Path, cfg.Target, strings.Join(cfg.Tags, ","))
		cfg.Target = "https://google.com"
		cfg.Tags = []string{}
		cfg, resp, err = client.Redirects.Update(cfg)
		require.NoError(t, err)
		require.Equal(t, 200, resp.StatusCode)
		require.NotNil(t, resp)
		require.NotNil(t, cfg)
		require.NotNil(t, cfg.ID)
		require.Equal(t, 0, len(cfg.Tags))
		t.Logf("Update [%s]=[%s]/[%s]/[%s]/[%s]", *cfg.ID, cfg.Domain, cfg.Path, cfg.Target, strings.Join(cfg.Tags, ","))
		cfg, resp, err = client.Redirects.Get(*cfg.ID)
		require.NoError(t, err)
		require.Equal(t, 200, resp.StatusCode)
		require.NotNil(t, resp)
		require.NotNil(t, cfg)
		require.NotNil(t, cfg.ID)
		require.Equal(t, 0, len(cfg.Tags))
		t.Logf("Get [%s]=[%s]/[%s]/[%s]/[%s]", *cfg.ID, cfg.Domain, cfg.Path, cfg.Target, strings.Join(cfg.Tags, ","))
		resp, err = client.Redirects.Delete(*cfg.ID)
		require.NoError(t, err)
		require.Equal(t, 204, resp.StatusCode)
		require.NotNil(t, resp)
	})
}

func TestCert(t *testing.T) {
	client := api.NewClient(
		&http.Client{Timeout: time.Second * 10},
		api.SetAPIKey("DoK9U6p6rddzUdL6KwOm"),
		api.SetEndpoint("https://api.nszero.com/v1/"),
	)

	t.Run("List", func(t *testing.T) {
		certList, resp, err := client.RedirectCertificates.List()
		require.NoError(t, err)
		require.Equal(t, 200, resp.StatusCode)
		require.GreaterOrEqual(t, len(certList), 1)
		t.Logf("+++ found %d certs", len(certList))
	})

	t.Run("Get", func(t *testing.T) {
		cfg, resp, err := client.RedirectCertificates.Get("b42dc47e-527f-46ac-b86d-14f6a76d3c58")
		require.NoError(t, err)
		require.Equal(t, 200, resp.StatusCode)
		require.NotNil(t, cfg)
		t.Logf("[%s]", cfg.Domain)
	})

	t.Run("Create/Update/Delete", func(t *testing.T) {
		cert, resp, err := client.RedirectCertificates.Create("test6.fformica.com")
		require.NoError(t, err)
		require.Equal(t, 200, resp.StatusCode)
		require.NotNil(t, resp)
		require.NotNil(t, cert)
		require.NotNil(t, cert.ID)
		t.Logf("[%s]=[%s]", *cert.ID, cert.Domain)
		time.Sleep(time.Duration(3))
		resp, err = client.RedirectCertificates.Update(*cert.ID)
		require.NoError(t, err)
		require.Equal(t, 204, resp.StatusCode)
		t.Logf("[%s]=[%s]", *cert.ID, cert.Domain)
		resp, err = client.RedirectCertificates.Delete(*cert.ID)
		require.NoError(t, err)
		require.Equal(t, 204, resp.StatusCode)
	})
}

func TestRecord(t *testing.T) {
	client := api.NewClient(
		&http.Client{Timeout: time.Second * 10},
		api.SetAPIKey("DoK9U6p6rddzUdL6KwOm"),
		api.SetEndpoint("https://api.nszero.com/v1/"),
	)

	client.Records.Delete("fformica.com", "my.test.a.fformica.com", "A")

	record := dns.NewRecord("fformica.com", "my.test.a.fformica.com", "A", map[string]string{"one": "two"}, nil)
	record.AddFilter(filter.NewGeofenceRegional(true))
	record.Regions = map[string]data.Region{"cal": {}}
	resp, err := client.Records.Create(record)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)

	record, resp, err = client.Records.Get("fformica.com", "my.test.a.fformica.com", "A")
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
	require.Equal(t, map[string]string{"one": "two"}, record.Tags)
	require.Equal(t, 1, len(record.Filters))
	require.Equal(t, 1, len(record.Regions))

	record.Tags = map[string]string{"three": "four"}
	record.Filters = nil
	record.Regions = nil
	resp, err = client.Records.Update(record)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)

	record, resp, err = client.Records.Get("fformica.com", "my.test.a.fformica.com", "A")
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
	require.Equal(t, map[string]string{"three": "four"}, record.Tags)
	require.Equal(t, 1, len(record.Filters))
	require.Equal(t, 1, len(record.Regions))

	record.Tags = nil
	record.Filters = []*filter.Filter{}
	record.Regions = map[string]data.Region{}
	resp, err = client.Records.Update(record)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
	require.Equal(t, 0, len(record.Filters))
	require.Equal(t, 0, len(record.Regions))

	record, resp, err = client.Records.Get("fformica.com", "my.test.a.fformica.com", "A")
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
	require.Equal(t, map[string]string{"three": "four"}, record.Tags)
}

func Test_Zone(t *testing.T) {
	client := api.NewClient(
		&http.Client{Timeout: time.Second * 10},
		api.SetAPIKey("DoK9U6p6rddzUdL6KwOm"),
		api.SetEndpoint("https://api.nszero.com/v1/"),
	)
	fqdn := "fftest.net"

	client.Zones.Delete(fqdn)

	zone := dns.NewZone(fqdn)
	resp, err := client.Zones.Create(zone)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
	zone, resp, err = client.Zones.Get(fqdn, true)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
	fmt.Println("create nil ->", zone.NetworkIDs)

	// add 0
	zone.NetworkIDs = []int{0}
	resp, err = client.Zones.Update(zone)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
	zone, resp, err = client.Zones.Get(fqdn, true)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
	fmt.Println("update [0] ->", zone.NetworkIDs)

	// no-op
	zone.NetworkIDs = nil
	resp, err = client.Zones.Update(zone)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
	zone, resp, err = client.Zones.Get(fqdn, true)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
	fmt.Println("update nil ->", zone.NetworkIDs)

	// empty
	zone.NetworkIDs = []int{}
	resp, err = client.Zones.Update(zone)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
	zone, resp, err = client.Zones.Get(fqdn, true)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
	fmt.Println("update [] ->", zone.NetworkIDs)

	resp, err = client.Zones.Delete(fqdn)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
}

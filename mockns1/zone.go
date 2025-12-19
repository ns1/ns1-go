package mockns1

import (
	"net/http"

	"gopkg.in/ns1/ns1-go.v2/rest/model/dns"
)

// AddZoneListTestCase sets up a test case for the api.Client.Zones.List()
// function
func (s *Service) AddZoneListTestCase(
	requestHeaders, responseHeaders http.Header,
	response []*dns.Zone,
) error {
	return s.AddTestCase(
		http.MethodGet, "/zones", http.StatusOK, requestHeaders,
		responseHeaders, "", response,
	)
}

// AddZoneGetTestCase sets up a test case for the api.Client.Zones.Get()
// function
func (s *Service) AddZoneGetTestCase(
	name string,
	requestHeaders, responseHeaders http.Header,
	response *dns.Zone,
	records bool,
) error {
	uri := "/zones/" + name
	if !records {
		uri = uri + "?records=false"
	}
	return s.AddTestCase(
		http.MethodGet, uri, http.StatusOK, requestHeaders,
		responseHeaders, "", response,
	)
}

// AddZoneCreateTestCase sets up a test case for the api.Client.Zones.Create()
// function
func (s *Service) AddZoneCreateTestCase(
	requestHeaders, responseHeaders http.Header,
	zone, response *dns.Zone,
) error {
	return s.AddTestCase(
		http.MethodPut, "/zones/"+zone.Zone, http.StatusCreated, requestHeaders,
		responseHeaders, zone, response,
	)
}

// AddZoneUpdateTestCase sets up a test case for the api.Client.Zones.Update()
// function
func (s *Service) AddZoneUpdateTestCase(
	requestHeaders, responseHeaders http.Header,
	zone, response *dns.Zone,
) error {
	return s.AddTestCase(
		http.MethodPost, "/zones/"+zone.Zone, http.StatusOK, requestHeaders,
		responseHeaders, zone, response,
	)
}

// AddZoneDeleteTestCase sets up a test case for the api.Client.Zones.Delete()
// function
func (s *Service) AddZoneDeleteTestCase(
	name string, requestHeaders, responseHeaders http.Header,
) error {
	return s.AddTestCase(
		http.MethodDelete, "/zones/"+name, http.StatusNoContent, requestHeaders,
		responseHeaders, "", "",
	)
}

// AddZoneExportZonefileTestCase sets up a test case for the api.Client.Zones.ExportZonefile()
// function
func (s *Service) AddZoneExportZonefileTestCase(
	name string,
	requestHeaders, responseHeaders http.Header,
	response *dns.ZoneFileExportStatus,
) error {
	return s.AddTestCase(
		http.MethodPut, "/zones/"+name+"/export/zonefile", http.StatusOK, requestHeaders,
		responseHeaders, map[string]interface{}{}, response,
	)
}

// AddGetExportZonefileStatusTestCase sets up a test case for the api.Client.Zones.GetExportZonefileStatus()
// function
func (s *Service) AddGetExportZonefileStatusTestCase(
	name string,
	requestHeaders, responseHeaders http.Header,
	response *dns.ZoneFileExportStatus,
) error {
	return s.AddTestCase(
		http.MethodGet, "/export/zonefile/"+name+"/status", http.StatusOK, requestHeaders,
		responseHeaders, "", response,
	)
}

// AddDownloadZonefileTestCase sets up a test case for the api.Client.Zones.DownloadZonefile()
// function
func (s *Service) AddDownloadZonefileTestCase(
	name string,
	requestHeaders, responseHeaders http.Header,
	response string,
) error {
	return s.AddTestCase(
		http.MethodGet, "/export/zonefile/"+name, http.StatusOK, requestHeaders,
		responseHeaders, "", response,
	)
}

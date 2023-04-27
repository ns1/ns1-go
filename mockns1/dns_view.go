package mockns1

import (
	"fmt"
	"net/http"

	"gopkg.in/ns1/ns1-go.v2/rest/model/dns"
)

// AddDNSViewListTestCase sets up a test case for the api.Client.View.List()
// function
func (s *Service) AddDNSViewListTestCase(
	requestHeaders, responseHeaders http.Header,
	response []*dns.DNSView,
) error {
	return s.AddTestCase(
		http.MethodGet, "views", http.StatusOK, requestHeaders,
		responseHeaders, "", response,
	)
}

// AddDNSViewGetTestCase sets up a test case for the api.Client.View.Get()
// function
func (s *Service) AddDNSViewGetTestCase(
	viewName string,
	requestHeaders, responseHeaders http.Header,
	response *dns.DNSView,
) error {
	return s.AddTestCase(
		http.MethodGet, fmt.Sprintf("views/%s", viewName), http.StatusOK, requestHeaders,
		responseHeaders, "", response,
	)
}

// AddDNSViewCreateTestCase sets up a test case for the api.Client.View.Create()
// function
func (s *Service) AddDNSViewCreateTestCase(
	requestHeaders, responseHeaders http.Header,
	dnsView, response *dns.DNSView,
) error {
	return s.AddTestCase(
		http.MethodPut, fmt.Sprintf("views/%s", dnsView.Name), http.StatusOK, requestHeaders,
		responseHeaders, dnsView, response,
	)
}

// AddDNSViewUpdateTestCase sets up a test case for the api.Client.View.Update()
// function
func (s *Service) AddDNSViewUpdateTestCase(
	requestHeaders, responseHeaders http.Header,
	dnsView, response *dns.DNSView,
) error {
	return s.AddTestCase(
		http.MethodPost, fmt.Sprintf("views/%s", dnsView.Name), http.StatusOK, requestHeaders,
		responseHeaders, dnsView, response,
	)
}

// AddDNSViewGetPreferencesTestCase sets up a test case for the api.Client.View.GetPreferences()
// function
func (s *Service) AddDNSViewGetPreferencesTestCase(
	requestHeaders, responseHeaders http.Header,
	response interface{},
) error {
	return s.AddTestCase(
		http.MethodGet, "config/views/preference", http.StatusOK, requestHeaders,
		responseHeaders, "", response,
	)
}

// AddDNSViewUpdatePreferencesTestCase sets up a test case for the api.Client.View.GetPreferences()
// function
func (s *Service) AddDNSViewUpdatePreferencesTestCase(
	requestHeaders, responseHeaders http.Header,
	body, response interface{},
) error {
	return s.AddTestCase(
		http.MethodPost, "config/views/preference", http.StatusOK, requestHeaders,
		responseHeaders, body, response,
	)
}

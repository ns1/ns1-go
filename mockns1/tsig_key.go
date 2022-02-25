package mockns1

import (
	"fmt"
	"net/http"

	"gopkg.in/ns1/ns1-go.v2/rest/model/dns"
)

// AddTsigKeyListTestCase sets up a test case for the api.Client.TSIG.List()
// function
func (s *Service) AddTsigKeyListTestCase(
	requestHeaders, responseHeaders http.Header,
	response []*dns.TSIGKey,
) error {
	return s.AddTestCase(
		http.MethodGet, "tsig", http.StatusOK, requestHeaders,
		responseHeaders, "", response,
	)
}

// AddTsigKeyGetTestCase sets up a test case for the api.Client.TSIG.Get()
// function
func (s *Service) AddTsigKeyGetTestCase(
	name string,
	requestHeaders, responseHeaders http.Header,
	response *dns.TSIGKey,
) error {
	return s.AddTestCase(
		http.MethodGet, fmt.Sprintf("/tsig/%s", name), http.StatusOK, requestHeaders,
		responseHeaders, "", response,
	)
}

// AddTsigKeyCreateTestCase sets up a test case for the api.Client.TSIG.Create()
// function
func (s *Service) AddTsigKeyCreateTestCase(
	requestHeaders, responseHeaders http.Header,
	tsigKey, response *dns.TSIGKey,
) error {
	return s.AddTestCase(
		http.MethodPut, fmt.Sprintf("tsig/%s", tsigKey.Name), http.StatusOK, requestHeaders,
		responseHeaders, tsigKey, response,
	)
}

// AddTsigKeyUpdateTestCase sets up a test case for the api.Client.TSIG.Update()
// function
func (s *Service) AddTsigKeyUpdateTestCase(
	requestHeaders, responseHeaders http.Header,
	tsigKey, response *dns.TSIGKey,
) error {
	return s.AddTestCase(
		http.MethodPost, fmt.Sprintf("tsig/%s", tsigKey.Name), http.StatusOK, requestHeaders,
		responseHeaders, tsigKey, response,
	)
}

// AddTsigKeyDeleteTestCase sets up a test case for the api.Client.TSIG.Delete()
// function
func (s *Service) AddTsigKeyDeleteTestCase(
	requestHeaders, responseHeaders http.Header,
	tsigKey, response *dns.TSIGKey,
) error {
	return s.AddTestCase(
		http.MethodDelete, fmt.Sprintf("tsig/%s", tsigKey.Name), http.StatusOK, requestHeaders,
		responseHeaders, "", "",
	)
}

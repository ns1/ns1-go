package mockns1

import (
	"net/http"

	"gopkg.in/ns1/ns1-go.v2/rest/model/dns"
)

// NetworkGetTestCase sets up a test case for the api.Client.Network.Get()
// function
func (s *Service) NetworkGetTestCase(
	requestHeaders, responseHeaders http.Header,
	response []*dns.Network,
) error {
	return s.AddTestCase(
		http.MethodGet, "/networks", http.StatusOK, requestHeaders,
		responseHeaders, "", response,
	)
}

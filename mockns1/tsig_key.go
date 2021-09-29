package mockns1

import (
	"net/http"

	"gopkg.in/ns1/ns1-go.v2/rest/model/dns"
)

// AddTsigKeyListTestCase sets up a test case for the api.Client.TSIG.List()
// function
func (s *Service) AddTsigKeyListTestCase(
	requestHeaders, responseHeaders http.Header,
	response []*dns.Tsig_key,
) error {
	return s.AddTestCase(
		http.MethodGet, "tsig", http.StatusOK, requestHeaders,
		responseHeaders, "", response,
	)
}

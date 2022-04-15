package mockns1

import (
	"fmt"
	"net/http"

	"gopkg.in/ns1/ns1-go.v2/rest/model/dns"
)

func (s *Service) AddRecordGetTestCase(
	zone, domain, recordType string,
	requestHeaders, responseHeaders http.Header,
	response *dns.Record,
) error {
	return s.AddTestCase(
		http.MethodGet,
		fmt.Sprintf("/zones/%s/%s/%s", zone, domain, recordType),
		http.StatusOK,
		requestHeaders,
		responseHeaders,
		"",
		response,
	)
}

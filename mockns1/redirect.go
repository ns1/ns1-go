package mockns1

import (
	"net/http"

	"gopkg.in/ns1/ns1-go.v2/rest/model/redirect"
)

// AddRedirectListTestCase sets up a test case for the api.Client.Redirects.List()
// function
func (s *Service) AddRedirectListTestCase(
	requestHeaders, responseHeaders http.Header,
	response *redirect.ConfigurationList,
) error {
	return s.AddTestCase(
		http.MethodGet, "/redirect", http.StatusOK, requestHeaders,
		responseHeaders, "", response,
	)
}

// AddRedirectGetTestCase sets up a test case for the api.Client.Redirects.Get()
// function
func (s *Service) AddRedirectGetTestCase(
	id string,
	requestHeaders, responseHeaders http.Header,
	response *redirect.Configuration,
) error {
	return s.AddTestCase(
		http.MethodGet, "/redirect/"+id, http.StatusOK, requestHeaders,
		responseHeaders, "", response,
	)
}

// AddRedirectCreateTestCase sets up a test case for the api.Client.Redirects.Create()
// function
func (s *Service) AddRedirectCreateTestCase(
	requestHeaders, responseHeaders http.Header,
	request, response *redirect.Configuration,
) error {
	return s.AddTestCase(
		http.MethodPut, "/redirect", http.StatusCreated, requestHeaders,
		responseHeaders, request, response,
	)
}

// AddRedirectUpdateTestCase sets up a test case for the api.Client.Redirects.Update()
// function
func (s *Service) AddRedirectUpdateTestCase(
	requestHeaders, responseHeaders http.Header,
	request, response *redirect.Configuration,
) error {
	return s.AddTestCase(
		http.MethodPost, "/redirect/"+*request.ID, http.StatusOK, requestHeaders,
		responseHeaders, request, response,
	)
}

// AddRedirectDeleteTestCase sets up a test case for the api.Client.Redirects.Delete()
// function
func (s *Service) AddRedirectDeleteTestCase(
	id string, requestHeaders, responseHeaders http.Header,
) error {
	return s.AddTestCase(
		http.MethodDelete, "/redirect/"+id, http.StatusNoContent, requestHeaders,
		responseHeaders, "", "",
	)
}

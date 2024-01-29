package mockns1

import (
	"net/http"

	"gopkg.in/ns1/ns1-go.v2/rest/model/redirect"
)

// AddRedirectCertificateListTestCase sets up a test case for the api.Client.RedirectCertificates.List()
// function
func (s *Service) AddRedirectCertificateListTestCase(
	requestHeaders, responseHeaders http.Header,
	response *redirect.CertificateList,
) error {
	return s.AddTestCase(
		http.MethodGet, "/redirect/certificates", http.StatusOK, requestHeaders,
		responseHeaders, "", response,
	)
}

// AddRedirectCertificateGetTestCase sets up a test case for the api.Client.RedirectCertificates.Get()
// function
func (s *Service) AddRedirectCertificateGetTestCase(
	id string,
	requestHeaders, responseHeaders http.Header,
	response *redirect.Certificate,
) error {
	return s.AddTestCase(
		http.MethodGet, "/redirect/certificates/"+id, http.StatusOK, requestHeaders,
		responseHeaders, "", response,
	)
}

// AddRedirectCertificateCreateTestCase sets up a test case for the api.Client.RedirectCertificates.Create()
// function
func (s *Service) AddRedirectCertificateCreateTestCase(
	requestHeaders, responseHeaders http.Header,
	request, response *redirect.Certificate,
) error {
	return s.AddTestCase(
		http.MethodPut, "/redirect/certificates", http.StatusCreated, requestHeaders,
		responseHeaders, request, response,
	)
}

// AddRedirectCertificateUpdateTestCase sets up a test case for the api.Client.RedirectCertificates.Update()
// function
func (s *Service) AddRedirectCertificateUpdateTestCase(
	requestHeaders, responseHeaders http.Header,
	id string, response *redirect.Certificate,
) error {
	return s.AddTestCase(
		http.MethodPost, "/redirect/certificates/"+id, http.StatusOK, requestHeaders,
		responseHeaders, "", response,
	)
}

// AddRedirectCertificateDeleteTestCase sets up a test case for the api.Client.RedirectCertificates.Delete()
// function
func (s *Service) AddRedirectCertificateDeleteTestCase(
	id string, requestHeaders, responseHeaders http.Header,
) error {
	return s.AddTestCase(
		http.MethodDelete, "/redirect/certificates/"+id, http.StatusNoContent, requestHeaders,
		responseHeaders, "", "",
	)
}

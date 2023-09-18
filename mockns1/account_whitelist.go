package mockns1

import (
	"fmt"
	"net/http"

	"gopkg.in/ns1/ns1-go.v2/rest/model/account"
)

// AddGlobalIPWhitelistListTestCase sets up a test case for the api.Client.GlobalIPWhitelist.List()
// function
func (s *Service) AddGlobalIPWhitelistListTestCase(
	requestHeaders, responseHeaders http.Header,
	response []*account.IPWhitelist,
) error {
	return s.AddTestCase(
		http.MethodGet, "/account/whitelist", http.StatusOK, requestHeaders,
		responseHeaders, "", response,
	)
}

// AddGlobalIPWhitelistGetTestCase sets up a test case for the api.Client.GlobalIPWhitelist.Get()
// function
func (s *Service) AddGlobalIPWhitelistGetTestCase(
	whitelistID string,
	requestHeaders, responseHeaders http.Header,
	response *account.IPWhitelist,
) error {
	return s.AddTestCase(
		http.MethodGet, fmt.Sprintf("/account/whitelist/%s", whitelistID), http.StatusOK, requestHeaders,
		responseHeaders, "", response,
	)
}

// AddGlobalIPWhitelistCreateTestCase sets up a test case for the api.Client.GlobalIPWhitelist.Create()
// function
func (s *Service) AddGlobalIPWhitelistCreateTestCase(
	requestHeaders, responseHeaders http.Header,
	whitelist, response *account.IPWhitelist,
) error {
	return s.AddTestCase(
		http.MethodPut, "/account/whitelist", http.StatusCreated, requestHeaders,
		responseHeaders, whitelist, response,
	)
}

// AddGlobalIPWhitelistUpdateTestCase sets up a test case for the api.Client.GlobalIPWhitelist.Update()
// function
func (s *Service) AddGlobalIPWhitelistUpdateTestCase(
	requestHeaders, responseHeaders http.Header,
	whitelist, response *account.IPWhitelist,
) error {
	return s.AddTestCase(
		http.MethodPost, fmt.Sprintf("/account/whitelist/%s", whitelist.ID), http.StatusOK, requestHeaders,
		responseHeaders, whitelist, response,
	)
}

// AddGlobalIPWhitelistDeleteTestCase sets up a test case for the api.Client.GlobalIPWhitelist.Delete()
// function
func (s *Service) AddGlobalIPWhitelistDeleteTestCase(
	whitelistID string,
	requestHeaders, responseHeaders http.Header,
) error {
	return s.AddTestCase(
		http.MethodDelete, fmt.Sprintf("/account/whitelist/%s", whitelistID), http.StatusOK, requestHeaders,
		responseHeaders, "", "",
	)
}

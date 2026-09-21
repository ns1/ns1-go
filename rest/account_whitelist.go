package rest

import (
	"errors"
	"fmt"
	"net/http"

	"gopkg.in/ns1/ns1-go.v2/rest/model/account"
)

// GlobalIPWhitelistService handles 'account/whitelist' endpoint.
type GlobalIPWhitelistService service

// List returns all global IP whitelists in the account.
// NS1 API docs: https://developer.ibm.com/apis/catalog/ns1--ibm-ns1-connect-api/api/API--ns1--ibm-ns1-connect-api#listAccountIPWhitelistRecords
func (s *GlobalIPWhitelistService) List() ([]*account.IPWhitelist, *http.Response, error) {
	req, err := s.client.NewRequest("GET", "account/whitelist", nil)
	if err != nil {
		return nil, nil, err
	}

	wl := []*account.IPWhitelist{}
	resp, err := s.client.Do(req, &wl)
	if err != nil {
		return nil, resp, err
	}

	return wl, resp, nil
}

// Get returns details of a single global IP whitelist.
// NS1 API docs: https://developer.ibm.com/apis/catalog/ns1--ibm-ns1-connect-api/api/API--ns1--ibm-ns1-connect-api#getAccountIPWhitelistRecord
func (s *GlobalIPWhitelistService) Get(id string) (*account.IPWhitelist, *http.Response, error) {
	path := fmt.Sprintf("account/whitelist/%s", id)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var wl account.IPWhitelist
	resp, err := s.client.Do(req, &wl)
	if err != nil {
		switch err := err.(type) {
		case *Error:
			if resourceMissingMatch(err.Message) {
				return nil, resp, ErrIPWhitelistMissing
			}
		}
		return nil, resp, err
	}

	return &wl, resp, nil
}

// Create takes a *IPWhitelist and creates a new global IP whitelist.
// NS1 API docs: https://developer.ibm.com/apis/catalog/ns1--ibm-ns1-connect-api/api/API--ns1--ibm-ns1-connect-api#createAccountIPWhitelistRecord
func (s *GlobalIPWhitelistService) Create(wl *account.IPWhitelist) (*http.Response, error) {
	req, err := s.client.NewRequest("PUT", "account/whitelist", wl)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, &wl)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Update changes the name or values for a global IP whitelist.
// NS1 API docs: https://developer.ibm.com/apis/catalog/ns1--ibm-ns1-connect-api/api/API--ns1--ibm-ns1-connect-api#modifyAccountIPWhitelistRecord
func (s *GlobalIPWhitelistService) Update(wl *account.IPWhitelist) (*http.Response, error) {
	path := fmt.Sprintf("account/whitelist/%s", wl.ID)

	req, err := s.client.NewRequest("POST", path, wl)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, &wl)
	if err != nil {
		switch err := err.(type) {
		case *Error:
			if resourceMissingMatch(err.Message) {
				return resp, ErrIPWhitelistMissing
			}
		}
		return resp, err
	}

	return resp, nil
}

// Delete deletes a global IP whitelist.
// NS1 API docs: https://developer.ibm.com/apis/catalog/ns1--ibm-ns1-connect-api/api/API--ns1--ibm-ns1-connect-api#removeAccountIPWhitelistRecord
func (s *GlobalIPWhitelistService) Delete(id string) (*http.Response, error) {
	path := fmt.Sprintf("account/whitelist/%s", id)

	req, err := s.client.NewRequest("DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		switch err := err.(type) {
		case *Error:
			if resourceMissingMatch(err.Message) {
				return resp, ErrIPWhitelistMissing
			}
		}
		return resp, err
	}

	return resp, nil
}

var (
	// ErrIPWhitelistMissing bundles GET/POST/DELETE error.
	ErrIPWhitelistMissing = errors.New("whitelist does not exist")
)

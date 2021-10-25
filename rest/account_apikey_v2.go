package rest

import (
	"fmt"
	"net/http"

	"gopkg.in/ns1/ns1-go.v2/rest/model/account"
)

// APIKeysServiceV2 handles 'account/apikeys' endpoint.
type APIKeysServiceV2 service

// List returns all api keys in the account.
//
// NS1 API docs: https://ns1.com/api/#apikeys-get
func (s *APIKeysServiceV2) List() ([]*account.APIKeyV2, *http.Response, error) {
	req, err := s.client.NewRequest("GET", "account/apikeys", nil)
	if err != nil {
		return nil, nil, err
	}

	kl := []*account.APIKeyV2{}
	resp, err := s.client.Do(req, &kl)
	if err != nil {
		return nil, resp, err
	}

	return kl, resp, nil
}

// Get returns details of an api key, including permissions, for a single API Key.
// Note: do not use the API Key itself as the keyid in the URL — use the id of the key.
//
// NS1 API docs: https://ns1.com/api/#apikeys-id-get
func (s *APIKeysServiceV2) Get(keyID string) (*account.APIKeyV2, *http.Response, error) {
	path := fmt.Sprintf("account/apikeys/%s", keyID)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var a account.APIKeyV2
	resp, err := s.client.Do(req, &a)
	if err != nil {
		switch err.(type) {
		case *Error:
			if err.(*Error).Message == "unknown api key" {
				return nil, resp, ErrKeyMissing
			}

		}
		return nil, resp, err
	}

	return &a, resp, nil
}

// Create takes a *APIKeyV2 and creates a new account apikey.
//
// NS1 API docs: https://ns1.com/api/#apikeys-put
func (s *APIKeysServiceV2) Create(a *account.APIKeyV2) (*http.Response, error) {
	var (
		req *http.Request
		err error
	)

	// If this is DDI then the permissions need to be transformed to DDI-compatible permissions.
	if s.client.DDI && a != nil {
		ddiAPIKey := apiKeyToDDIAPIKeyV2(a)
		req, err = s.client.NewRequest("PUT", "account/apikeys", ddiAPIKey)
		if err != nil {
			return nil, err
		}
	} else {
		req, err = s.client.NewRequest("PUT", "account/apikeys", a)
		if err != nil {
			return nil, err
		}
	}

	// Update account fields with data from api(ensure consistent)
	resp, err := s.client.Do(req, &a)
	if err != nil {
		switch err.(type) {
		case *Error:
			if err.(*Error).Message == fmt.Sprintf("api key with name \"%s\" exists", a.Name) {
				return resp, ErrKeyExists
			}
		}
		return resp, err
	}

	return resp, nil
}

// Update changes the name or access rights for an API Key.
//
// NS1 API docs: https://ns1.com/api/#apikeys-id-post
func (s *APIKeysServiceV2) Update(a *account.APIKeyV2) (*http.Response, error) {
	path := fmt.Sprintf("account/apikeys/%s", a.ID)

	var (
		req *http.Request
		err error
	)

	// If this is DDI then the permissions need to be transformed to DDI-compatible permissions.
	if s.client.DDI && a != nil {
		ddiAPIKey := apiKeyToDDIAPIKeyV2(a)
		req, err = s.client.NewRequest("POST", path, ddiAPIKey)
		if err != nil {
			return nil, err
		}
	} else {
		req, err = s.client.NewRequest("POST", path, a)
		if err != nil {
			return nil, err
		}
	}

	// Update apikey fields with data from api(ensure consistent)
	resp, err := s.client.Do(req, &a)
	if err != nil {
		switch err.(type) {
		case *Error:
			if err.(*Error).Message == "unknown api key" {
				return resp, ErrKeyMissing
			}
		}
		return resp, err
	}

	return resp, nil
}

// Delete deletes an apikey.
//
// NS1 API docs: https://ns1.com/api/#apikeys-id-delete
func (s *APIKeysServiceV2) Delete(keyID string) (*http.Response, error) {
	path := fmt.Sprintf("account/apikeys/%s", keyID)

	req, err := s.client.NewRequest("DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		switch err.(type) {
		case *Error:
			if err.(*Error).Message == "unknown api key" {
				return resp, ErrKeyMissing
			}
		}
		return resp, err
	}

	return resp, nil
}

// Same as v1
// var (
// 	// ErrKeyExists bundles PUT create error.
// 	ErrKeyExists = errors.New("key already exists")
// 	// ErrKeyMissing bundles GET/POST/DELETE error.
// 	ErrKeyMissing = errors.New("key does not exist")
// )

func apiKeyToDDIAPIKeyV2(k *account.APIKeyV2) *ddiAPIKeyV2 {
	ddiAPIKey := &ddiAPIKeyV2{
		ID:                k.ID,
		Key:               k.Key,
		LastAccess:        k.LastAccess,
		Name:              k.Name,
		TeamIDs:           k.TeamIDs,
		IPWhitelist:       k.IPWhitelist,
		IPWhitelistStrict: k.IPWhitelistStrict,
	}

	ddiAPIKey.Permissions = convertDDIPermissionsV2(k.Permissions)

	return ddiAPIKey
}

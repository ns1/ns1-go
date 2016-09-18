package rest

import (
	"net/http"

	"gopkg.in/ns1/ns1-go.v2/rest/model/account"
)

// APIKeysService handles 'account/apikeys' endpoint.
type APIKeysService service

// List returns all api keys in the account.
//
// NS1 API docs: https://ns1.com/api/#apikeys-get
func (s *APIKeysService) List() ([]*account.APIKey, *http.Response, error) {
	req, err := s.client.NewRequest("GET", "account/apikeys", nil)
	if err != nil {
		return nil, nil, err
	}

	kl := []*account.APIKey{}
	resp, err := s.client.Do(req, &kl)
	if err != nil {
		return nil, resp, err
	}

	return kl, resp, nil
}

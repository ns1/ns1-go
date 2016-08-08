package rest

import (
	"github.com/ns1/ns1-go/account"
)

const (
	apiKeyPath = "account/apikeys"
)

// APIKeysService handles 'account/apikeys' endpoint.
type APIKeysService service

// List returns all api keys in the account.
//
// NS1 API docs: https://ns1.com/api/#apikeys-get
func (s *APIKeysService) List() ([]*account.APIKey, error) {
	req, err := s.client.NewRequest("GET", apiKeyPath, nil)
	if err != nil {
		return nil, err
	}

	kl := []*account.APIKey{}
	_, err = s.client.Do(req, &kl)
	if err != nil {
		return nil, err
	}

	return kl, nil
}

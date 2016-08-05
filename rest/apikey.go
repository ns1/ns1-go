package rest

import ns1 "github.com/ns1/ns1-go"

const (
	apiKeyPath = "account/apikeys"
)

// APIKeysService handles 'account/apikeys' endpoint.
type APIKeysService service

// List returns all api keys in the account.
//
// NS1 API docs: https://ns1.com/api/
func (s *APIKeysService) List() ([]*ns1.APIKey, error) {
	req, err := s.client.NewRequest("GET", apiKeyPath, nil)
	if err != nil {
		return nil, err
	}

	kl := []*ns1.APIKey{}
	_, err = s.client.Do(req, &kl)
	if err != nil {
		return nil, err
	}

	return kl, nil
}

package rest

import (
	"errors"
	"fmt"
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

// Get returns details of an api key, including permissions, for a single API Key.
// Note: do not use the API Key itself as the keyid in the URL — use the id of the key.
//
// NS1 API docs: https://ns1.com/api/#apikeys-id-get
func (s *APIKeysService) Get(keyID string) (*account.APIKey, *http.Response, error) {
	path := fmt.Sprintf("account/apikeys/%s", keyID)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var a account.APIKey
	resp, err := s.client.Do(req, &a)
	if err != nil {
		switch err.(type) {
		case *Error:
			if resourceMissingMatch(err.(*Error).Message) {
				return nil, resp, ErrKeyMissing
			}

		}
		return nil, resp, err
	}

	return &a, resp, nil
}

// Create takes a *APIKey and creates a new account apikey.
//
// NS1 API docs: https://ns1.com/api/#apikeys-put
func (s *APIKeysService) Create(a *account.APIKey) (*http.Response, error) {
	var (
		req *http.Request
		err error
	)

	req, err = s.client.NewRequest("PUT", "account/apikeys", a)
	if err != nil {
		return nil, err
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
func (s *APIKeysService) Update(a *account.APIKey) (*http.Response, error) {
	path := fmt.Sprintf("account/apikeys/%s", a.ID)

	var (
		req *http.Request
		err error
	)

	req, err = s.client.NewRequest("POST", path, a)
	if err != nil {
		return nil, err
	}

	// Update apikey fields with data from api(ensure consistent)
	resp, err := s.client.Do(req, &a)
	if err != nil {
		switch err.(type) {
		case *Error:
			if resourceMissingMatch(err.(*Error).Message) {
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
func (s *APIKeysService) Delete(keyID string) (*http.Response, error) {
	path := fmt.Sprintf("account/apikeys/%s", keyID)

	req, err := s.client.NewRequest("DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		switch err.(type) {
		case *Error:
			if resourceMissingMatch(err.(*Error).Message) {
				return resp, ErrKeyMissing
			}
		}
		return resp, err
	}

	return resp, nil
}

// UpdateSecret updates an API key secret's enabled status or expiration date.
//
// NS1 API docs: https://ns1.com/api/#apikeys-v1-secrets-secretid-put
func (s *APIKeysService) UpdateSecret(secret *account.APIKeySecret) (*http.Response, error) {
	path := fmt.Sprintf("apikeys/v1/secrets/%s", secret.ID)

	req, err := s.client.NewRequest("PUT", path, secret)
	if err != nil {
		return nil, err
	}

	// Update secret fields with data from api(ensure consistent)
	resp, err := s.client.Do(req, &secret)
	if err != nil {
		switch err.(type) {
		case *Error:
			if resourceMissingMatch(err.(*Error).Message) {
				return resp, ErrSecretMissing
			}
		}
		return resp, err
	}

	return resp, nil
}

// DeleteSecret deletes an API key secret.

func (s *APIKeysService) DeleteSecret(secretID string) (*http.Response, error) {
	path := fmt.Sprintf("apikeys/v1/secrets/%s", secretID)

	req, err := s.client.NewRequest("DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		switch err.(type) {
		case *Error:
			if resourceMissingMatch(err.(*Error).Message) {
				return resp, ErrSecretMissing
			}
		}
		return resp, err
	}

	return resp, nil
}

// GetSecret retrieves details of a specific API key secret by its ID.
//
// NS1 API docs: https://ns1.com/api/#apikeys-v1-secrets-secretid-get
func (s *APIKeysService) GetSecret(secretID string) (*account.APIKeySecret, *http.Response, error) {
	path := fmt.Sprintf("apikeys/v1/secrets/%s", secretID)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var secret account.APIKeySecret
	resp, err := s.client.Do(req, &secret)
	if err != nil {
		switch err.(type) {
		case *Error:
			if resourceMissingMatch(err.(*Error).Message) {
				return nil, resp, ErrSecretMissing
			}
		}
		return nil, resp, err
	}

	return &secret, resp, nil
}

// GetSecretSelf retrieves details of the API key secret used in the current request.
// This allows an API key to query its own secret information without needing manage_apikeys permission.
//
// NS1 API docs: https://ns1.com/api/#apikeys-v1-secrets-self-get
func (s *APIKeysService) GetSecretSelf() (*account.APIKeySecret, *http.Response, error) {
	return s.GetSecret("self")
}

// RenewSecret creates a new secret for an API key specified by its secret ID.
// This generates a new secret value with an updated expiration date.
// The API key must have an expiry_duration set, and cannot have more than 2 active secrets.
// Returns the new secret with the plaintext secret value (only time it's visible).
//
// NS1 API docs: https://ns1.com/api/#apikeys-v1-secrets-secretid-renew-post
func (s *APIKeysService) RenewSecret(secretID string) (*account.APIKeySecret, *http.Response, error) {
	path := fmt.Sprintf("apikeys/v1/secrets/%s/renew", secretID)

	req, err := s.client.NewRequest("POST", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var secret account.APIKeySecret
	resp, err := s.client.Do(req, &secret)
	if err != nil {
		switch err.(type) {
		case *Error:
			if resourceMissingMatch(err.(*Error).Message) {
				return nil, resp, ErrSecretMissing
			}
		}
		return nil, resp, err
	}

	return &secret, resp, nil
}

// RenewSecretSelf renews the API key secret used in the current request.
// This generates a new secret value with an updated expiration date.
// This allows an API key to renew itself without needing manage_apikeys permission.
// The API key must have an expiry_duration set, and cannot have more than 2 active secrets.
// Returns the new secret with the plaintext secret value (only time it's visible).
//
// NS1 API docs: https://ns1.com/api/#apikeys-v1-secrets-self-renew-post
func (s *APIKeysService) RenewSecretSelf() (*account.APIKeySecret, *http.Response, error) {
	return s.RenewSecret("self")
}

var (
	// ErrKeyExists bundles PUT create error.
	ErrKeyExists = errors.New("key already exists")
	// ErrKeyMissing bundles GET/POST/DELETE error.
	ErrKeyMissing = errors.New("key does not exist")
	// ErrSecretMissing bundles secret GET/PUT/DELETE error.
	ErrSecretMissing = errors.New("secret does not exist")
)

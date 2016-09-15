package rest

import (
	"net/http"

	"github.com/ns1/ns1-go/rest/model/account"
)

// UsersService handles 'account/users' endpoint.
type UsersService service

// List returns all users in the account.
//
// NS1 API docs: https://ns1.com/api/#users-get
func (s *UsersService) List() ([]*account.User, *http.Response, error) {
	req, err := s.client.NewRequest("GET", "account/users", nil)
	if err != nil {
		return nil, nil, err
	}

	ul := []*account.User{}
	resp, err := s.client.Do(req, &ul)
	if err != nil {
		return nil, resp, err
	}

	return ul, resp, nil
}

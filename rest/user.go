package rest

import "github.com/ns1/ns1-go/account"

const (
	userPath = "account/users"
)

// UsersService handles 'account/users' endpoint.
type UsersService service

// List returns all users in the account.
//
// NS1 API docs: https://ns1.com/api/#users-get
func (s *UsersService) List() ([]*account.User, error) {
	req, err := s.client.NewRequest("GET", userPath, nil)
	if err != nil {
		return nil, err
	}

	ul := []*account.User{}
	_, err = s.client.Do(req, &ul)
	if err != nil {
		return nil, err
	}

	return ul, nil
}

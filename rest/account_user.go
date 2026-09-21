package rest

import (
	"errors"
	"fmt"
	"net/http"

	"gopkg.in/ns1/ns1-go.v2/rest/model/account"
)

// UsersService handles 'account/users' endpoint.
type UsersService service

// List returns all users in the account.
//
// NS1 API docs: https://developer.ibm.com/apis/catalog/ns1--ibm-ns1-connect-api/api/API--ns1--ibm-ns1-connect-api#listUsers
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

// Get returns details of a single user.
//
// NS1 API docs: https://developer.ibm.com/apis/catalog/ns1--ibm-ns1-connect-api/api/API--ns1--ibm-ns1-connect-api#getUser
func (s *UsersService) Get(username string) (*account.User, *http.Response, error) {
	path := fmt.Sprintf("account/users/%s", username)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var u account.User
	resp, err := s.client.Do(req, &u)
	if err != nil {
		switch err := err.(type) {
		case *Error:
			if resourceMissingMatch(err.Message) {
				return nil, resp, ErrUserMissing
			}
		}
		return nil, resp, err
	}

	return &u, resp, nil
}

// Create takes a *User and creates a new account user.
//
// NS1 API docs: https://developer.ibm.com/apis/catalog/ns1--ibm-ns1-connect-api/api/API--ns1--ibm-ns1-connect-api#createUser
func (s *UsersService) Create(u *account.User) (*http.Response, error) {
	var (
		req *http.Request
		err error
	)

	req, err = s.client.NewRequest("PUT", "account/users", u)
	if err != nil {
		return nil, err
	}

	// Update user fields with data from api(ensure consistent)
	resp, err := s.client.Do(req, &u)
	if err != nil {
		switch err := err.(type) {
		case *Error:
			if err.Message == "request failed:Login Name is already in use." {
				return resp, ErrUserExists
			}
		}
		return resp, err
	}

	return resp, nil
}

// Update change contact details, notification settings, or access rights for a user.
//
// NS1 API docs: https://developer.ibm.com/apis/catalog/ns1--ibm-ns1-connect-api/api/API--ns1--ibm-ns1-connect-api#modifyUser
func (s *UsersService) Update(u *account.User) (*http.Response, error) {
	path := fmt.Sprintf("account/users/%s", u.Username)

	var (
		req *http.Request
		err error
	)

	req, err = s.client.NewRequest("POST", path, u)
	if err != nil {
		return nil, err
	}

	// Update user fields with data from api(ensure consistent)
	resp, err := s.client.Do(req, &u)
	if err != nil {
		switch err := err.(type) {
		case *Error:
			if err.Message == "Unknown user" {
				return resp, ErrUserMissing
			}
		}
		return resp, err
	}

	return resp, nil
}

// Delete deletes a user.
//
// NS1 API docs: https://developer.ibm.com/apis/catalog/ns1--ibm-ns1-connect-api/api/API--ns1--ibm-ns1-connect-api#removeUser
func (s *UsersService) Delete(username string) (*http.Response, error) {
	path := fmt.Sprintf("account/users/%s", username)

	req, err := s.client.NewRequest("DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		switch err := err.(type) {
		case *Error:
			if err.Message == "Unknown user" {
				return resp, ErrUserMissing
			}
		}
		return resp, err
	}

	return resp, nil
}

var (
	// ErrUserExists bundles PUT create error.
	ErrUserExists = errors.New("user already exists")
	// ErrUserMissing bundles GET/POST/DELETE error.
	ErrUserMissing = errors.New("user does not exist")
)

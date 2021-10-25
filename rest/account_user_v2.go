package rest

import (
	"fmt"
	"net/http"

	"gopkg.in/ns1/ns1-go.v2/rest/model/account"
)

// UsersServiceV2 handles 'account/users' endpoint.
type UsersServiceV2 service

// List returns all users in the account.
//
// NS1 API docs: https://ns1.com/api/#users-get
func (s *UsersServiceV2) List() ([]*account.UserV2, *http.Response, error) {
	req, err := s.client.NewRequest("GET", "account/users", nil)
	if err != nil {
		return nil, nil, err
	}

	ul := []*account.UserV2{}
	resp, err := s.client.Do(req, &ul)
	if err != nil {
		return nil, resp, err
	}

	return ul, resp, nil
}

// Get returns details of a single user.
//
// NS1 API docs: https://ns1.com/api/#users-user-get
func (s *UsersServiceV2) Get(username string) (*account.UserV2, *http.Response, error) {
	path := fmt.Sprintf("account/users/%s", username)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var u account.UserV2
	resp, err := s.client.Do(req, &u)
	if err != nil {
		switch err.(type) {
		case *Error:
			if err.(*Error).Message == "Unknown user" {
				return nil, resp, ErrUserMissing
			}
		}
		return nil, resp, err
	}

	return &u, resp, nil
}

// Create takes a *UserV2 and creates a new account user.
//
// NS1 API docs: https://ns1.com/api/#users-put
func (s *UsersServiceV2) Create(u *account.UserV2) (*http.Response, error) {
	var (
		req *http.Request
		err error
	)

	// If user u is in a team, its permissions are invalid and will use the team's permissions
	if len(u.TeamIDs) != 0 {
		u.Permissions = nil
	}

	// If this is DDI then the permissions need to be transformed to DDI-compatible permissions.
	if s.client.DDI && u != nil {
		ddiUser := userToDDIUserV2(u)
		req, err = s.client.NewRequest("PUT", "account/users", ddiUser)
		if err != nil {
			return nil, err
		}
	} else {
		req, err = s.client.NewRequest("PUT", "account/users", u)
		if err != nil {
			return nil, err
		}
	}

	// Update user fields with data from api(ensure consistent)
	resp, err := s.client.Do(req, &u)
	if err != nil {
		switch err.(type) {
		case *Error:
			if err.(*Error).Message == "request failed:Login Name is already in use." {
				return resp, ErrUserExists
			}
		}
		return resp, err
	}

	return resp, nil
}

// Update change contact details, notification settings, or access rights for a user.
//
// NS1 API docs: https://ns1.com/api/#users-user-post
func (s *UsersServiceV2) Update(u *account.UserV2) (*http.Response, error) {
	path := fmt.Sprintf("account/users/%s", u.Username)

	var (
		req *http.Request
		err error
	)

	// If this is DDI then the permissions need to be transformed to DDI-compatible permissions.
	if s.client.DDI && u != nil {
		ddiUser := userToDDIUserV2(u)
		req, err = s.client.NewRequest("POST", path, ddiUser)
		if err != nil {
			return nil, err
		}
	} else {
		req, err = s.client.NewRequest("POST", path, u)
		if err != nil {
			return nil, err
		}
	}

	// Update user fields with data from api(ensure consistent)
	resp, err := s.client.Do(req, &u)
	if err != nil {
		switch err.(type) {
		case *Error:
			if err.(*Error).Message == "Unknown user" {
				return resp, ErrUserMissing
			}
		}
		return resp, err
	}

	return resp, nil
}

// Delete deletes a user.
//
// NS1 API docs: https://ns1.com/api/#users-user-delete
func (s *UsersServiceV2) Delete(username string) (*http.Response, error) {
	path := fmt.Sprintf("account/users/%s", username)

	req, err := s.client.NewRequest("DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		switch err.(type) {
		case *Error:
			if err.(*Error).Message == "Unknown user" {
				return resp, ErrUserMissing
			}
		}
		return resp, err
	}

	return resp, nil
}

// Same as v1
// var (
// 	// ErrUserExists bundles PUT create error.
// 	ErrUserExists = errors.New("user already exists")
// 	// ErrUserMissing bundles GET/POST/DELETE error.
// 	ErrUserMissing = errors.New("user does not exist")
// )

func userToDDIUserV2(u *account.UserV2) *ddiUserV2 {
	ddiUser := &ddiUserV2{
		LastAccess:        u.LastAccess,
		Name:              u.Name,
		Username:          u.Username,
		Email:             u.Email,
		TeamIDs:           u.TeamIDs,
		Notify:            u.Notify,
		IPWhitelist:       u.IPWhitelist,
		IPWhitelistStrict: u.IPWhitelistStrict,
	}

	ddiUser.Permissions = convertDDIPermissionsV2(u.Permissions)

	return ddiUser
}

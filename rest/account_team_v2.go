package rest

import (
	"fmt"
	"net/http"

	"gopkg.in/ns1/ns1-go.v2/rest/model/account"
)

// TeamsServiceV2 handles 'account/teams' endpoint.
type TeamsServiceV2 service

// List returns all teams in the account.
//
// NS1 API docs: https://ns1.com/api/#teams-get
func (s *TeamsServiceV2) List() ([]*account.TeamV2, *http.Response, error) {
	req, err := s.client.NewRequest("GET", "account/teams", nil)
	if err != nil {
		return nil, nil, err
	}

	tl := []*account.TeamV2{}
	resp, err := s.client.Do(req, &tl)
	if err != nil {
		return nil, resp, err
	}

	return tl, resp, nil
}

// Get returns details of a single team.
//
// NS1 API docs: https://ns1.com/api/#teams-id-get
func (s *TeamsServiceV2) Get(id string) (*account.TeamV2, *http.Response, error) {
	path := fmt.Sprintf("account/teams/%s", id)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var t account.TeamV2
	resp, err := s.client.Do(req, &t)
	if err != nil {
		switch err.(type) {
		case *Error:
			if err.(*Error).Message == "Unknown team id" {
				return nil, resp, ErrTeamMissing
			}
		}
		return nil, resp, err
	}

	return &t, resp, nil
}

// Create takes a *Team and creates a new account team.
//
// NS1 API docs: https://ns1.com/api/#teams-put
func (s *TeamsServiceV2) Create(t *account.TeamV2) (*http.Response, error) {
	var (
		req *http.Request
		err error
	)

	// If this is DDI then the permissions need to be transformed to DDI-compatible permissions.
	if s.client.DDI && t != nil {
		ddiTeam := teamToDDITeamV2(t)
		req, err = s.client.NewRequest("PUT", "account/teams", ddiTeam)
		if err != nil {
			return nil, err
		}
	} else {
		req, err = s.client.NewRequest("PUT", "account/teams", t)
		if err != nil {
			return nil, err
		}
	}

	// Update team fields with data from api(ensure consistent)
	resp, err := s.client.Do(req, &t)
	if err != nil {
		switch err.(type) {
		case *Error:
			if err.(*Error).Message == fmt.Sprintf("team with name \"%s\" exists", t.Name) {
				return resp, ErrTeamExists
			}
		}
		return resp, err
	}

	return resp, nil
}

// Update changes the name or access rights for a team.
//
// NS1 API docs: https://ns1.com/api/#teams-id-post
func (s *TeamsServiceV2) Update(t *account.TeamV2) (*http.Response, error) {
	path := fmt.Sprintf("account/teams/%s", t.ID)

	var (
		req *http.Request
		err error
	)

	// If this is DDI then the permissions need to be transformed to DDI-compatible permissions.
	if s.client.DDI && t != nil {
		ddiTeam := teamToDDITeamV2(t)
		req, err = s.client.NewRequest("POST", path, ddiTeam)
		if err != nil {
			return nil, err
		}
	} else {
		req, err = s.client.NewRequest("POST", path, t)
		if err != nil {
			return nil, err
		}
	}

	// Update team fields with data from api(ensure consistent)
	resp, err := s.client.Do(req, &t)
	if err != nil {
		switch err.(type) {
		case *Error:
			if err.(*Error).Message == "unknown team id" {
				return resp, ErrTeamMissing
			}
		}
		return resp, err
	}

	return resp, nil
}

// Delete deletes a team.
//
// NS1 API docs: https://ns1.com/api/#teams-id-delete
func (s *TeamsServiceV2) Delete(id string) (*http.Response, error) {
	path := fmt.Sprintf("account/teams/%s", id)

	req, err := s.client.NewRequest("DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		switch err.(type) {
		case *Error:
			if err.(*Error).Message == "unknown team id" {
				return resp, ErrTeamMissing
			}
		}
		return resp, err
	}

	return resp, nil
}

// Same as v1
// var (
// 	// ErrTeamExists bundles PUT create error.
// 	ErrTeamExists = errors.New("team already exists")
// 	// ErrTeamMissing bundles GET/POST/DELETE error.
// 	ErrTeamMissing = errors.New("team does not exist")
// )

func teamToDDITeamV2(t *account.TeamV2) *ddiTeamV2 {
	ddiTeam := &ddiTeamV2{
		ID:          t.ID,
		Name:        t.Name,
		IPWhitelist: t.IPWhitelist,
	}

	ddiTeam.Permissions = convertDDIPermissionsV2(t.Permissions)

	return ddiTeam
}

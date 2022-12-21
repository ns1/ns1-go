package rest

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"gopkg.in/ns1/ns1-go.v2/rest/model/account"
)

// TeamsService handles 'account/teams' endpoint.
type TeamsService service

// List returns all teams in the account.
//
// NS1 API docs: https://ns1.com/api/#teams-get
func (s *TeamsService) List() ([]*account.Team, *http.Response, error) {
	return s.ListWithContext(context.Background())
}

// ListWithContext is the same as List, but takes a context.
func (s *TeamsService) ListWithContext(ctx context.Context) ([]*account.Team, *http.Response, error) {
	req, err := s.client.NewRequest("GET", "account/teams", nil)
	if err != nil {
		return nil, nil, err
	}
	req = req.WithContext(ctx)

	tl := []*account.Team{}
	resp, err := s.client.Do(req, &tl)
	if err != nil {
		return nil, resp, err
	}

	return tl, resp, nil
}

// Get returns details of a single team.
//
// NS1 API docs: https://ns1.com/api/#teams-id-get
func (s *TeamsService) Get(id string) (*account.Team, *http.Response, error) {
	return s.GetWithContext(context.Background(), id)
}

// GetWithContext is the same as Get, but takes a context.
func (s *TeamsService) GetWithContext(ctx context.Context, id string) (*account.Team, *http.Response, error) {
	path := fmt.Sprintf("account/teams/%s", id)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}
	req = req.WithContext(ctx)

	var t account.Team
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
func (s *TeamsService) Create(t *account.Team) (*http.Response, error) {
	return s.CreateWithContext(context.Background(), t)
}

// CreateWithContext is the same as Create, but takes a context.
func (s *TeamsService) CreateWithContext(ctx context.Context, t *account.Team) (*http.Response, error) {
	var (
		req *http.Request
		err error
	)

	// If this is DDI then the permissions need to be transformed to DDI-compatible permissions.
	if s.client.DDI && t != nil {
		ddiTeam := teamToDDITeam(t)
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
	req = req.WithContext(ctx)

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
func (s *TeamsService) Update(t *account.Team) (*http.Response, error) {
	return s.UpdateWithContext(context.Background(), t)
}

// UpdateWithContext is the same as Update, but takes a context.
func (s *TeamsService) UpdateWithContext(ctx context.Context, t *account.Team) (*http.Response, error) {
	path := fmt.Sprintf("account/teams/%s", t.ID)

	var (
		req *http.Request
		err error
	)

	// If this is DDI then the permissions need to be transformed to DDI-compatible permissions.
	if s.client.DDI && t != nil {
		ddiTeam := teamToDDITeam(t)
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
	req = req.WithContext(ctx)

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
func (s *TeamsService) Delete(id string) (*http.Response, error) {
	return s.DeleteWithContext(context.Background(), id)
}

// DeleteWithContext is the same as Delete, but takes a context.
func (s *TeamsService) DeleteWithContext(ctx context.Context, id string) (*http.Response, error) {
	path := fmt.Sprintf("account/teams/%s", id)

	req, err := s.client.NewRequest("DELETE", path, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

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

var (
	// ErrTeamExists bundles PUT create error.
	ErrTeamExists = errors.New("team already exists")
	// ErrTeamMissing bundles GET/POST/DELETE error.
	ErrTeamMissing = errors.New("team does not exist")
)

func teamToDDITeam(t *account.Team) *ddiTeam {
	ddiTeam := &ddiTeam{
		ID:          t.ID,
		Name:        t.Name,
		IPWhitelist: t.IPWhitelist,
		Permissions: ddiPermissionsMap{
			DNS:  t.Permissions.DNS,
			Data: t.Permissions.Data,
			Account: permissionsDDIAccount{
				ManageUsers:           t.Permissions.Account.ManageUsers,
				ManageTeams:           t.Permissions.Account.ManageTeams,
				ManageApikeys:         t.Permissions.Account.ManageApikeys,
				ManageAccountSettings: t.Permissions.Account.ManageAccountSettings,
				ViewActivityLog:       t.Permissions.Account.ViewActivityLog,
			},
		},
	}

	if t.Permissions.Security != nil {
		ddiTeam.Permissions.Security = permissionsDDISecurity(*t.Permissions.Security)
	}

	if t.Permissions.DHCP != nil {
		ddiTeam.Permissions.DHCP = *t.Permissions.DHCP
	}

	if t.Permissions.IPAM != nil {
		ddiTeam.Permissions.IPAM = *t.Permissions.IPAM
	}

	return ddiTeam
}

package rest

import (
	"net/http"

	"gopkg.in/ns1/ns1-go.v2/rest/model/account"
)

// TeamsService handles 'account/teams' endpoint.
type TeamsService service

// List returns all teams in the account.
//
// NS1 API docs: https://ns1.com/api/#teams-get
func (s *TeamsService) List() ([]*account.Team, *http.Response, error) {
	req, err := s.client.NewRequest("GET", "account/teams", nil)
	if err != nil {
		return nil, nil, err
	}

	tl := []*account.Team{}
	resp, err := s.client.Do(req, &tl)
	if err != nil {
		return nil, resp, err
	}

	return tl, resp, nil
}

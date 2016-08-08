package rest

import "github.com/ns1/ns1-go/account"

const (
	teamPath = "account/teams"
)

// TeamsService handles 'account/teams' endpoint.
type TeamsService service

// List returns all teams in the account.
//
// NS1 API docs: https://ns1.com/api/#teams-get
func (s *TeamsService) List() ([]*account.Team, error) {
	req, err := s.client.NewRequest("GET", teamPath, nil)
	if err != nil {
		return nil, err
	}

	tl := []*account.Team{}
	_, err = s.client.Do(req, &tl)
	if err != nil {
		return nil, err
	}

	return tl, nil
}

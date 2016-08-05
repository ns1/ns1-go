package rest

import ns1 "github.com/ns1/ns1-go"

const (
	teamPath = "account/teams"
)

// TeamsService handles 'account/teams' endpoint.
type TeamsService service

// List returns all teams in the account.
//
// NS1 API docs: https://ns1.com/api/#teams-get
func (s *TeamsService) List() ([]*ns1.Team, error) {
	req, err := s.client.NewRequest("GET", teamPath, nil)
	if err != nil {
		return nil, err
	}

	tl := []*ns1.Team{}
	_, err = s.client.Do(req, &tl)
	if err != nil {
		return nil, err
	}

	return tl, nil
}

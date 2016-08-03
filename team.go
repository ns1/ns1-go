package ns1

// Team wraps an NS1 /accounts/teams resource
type Team struct {
	Id          string         `json:"id,omitempty"`
	Name        string         `json:"name"`
	Permissions PermissionsMap `json:"permissions"`
}

// // GetTeams returns a list of all teams under the account
// func (c APIClient) GetTeams() ([]Team, error) {
// 	var teams []Team
// 	_, err := c.doHTTPUnmarshal("GET", "https://api.nsone.net/v1/account/teams", nil, &teams)
// 	return teams, err
// }

// // GetTeam takes an ID and returns details, including permissions, for a single team
// func (c APIClient) GetTeam(id string) (Team, error) {
// 	var t Team
// 	status, err := c.doHTTPUnmarshal("GET", fmt.Sprintf("https://api.nsone.net/v1/account/teams/%s", id), nil, &t)
// 	if status == 404 {
// 		t.Id = ""
// 		t.Name = ""
// 		return t, nil
// 	}
// 	return t, err
// }

// // CreateTeam takes a *Team and creates a new team
// func (c APIClient) CreateTeam(t *Team) error {
// 	return c.doHTTPBoth("PUT", "https://api.nsone.net/v1/account/teams", &t)
// }

// // DeleteTeam takes an ID and deletes a team. Any users of API keys that belong to the team will be removed from the team.
// func (c APIClient) DeleteTeam(id string) error {
// 	return c.doHTTPDelete(fmt.Sprintf("https://api.nsone.net/v1/account/teams/%s", id))
// }

// // UpdateTeam takes a *Team and change name or access rights for a team
// func (c APIClient) UpdateTeam(t *Team) error {
// 	return c.doHTTPBoth("POST", fmt.Sprintf("https://api.nsone.net/v1/account/teams/%s", t.Id), t)
// }

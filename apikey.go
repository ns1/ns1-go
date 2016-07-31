package nsone

import "fmt"

// Apikey wraps an NS1 /account/apikeys resource
type Apikey struct {
	Id          string         `json:"id,omitempty"`
	Name        string         `json:"name"`
	Key         string         `json:"key,omitempty"`
	LastAccess  int            `json:"last_access,omitempty"`
	Teams       []string       `json:"teams"`
	Permissions PermissionsMap `json:"permissions"`
}

// GetApikeys returns a list of all API keys under the account
func (c APIClient) GetApikeys() ([]Apikey, error) {
	var apikeys []Apikey
	_, err := c.doHTTPUnmarshal("GET", "https://api.nsone.net/v1/account/apikeys", nil, &apikeys)
	return apikeys, err
}

// GetApikey takes an ID and returns details, including permissions, for a single API key
func (c APIClient) GetApikey(id string) (Apikey, error) {
	var k Apikey
	status, err := c.doHTTPUnmarshal("GET", fmt.Sprintf("https://api.nsone.net/v1/account/apikeys/%s", id), nil, &k)
	if status == 404 {
		k.Id = ""
		k.Key = ""
		k.Name = ""
		return k, nil
	}
	return k, err
}

// CreateApikey takes an *Apikey and creates a new API key
func (c APIClient) CreateApikey(k *Apikey) error {
	return c.doHTTPBoth("PUT", fmt.Sprintf("https://api.nsone.net/v1/account/apikeys/%s", k.Id), &k)
}

// DeleteApikey takes an ID and deletes and API key
func (c APIClient) DeleteApikey(id string) error {
	return c.doHTTPDelete(fmt.Sprintf("https://api.nsone.net/v1/account/apikeys/%s", id))
}

// UpdateApikey takes an *Apikey and change name or access rights for an API key
func (c APIClient) UpdateApikey(k *Apikey) error {
	return c.doHTTPBoth("POST", fmt.Sprintf("https://api.nsone.net/v1/account/apikeys/%s", k.Id), k)
}

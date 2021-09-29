package rest

import (
	"net/http"

	"gopkg.in/ns1/ns1-go.v2/rest/model/dns"
)

// TsigService handles 'tsig' endpoint.
type TsigService service

// List returns all tsig keys and basic tsig keys configuration details for each.
//
// NS1 API docs: https://ns1.com/api/#getlist-tsig-keys
func (s *TsigService) List() ([]*dns.Tsig_key, *http.Response, error) {
	req, err := s.client.NewRequest("GET", "tsig", nil)
	if err != nil {
		return nil, nil, err
	}

	tsigKeyList := []*dns.Tsig_key{}
	var resp *http.Response
	resp, err = s.client.Do(req, &tsigKeyList)
	if err != nil {
		return nil, resp, err
	}

	return tsigKeyList, resp, nil
}

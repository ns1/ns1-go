package rest

import (
	"net/http"

	"gopkg.in/ns1/ns1-go.v2/rest/model/dns"
)

// NetworkService handles the 'networks' endpoint
type NetworkService service

// Get returns a list of all available NS1 DNS networks associated
// with your account.
// NS1 API docs: https://developer.ibm.com/apis/catalog/ns1--ibm-ns1-connect-api/api/API--ns1--ibm-ns1-connect-api#listNetworks
func (s *NetworkService) Get() ([]*dns.Network, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "networks", nil)
	if err != nil {
		return nil, nil, err
	}

	networks := []*dns.Network{}
	var resp *http.Response

	resp, err = s.client.Do(req, &networks)
	if err != nil {
		return nil, resp, err
	}

	return networks, resp, nil
}

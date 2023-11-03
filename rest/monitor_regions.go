package rest

import (
	"net/http"

	"gopkg.in/ns1/ns1-go.v2/rest/model/monitor"
)

// RegionsService handles 'monitoring/regions' endpoint.
type RegionsService service

// List returns all available monitoring regions for the account.
//
// NS1 API docs: https://ns1.com/api?docId=2247
func (s *RegionsService) List() ([]*monitor.Region, *http.Response, error) {
	req, err := s.client.NewRequest("GET", "monitoring/regions", nil)
	if err != nil {
		return nil, nil, err
	}

	mrl := []*monitor.Region{}
	resp, err := s.client.Do(req, &mrl)
	if err != nil {
		return nil, resp, err
	}

	return mrl, resp, nil
}

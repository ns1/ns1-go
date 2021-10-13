package rest

import (
	"errors"
	"fmt"
	"net/http"

	"gopkg.in/ns1/ns1-go.v2/rest/model/dns"
)

// DNSViewService handles 'views/' endpoint.
type DNSViewService service

// List returns all DNS Views
//
// NS1 API docs: https://ns1.com/api#getlist-all-dns-views
func (s *DNSViewService) List() ([]*dns.DNSView, *http.Response, error) {
	req, err := s.client.NewRequest("GET", "views", nil)
	if err != nil {
		return nil, nil, err
	}

	var vl []*dns.DNSView
	resp, err := s.client.Do(req, &vl)
	if err != nil {
		return nil, resp, err
	}

	return vl, resp, nil
}

// Create takes a *dns.DNSView and creates a new DNS View.
//
// The given DNSView must have at least the name
// NS1 API docs: https://ns1.com/api#putcreate-a-dns-view
func (s *DNSViewService) Create(v *dns.DNSView) (*http.Response, error) {
	req, err := s.client.NewRequest("PUT", fmt.Sprintf("/v1/views/%s", v.Name), v)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		switch errType := err.(type) {
		case *Error:
			if errType.Resp.StatusCode == 409 {
				return nil, ErrViewExists
			}
		}

		return resp, err
	}

	return resp, nil
}

var (
	// ErrViewExists bundles CREATE error.
	ErrViewExists = errors.New("DNS view already exists")
)

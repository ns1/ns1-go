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

// Get takes a DNS view name and returns DNSView struct.
//
// NS1 API docs: https://ns1.com/api#getview-dns-view-details
func (s *DNSViewService) Get(view_name string) (*dns.DNSView, *http.Response, error) {
	path := fmt.Sprintf("views/%s", view_name)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var v dns.DNSView
	resp, err := s.client.Do(req, &v)
	if err != nil {
		switch errType := err.(type) {
		case *Error:
			if errType.Resp.StatusCode == 404 {
				return nil, resp, ErrViewMissing
			}
		}
		return nil, resp, err
	}

	return &v, resp, nil
}

// Update takes a *dns.DNSView and updates the DNS view with same name on NS1.
//
// NS1 API docs: https://ns1.com/api#postedit-a-dns-view
func (s *DNSViewService) Update(v *dns.DNSView) (*http.Response, error) {
	path := fmt.Sprintf("views/%s", v.Name)

	req, err := s.client.NewRequest("POST", path, &v)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, &v)
	if err != nil {
		switch errType := err.(type) {
		case *Error:
			if errType.Resp.StatusCode == 404 {
				return resp, ErrViewMissing
			}
		}
		return resp, err
	}

	return resp, nil
}

// Delete takes a DNS view name, and removes an existing DNS view
//
// NS1 API docs: https://ns1.com/api#deletedelete-a-dns-view
func (s *DNSViewService) Delete(view_name string) (*http.Response, error) {
	path := fmt.Sprintf("views/%s", view_name)

	req, err := s.client.NewRequest("DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		switch errType := err.(type) {
		case *Error:
			if errType.Resp.StatusCode == 404 {
				return resp, ErrViewMissing
			}
		}
		return resp, err
	}

	return resp, nil
}

var (
	// ErrViewExists bundles CREATE error.
	ErrViewExists = errors.New("DNS view already exists")

	// ErrViewExists bundles GET error.
	ErrViewMissing = errors.New("DNS view not found")
)

package rest

import (
	"errors"
	"fmt"
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

// Get takes a TSIG key name and returns a single TSIG key and its basic configuration details.
//
// NS1 API docs: https://ns1.com/api/#getview-tsig-key-details
func (s *TsigService) Get(name string) (*dns.Tsig_key, *http.Response, error) {
	path := fmt.Sprintf("tsig/%s", name)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var tk dns.Tsig_key
	var resp *http.Response
	resp, err = s.client.Do(req, &tk)

	if err != nil {
		switch errType := err.(type) {
		case *Error:
			if errType.Message == fmt.Sprintf("TSIG %s. was not found", name) {
				return nil, resp, ErrTsigKeyMissing
			}
		}
		return nil, resp, err
	}

	return &tk, resp, nil
}

// Create takes a *Tsig_key and creates a new TSIG key.
//
// NS1 API docs: https://ns1.com/api/#putcreate-a-tsig-key
func (s *TsigService) Create(tk *dns.Tsig_key) (*http.Response, error) {
	path := fmt.Sprintf("tsig/%s", tk.Name)

	req, err := s.client.NewRequest("PUT", path, &tk)
	if err != nil {
		return nil, err
	}

	// Update TSIG key fields with data from api(ensure consistent)
	resp, err := s.client.Do(req, &tk)
	if err != nil {
		switch errType := err.(type) {
		case *Error:
			if errType.Message == fmt.Sprintf("TSIG %s. already exists", tk.Name) {
				return resp, ErrTsigKeyExists
			}
		}
		return resp, err
	}

	return resp, nil
}

// Update takes a *Tsug_key and modifies basic details of a TSIG key.
//
// NS1 API docs: https://ns1.com/api/#postmodify-a-tsig-key
func (s *TsigService) Update(tk *dns.Tsig_key) (*http.Response, error) {
	path := fmt.Sprintf("tsig/%s", tk.Name)

	req, err := s.client.NewRequest("POST", path, &tk)
	if err != nil {
		return nil, err
	}

	// Update TSIG key fields with data from api(ensure consistent)
	resp, err := s.client.Do(req, &tk)
	if err != nil {
		switch errType := err.(type) {
		case *Error:
			if errType.Message == fmt.Sprintf("TSIG %s. was not found", tk.Name) {
				return resp, ErrTsigKeyMissing
			}
		}
		return resp, err
	}

	return resp, nil
}

var (
	// ErrTsigKeyExists bundles PUT create error.
	ErrTsigKeyExists = errors.New("TSIG key already exists")
	// ErrTsigKeyMissing bundles GET/POST/DELETE error.
	ErrTsigKeyMissing = errors.New("TSIG key does not exist")
)

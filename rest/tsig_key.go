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

var (
	// ErrTsigKeyMissing bundles GET/POST/DELETE error.
	ErrTsigKeyMissing = errors.New("TSIG key does not exist")
)

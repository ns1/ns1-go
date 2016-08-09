package rest

import (
	"fmt"

	"github.com/ns1/ns1-go/model/dns"
)

// RecordsService handles 'zones/ZONE/DOMAIN/TYPE' endpoint.
type RecordsService service

// Get takes a zone, domain and record type t and returns full configuration for a DNS record.
//
// NS1 API docs: https://ns1.com/api/#record-get
func (s *RecordsService) Get(zone string, domain string, t string) (*dns.Record, *Response, error) {
	path := fmt.Sprintf("zones/%s/%s/%s", zone, domain, t)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var r dns.Record
	resp, err = s.client.Do(req, &r)
	if err != nil {
		return nil, resp, err
	}

	return &r, resp, nil
}

// Create takes a *Record and creates a new DNS record in the specified zone, for the specified domain, of the given record type.
//
// NS1 API docs: https://ns1.com/api/#record-put
func (s *RecordsService) Create(r *dns.Record) (*Response, error) {
	path := fmt.Sprintf("zones/%s/%s/%s", r.Zone, r.Domain, r.Type)

	req, err := s.client.NewRequest("PUT", path, &r)
	if err != nil {
		return nil, err
	}

	// Update record fields with data from api(ensure consistent)
	resp, err = s.client.Do(req, &r)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// UpdateRecord takes a *Record and modifies configuration details for an existing DNS record.
//
// NS1 API docs: https://ns1.com/api/#record-post
func (s *RecordsService) UpdateRecord(r *dns.Record) (*Response, error) {
	path := fmt.Sprintf("zones/%s/%s/%s", r.Zone, r.Domain, r.Type)

	req, err := s.client.NewRequest("POST", path, &r)
	if err != nil {
		return nil, err
	}

	// Update records fields with data from api(ensure consistent)
	resp, err = s.client.Do(req, &r)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// DeleteRecord takes a zone, domain and record type t and removes an existing record and all associated answers and configuration details.
//
// NS1 API docs: https://ns1.com/api/#record-delete
func (s *RecordsService) DeleteRecord(zone string, domain string, t string) (*Response, error) {
	path := fmt.Sprintf("zones/%s/%s/%s", zone, domain, t)

	req, err := s.client.NewRequest("DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err = s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

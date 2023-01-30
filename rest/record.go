package rest

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"gopkg.in/ns1/ns1-go.v2/rest/model/dns"
)

// RecordsService handles 'zones/ZONE/DOMAIN/TYPE' endpoint.
type RecordsService service

// Get takes a zone, domain and record type t and returns full configuration for a DNS record.
//
// NS1 API docs: https://ns1.com/api/#record-get
func (s *RecordsService) Get(ctx context.Context, zone, domain, t string) (*dns.Record, *http.Response, error) {
	path := fmt.Sprintf("zones/%s/%s/%s", zone, domain, t)

	req, err := s.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var r dns.Record
	resp, err := s.client.Do(req, &r)
	if err != nil {
		switch err.(type) {
		case *Error:
			if err.(*Error).Message == "record not found" {
				return nil, resp, ErrRecordMissing
			}
		}
		return nil, resp, err
	}

	return &r, resp, nil
}

// Create takes a *Record and creates a new DNS record in the specified zone, for the specified domain, of the given record type.
//
// The given record must have at least one answer.
// NS1 API docs: https://ns1.com/api/#record-put
func (s *RecordsService) Create(ctx context.Context, r *dns.Record) (*http.Response, error) {
	path := fmt.Sprintf("zones/%s/%s/%s", r.Zone, r.Domain, r.Type)

	req, err := s.client.NewRequest(ctx, "PUT", path, &r)
	if err != nil {
		return nil, err
	}

	// Update record fields with data from api(ensure consistent)
	resp, err := s.client.Do(req, &r)
	if err != nil {
		switch err.(type) {
		case *Error:
			switch err.(*Error).Message {
			case "zone not found":
				return resp, ErrZoneMissing
			case "record already exists":
				return resp, ErrRecordExists
			}
		}
		return resp, err
	}

	return resp, nil
}

// Update takes a *Record and modifies configuration details for an existing DNS record.
//
// Only the fields to be updated are required in the given record.
// NS1 API docs: https://ns1.com/api/#record-post
func (s *RecordsService) Update(ctx context.Context, r *dns.Record) (*http.Response, error) {
	path := fmt.Sprintf("zones/%s/%s/%s", r.Zone, r.Domain, r.Type)

	req, err := s.client.NewRequest(ctx, "POST", path, &r)
	if err != nil {
		return nil, err
	}

	// Update records fields with data from api(ensure consistent)
	resp, err := s.client.Do(req, &r)
	if err != nil {
		switch err.(type) {
		case *Error:
			switch err.(*Error).Message {
			case "zone not found":
				return resp, ErrZoneMissing
			case "record not found":
				return resp, ErrRecordMissing
			case "record already exists":
				return resp, ErrRecordExists
			}
		}
		return resp, err
	}

	return resp, nil
}

// Delete takes a zone, domain and record type t and removes an existing record and all associated answers and configuration details.
//
// NS1 API docs: https://ns1.com/api/#record-delete
func (s *RecordsService) Delete(ctx context.Context, zone string, domain string, t string) (*http.Response, error) {
	path := fmt.Sprintf("zones/%s/%s/%s", zone, domain, t)

	req, err := s.client.NewRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		switch err.(type) {
		case *Error:
			if err.(*Error).Message == "record not found" {
				return resp, ErrRecordMissing
			}
		}
		return resp, err
	}

	return resp, nil
}

var (
	// ErrRecordExists bundles PUT create error.
	ErrRecordExists = errors.New("record already exists")
	// ErrRecordMissing bundles GET/POST/DELETE error.
	ErrRecordMissing = errors.New("record does not exist")
)

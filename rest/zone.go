package rest

import (
	"fmt"

	"github.com/ns1/ns1-go/dns"
)

const (
	zonePath = "zones"
)

// ZonesService handles 'zones' endpoint.
type ZonesService service

// List returns all active zones and basic zone configuration details for each.
//
// NS1 API docs: https://ns1.com/api/#zones-get
func (s *ZonesService) List() ([]*dns.Zone, error) {
	req, err := s.client.NewRequest("GET", zonePath, nil)
	if err != nil {
		return nil, err
	}

	zl := []*dns.Zone{}
	_, err = s.client.Do(req, &zl)
	if err != nil {
		return nil, err
	}

	return zl, nil
}

// Get takes a zone name and returns a single active zone and its basic configuration details.
//
// NS1 API docs: https://ns1.com/api/#zones-zone-get
func (s *ZonesService) Get(zone string) (*dns.Zone, error) {
	path := fmt.Sprintf("%s/%s", zonePath, zone)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var z dns.Zone
	_, err = s.client.Do(req, &z)
	if err != nil {
		return nil, err
	}

	return &z, nil
}

// Create takes a *Zone and creates a new DNS zone.
//
// NS1 API docs: https://ns1.com/api/#zones-put
func (s *ZonesService) Create(z *dns.Zone) error {
	path := fmt.Sprintf("%s/%s", zonePath, z.Zone)

	req, err := s.client.NewRequest("PUT", path, &z)
	if err != nil {
		return err
	}

	// Update zones fields with data from api(ensure consistent)
	_, err = s.client.Do(req, &z)
	if err != nil {
		return err
	}

	return nil
}

// Update takes a *Zone and modifies basic details of a DNS zone.
//
// NS1 API docs: https://ns1.com/api/#zones-post
func (s *ZonesService) Update(z *dns.Zone) error {
	path := fmt.Sprintf("%s/%s", zonePath, z.Zone)

	req, err := s.client.NewRequest("POST", path, &z)
	if err != nil {
		return err
	}

	// Update zones fields with data from api(ensure consistent)
	_, err = s.client.Do(req, &z)
	if err != nil {
		return err
	}

	return nil
}

// Delete takes a zone and destroys an existing DNS zone and all records in the zone.
//
// NS1 API docs: https://ns1.com/api/#zones-delete
func (s *ZonesService) Delete(zone string) error {
	path := fmt.Sprintf("%s/%s", zonePath, zone)

	req, err := s.client.NewRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(req, nil)
	if err != nil {
		return err
	}

	return nil
}

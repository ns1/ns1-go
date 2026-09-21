package rest

import (
	"net/http"

	"gopkg.in/ns1/ns1-go.v2/rest/model/account"
)

// WarningsService handles 'account/usagewarnings' endpoint.
//
// Deprecated: This service is deprecated and will be removed in a future release.
type WarningsService service

// Get returns toggles and thresholds used when sending overage warning
// alert messages to users with billing notifications enabled.
//
// Deprecated: This method is deprecated and will be removed in a future release.
func (s *WarningsService) Get() (*account.UsageWarning, *http.Response, error) {
	req, err := s.client.NewRequest("GET", "account/usagewarnings", nil)
	if err != nil {
		return nil, nil, err
	}

	var uw account.UsageWarning
	resp, err := s.client.Do(req, &uw)
	if err != nil {
		return nil, resp, err
	}

	return &uw, resp, nil
}

// Update changes alerting toggles and thresholds for overage warning alert messages.
//
// Deprecated: This method is deprecated and will be removed in a future release.
func (s *WarningsService) Update(uw *account.UsageWarning) (*http.Response, error) {
	req, err := s.client.NewRequest("POST", "account/usagewarnings", &uw)
	if err != nil {
		return nil, err
	}

	// Update usagewarnings fields with data from api(ensure consistent)
	resp, err := s.client.Do(req, &uw)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

package rest

import (
	"context"
	"net/http"

	"gopkg.in/ns1/ns1-go.v2/rest/model/account"
)

// WarningsService handles 'account/usagewarnings' endpoint.
type WarningsService service

// Get returns toggles and thresholds used when sending overage warning
// alert messages to users with billing notifications enabled.
//
// NS1 API docs: https://ns1.com/api/#usagewarnings-get
func (s *WarningsService) Get() (*account.UsageWarning, *http.Response, error) {
	return s.GetWithContext(context.Background())
}

// GetWithContext is the same as Get, but takes a context.
func (s *WarningsService) GetWithContext(ctx context.Context) (*account.UsageWarning, *http.Response, error) {
	req, err := s.client.NewRequest("GET", "account/usagewarnings", nil)
	if err != nil {
		return nil, nil, err
	}
	req = req.WithContext(ctx)

	var uw account.UsageWarning
	resp, err := s.client.Do(req, &uw)
	if err != nil {
		return nil, resp, err
	}

	return &uw, resp, nil
}

// Update changes alerting toggles and thresholds for overage warning alert messages.
//
// NS1 API docs: https://ns1.com/api/#usagewarnings-post
func (s *WarningsService) Update(uw *account.UsageWarning) (*http.Response, error) {
	return s.UpdateWithContext(context.Background(), uw)
}

// UpdateWithContext is the same as Update, but takes a context.
func (s *WarningsService) UpdateWithContext(ctx context.Context, uw *account.UsageWarning) (*http.Response, error) {
	req, err := s.client.NewRequest("POST", "account/usagewarnings", &uw)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	// Update usagewarnings fields with data from api(ensure consistent)
	resp, err := s.client.Do(req, &uw)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

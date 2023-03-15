package rest

import (
	"context"
	"fmt"
	"net/http"
)

const statsQPSEndpoint = "stats/qps"

// StatsService handles 'stats/qps' endpoint.
type StatsService service

// GetQPS returns current queries per second (QPS) for the account.
// The QPS number is lagged by approximately 30 seconds for statistics collection;
// and the rate is computed over the preceding minute.
func (s *StatsService) GetQPS(ctx context.Context) (float32, *http.Response, error) {
	return s.getQPS(ctx, statsQPSEndpoint)
}

// GetZoneQPS returns current queries per second (QPS) for a specific zone.
// The QPS number is lagged by approximately 30 seconds for statistics collection;
// and the rate is computed over the preceding minute.
func (s *StatsService) GetZoneQPS(ctx context.Context, zone string) (float32, *http.Response, error) {
	path := fmt.Sprintf("%s/%s", statsQPSEndpoint, zone)
	return s.getQPS(ctx, path)
}

// GetRecordQPS returns current queries per second (QPS) for a specific record.
// The QPS number is lagged by approximately 30 seconds for statistics collection;
// and the rate is computed over the preceding minute.
func (s *StatsService) GetRecordQPS(ctx context.Context, zone, record, t string) (float32, *http.Response, error) {
	path := fmt.Sprintf("%s/%s/%s/%s", statsQPSEndpoint, zone, record, t)
	return s.getQPS(ctx, path)
}

func (s *StatsService) getQPS(ctx context.Context, path string) (float32, *http.Response, error) {
	req, err := s.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return 0, nil, err
	}

	var value struct {
		/* by default unmartial will ignore any extra fields so we don't need these
		Networks []struct {
			Network int
			Qps float32
		}
		*/
		Qps float32
	}
	resp, err := s.client.Do(req, &value)

	if err != nil {
		switch err.(type) {
		case *Error:
			switch err.(*Error).Message {
			case "zone not found":
				return 0, nil, ErrZoneMissing
			case "record not found":
				return 0, nil, ErrRecordMissing
			}
		}
		return 0, resp, err
	}
	return value.Qps, resp, nil
}

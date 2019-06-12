package rest

import (
	"errors"
	"fmt"
	"net/http"
)

const statsQpsEndpoint = "stats/qps"

// StatsService handles 'stats/qps' endpoint.
type StatsService service

//Returns current queries per second (QPS) for the account.
// The QPS number is lagged by approximately 30 seconds for statistics collection;
// and the rate is computed over the preceding minute.
func (s StatsService) GetQps() (float32, *http.Response, error) {
	return s.getQps(statsQpsEndpoint)
}

// Returns current queries per second (QPS) for a specific zone.
// The QPS number is lagged by approximately 30 seconds for statistics collection;
// and the rate is computed over the preceding minute.
func (s StatsService) GetZoneQps(zone string) (float32, *http.Response, error) {
	path := fmt.Sprintf("%s/%s", statsQpsEndpoint, zone)
	return s.getQps(path)
}

// Returns current queries per second (QPS) for a specific record.
// The QPS number is lagged by approximately 30 seconds for statistics collection;
// and the rate is computed over the preceding minute.
func (s StatsService) GetRecordQps(zone, record, t string) (float32, *http.Response, error) {
	path := fmt.Sprintf("%s/%s/%s/%s", statsQpsEndpoint, zone, record, t)
	return s.getQps(path)
}

func (s StatsService) getQps(path string) (float32, *http.Response, error) {
	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return 0, nil, err
	}

	var r map[string]float32
	resp, err := s.client.Do(req, &r)
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
		return 0, nil, err
	}

	qps, ok := r["qps"]
	if !ok {
		return 0, nil, errors.New(fmt.Sprintf("Could not find 'qps' key in returned data: %v", resp))
	}
	return qps, resp, nil
}

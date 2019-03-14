package rest

import (
	"errors"
	"fmt"
)

const statsQpsEndpoint = "stats/qps"

// StatsService handles 'stats/qps' endpoint.
type StatsService service

//Returns current queries per second (QPS) for the account.
// The QPS number is lagged by approximately 30 seconds for statistics collection;
// and the rate is computed over the preceding minute.
func (s StatsService) GetQps() (float32, error) {
	return s.getQps(statsQpsEndpoint)
}

// Returns current queries per second (QPS) for a specific zone.
// The QPS number is lagged by approximately 30 seconds for statistics collection;
// and the rate is computed over the preceding minute.
func (s StatsService) GetZoneQps(zone string) (float32, error) {
	path := fmt.Sprintf("%s/%s", statsQpsEndpoint, zone)
	return s.getQps(path)
}

// Returns current queries per second (QPS) for a specific record.
// The QPS number is lagged by approximately 30 seconds for statistics collection;
// and the rate is computed over the preceding minute.
func (s StatsService) GetRecordQps(zone, record, t string) (float32, error) {
	path := fmt.Sprintf("%s/%s/%s/%s", statsQpsEndpoint, zone, record, t)
	return s.getQps(path)
}

func (s StatsService) getQps(path string) (float32, error) {
	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return 0, err
	}

	var r map[string]float32
	resp, err := s.client.Do(req, &r)
	if err != nil {
		switch err.(type) {
		case *Error:
			switch err.(*Error).Message {
			case "zone not found":
				return 0, ErrZoneMissing
			}
		}
		return 0, err
	}

	qps, ok := r["qps"]
	if !ok {
		return 0, errors.New(fmt.Sprintf("Could not find 'qps' key in returned data: %v", resp))
	}
	return qps, nil

}

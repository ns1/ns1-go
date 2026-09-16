package rest

import (
	"fmt"
	"net/http"
)

const statsQPSEndpoint = "stats/qps"

// StatsService handles 'stats/qps' endpoint.
type StatsService service

// GetQPS returns current queries per second (QPS) for the account.
// The QPS number is lagged by approximately 30 seconds for statistics collection;
// and the rate is computed over the preceding minute.
func (s *StatsService) GetQPS() (float32, *http.Response, error) {
	return s.getQPS(statsQPSEndpoint)
}

// GetZoneQPS returns current queries per second (QPS) for a specific zone.
// The QPS number is lagged by approximately 30 seconds for statistics collection;
// and the rate is computed over the preceding minute.
// NS1 API docs: https://developer.ibm.com/apis/catalog/ns1--ibm-ns1-connect-api/api/API--ns1--ibm-ns1-connect-api#getQpsByZoneName
func (s *StatsService) GetZoneQPS(zone string) (float32, *http.Response, error) {
	path := fmt.Sprintf("%s/%s", statsQPSEndpoint, zone)
	return s.getQPS(path)
}

// GetRecordQPS returns current queries per second (QPS) for a specific record.
// The QPS number is lagged by approximately 30 seconds for statistics collection;
// and the rate is computed over the preceding minute.
// NS1 API docs: https://developer.ibm.com/apis/catalog/ns1--ibm-ns1-connect-api/api/API--ns1--ibm-ns1-connect-api#getQpsByZoneAndDomainAndRecord
func (s *StatsService) GetRecordQPS(zone, record, t string) (float32, *http.Response, error) {
	path := fmt.Sprintf("%s/%s/%s/%s", statsQPSEndpoint, zone, record, t)
	return s.getQPS(path)
}

// NS1 API docs: https://developer.ibm.com/apis/catalog/ns1--ibm-ns1-connect-api/api/API--ns1--ibm-ns1-connect-api#getQps
func (s *StatsService) getQPS(path string) (float32, *http.Response, error) {
	req, err := s.client.NewRequest("GET", path, nil)
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
		QPS float32
	}
	resp, err := s.client.Do(req, &value)

	if err != nil {
		switch err := err.(type) {
		case *Error:
			switch err.Message {
			case "zone not found":
				return 0, nil, ErrZoneMissing
			case "record not found":
				return 0, nil, ErrRecordMissing
			}
		}
		return 0, resp, err
	}
	return value.QPS, resp, nil
}

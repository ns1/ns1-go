package mockns1

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	api "gopkg.in/ns1/ns1-go.v2/rest"
)

// Service is a controller for a mock server suitable for responding to
// gopkg.in/ns1/ns1-go.v2 client requests. This object should always be
// initialized via the New() function.
type Service struct {
	// Address is set by New() to the listen address of the mock server
	Address string

	server *httptest.Server
	tests  map[string]map[string][]*testCase // method, uri
	tb     testing.TB
}

// New creates and starts a new TLS based *httptest.Server instance. As a
// self signed certificate is being used by the server, your HTTP client
// being used with your NS1 client needs to disable TLS verification. The
// returned Doer is a *http.Client with TLS verification disable that
// is appropriate for use.
//
// The mock service is compatible with both testing.B and testing.T instances
// so that it may be used in benchmarking as well as testing. If a
// testing.B instance is supplied the timer will be stopped for the
// duration of the ServeHTTP() and AddTestCase() methods to minimise
// the impact on your benchmark statistics.
func New(tb testing.TB) (*Service, api.Doer, error) {
	s := &Service{
		tb:    tb,
		tests: map[string]map[string][]*testCase{},
	}

	hc := &http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		}},
	}

	s.server = httptest.NewTLSServer(s)

	u, err := url.Parse(s.server.URL)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to parse server URL for address: %s", err)
	}
	s.Address = u.Host

	return s, hc, nil
}

// Shutdown cleans up the mock instance
func (s *Service) Shutdown() {
	s.server.Close()
}

func (s *Service) startTimer() {
	b, ok := s.tb.(*testing.B)
	if !ok {
		return
	}

	b.StartTimer()
}

func (s *Service) stopTimer() {
	b, ok := s.tb.(*testing.B)
	if !ok {
		return
	}

	b.StopTimer()
}

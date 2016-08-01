package nsone

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	clientVersion    = "0.9.0"
	defaultEndpoint  = "https://api.nsone.net/v1/"
	defaultUserAgent = "Golang client-" + clientVersion
)

// Doer is a single method interface that allows a user to extend/augment an http.Client instance.
// Note: http.Client satisfies the Doer interface.
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

// APIClient stores NS1 client state
type APIClient struct {
	// client handles all http communication.
	// The default value is *http.Client
	client Doer

	// NS1 rest endpoint, overrides default if given.
	Endpoint *url.URL

	// NS1 api key (value for http request header 'X-NSONE-Key').
	ApiKey string

	// NS1 go rest user agent (value for http request header 'User-Agent').
	UserAgent string

	// Rate limiting strategy for the APIClient instance.
	RateLimitFunc func(RateLimit)

	// Enables verbose logs.
	debug bool
}

// New takes an API Key and creates an *APIClient
func New(k string) *APIClient {
	endpoint, _ := url.Parse(defaultEndpoint)
	return &APIClient{
		client:        http.DefaultClient,
		Endpoint:      endpoint,
		ApiKey:        k,
		RateLimitFunc: defaultRateLimitFunc,
		UserAgent:     defaultUserAgent,
	}
}

func NewAPIClient(httpClient Doer, options ...APIClientOption) *APIClient {
	endpoint, _ := url.Parse(defaultEndpoint)

	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	c := &APIClient{
		client:        httpClient,
		Endpoint:      endpoint,
		RateLimitFunc: defaultRateLimitFunc,
		UserAgent:     defaultUserAgent,
	}

	for _, option := range options {
		option(c)
	}
	return c
}

// Debug enables debug logging
func (c *APIClient) Debug() {
	c.debug = true
}

type APIClientOption func(*APIClient)

func SetClient(client Doer) APIClientOption {
	return func(c *APIClient) { c.client = client }
}

func SetApiKey(key string) APIClientOption {
	return func(c *APIClient) { c.ApiKey = key }
}

func SetEndpoint(endpoint string) APIClientOption {
	return func(c *APIClient) { c.Endpoint, _ = url.Parse(endpoint) }
}

func SetUserAgent(ua string) APIClientOption {
	return func(c *APIClient) { c.UserAgent = ua }
}

// Contains all http responses outside the 2xx range.
type RestError struct {
	Resp    *http.Response
	Message string
}

// Satisfy std lib error interface.
func (re *RestError) Error() string {
	return fmt.Sprintf("%v %v: %d %v", re.Resp.Request.Method, re.Resp.Request.URL, re.Resp.StatusCode, re.Message)
}

// Handles parsing of rest api errors. Returns nil if no error.
func CheckResponse(resp *http.Response) error {
	if c := resp.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	restError := &RestError{Resp: resp}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if len(b) == 0 {
		return restError
	}

	json.Unmarshal(b, restError)
	return restError
}

func (c APIClient) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = CheckResponse(resp)
	if err != nil {
		return resp, err
	}

	// Update the clients' rate limit.
	if len(resp.Header["X-Ratelimit-Limit"]) > 0 {
		var remaining int
		var period int
		limit, err := strconv.Atoi(resp.Header["X-Ratelimit-Limit"][0])
		if err == nil {
			remaining, err = strconv.Atoi(resp.Header["X-Ratelimit-Remaining"][0])
			if err == nil {
				period, err = strconv.Atoi(resp.Header["X-Ratelimit-Period"][0])
			}
		}
		if err == nil {
			c.RateLimitFunc(RateLimit{
				Limit:     limit,
				Remaining: remaining,
				Period:    period,
			})
		}
	}

	if v != nil {
		// Try to decode body into the given type.
		err := json.NewDecoder(resp.Body).Decode(&v)
		if err != nil {
			return nil, err
		}
	}

	return resp, err
}

func (c *APIClient) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	uri := c.Endpoint.ResolveReference(rel)

	if c.debug {
		log.Printf("[DEBUG] %s: %s (%s)", method, uri.String(), body)
	}

	// Encode body as json
	buf := new(bytes.Buffer)
	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, uri.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-NSONE-Key", c.ApiKey)
	req.Header.Add("User-Agent", c.UserAgent)
	return req, nil
}

// RateLimit stores X-Ratelimit-* headers
type RateLimit struct {
	Limit     int
	Remaining int
	Period    int
}

var defaultRateLimitFunc = func(rl RateLimit) {}

// PercentageLeft returns the ratio of Remaining to Limit as a percentage
func (rl RateLimit) PercentageLeft() int {
	return rl.Remaining * 100 / rl.Limit
}

// WaitTime returns the time.Duration ratio of Period to Limit
func (rl RateLimit) WaitTime() time.Duration {
	return (time.Second * time.Duration(rl.Period)) / time.Duration(rl.Limit)
}

// WaitTimeRemaining returns the time.Duration ratio of Period to Remaining
func (rl RateLimit) WaitTimeRemaining() time.Duration {
	return (time.Second * time.Duration(rl.Period)) / time.Duration(rl.Remaining)
}

// RateLimitStrategyNone sets RateLimitFunc to an empty func
func (c *APIClient) RateLimitStrategyNone() {
	c.RateLimitFunc = defaultRateLimitFunc
}

// RateLimitStrategySleep sets RateLimitFunc to sleep by WaitTimeRemaining
func (c *APIClient) RateLimitStrategySleep() {
	c.RateLimitFunc = func(rl RateLimit) {
		remaining := rl.WaitTimeRemaining()
		if c.debug {
			log.Printf("Rate limiting - Limit %d Remaining %d in period %d: Sleeping %dns", rl.Limit, rl.Remaining, rl.Period, remaining)
		}
		time.Sleep(remaining)
	}
}

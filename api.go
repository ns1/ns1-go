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
	defaultEndpoint = "https://api.nsone.net/v1/"
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
	client *Doer

	// NS1 rest endpoint, overrides default if given.
	Endpoint *url.URL

	// NS1 api key (value for http request header 'X-NSONE-Key')
	ApiKey string

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
	}
}

// Debug enables debug logging
func (c *APIClient) Debug() {
	c.debug = true
}

func (c APIClient) doHTTP(method string, uri string, rbody []byte) ([]byte, int, error) {
	var body []byte
	r := bytes.NewReader(rbody)
	if c.debug {
		log.Printf("[DEBUG] %s: %s (%s)", method, uri, string(rbody))
	}
	req, err := http.NewRequest(method, uri, r)
	if err != nil {
		return body, 510, err
	}
	req.Header.Add("X-NSONE-Key", c.ApiKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return body, 510, err
	}
	if c.debug {
		log.Println(resp)
	}
	body, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
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
	if resp.StatusCode != 200 {
		return body, resp.StatusCode, fmt.Errorf("%s: %s", resp.Status, string(body))
	}
	if c.debug {
		log.Println(fmt.Sprintf("Response body: %s", string(body)))
	}
	return body, resp.StatusCode, nil

}

func (c APIClient) doHTTPUnmarshal(method string, uri string, rbody []byte, unpackInto interface{}) (int, error) {
	body, status, err := c.doHTTP(method, uri, rbody)
	if err != nil {
		return status, err
	}
	return status, json.Unmarshal(body, unpackInto)
}

func (c APIClient) doHTTPBoth(method string, uri string, s interface{}) error {
	rbody, err := json.Marshal(s)
	if err != nil {
		return err
	}
	_, err = c.doHTTPUnmarshal(method, uri, rbody, s)
	return err
}

func (c APIClient) doHTTPDelete(uri string) error {
	_, _, err := c.doHTTP("DELETE", uri, nil)
	return err
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

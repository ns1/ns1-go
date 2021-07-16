package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_RateLimit(t *testing.T) {
	r := RateLimit{
		Limit:     10,
		Remaining: 10,
		Period:    10,
	}
	if r.PercentageLeft() != 100 {
		t.Error("PercentLeft != 100")
	}
	if r.WaitTime() != time.Second {
		t.Error("WaitTime is wrong duration ", r.WaitTime())
	}
	if r.WaitTimeRemaining() != (time.Duration(1) * time.Second) {
		t.Error("WaitTimeRemaining is wrong duration ", r.WaitTimeRemaining())
	}

	r.Remaining = 5
	if r.PercentageLeft() != 50 {
		t.Error("PercentLeft != 50")
	}
	if r.WaitTime() != time.Second {
		t.Error("WaitTime is wrong duration ", r.WaitTime())
	}
	if r.WaitTimeRemaining() != (time.Duration(2) * time.Second) {
		t.Error("WaitTimeRemaining is wrong duration ", r.WaitTimeRemaining())
	}

	r.Remaining = 0
	if r.PercentageLeft() != 0 {
		t.Error("PercentLeft != 0")
	}
	if r.WaitTime() != time.Second {
		t.Error("WaitTime is wrong duration ", r.WaitTime())
	}
	if r.WaitTimeRemaining() != (time.Duration(10) * time.Second) {
		t.Error("WaitTimeRemaining is wrong duration ", r.WaitTimeRemaining())
	}
}

func TestClient_Do(t *testing.T) {
	// It should return the response without error
	httpClient := mockHTTPClient{}
	client := NewClient(&httpClient, SetEndpoint(""))
	req, _ := http.NewRequest("GET", "http://example.com", new(bytes.Buffer))

	mockResp := http.Response{
		Body:       ioutil.NopCloser(bytes.NewBufferString("{}")),
		StatusCode: 200,
	}
	httpClient.On("Do", req).Return(&mockResp, nil)

	resp, err := client.Do(req, "")

	httpClient.AssertExpectations(t)

	assert.Equal(t, &mockResp, resp)
	assert.Nil(t, err)
}

func TestClient_DoWithHTTPClientError(t *testing.T) {
	// It should return nil response and the error
	httpClient := mockHTTPClient{}
	client := NewClient(&httpClient, SetEndpoint(""))
	req, _ := http.NewRequest("GET", "http://example.com", new(bytes.Buffer))

	mockError := errors.New("Some Error")
	httpClient.On("Do", req).Return(nil, mockError)

	resp, err := client.Do(req, "")

	httpClient.AssertExpectations(t)

	assert.Nil(t, resp)
	assert.Equal(t, mockError, err)
}

func TestClient_DoWithHTTPClientErrorRateLimit(t *testing.T) {
	// It should return nil response and the error, also should run the rate limit func
	httpClient := mockHTTPClient{}
	client := NewClient(&httpClient, SetEndpoint(""))
	var rateLimit *RateLimit = nil
	client.RateLimitFunc = func(r RateLimit) {
		rateLimit = &r
	}
	req, _ := http.NewRequest("GET", "http://example.com", new(bytes.Buffer))

	headerResponse := make(http.Header)
	headerResponse.Add(headerRateLimit, "10")
	headerResponse.Add(headerRateRemaining, "10")
	headerResponse.Add(headerRatePeriod, "10")
	mockResp := http.Response{
		Body:       ioutil.NopCloser(bytes.NewBufferString("{}")),
		Header:     headerResponse,
		StatusCode: 400,
	}
	httpClient.On("Do", req).Return(&mockResp, nil)

	client.Do(req, "")

	httpClient.AssertExpectations(t)

	assert.NotNil(t, rateLimit)
}

func TestClient_DoWithNon2XXResponse(t *testing.T) {
	// It should return a pointer to the response, and a pointer to Error (with
	// the response)
	httpClient := mockHTTPClient{}
	client := NewClient(&httpClient, SetEndpoint(""))
	req, _ := http.NewRequest("GET", "http://example.com", new(bytes.Buffer))

	mockResp := http.Response{
		Body:       ioutil.NopCloser(bytes.NewBufferString("")),
		StatusCode: 404,
	}
	httpClient.On("Do", req).Return(&mockResp, nil)

	resp, err := client.Do(req, "")

	httpClient.AssertExpectations(t)

	assert.Equal(t, &mockResp, resp)
	assert.Equal(t, &Error{Resp: &mockResp}, err)
}

func TestClient_DoWithNonJSONResponse(t *testing.T) {
	// It should return a nil response, and the error from JSON Decoder
	httpClient := mockHTTPClient{}
	client := NewClient(&httpClient, SetEndpoint(""))
	req, _ := http.NewRequest("GET", "http://example.com", new(bytes.Buffer))

	mockResp := http.Response{
		Body:       ioutil.NopCloser(bytes.NewBufferString("INVALID")),
		StatusCode: 200,
	}
	httpClient.On("Do", req).Return(&mockResp, nil)

	resp, err := client.Do(req, "")

	httpClient.AssertExpectations(t)

	assert.Nil(t, resp)
	assert.IsType(t, &json.SyntaxError{}, errors.Unwrap(err))
}

func TestClient_DoWithPagination(t *testing.T) {
	// It should call nextFunc
	// It should return the last response without error
	// It should not set HTTPS on the Link when endpoint is HTTP
	httpClient := mockHTTPClient{}
	client := NewClient(&httpClient, SetEndpoint("http://"))

	req, _ := http.NewRequest("GET", "http://example.com", new(bytes.Buffer))
	firstResp := http.Response{
		Body:       ioutil.NopCloser(bytes.NewBufferString("{}")),
		Header:     http.Header{"Link": []string{`<http://example.com/?after=1&limit=2>; rel="next"`}},
		StatusCode: 200,
	}
	nextResp := http.Response{
		Body:       ioutil.NopCloser(bytes.NewBufferString("{}")),
		StatusCode: 200,
	}
	var v interface{}

	httpClient.On("Do", req).Return(&firstResp, nil)
	httpClient.On("nextFunc", &v, "http://example.com/?after=1&limit=2").Return(&nextResp, nil)

	resp, err := client.DoWithPagination(req, v, httpClient.nextFunc)

	httpClient.AssertExpectations(t)

	assert.Equal(t, &nextResp, resp)
	assert.Nil(t, err)
}

func TestClient_DoWithPaginationForceHTTPS(t *testing.T) {
	// When the endpoint is HTTPS, it should Force HTTPS when following Link
	// headers
	httpClient := mockHTTPClient{}
	client := NewClient(&httpClient, SetEndpoint("https://"))

	req, _ := http.NewRequest("GET", "https://example.com", new(bytes.Buffer))
	firstResp := http.Response{
		Body:       ioutil.NopCloser(bytes.NewBufferString("{}")),
		Header:     http.Header{"Link": []string{`<http://example.com/?after=1&limit=2>; rel="next"`}},
		StatusCode: 200,
	}
	nextResp := http.Response{
		Body:       ioutil.NopCloser(bytes.NewBufferString("{}")),
		StatusCode: 200,
	}
	var v interface{}

	httpClient.On("Do", req).Return(&firstResp, nil)
	httpClient.On("nextFunc", &v, "https://example.com/?after=1&limit=2").Return(&nextResp, nil)

	resp, err := client.DoWithPagination(req, v, httpClient.nextFunc)

	httpClient.AssertExpectations(t)

	assert.Equal(t, &nextResp, resp)
	assert.Nil(t, err)
}

func TestClient_getURI(t *testing.T) {
	// It should delegate to client.Do
	httpClient := mockHTTPClient{}
	client := NewClient(&httpClient, SetEndpoint(""))

	mockResp := http.Response{
		Body:       ioutil.NopCloser(bytes.NewBufferString("{}")),
		StatusCode: 200,
	}
	httpClient.On("Do", mock.Anything).Return(&mockResp, nil)

	var v interface{}
	resp, err := client.getURI(v, "http://example.com")

	assert.Equal(t, &mockResp, resp)
	assert.Nil(t, err)
}

func TestClient_getURIWithNon2XXResponse(t *testing.T) {
	// It should return a pointer to the response, and a pointer to Error (with
	// the response
	httpClient := mockHTTPClient{}
	client := NewClient(&httpClient, SetEndpoint(""))

	mockResp := http.Response{
		Body:       ioutil.NopCloser(bytes.NewBufferString("{}")),
		StatusCode: 418,
	}
	httpClient.On("Do", mock.Anything).Return(&mockResp, nil)

	var v interface{}
	resp, err := client.getURI(v, "http://example.com")

	assert.Equal(t, &mockResp, resp)
	assert.Equal(t, &Error{Resp: &mockResp}, err)
}

type mockHTTPClient struct {
	mock.Mock
}

func (c *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := c.Called(req)
	if res := args.Get(0); res == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*http.Response), args.Error(1)
}

// Hanging this off of mockHTTPClient for convent access to mocking stuff
func (c *mockHTTPClient) nextFunc(v *interface{}, uri string) (*http.Response, error) {
	args := c.Called(v, uri)
	return args.Get(0).(*http.Response), args.Error(1)
}

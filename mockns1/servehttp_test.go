package mockns1_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/ns1/ns1-go.v2/mockns1"
)

func TestServeHTTP(t *testing.T) {
	t.Run("ServeHTTP", func(t *testing.T) {
		mock, _, err := mockns1.New(t)
		require.Nil(t, err)
		defer mock.Shutdown()

		t.Run("Errors", func(t *testing.T) {
			t.Run("Method", func(t *testing.T) {
				mw := &mockWriter{buf: bytes.NewBufferString("")}
				req := &http.Request{Method: http.MethodOptions}

				mock.ServeHTTP(mw, req)
				require.Equal(t, http.StatusNotFound, mw.status, mw.buf.String())
				require.Equal(t, `{"message": "request not found: method"}`, mw.buf.String())
			})

			t.Run("URI", func(t *testing.T) {
				mw := &mockWriter{buf: bytes.NewBufferString("")}
				req := &http.Request{
					Method:     http.MethodGet,
					RequestURI: "/test",
				}

				require.Nil(t, mock.AddTestCase(
					http.MethodGet, "/test", http.StatusTeapot, nil, nil, "", "",
				))

				mock.ServeHTTP(mw, req)
				require.Equal(t, http.StatusNotFound, mw.status, mw.buf.String())
				require.Equal(t, `{"message": "request not found: uri"}`, mw.buf.String())
			})

			t.Run("Test", func(t *testing.T) {
				mw := &mockWriter{buf: bytes.NewBufferString("")}
				req := &http.Request{
					Method:     http.MethodGet,
					RequestURI: "/v1/bad/body",
					Body:       ioutil.NopCloser(bytes.NewReader([]byte("body"))),
				}

				require.Nil(t, mock.AddTestCase(
					http.MethodGet, "/bad/body", http.StatusTeapot, nil, nil, "", "",
				))

				mock.ServeHTTP(mw, req)
				require.Equal(t, http.StatusNotFound, mw.status, mw.buf.String())
				require.Equal(t, `{"message": "request not found: no test"}`, mw.buf.String())
			})
		})

		t.Run("Success", func(t *testing.T) {
			t.Run("Header", func(t *testing.T) {
				mw := &mockWriter{buf: bytes.NewBufferString("")}
				req := &http.Request{
					Method:     http.MethodGet,
					RequestURI: "/v1/request/header",
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
					Header:     http.Header{},
				}
				req.Header.Set("x-test-header", "test-value")

				hdrs := http.Header{}
				hdrs.Set("x-test-header", "test-value")
				require.Nil(t, mock.AddTestCase(
					http.MethodGet, "/request/header", http.StatusOK, hdrs, nil, "",
					"header match",
				))

				mock.ServeHTTP(mw, req)
				require.Equal(t, http.StatusOK, mw.status, mw.buf.String())
				require.Equal(t, "header match", mw.buf.String())
			})

			t.Run("Body", func(t *testing.T) {
				mw := &mockWriter{buf: bytes.NewBufferString("")}
				req := &http.Request{
					Method:     http.MethodGet,
					RequestURI: "/v1/request/body",
					Body:       ioutil.NopCloser(bytes.NewReader([]byte("body"))),
				}

				require.Nil(t, mock.AddTestCase(
					http.MethodGet, "/request/body", http.StatusOK, nil, nil, "body",
					"body match",
				))

				mock.ServeHTTP(mw, req)
				require.Equal(t, http.StatusOK, mw.status, mw.buf.String())
				require.Equal(t, "body match", mw.buf.String())
			})
		})
	})
}

type mockWriter struct {
	buf     *bytes.Buffer
	headers http.Header
	status  int
}

func (mw *mockWriter) Header() http.Header {
	return mw.headers
}

func (mw *mockWriter) Write(data []byte) (int, error) {
	return mw.buf.Write(data)
}

func (mw *mockWriter) WriteHeader(status int) {
	mw.status = status
}

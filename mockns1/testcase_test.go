package mockns1_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/ns1/ns1-go.v2/mockns1"
)

func TestTestCase(t *testing.T) {
	mock, _, err := mockns1.New(t)
	require.Nil(t, err)
	defer mock.Shutdown()

	t.Run("AddTestCase", func(t *testing.T) {
		t.Run("Error", func(t *testing.T) {
			t.Run("Duplicate", func(t *testing.T) {
				require.Nil(t, mock.AddTestCase(http.MethodGet, "test/header", http.StatusOK,
					nil, nil, "", ""))
				err := mock.AddTestCase(http.MethodGet, "test/header", http.StatusOK,
					nil, nil, "", "")
				require.NotNil(t, err)
				require.Equal(t, "test case already registered", err.Error())
			})
		})

		t.Run("Success", func(t *testing.T) {
			require.Nil(t, mock.AddTestCase(http.MethodGet, "test/case", http.StatusOK,
				nil, nil, "", ""))
		})
	})

	t.Run("ClearTestCases", func(t *testing.T) {
		require.Nil(t, mock.AddTestCase(http.MethodGet, "test/clear", http.StatusOK,
			nil, nil, "", ""))

		mw := &mockWriter{buf: bytes.NewBufferString("")}
		req := &http.Request{
			Method:     http.MethodGet,
			RequestURI: "/v1/test/clear",
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
			Header:     http.Header{},
		}

		mock.ServeHTTP(mw, req)
		require.Equal(t, http.StatusOK, mw.status, mw.buf.String())

		mock.ClearTestCases()
		mock.ServeHTTP(mw, req)
		require.Equal(t, http.StatusNotFound, mw.status, mw.buf.String())
	})
}

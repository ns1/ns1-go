package mockns1_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/ns1/ns1-go.v2/mockns1"
)

func TestTestCase(t *testing.T) {
	t.Run("AddTestCase", func(t *testing.T) {
		mock, _, err := mockns1.New(t)
		require.Nil(t, err)
		defer mock.Shutdown()

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
}

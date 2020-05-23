package mockns1_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/ns1/ns1-go.v2/mockns1"
)

func TestService(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		mock, doer, err := mockns1.New(t)
		require.Nil(t, err)
		require.IsType(t, new(http.Client), doer)
		require.NotNil(t, mock)
		require.NotEmpty(t, mock.Address)
		mock.Shutdown()
	})

	t.Run("Shutdown", func(t *testing.T) {
		mock, _, err := mockns1.New(t)
		require.Nil(t, err)
		require.NotPanics(t, mock.Shutdown)
	})
}

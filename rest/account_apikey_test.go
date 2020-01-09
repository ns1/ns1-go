package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/ns1/ns1-go.v2/rest/model/account"
)

func TestCreateAPIKey(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)

		var k account.APIKey
		require.NoError(t, json.Unmarshal(b, &k))
		assert.Nil(t, k.Permissions.Security)
		assert.Nil(t, k.Permissions.DHCP)
		assert.Nil(t, k.Permissions.IPAM)

		w.Write(b)
	}))
	defer ts.Close()
	c := NewClient(nil, SetEndpoint(ts.URL))

	k := &account.APIKey{
		ID:          "id-1",
		Key:         "key-1",
		Name:        "name-1",
		Permissions: account.PermissionsMap{},
	}

	_, err := c.APIKeys.Create(k)
	require.NoError(t, err)
}

func TestCreateDDIAPIKey(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)

		var k account.APIKey
		require.NoError(t, json.Unmarshal(b, &k))
		assert.NotNil(t, k.Permissions.Security)
		assert.NotNil(t, k.Permissions.DHCP)
		assert.NotNil(t, k.Permissions.IPAM)

		w.Write(b)
	}))
	defer ts.Close()
	c := NewClient(nil, SetEndpoint(ts.URL), SetDDIAPI())

	k := &account.APIKey{
		ID:          "id-1",
		Key:         "key-1",
		Name:        "name-1",
		Permissions: account.PermissionsMap{},
	}

	_, err := c.APIKeys.Create(k)
	require.NoError(t, err)
}

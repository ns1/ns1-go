package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/ns1/ns1-go.v2/common/conv"
	"gopkg.in/ns1/ns1-go.v2/rest/model/account"
)

func TestCreateUserV2(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)

		var u account.UserV2
		require.NoError(t, json.Unmarshal(b, &u))
		assert.Nil(t, u.Permissions)

		w.Write(b)
	}))
	defer ts.Close()
	c := NewClient(nil, SetEndpoint(ts.URL))

	u := &account.UserV2{
		Name:     "name-1",
		Username: "user-1",
		Email:    "email-1",
	}

	_, err := c.UsersV2.Create(u)
	require.NoError(t, err)
}

func TestCreateDDIUserV2(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)

		var u account.UserV2
		require.NoError(t, json.Unmarshal(b, &u))
		assert.NotNil(t, u.Permissions.Security)
		assert.NotNil(t, u.Permissions.DHCP)
		assert.NotNil(t, u.Permissions.IPAM)
		assert.Nil(t, u.Permissions.DNS)
		assert.Nil(t, u.Permissions.Account)
		assert.Nil(t, u.Permissions.Monitoring)
		assert.NotNil(t, u.IPWhitelist)
		assert.True(t, u.IPWhitelistStrict)

		w.Write(b)
	}))
	defer ts.Close()
	c := NewClient(nil, SetEndpoint(ts.URL), SetDDIAPI())
	u := &account.UserV2{
		Name:              "name-1",
		Username:          "user-1",
		Email:             "email-1",
		IPWhitelist:       []string{"1.1.1.1"},
		IPWhitelistStrict: true,
		Permissions: &account.PermissionsMapV2{
			Security: &account.PermissionsSecurityV2{
				ManageGlobal2FA: conv.BoolPtrFrom(true),
			},
			DHCP: &account.PermissionsDHCPV2{
				ManageDHCP: conv.BoolPtrFrom(true),
			},
			IPAM: &account.PermissionsIPAMV2{
				ManageIPAM: conv.BoolPtrFrom(true),
			},
		},
	}

	_, err := c.UsersV2.Create(u)
	require.NoError(t, err)
}

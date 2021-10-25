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

func TestCreateTeamV2(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)

		var tm account.TeamV2
		require.NoError(t, json.Unmarshal(b, &tm))
		assert.Nil(t, tm.Permissions.Security)
		assert.Nil(t, tm.Permissions.DHCP)
		assert.Nil(t, tm.Permissions.IPAM)

		w.Write(b)
	}))
	defer ts.Close()
	c := NewClient(nil, SetEndpoint(ts.URL))

	tm := &account.TeamV2{
		ID:          "id-1",
		Name:        "team-1",
		Permissions: &account.PermissionsMapV2{},
	}

	_, err := c.TeamsV2.Create(tm)
	require.NoError(t, err)
}

func TestCreateDDITeamV2(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)

		var tm account.TeamV2
		require.NoError(t, json.Unmarshal(b, &tm))
		assert.NotNil(t, tm.Permissions.Security)
		assert.NotNil(t, tm.Permissions.DHCP)
		assert.NotNil(t, tm.Permissions.IPAM)
		assert.NotNil(t, tm.IPWhitelist)

		w.Write(b)
	}))
	defer ts.Close()
	c := NewClient(nil, SetEndpoint(ts.URL), SetDDIAPI())

	tm := &account.TeamV2{
		ID:   "id-1",
		Name: "team-1",
		IPWhitelist: []account.IPWhitelist{
			{Name: "whitelist", Values: []string{"1.1.1.1"}},
		},
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

	_, err := c.TeamsV2.Create(tm)
	require.NoError(t, err)
}

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

func TestCreateAPIKeyV2(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)

		var k account.APIKeyV2
		require.NoError(t, json.Unmarshal(b, &k))
		assert.Nil(t, k.Permissions.Security)
		assert.Nil(t, k.Permissions.DHCP)
		assert.Nil(t, k.Permissions.IPAM)

		_, err = w.Write(b)
		require.NoError(t, err)
	}))
	defer ts.Close()
	c := NewClient(nil, SetEndpoint(ts.URL))

	k := &account.APIKeyV2{
		ID:          "id-1",
		Key:         "key-1",
		Name:        "name-1",
		Permissions: &account.PermissionsMapV2{},
	}

	_, err := c.APIKeysV2.Create(k)
	require.NoError(t, err)
}

func TestCreateDDIAPIKeyV2(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)

		var k account.APIKeyV2
		require.NoError(t, json.Unmarshal(b, &k))
		switch k.ID {
		case "ddi-no-authtags":
			assert.NotNil(t, k.Permissions.Security)
			assert.NotNil(t, k.Permissions.DHCP)
			assert.NotNil(t, k.Permissions.IPAM)
			assert.NotNil(t, k.IPWhitelist)
			assert.True(t, k.IPWhitelistStrict)
			// ensure auth tag permissions are not included by default to maintain backwards compatibility
			assert.Nil(t, k.Permissions.DHCP.TagsAllow)
			assert.Nil(t, k.Permissions.DHCP.TagsDeny)
			assert.Nil(t, k.Permissions.IPAM.TagsAllow)
			assert.Nil(t, k.Permissions.IPAM.TagsDeny)
		case "ddi-authtags":
			assert.NotNil(t, k.Permissions.DHCP)
			assert.NotNil(t, k.Permissions.IPAM)
			assert.Equal(t, "auth:dhcpallow", (*k.Permissions.DHCP.TagsAllow)[0].Name)
			assert.Equal(t, "", (*k.Permissions.DHCP.TagsAllow)[0].Value)
			assert.Equal(t, "auth:dhcpdeny", (*k.Permissions.DHCP.TagsDeny)[0].Name)
			assert.Equal(t, "denyme", (*k.Permissions.DHCP.TagsDeny)[0].Value)
			assert.Equal(t, "auth:ipamallow", (*k.Permissions.IPAM.TagsAllow)[0].Name)
			assert.Equal(t, "", (*k.Permissions.IPAM.TagsAllow)[0].Value)
			assert.Equal(t, "auth:ipamdeny", (*k.Permissions.IPAM.TagsDeny)[0].Name)
			assert.Equal(t, "denyme", (*k.Permissions.IPAM.TagsDeny)[0].Value)
		case "ddi-empty-authtags":
			assert.NotNil(t, k.Permissions.DHCP)
			assert.NotNil(t, k.Permissions.IPAM)
			assert.Equal(t, []account.AuthTag{}, *k.Permissions.DHCP.TagsAllow)
			assert.Equal(t, []account.AuthTag{}, *k.Permissions.DHCP.TagsDeny)
			assert.Equal(t, []account.AuthTag{}, *k.Permissions.IPAM.TagsAllow)
			assert.Equal(t, []account.AuthTag{}, *k.Permissions.IPAM.TagsDeny)
		}

		_, err = w.Write(b)
		require.NoError(t, err)
	}))
	defer ts.Close()
	c := NewClient(nil, SetEndpoint(ts.URL), SetDDIAPI())

	// Create a key without auth tags
	k := &account.APIKeyV2{
		ID:                "ddi-no-authtags",
		Key:               "key-1",
		Name:              "name-1",
		IPWhitelist:       []string{"1.1.1.1"},
		IPWhitelistStrict: true,
		Permissions: &account.PermissionsMapV2{
			Security: &account.PermissionsSecurityV2{
				ManageGlobal2FA: conv.BoolPtrFrom(true),
			},
			IPAM: &account.PermissionsIPAMV2{
				ManageIPAM: conv.BoolPtrFrom(false),
			},
			DHCP: &account.PermissionsDHCPV2{
				ManageDHCP: conv.BoolPtrFrom(false),
			},
		},
	}

	_, err := c.APIKeysV2.Create(k)
	require.NoError(t, err)
	// Create a key with auth tags
	k = &account.APIKeyV2{
		ID:   "ddi-authtags",
		Key:  "key-2",
		Name: "name-2",
		Permissions: &account.PermissionsMapV2{
			DHCP: &account.PermissionsDHCPV2{
				TagsAllow: &[]account.AuthTag{
					{
						Name:  "auth:dhcpallow",
						Value: "",
					},
				},
				TagsDeny: &[]account.AuthTag{
					{
						Name:  "auth:dhcpdeny",
						Value: "denyme",
					},
				},
			},
			IPAM: &account.PermissionsIPAMV2{
				TagsAllow: &[]account.AuthTag{
					{
						Name:  "auth:ipamallow",
						Value: "",
					},
				},
				TagsDeny: &[]account.AuthTag{
					{
						Name:  "auth:ipamdeny",
						Value: "denyme",
					},
				},
			},
		},
	}
	// Create a key with empty auth tags
	k = &account.APIKeyV2{
		ID:   "ddi-empty-authtags",
		Key:  "key-3",
		Name: "name-3",
		Permissions: &account.PermissionsMapV2{
			DHCP: &account.PermissionsDHCPV2{
				TagsAllow: &[]account.AuthTag{},
				TagsDeny:  &[]account.AuthTag{},
			},
			IPAM: &account.PermissionsIPAMV2{
				TagsAllow: &[]account.AuthTag{},
				TagsDeny:  &[]account.AuthTag{},
			},
		},
	}

	_, err = c.APIKeysV2.Create(k)
	require.NoError(t, err)
}

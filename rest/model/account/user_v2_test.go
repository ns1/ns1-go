package account

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/ns1/ns1-go.v2/common/conv"
)

func TestUnmarshalUsersV2(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
		out  UserV2
	}{
		{
			"basic",
			[]byte(`{
"permissions": {
	"dns": {
		"view_zones": true,
		"manage_zones": true,
		"zones_allow_by_default": false,
		"zones_deny": [],
		"zones_allow": ["example.com"]
	},
	"data": {
		"push_to_datafeeds": false,
		"manage_datasources": false,
		"manage_datafeeds": false
	},
	"account": {
		"manage_payment_methods": false,
		"manage_plan": false,
		"manage_teams": false,
		"manage_apikeys": false,
		"manage_account_settings": false,
		"view_activity_log": false,
		"view_invoices": false,
		"manage_users": false,
		"manage_ip_whitelist": true
	},
	"monitoring": {
		"manage_lists": false,
		"manage_jobs": false,
		"view_jobs": false
	}
	},
    "teams": [],
    "email": "support@nsone.net",
    "last_access": 1376325771.0,
    "created": 1376325771.0,
    "notify": {
      "billing": true
    },
    "name": "API Example",
    "username": "apiexample",
	"ip_whitelist": ["1.1.1.1", "2.2.2.2"],
	"ip_whitelist_strict": true,
	"2fa_enabled": true,
	"shared_auth": {
        "saml": {
            "sso": true,
            "idp": {
                "use_metadata_url": true,
                "metadata_url": null,
                "metadata_file": null,
                "provider": "okta"
            }
        }
    }
  }`),
			UserV2{
				Username:             "apiexample",
				Name:                 "API Example",
				Email:                "support@nsone.net",
				LastAccess:           1376325771.0,
				Created:              1376325771.0,
				Notify:               NotificationSettings{true},
				TeamIDs:              []string{},
				IPWhitelist:          []string{"1.1.1.1", "2.2.2.2"},
				IPWhitelistStrict:    true,
				TwoFactorAuthEnabled: true,
				SharedAuth: SharedAuth{
					SAML: SAML{
						SSO: true,
						IDP: IDP{
							UseMetadataURL: pointerBool(true),
							MetadataURL:    nil,
							MetadataFile:   nil,
							Provider:       pointerString("okta"),
						},
					},
				},
				Permissions: &PermissionsMapV2{
					DNS: &PermissionsDNSV2{
						ViewZones:           conv.BoolPtrFrom(true),
						ManageZones:         conv.BoolPtrFrom(true),
						ZonesAllowByDefault: conv.BoolPtrFrom(false),
						ZonesDeny:           []string{},
						ZonesAllow:          []string{"example.com"},
					},
					Data: &PermissionsDataV2{
						PushToDatafeeds:   conv.BoolPtrFrom(false),
						ManageDatasources: conv.BoolPtrFrom(false),
						ManageDatafeeds:   conv.BoolPtrFrom(false),
					},
					Account: &PermissionsAccountV2{
						ManagePaymentMethods:  conv.BoolPtrFrom(false),
						ManagePlan:            conv.BoolPtrFrom(false),
						ManageTeams:           conv.BoolPtrFrom(false),
						ManageApikeys:         conv.BoolPtrFrom(false),
						ManageAccountSettings: conv.BoolPtrFrom(false),
						ViewActivityLog:       conv.BoolPtrFrom(false),
						ViewInvoices:          conv.BoolPtrFrom(false),
						ManageUsers:           conv.BoolPtrFrom(false),
						ManageIPWhitelist:     conv.BoolPtrFrom(true),
					},
					Monitoring: &PermissionsMonitoringV2{
						ManageLists: conv.BoolPtrFrom(false),
						ManageJobs:  conv.BoolPtrFrom(false),
						ViewJobs:    conv.BoolPtrFrom(false),
					},
				},
			},
		},
		{
			"perms",
			[]byte(`{
    "permissions": {
      "dns": {
        "view_zones": true,
        "manage_zones": true,
        "zones_allow_by_default": false,
        "zones_deny": [],
        "zones_allow": ["example.com"]
      },
      "data": {
        "push_to_datafeeds": false,
        "manage_datasources": false,
        "manage_datafeeds": false
      },
      "account": {
        "manage_payment_methods": false,
        "manage_plan": false,
        "manage_teams": false,
        "manage_apikeys": false,
        "manage_account_settings": false,
        "view_activity_log": false,
        "view_invoices": false,
        "manage_users": false,
        "manage_ip_whitelist": true
      },
      "monitoring": {
        "manage_lists": false,
        "manage_jobs": false,
        "view_jobs": false
      }
    },
    "teams": ["520422919f782d37dffb588a"],
    "email": "newuser@example.com",
    "last_access": null,
    "notify": {
      "billing": true
    },
    "name": "New User",
    "username": "newuser"
  }`),
			UserV2{
				Username:   "newuser",
				Name:       "New User",
				Email:      "newuser@example.com",
				TeamIDs:    []string{"520422919f782d37dffb588a"},
				LastAccess: 0.0,
				Notify:     NotificationSettings{true},
				Permissions: &PermissionsMapV2{
					DNS: &PermissionsDNSV2{
						ViewZones:           conv.BoolPtrFrom(true),
						ManageZones:         conv.BoolPtrFrom(true),
						ZonesAllowByDefault: conv.BoolPtrFrom(false),
						ZonesDeny:           []string{},
						ZonesAllow:          []string{"example.com"},
					},
					Data: &PermissionsDataV2{
						PushToDatafeeds:   conv.BoolPtrFrom(false),
						ManageDatasources: conv.BoolPtrFrom(false),
						ManageDatafeeds:   conv.BoolPtrFrom(false),
					},
					Account: &PermissionsAccountV2{
						ManagePaymentMethods:  conv.BoolPtrFrom(false),
						ManagePlan:            conv.BoolPtrFrom(false),
						ManageTeams:           conv.BoolPtrFrom(false),
						ManageApikeys:         conv.BoolPtrFrom(false),
						ManageAccountSettings: conv.BoolPtrFrom(false),
						ViewActivityLog:       conv.BoolPtrFrom(false),
						ViewInvoices:          conv.BoolPtrFrom(false),
						ManageUsers:           conv.BoolPtrFrom(false),
						ManageIPWhitelist:     conv.BoolPtrFrom(true),
					},
					Monitoring: &PermissionsMonitoringV2{
						ManageLists: conv.BoolPtrFrom(false),
						ManageJobs:  conv.BoolPtrFrom(false),
						ViewJobs:    conv.BoolPtrFrom(false),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var u UserV2
			require.NoError(t, json.Unmarshal(tt.in, &u))
			require.Equal(t, tt.out, u)
		})
	}
}

// Same as v1
// func pointerBool(b bool) *bool {
// 	return &b
// }

// func pointerString(s string) *string {
// 	return &s
// }

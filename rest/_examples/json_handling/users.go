package main

import (
	"encoding/json"
	"fmt"

	"github.com/ns1/ns1-go/model/account"
)

func prettyPrint(header string, v interface{}) {
	fmt.Println(header)
	fmt.Printf("%#v \n", v)
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}

func main() {
	// d := []byte(`[
	// {
	// "permissions": {},
	// "teams": [],
	// "email": "support@nsone.net",
	// "last_access": 1376325771,
	// "notify": {
	// "billing": true
	// },
	// "name": "API Example",
	// "username": "apiexample"
	// },
	// {
	// "permissions": {
	// "dns": {
	// "view_zones": true,
	// "manage_zones": true,
	// "zones_allow_by_default": false,
	// "zones_deny": [],
	// "zones_allow": ["example.com"]
	// },
	// "data": {
	// "push_to_datafeeds": false,
	// "manage_datasources": false,
	// "manage_datafeeds": false
	// },
	// "account": {
	// "manage_payment_methods": false,
	// "manage_plan": false,
	// "manage_teams": false,
	// "manage_apikeys": false,
	// "manage_account_settings": false,
	// "view_activity_log": false,
	// "view_invoices": false,
	// "manage_users": false
	// },
	// "monitoring": {
	// "manage_lists": false,
	// "manage_jobs": false,
	// "view_jobs": false
	// }
	// },
	// "teams": ["520422919f782d37dffb588a"],
	// "email": "newuser@example.com",
	// "last_access": null,
	// "notify": {
	// "billing": true
	// },
	// "name": "New User",
	// "username": "newuser"
	// }
	// ]
	// `)
	// ul := []*account.User{}
	// if err := json.Unmarshal(d, &ul); err != nil {
	// 	log.Fatal(err)
	// }

	permMap := account.PermissionsMap{
		DNS: account.PermissionsDNS{
			ViewZones:           true,
			ManageZones:         true,
			ZonesAllowByDefault: false,
			ZonesDeny:           []string{},
			ZonesAllow:          []string{"example.com"},
		},
		Data: account.PermissionsData{
			PushToDatafeeds:   false,
			ManageDatasources: false,
			ManageDatafeeds:   false,
		},
		Account: account.PermissionsAccount{
			ManagePaymentMethods:  false,
			ManagePlan:            false,
			ManageTeams:           false,
			ManageApikeys:         false,
			ManageAccountSettings: false,
			ViewActivityLog:       false,
			ViewInvoices:          false,
			ManageUsers:           false,
		},
		Monitoring: account.PermissionsMonitoring{
			ManageLists: false,
			ManageJobs:  false,
			ViewJobs:    false,
		},
	}
	prettyPrint("perms:", permMap.DNS.ZonesAllow)

}

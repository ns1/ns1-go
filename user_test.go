package nsone

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestUnmarshalUsers(t *testing.T) {
	data, err := ioutil.ReadFile("test/data/account_users.json")
	if err != nil {
		t.Error(err)
	}
	var users []User
	if err = json.Unmarshal(data, &users); err != nil {
		t.Error(err)
	}
	if len(users) != 1 {
		t.Error("user not found")
	}
	user := users[0]
	if user.Username != "test" {
		t.Error("username")
	}
	if user.Name != "Test User" {
		t.Error("name")
	}
	if len(user.Teams) != 0 {
		t.Error("teams")
	}
	// FIXME user.Notify
	if user.LastAccess != 1428457118 {
		t.Error("last access")
	}
	if user.Email != "test@test.com" {
		t.Error("email")
	}
	p := user.Permissions
	pm := p.Monitoring
	if !pm.ManageJobs {
		t.Error("manage jobs")
	}
	if !pm.ViewJobs {
		t.Error("view jobs")
	}
	if !pm.ManageLists {
		t.Error("manage lists")
	}
	pa := p.Account
	if !pa.ManagePlan {
		t.Error("manage plan")
	}
	if !pa.ManageUsers {
		t.Error("manage users")
	}
	if !pa.ViewInvoices {
		t.Error("view invoices")
	}
	if !pa.ManageTeams {
		t.Error("manage teams")
	}
	if !pa.ManagePaymentMethods {
		t.Error("manage payment methods")
	}
	if !pa.ManageAccountSettings {
		t.Error("manage account settings")
	}
	if !pa.ManageApikeys {
		t.Error("manage api keys")
	}
	if !pa.ViewActivityLog {
		t.Error("view activity log")
	}
	pdata := p.Data
	if !pdata.PushToDatafeeds {
		t.Error("push to data feeds")
	}
	if !pdata.ManageDatasources {
		t.Error("manage data sources")
	}
	if !pdata.ManageDatafeeds {
		t.Error("manage data feeds")
	}
	pdns := p.Dns
	if !pdns.ManageZones {
		t.Error("manage zones")
	}
	if !pdns.ViewZones {
		t.Error("view zones")
	}
	if len(pdns.ZonesAllow) != 1 {
		t.Error("zones allow size")
	}
	if pdns.ZonesAllow[0] != "foo.com" {
		t.Error("zones allow")
	}
	if len(pdns.ZonesDeny) != 1 {
		t.Error("zone deny size")
	}
	if pdns.ZonesDeny[0] != "bar.com" {
		t.Error("zone deny")
	}
}

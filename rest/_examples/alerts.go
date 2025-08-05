package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	api "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/alerting"
	"gopkg.in/ns1/ns1-go.v2/rest/model/dns"
	"gopkg.in/ns1/ns1-go.v2/rest/model/monitor"
)

var client *api.Client

// Helper that initializes rest api client from environment variable.
func init() {
	k := os.Getenv("NS1_APIKEY")
	if k == "" {
		fmt.Println("NS1_APIKEY environment variable is not set, giving up")
	}

	httpClient := &http.Client{Timeout: time.Second * 10}
	// Adds logging to each http request.
	doer := api.Decorate(httpClient, api.Logging(log.New(os.Stdout, "", log.LstdFlags)))
	client = api.NewClient(doer, api.SetAPIKey(k))
}

func prettyPrint(header string, v interface{}) {
	fmt.Println(header)
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}

func main() {
	alerts, _, err := client.Alerts.List()
	if err != nil {
		log.Fatal(err)
	}
	for _, a := range alerts {
		prettyPrint("alert:", a)
	}

	webhook := monitor.NewWebNotification("test.com/test", map[string]string{})
	webhookList := monitor.NewNotifyList("my webhook list", webhook)
	_, err = client.Notifications.Create(webhookList)
	if err != nil {
		log.Fatal(err)
	}
	prettyPrint("Webhook NotifyList:", webhookList)

	email := monitor.NewEmailNotification("test@test.com")
	emailList := monitor.NewNotifyList("my email list", email)
	_, err = client.Notifications.Create(emailList)
	if err != nil {
		log.Fatal(err)
	}
	prettyPrint("Email NotifyList:", emailList)

	// Construct/Create a zone.
	domain := "myalerttest.com"

	z := dns.NewZone(domain)
	z.NxTTL = 3600
	z.Secondary = &dns.ZoneSecondary{
		Enabled:     true,
		PrimaryIP:   "127.0.0.1",
		PrimaryPort: 67,
	}
	_, err = client.Zones.Create(z)
	if err != nil {
		// Ignore if zone already exists
		if err != api.ErrZoneExists {
			log.Fatal(err)
		} else {
			log.Println("Zone already exists, continuing...")
		}
	}
	prettyPrint("Zone:", z)

	alert := alerting.NewZoneAlert("myalerttest.com - transfer failed", "transfer_failed", []string{webhookList.ID}, []string{domain})
	prettyPrint("Creating alert: ", alert)
	_, err = client.Alerts.Create(alert)
	if err != nil {
		if err == api.ErrAlertExists {
			// This is fatal as we need the id returned on create.
			log.Println("Alert already exists.")
		}
		log.Fatal(err)

	}
	alertID := *alert.ID

	ssoAlert := alerting.NewSSOAlert("my sso alert", []string{emailList.ID})
	prettyPrint("Creating SSO alert: ", ssoAlert)
	_, err = client.Alerts.Create(ssoAlert)
	if err != nil && err != api.ErrAlertExists {
		log.Fatal(err)
	}
	if ssoAlert.ID != nil {
		defer client.Alerts.Delete(*ssoAlert.ID)
	}

	redirectAlert := alerting.NewRedirectAlert("my redirect alert", []string{emailList.ID})
	prettyPrint("Creating Redirect alert: ", redirectAlert)
	_, err = client.Alerts.Create(redirectAlert)
	if err != nil && err != api.ErrAlertExists {
		log.Fatal(err)
	}
	if redirectAlert.ID != nil {
		defer client.Alerts.Delete(*redirectAlert.ID)
	}

	// Pass the id and the field(s) to change on Update.
	updatedName := "myalerttest.com - updated"
	alertUpdate := &alerting.Alert{
		ID:              &alertID,
		Name:            &updatedName,
		NotifierListIds: []string{webhookList.ID, emailList.ID},
	}
	_, err = client.Alerts.Update(alertUpdate)
	if err != nil {
		log.Fatal(err)
	}

	prettyPrint("Updated Alert:", alertUpdate)

	// To pass the whole alert object on Replace, retrieve it by ID it first.
	alertToReplace, _, err := client.Alerts.Get(alertID)
	if err != nil {
		log.Fatal(err)
	}

	// Replace values in retrieved struct with new values.
	// e.g. Change name and clear list.
	replacedName := "myalerttest.com - replaced"
	alertToReplace.Name = &replacedName
	alertToReplace.NotifierListIds = []string{}

	// Pass the whole alert object
	_, err = client.Alerts.Replace(alertToReplace)
	if err != nil {
		log.Fatal(err)
	}

	prettyPrint("Replaced Alert:", alertToReplace)

	// Delete the alert.
	_, err = client.Alerts.Delete(*alertToReplace.ID)
	if err != nil {
		log.Fatal(err)
	}
}

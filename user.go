package nsone

import "fmt"

// User wraps an NS1 /account/users resource
type User struct {
	Name        string               `json:"name"`
	Username    string               `json:"username"`
	Notify      NotificationSettings `json:"notify"`
	LastAccess  int                  `json:"last_access,omitempty"`
	Email       string               `json:"email"`
	Teams       []string             `json:"teams"`
	Permissions PermissionsMap       `json:"permissions"`
}

// NotificationSettings wraps a User's "notify" attribute
type NotificationSettings struct {
	Billing bool `json:"billing"`
}

// PermissionsMap wraps a User's "permissions" attribute
type PermissionsMap struct {
	Dns        PermissionsDns        `json:"dns"`
	Data       PermissionsData       `json:"data"`
	Account    PermissionsAccount    `json:"account"`
	Monitoring PermissionsMonitoring `json:"monitoring"`
}

// PermissionsDns wraps a User's "permissions.dns" attribute
type PermissionsDns struct {
	ViewZones           bool     `json:"view_zones"`
	ManageZones         bool     `json:"manage_zones"`
	ZonesAllowByDefault bool     `json:"zones_allow_by_default"`
	ZonesDeny           []string `json:"zones_deny"`
	ZonesAllow          []string `json:"zones_allow"`
}

// PermissionsData wraps a User's "permissions.data" attribute
type PermissionsData struct {
	PushToDatafeeds   bool `json:"push_to_datafeeds"`
	ManageDatasources bool `json:"manage_datasources"`
	ManageDatafeeds   bool `json:"manage_datafeeds"`
}

// PermissionsAccount wraps a User's "permissions.account" attribute
type PermissionsAccount struct {
	ManageUsers           bool `json:"manage_users"`
	ManagePaymentMethods  bool `json:"manage_payment_methods"`
	ManagePlan            bool `json:"manage_plan"`
	ManageTeams           bool `json:"manage_teams"`
	ManageApikeys         bool `json:"manage_apikeys"`
	ManageAccountSettings bool `json:"manage_account_settings"`
	ViewActivityLog       bool `json:"view_activity_log"`
	ViewInvoices          bool `json:"view_invoices"`
}

// PermissionsMonitoring wraps a User's "permissions.monitoring" attribute
type PermissionsMonitoring struct {
	ManageLists bool `json:"manage_lists"`
	ManageJobs  bool `json:"manage_jobs"`
	ViewJobs    bool `json:"view_jobs"`
}

// GetUsers returns a list of all users with access to the account
func (c APIClient) GetUsers() ([]User, error) {
	req, err := c.NewRequest("GET", "account/users", nil)
	if err != nil {
		return nil, err
	}

	var users []User
	_, err = c.Do(req, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// GetUser takes a username and returns the details for a single user
func (c APIClient) GetUser(username string) (*User, error) {
	path := fmt.Sprintf("account/users/%s", username)
	req, err := c.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var u User
	_, err = c.Do(req, &u)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// // CreateUser takes a *User and creates a new user
// func (c APIClient) CreateUser(u *User) error {
// 	return c.doHTTPBoth("PUT", fmt.Sprintf("https://api.nsone.net/v1/account/users/%s", u.Username), &u)
// }

// // DeleteUser takes a username and deletes a user from the account
// func (c APIClient) DeleteUser(username string) error {
// 	return c.doHTTPDelete(fmt.Sprintf("https://api.nsone.net/v1/account/users/%s", username))
// }

// // UpdateUser takes a *User and change contact details, notification settings or access rights for a user
// func (c APIClient) UpdateUser(user *User) error {
// 	return c.doHTTPBoth("POST", fmt.Sprintf("https://api.nsone.net/v1/account/users/%s", user.Username), user)
// }

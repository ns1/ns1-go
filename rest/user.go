package rest

// // GetUsers returns a list of all users with access to the account
// func (c APIClient) GetUsers() ([]User, error) {
// 	req, err := c.NewRequest("GET", "account/users", nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var users []User
// 	_, err = c.Do(req, &users)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return users, nil
// }

// // GetUser takes a username and returns the details for a single user
// func (c APIClient) GetUser(username string) (*User, error) {
// 	path := fmt.Sprintf("account/users/%s", username)
// 	req, err := c.NewRequest("GET", path, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var u User
// 	_, err = c.Do(req, &u)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &u, nil
// }

// // // CreateUser takes a *User and creates a new user
// // func (c APIClient) CreateUser(u *User) error {
// // 	return c.doHTTPBoth("PUT", fmt.Sprintf("https://api.nsone.net/v1/account/users/%s", u.Username), &u)
// // }

// // // DeleteUser takes a username and deletes a user from the account
// // func (c APIClient) DeleteUser(username string) error {
// // 	return c.doHTTPDelete(fmt.Sprintf("https://api.nsone.net/v1/account/users/%s", username))
// // }

// // // UpdateUser takes a *User and change contact details, notification settings or access rights for a user
// // func (c APIClient) UpdateUser(user *User) error {
// // 	return c.doHTTPBoth("POST", fmt.Sprintf("https://api.nsone.net/v1/account/users/%s", user.Username), user)
// // }

// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package stormpath

import (
	"encoding/base64"
	"github.com/jmcvetta/napping"
	"log"
	"net/url"
)

// An Application in represents a real world application that can communicate
// with and be provisioned by Stormpath. After it is defined, an application is
// mapped to one or more directories or groups, whose users are then granted
// access to the application.
type Application struct {
	Href           string // Stormpath URL for this application
	ApiId          string // Stormpath API key ID
	ApiSecret      string // Stormpath API key secret
	LoadCustomData bool   // Whether it should load custom data on login
}

func (a *Application) userinfo() *url.Userinfo {
	return url.UserPassword(a.ApiId, a.ApiSecret)
}

// CreateAccount creates a new account accessible to the application.
func (app *Application) CreateAccount(template Account) (Account, error) {
	/*
		data := &map[string]string{
			"username":  username,
			"password":  password,
			"email":     email,
			"surname":   surname,
			"givenName": givenName,
		}
	*/
	url := app.Href + "/accounts"
	acct := Account{}
	e := new(StormpathError)
	req := &napping.Request{
		Userinfo: app.userinfo(),
		Url:      url,
		Method:   "POST",
		Payload:  &template,
		Result:   &acct,
		Error:    e,
	}
	res, err := napping.Send(req)
	if err != nil {
		return acct, err
	}
	acct.app = app
	if res.Status() != 201 {
		log.Println(res.Status())
		log.Println(e)
		return acct, BadResponse
	}
	return acct, nil
}

// Authenticate with Stormpath using supplied credentials.  Username may be
// either a username or the user's email.
func (app *Application) Authenticate(username, password string) (Account, error) {
	acct := Account{}
	s := username + ":" + password
	value := base64.URLEncoding.EncodeToString([]byte(s))
	m := map[string]string{
		"type":  "basic",
		"value": value,
	}
	loginUrl := app.Href + "/loginAttempts"
	var resp struct {
		Account struct {
			Href string `json:"href"`
		} `json:"account"`
	}
	e := new(StormpathError)
	req := &napping.Request{
		Userinfo: app.userinfo(),
		Url:      loginUrl,
		Method:   "POST",
		Payload:  &m,
		Result:   &resp,
		Error:    &e,
	}
	res, err := napping.Send(req)
	if err != nil {
		return acct, err
	}
	if res.Status() != 200 {
		log.Println(res.Status())
		log.Println(res)
		log.Println(e)
		return acct, InvalidUsernamePassword
	}
	return app.GetAccount(resp.Account.Href)
}

// GetAccount returns the specified account object, if it exists.
func (app *Application) GetAccount(href string) (Account, error) {
	acct := Account{}
	e := new(StormpathError)
	if app.LoadCustomData {
		href = href + "?expand=customData"
	}
	req := &napping.Request{
		Userinfo: app.userinfo(),
		Url:      href,
		Method:   "GET",
		Result:   &acct,
		Error:    e,
	}
	res, err := napping.Send(req)
	if err != nil {
		return acct, err
	}
	// remove additional metadata because it will prevent us from saving
	if app.LoadCustomData && acct.CustomData != nil {
		delete(acct.CustomData, "createdAt")
		delete(acct.CustomData, "modifiedAt")
	}
	acct.app = app
	if res.Status() != 200 {
		log.Println(res.Status())
		log.Println(e)
		return acct, BadResponse
	}
	return acct, nil
}

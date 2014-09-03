// Copyright (C) 2013 Jason McVetta, all rights reserved.

package stormpath

import (
	"encoding/json"
	"github.com/bmizerany/assert"
	"github.com/darkhelmet/env"
	"github.com/jmcvetta/randutil"
	"log"
	"testing"
)

func setupApplication(t *testing.T) *Application {
	log.SetFlags(log.Ltime | log.Ldate | log.Lshortfile)
	spApp := env.String("STORMPATH_APP")
	apiId := env.String("STORMPATH_API_ID")
	apiSecret := env.String("STORMPATH_API_SECRET")
	s := Application{
		Href:      spApp,
		ApiId:     apiId,
		ApiSecret: apiSecret,
	}
	return &s
}

func createAccountTemplate(t *testing.T) Account {
	rnd, err := randutil.AlphaString(8)
	if err != nil {
		t.Error(err)
	}
	email := "jason.mcvetta+" + rnd + "@gmail.com"
	password := rnd + "Xy123" // Ensure we meet password requirements
	tmpl := Account{
		Username:   rnd,
		Email:      email,
		Password:   password,
		GivenName:  "James",
		MiddleName: "T",
		Surname:    "Kirk",
	}
	return tmpl
}

func TestCreateAccount(t *testing.T) {
	app := setupApplication(t)
	tmpl := createAccountTemplate(t)
	acct, err := app.CreateAccount(tmpl)
	if err != nil {
		t.Fatal(err)
	}
	//
	// Cleanup
	//
	acct.Delete()
}

func TestDeleteAccount(t *testing.T) {
	app := setupApplication(t)
	tmpl := createAccountTemplate(t)
	acct, _ := app.CreateAccount(tmpl)
	err := acct.Delete()
	if err != nil {
		t.Error(err)
	}
}

func TestGetAccount(t *testing.T) {
	app := setupApplication(t)
	tmpl := createAccountTemplate(t)
	acct0, _ := app.CreateAccount(tmpl)
	defer acct0.Delete()
	acct1, err := app.GetAccount(acct0.Href)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, acct0, acct1)
}

func TestUpdateAccount(t *testing.T) {
	app := setupApplication(t)
	tmpl := createAccountTemplate(t)
	acct0, _ := app.CreateAccount(tmpl)
	defer acct0.Delete()
	acct0.GivenName = "Mister"
	acct0.MiddleName = ""
	acct0.Surname = "Spock"
	err := acct0.Update()
	if err != nil {
		t.Fatal(err)
	}
	acct1, _ := app.GetAccount(acct0.Href)
	assert.Equal(t, "Mister", acct1.GivenName)
	assert.Equal(t, "Spock", acct1.Surname)
}

func TestSaveCustomData(t *testing.T) {
	app := setupApplication(t)
	app.LoadCustomData = true
	tmpl := createAccountTemplate(t)
	acct0, _ := app.CreateAccount(tmpl)
	defer acct0.Delete()
	msg := json.RawMessage("\"foo\"")
	acct0.CustomData = map[string]*json.RawMessage{
		"field1": &msg,
	}
	err := acct0.Update()
	if err != nil {
		t.Fatal(err)
	}
	acct1, _ := app.GetAccount(acct0.Href)
	resp := acct1.CustomData["field1"]
	b, _ := resp.MarshalJSON()
	s := string(b)
	assert.Equal(t, "\"foo\"", s)
}

func TestDontLoadCustomData(t *testing.T) {
	app := setupApplication(t)
	tmpl := createAccountTemplate(t)
	acct0, _ := app.CreateAccount(tmpl)
	defer acct0.Delete()
	msg := json.RawMessage("\"foo\"")
	acct0.CustomData = map[string]*json.RawMessage{
		"field1": &msg,
	}
	err := acct0.Update()
	if err != nil {
		t.Fatal(err)
	}
	acct1, _ := app.GetAccount(acct0.Href)
	_, ok := acct1.CustomData["field1"]
	assert.Equal(t, ok, false)
}

func TestManualLoadCustomData(t *testing.T) {
	app := setupApplication(t)
	tmpl := createAccountTemplate(t)
	acct0, _ := app.CreateAccount(tmpl)
	defer acct0.Delete()
	msg := json.RawMessage("\"foo\"")
	acct0.CustomData = map[string]*json.RawMessage{
		"field1": &msg,
	}
	err := acct0.Update()
	if err != nil {
		t.Fatal(err)
	}
	acct1, _ := app.GetAccount(acct0.Href)
	if err := acct1.LoadCustomData(); err != nil {
		t.Fatal(err)
	}
	resp, ok := acct1.CustomData["field1"]
	assert.Equal(t, ok, true)
	b, _ := resp.MarshalJSON()
	s := string(b)
	assert.Equal(t, "\"foo\"", s)
}

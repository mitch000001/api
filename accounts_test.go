package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAccountValidations(t *testing.T) {
	account := Account{}

	if account.IsValid() {
		t.Fatalf("expected empty account to be invalid")
	}
}

func TestAccountsIndex(t *testing.T) {
	app = &App{SetupDb(), I18nInit()}
	app.ClearDb()

	if err := app.Db.Query(`INSERT INTO accounts (code, label) VALUES ('5900', 'Fremdleistungen') RETURNING *`).Run(); err != nil {
		t.Fatalf("Unable to insert fake account: %v", err)
	}

	request, _ := http.NewRequest("GET", "/accounts", strings.NewReader(""))
	response := httptest.NewRecorder()

	app.AccountIndexHandler(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("Non-expected status code%v:\n\tbody: %+v", "200", response.Code)
	}

	decoder := json.NewDecoder(response.Body)
	var accounts []Account

	_ = decoder.Decode(&accounts)
	if len(accounts) != 1 {
		t.Fatalf("Received wrong number of accounts: %v - '%v'", accounts, response.Body)
	}
}

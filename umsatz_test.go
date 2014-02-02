package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
  "testing"
  // "io/ioutil"
  "fmt"
)

func init() {
	jetDb = SetupDb()
}

func ClearDb() {
	jetDb.Query("DELETE FROM positions").Run()
	jetDb.Query("DELETE FROM fiscal_periods").Run()
}

func TestFiscalPeriodsIndex(t *testing.T) {
	ClearDb()
	jetDb.Query("INSERT INTO fiscal_periods (year) VALUES (2014)").Run()

	request, _ := http.NewRequest("GET", "/fiscalPeriods", strings.NewReader(""))
	response := httptest.NewRecorder()

	FiscalPeriodIndexHandler(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("Non-expected status code%v:\n\tbody: %+v", "200", response.Code)
	}

	decoder := json.NewDecoder(response.Body)

	var fiscalPeriods []FiscalPeriod
	_ = decoder.Decode(&fiscalPeriods)
	if len(fiscalPeriods) != 1 {
		t.Fatalf("Received wrong number of fiscalPeriods: %v - %v", fiscalPeriods, response.Body)
	}
}

func TestFiscalPeriodsPositionCreation(t *testing.T) {
	ClearDb()
	jetDb.Query("INSERT INTO fiscal_periods (year) VALUES (2014)").Run()

	request, _ := http.NewRequest("POST", "/fiscalPeriods/2014/positions", strings.NewReader(`
      {
        "category": "Freelance",
        "account": "5900",
        "type": "income",
        "invoiceDate": "2014-02-02",
        "invoiceNumber": "20140201",
        "totalAmount": 2099,
        "currency": "EUR",
        "tax": 700,
        "description": "Kunde A Februar"
      }`,
	))
	response := httptest.NewRecorder()

	FiscalPeriodCreatePositionHandler(response, request, map[string]string{"year": "2014"})

	if response.Code != http.StatusOK {
		t.Fatalf("Non-expected status code%v:\n\tbody: %+v", "200", response.Code)
	}

  decoder := json.NewDecoder(response.Body)

  var position Position
  _ = decoder.Decode(&position)

  // b, _ := ioutil.ReadAll(response.Body)
  // fmt.Printf("'%#v'", string(b))
  // fmt.Printf("%#v", position)

  if len(position.Errors) != 0 {
    t.Fatalf("payload should have been valid, got '%v'", position.Errors)
  }

  if position.Category != "Freelance" {
    t.Fatalf("did not persist category correctly, expected 'Freelance', got %#v", position.Category)
  }
  if position.Account != "5900" {
    t.Fatalf("did not persist account correctly, got %v", position.Account)
  }
  if position.PositionType != "income" {
    t.Fatalf("did not persist type correctly, got %v", position.PositionType)
  }
}

func TestFiscalPeriodsPositionCreationWithMissingPositionAttributes(t *testing.T) {
  ClearDb()
  jetDb.Query("INSERT INTO fiscal_periods (year) VALUES (2014)").Run()

  request, _ := http.NewRequest("POST", "/fiscalPeriods/2014/positions", strings.NewReader(`
      {
        "category": "Freelance",
        "account": "5900",
        "invoiceDate": "2014-02-02",
        "invoiceNumber": "20140201",
        "totalAmount": 2099,
        "currency": "EUR",
        "tax": 700,
        "description": "Kunde A Februar"
      }`,
  ))
  response := httptest.NewRecorder()

  FiscalPeriodCreatePositionHandler(response, request, map[string]string{"year": "2014"})

  if response.Code != http.StatusBadRequest {
    t.Fatalf("Non-expected status code%v:\n\tbody: %+v", "400", response.Code)
  }

  decoder := json.NewDecoder(response.Body)

  var position Position
  _ = decoder.Decode(&position)

  if len(position.Errors) == 0 {
    t.Fatalf("payload should have been invalid")
  }

  if false {
    fmt.Printf("%#v", position.Errors)
  }
}
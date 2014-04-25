package controllers

import (
	"github.com/umsatz/api/models"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestFiscalPeriodsPositionsIndex(t *testing.T) {
  app := &App{}
  app.SetupDb()
  app.ClearDb()

  var fiscalPeriod models.FiscalPeriod
  app.Db.Query("INSERT INTO fiscal_periods (year) VALUES (2014) RETURNING *").Rows(&fiscalPeriod)
  app.Db.Query(`INSERT INTO "positions"
        (account_code_from, account_code_to, type, invoice_date, booking_date, invoice_number, total_amount_cents, currency, tax, fiscal_period_id, description, attachment_path)
      VALUES
        ('5900', '1100', 'expense', NOW(), NOW(), '2001312', 0, 'EUR', 0, $1, '', '')`, fiscalPeriod.Id).Run()

  request, _ := http.NewRequest("GET", "/fiscalPeriods/2014/positions", strings.NewReader(""))
  response := httptest.NewRecorder()

  app.FiscalPeriodPositionIndexHandler(response, request, map[string]string{"year": "2014"})

  if response.Code != http.StatusOK {
    t.Fatalf("Non-expected status code%v:\n\tbody: %+v", "200", response.Code)
  }

  decoder := json.NewDecoder(response.Body)
  var positions []models.Position
  _ = decoder.Decode(&positions)
  if len(positions) != 1 {
    t.Fatalf("Received wrong number of positions: %v - %v", positions, response.Body)
  }
}


func TestFiscalPeriodsPositionCreation(t *testing.T) {
	app := &App{}
	app.SetupDb()
	app.ClearDb()

	app.Db.Query("INSERT INTO fiscal_periods (year) VALUES (2014)").Run()

	request, _ := http.NewRequest("POST", "/fiscalPeriods/2014/positions", strings.NewReader(`
      {
        "accountCodeFrom": "5900",
        "accountCodeTo": "1100",
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

	app.FiscalPeriodCreatePositionHandler(response, request, map[string]string{"year": "2014"})

	if response.Code != http.StatusOK {
		t.Fatalf("Non-expected status code%v:\n\tbody: %+v", "200", response.Code)
	}

	decoder := json.NewDecoder(response.Body)

	var position models.Position
	_ = decoder.Decode(&position)

	if position.AccountCodeFrom != "5900" {
		t.Fatalf("did not persist accountCodeFrom correctly, got '%v'", position.AccountCodeFrom)
	}
  if position.AccountCodeTo != "1100" {
    t.Fatalf("did not persist accountCodeTo correctly, got '%v'", position.AccountCodeTo)
  }
	if position.PositionType != "income" {
		t.Fatalf("did not persist type correctly, got %v", position.PositionType)
	}

	updateRequest, _ := http.NewRequest("PUT", ("/fiscalPeriods/2014/positions/" + strconv.Itoa(position.Id)), strings.NewReader(`
      {
        "type": "expense"
      }`,
	))
	updateResponse := httptest.NewRecorder()
	app.FiscalPeriodUpdatePositionHandler(updateResponse, updateRequest, map[string]string{"year": "2014", "id": strconv.Itoa(position.Id)})

	if response.Code != http.StatusOK {
		t.Fatalf("Non-expected status code%v:\n\tbody: %+v", "200", response.Code)
	}

	decoder = json.NewDecoder(updateResponse.Body)
	var updatedPosition models.Position
	updateErr := decoder.Decode(&updatedPosition)

	if updateErr != nil {
		t.Fatalf("error decoding update '%v'", updateErr)
	}

	if updatedPosition.PositionType != "expense" {
		t.Fatalf("position should have been expense now, got '%v'", updatedPosition.PositionType)
	}

}

func TestFiscalPeriodsPositionCreationWithMissingPositionAttributes(t *testing.T) {
	app := &App{}
	app.SetupDb()
	app.ClearDb()
	app.Db.Query("INSERT INTO fiscal_periods (year) VALUES (2014)").Run()

	request, _ := http.NewRequest("POST", "/fiscalPeriods/2014/positions", strings.NewReader(`
      {
        "accountCodeFrom": "5900",
        "accountCodeTo": "5900",
        "invoiceDate": "2014-02-02",
        "invoiceNumber": "20140201",
        "totalAmount": 2099,
        "currency": "EUR",
        "tax": 705,
        "description": "Kunde A Februar"
      }`,
	))
	response := httptest.NewRecorder()

	app.FiscalPeriodCreatePositionHandler(response, request, map[string]string{"year": "2014"})

	if response.Code != http.StatusBadRequest {
		t.Fatalf("Non-expected status code%v:\n\tbody: %+v", "400", response.Code)
	}

	decoder := json.NewDecoder(response.Body)

	var position models.Position
	_ = decoder.Decode(&position)

	if len(position.Errors) == 0 {
		t.Fatalf("payload should have been invalid")
	}

	if false {
		fmt.Printf("%#v", position.Errors)
	}
}

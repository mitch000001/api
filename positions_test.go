package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/melvinmt/gt"
)

func TestPositionValidations(t *testing.T) {
	position := Position{}
	var locales gt.Strings
	g := &gt.Build{
		Index:  locales,
		Origin: "en",
	}

	if position.IsValid(g) {
		t.Fatalf("expected empty position to be invalid")
	}
	position.Currency = "EUR"
	if position.IsValid(g) {
		t.Fatalf("expected empty position to be invalid")
	}
	position.PositionType = "income"
	if position.IsValid(g) {
		t.Fatalf("expected empty position to be invalid")
	}
	position.AccountCodeFrom = "5900"
	if position.IsValid(g) {
		t.Fatalf("expected empty position to be invalid")
	}
	position.AccountCodeTo = "5900"
	if position.IsValid(g) {
		t.Fatalf("expected empty position to be invalid")
	}
	position.InvoiceDate = ShortDate(time.Now())
	if position.IsValid(g) {
		t.Fatalf("expected empty position to be invalid")
	}
	position.InvoiceNumber = "20140101"

	if !position.IsValid(g) {
		t.Fatalf("expect position to be valid")
	}
}

func TestPositionUnmarshal(t *testing.T) {
	payload := `{
    "fiscalPeriodId": null,
    "accountCodeFrom": "5900",
    "accountCodeTo": "1100",
    "type":"expense",
    "invoiceDate":"2014-01-01",
    "invoiceNumber":"20140101",
    "totalAmountCents":4255,
    "tax":7,
    "description":"",
    "currency":"EUR"
  }`

	decoder := json.NewDecoder(strings.NewReader(payload))
	var position Position
	if err := decoder.Decode(&position); err != nil && err != io.EOF {
		t.Fatalf("error decoding JSON", err)
	}
}

func TestFiscalPeriodsPositionsIndex(t *testing.T) {
	app = &App{SetupDb(), I18nInit()}
	app.ClearDb()

	var fiscalPeriod FiscalPeriod
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
	var positions []Position
	_ = decoder.Decode(&positions)
	if len(positions) != 1 {
		t.Fatalf("Received wrong number of positions: %v - %v", positions, response.Body)
	}
}

func TestFiscalPeriodsPositionCreation(t *testing.T) {
	app = &App{SetupDb(), I18nInit()}
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

	var position Position
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
	var updatedPosition Position
	updateErr := decoder.Decode(&updatedPosition)

	if updateErr != nil {
		t.Fatalf("error decoding update '%v'", updateErr)
	}

	if updatedPosition.PositionType != "expense" {
		t.Fatalf("position should have been expense now, got '%v'", updatedPosition.PositionType)
	}

}

func TestFiscalPeriodsPositionCreationWithMissingPositionAttributes(t *testing.T) {
	app = &App{SetupDb(), I18nInit()}
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

	var position Position
	_ = decoder.Decode(&position)

	if len(position.Errors) == 0 {
		t.Fatalf("payload should have been invalid")
	}

	if false {
		fmt.Printf("%#v", position.Errors)
	}
}

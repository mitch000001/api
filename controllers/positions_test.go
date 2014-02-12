package controllers

import (
  "../models"
  "encoding/json"
  "net/http"
  "net/http/httptest"
  "strconv"
  "strings"
  "testing"
  "fmt"
)

func TestFiscalPeriodsPositionCreation(t *testing.T) {
  app := &App{}
  app.SetupDb()
  app.ClearDb()

  app.Db.Query("INSERT INTO fiscal_periods (year) VALUES (2014)").Run()

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

  app.FiscalPeriodCreatePositionHandler(response, request, map[string]string{"year": "2014"})

  if response.Code != http.StatusOK {
    t.Fatalf("Non-expected status code%v:\n\tbody: %+v", "200", response.Code)
  }

  decoder := json.NewDecoder(response.Body)

  var position models.Position
  _ = decoder.Decode(&position)

  if position.Category != "Freelance" {
    t.Fatalf("did not persist category correctly, expected 'Freelance', got %#v", position.Category)
  }
  if position.Account != "5900" {
    t.Fatalf("did not persist account correctly, got '%v'", position.Account)
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

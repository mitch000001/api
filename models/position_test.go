package models

import (
  "time"
  "testing"
  "strings"
  "encoding/json"
  "io"
)

func TestPositionValidations(t *testing.T) {
  position := Position{}

  if position.IsValid() {
    t.Fatalf("expected empty position to be invalid")
  }
  position.Currency = "EUR"
  if position.IsValid() {
    t.Fatalf("expected empty position to be invalid")
  }
  position.PositionType = "income"
  if position.IsValid() {
    t.Fatalf("expected empty position to be invalid")
  }
  position.AccountCodeFrom = "5900"
  if position.IsValid() {
    t.Fatalf("expected empty position to be invalid")
  }
  position.AccountCodeTo = "5900"
  if position.IsValid() {
    t.Fatalf("expected empty position to be invalid")
  }
  position.InvoiceDate = ShortDate(time.Now())
  if position.IsValid() {
    t.Fatalf("expected empty position to be invalid")
  }
  position.InvoiceNumber = "20140101"

  if !position.IsValid() {
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
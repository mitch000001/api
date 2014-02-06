package models

import (
  "time"
  "testing"
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
  position.Account = "5900"
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
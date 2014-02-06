package umsatz

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
  position.PositionType = "income"
  position.Account = "5900"
  position.InvoiceDate = ShortDate(time.Now())
  position.InvoiceNumber = "20140101"

  if !position.IsValid() {
    t.Fatalf("expect position to be valid")
  }
}
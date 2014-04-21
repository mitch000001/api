package models

import (
  // "time"
  "testing"
  // "strings"
  // "encoding/json"
  // "io"
)

func TestAccountValidations(t *testing.T) {
  account := Account{}

  if account.IsValid() {
    t.Fatalf("expected empty account to be invalid")
  }
}
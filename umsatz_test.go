package main

import (
  _ "encoding/json"
  _ "github.com/eaigner/hood"
  "net/http"
  "net/http/httptest"
  "strings"
  "testing"
  _ "time"
)

func TestFiscalPeriodsIndex(t *testing.T) {
  http.NewRequest("POST", "/fiscalPeriods", strings.NewReader(""))
  response := httptest.NewRecorder()

  if response.Code != http.StatusOK {
    t.Fatalf("Non-expected status code%v:\n\tbody: %v", "200", response.Code)
  }
}
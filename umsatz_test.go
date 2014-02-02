package main

import (
  "encoding/json"
  _ "github.com/eaigner/hood"
  "net/http"
  "net/http/httptest"
  "strings"
  "testing"
)

func init() {
  hd = SetupHood()
}

func TestFiscalPeriodsIndex(t *testing.T) {
  hd.Exec("DELETE FROM fiscal_periods")

  var fiscalPeriod FiscalPeriods = FiscalPeriods{Year: 2014}
  hd.Save(&fiscalPeriod)

  http.NewRequest("POST", "/fiscalPeriods", strings.NewReader(""))
  response := httptest.NewRecorder()

  if response.Code != http.StatusOK {
    t.Fatalf("Non-expected status code%v:\n\tbody: %v", "200", response.Code)
  }

  decoder := json.NewDecoder(response.Body)

  var fiscalPeriods []FiscalPeriods
  _ = decoder.Decode(&fiscalPeriods)
  if len(fiscalPeriods) > 1 {
    t.Fatalf("Received wrong number of fiscalPeriods: %v", fiscalPeriods)
  }
}
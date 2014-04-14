package controllers

import (
	. "github.com/umsatz/api/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFiscalPeriodsIndex(t *testing.T) {
	app := &App{}
	app.SetupDb()
	app.ClearDb()

	app.Db.Query("INSERT INTO fiscal_periods (year) VALUES (2014)").Run()

	request, _ := http.NewRequest("GET", "/fiscalPeriods", strings.NewReader(""))
	response := httptest.NewRecorder()

	app.FiscalPeriodIndexHandler(response, request)

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

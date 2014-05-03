package controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/umsatz/api/models"
)

func (app *App) FiscalPeriodIndexHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var fiscalPeriods []models.FiscalPeriod
	err := app.Db.Query(`SELECT * FROM "fiscal_periods" ORDER BY year ASC`).Rows(&fiscalPeriods)

	if err != nil {
		log.Println("database error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bytes, err := json.Marshal(fiscalPeriods)
	if err != nil {
		log.Println("json marshal error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	io.WriteString(w, string(bytes))
}

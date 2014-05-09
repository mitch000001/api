package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type FiscalPeriod struct {
	Id        int        `json:"-"`
	Year      int        `json:"year"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Positions []Position `json:"positions"`
}

func (app *App) FiscalPeriodIndexHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var fiscalPeriods []FiscalPeriod
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

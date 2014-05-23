package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/splicers/jet"
)

type FiscalPeriod struct {
	Id        int        `json:"-"`
	Year      int        `json:"year"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Positions []Position `json:"-"`
}

type fiscalPeriodResponse struct {
	FiscalPeriod
	PositionsCount int    `json:"positions_count"`
	Links          []Link `json:"links"`
}

func loadFiscalPeriods(db *jet.Db) ([]fiscalPeriodResponse, error) {
	var fiscalPeriods []fiscalPeriodResponse

	if err := db.Query(`
		SELECT
			fiscal_periods.*,
			(SELECT count(*) FROM positions WHERE fiscal_period_id = fiscal_periods.id) AS positions_count
		FROM fiscal_periods
		ORDER BY year ASC`).Rows(&fiscalPeriods); err != nil {
		return nil, err
	}

	for i, fiscalPeriod := range fiscalPeriods {
		link := fmt.Sprintf("/fiscalPeriods/%v/positions", fiscalPeriod.Year)
		links := []Link{
			Link{"positions", link},
		}
		fiscalPeriods[i].Links = links
	}

	return fiscalPeriods, nil
}

func (app *App) FiscalPeriodIndexHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/vnd.umsatz+json; charset=utf-8")

	fiscalPeriods, err := loadFiscalPeriods(app.Db)

	if err != nil {
		log.Println("database error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bytes, err := json.Marshal(fiscalPeriods)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	io.WriteString(w, string(bytes))
}

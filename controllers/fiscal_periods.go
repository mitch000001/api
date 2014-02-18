package controllers

import (
  "../models"
  "encoding/json"
  "log"
  "io"
  "net/http"
)

func (app *App) FiscalPeriodIndexHandler(w http.ResponseWriter, req *http.Request) {
  w.Header().Set("Content-Type", "application/json; charset=utf-8")

  log.Println("GET /fiscalPeriods")
  var fiscalPeriods []models.FiscalPeriod
  err := app.Db.Query(`SELECT * FROM "fiscal_periods" ORDER BY year ASC`).Rows(&fiscalPeriods)

  if err != nil {
    log.Println("database error", err)
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  for i, fiscalPeriod := range fiscalPeriods {
    var positions []models.Position
    err = app.Db.Query(`SELECT *, type as position_type FROM positions WHERE fiscal_period_id = $1`, fiscalPeriod.Id).Rows(&positions)
    fiscalPeriods[i].Positions = positions
  }

  bytes, err := json.Marshal(fiscalPeriods)
  if err != nil {
    log.Println("json marshal error", err)
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  io.WriteString(w, string(bytes))
}
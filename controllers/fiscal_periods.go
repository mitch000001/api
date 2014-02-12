package controllers

import (
  "../models"
  "encoding/json"
  "fmt"
  "log"
  "io"
  "net/http"
)

func (app *App) FiscalPeriodIndexHandler(w http.ResponseWriter, req *http.Request) {
  log.Println("GET /fiscalPeriods")
  var fiscalPeriods []models.FiscalPeriod
  err := app.Db.Query(`SELECT * FROM "fiscal_periods" ORDER BY year ASC`).Rows(&fiscalPeriods)

  if err != nil {
    log.Fatal("unable to load fiscalPeriods", err)
  }

  for i, fiscalPeriod := range fiscalPeriods {
    var positions []models.Position
    err = app.Db.Query(`SELECT *, type as position_type FROM positions WHERE fiscal_period_id = $1`, fiscalPeriod.Id).Rows(&positions)
    fiscalPeriods[i].Positions = positions
    fmt.Println("%v", positions)
  }

  b, err := json.Marshal(fiscalPeriods)
  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  if err == nil {
    if string(b) == "null" {
      io.WriteString(w, "[]")
    } else {
      io.WriteString(w, string(b))
    }
  } else {
    io.WriteString(w, "[]")
  }
}
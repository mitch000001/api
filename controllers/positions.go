package controllers

import (
	. "github.com/umsatz/api/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func (app *App) FiscalPeriodDeletePositionHandler(w http.ResponseWriter, req *http.Request, vars map[string]string) {
	log.Println("DELETE /fiscalPeriods/%v/positions/%v", vars["year"], vars["id"])

	var fiscalPeriods []FiscalPeriod
	app.Db.Query(`SELECT * FROM "fiscal_periods" WHERE year = $1 LIMIT 1`, vars["year"]).Rows(&fiscalPeriods)
	var fiscalPeriod FiscalPeriod = fiscalPeriods[0]

	var positions []Position
	if err := app.Db.Query(`SELECT * FROM positions WHERE fiscal_period_id = $1 AND id = $2`, fiscalPeriod.Id, vars["id"]).Rows(&positions); err != nil {
		log.Println("database error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var position Position = positions[0]

	if err := app.Db.Query(`DELETE FROM positions
    WHERE fiscal_period_id = (SELECT id FROM fiscal_periods WHERE year = $1 LIMIT 1) AND positions.id = $2`, vars["year"], vars["id"]).Run(); err != nil {
		log.Println("database error", err)
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprintf(`{ "errors": "%v" }`, err))
		return
	}
	os.Remove(position.AttachmentPath)

	io.WriteString(w, "")
}

func (app *App) FiscalPeriodUpdatePositionHandler(w http.ResponseWriter, req *http.Request, vars map[string]string) {
	log.Println("PUT /fiscalPeriods/%v/positions/%v", vars["year"], vars["id"])
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var fiscalPeriods []FiscalPeriod
	app.Db.Query(`SELECT * FROM "fiscal_periods" WHERE year = $1 LIMIT 1`, vars["year"]).Rows(&fiscalPeriods)
	var fiscalPeriod FiscalPeriod = fiscalPeriods[0]

	var positions []Position
	err := app.Db.Query(`SELECT * FROM positions WHERE fiscal_period_id = $1 AND id = $2`, fiscalPeriod.Id, vars["id"]).Rows(&positions)
	var position Position = positions[0]
	if err != nil {
		log.Fatal("unknown position", err)
	}

	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&position); err != nil && err != io.EOF {
		log.Fatal("decode error", err)
	}

	if !position.IsValid() {
		log.Println("INFO: unable to update position due to validation errors: %v", position.Errors)
		w.WriteHeader(http.StatusBadRequest)

		if b, err := json.Marshal(position); err == nil {
			io.WriteString(w, string(b))
		}
	}

	if err := position.StoreAttachment(fmt.Sprintf("uploads/%v", vars["year"])); err != nil {
		log.Println("INFO: unable to store attachment: %v", err)
	}

	updateError := app.Db.Query(`UPDATE "positions" SET
        category = $1,
        account_code = $2,
        type = $3,
        invoice_date = $4,
        invoice_number = $5,
        total_amount_cents = $6,
        currency = $7,
        tax = $8,
        fiscal_period_id = $9,
        description = $10,
        attachment_path = $11
        WHERE ID = $12`,
		position.Category,
		position.AccountCode,
		position.PositionType,
		time.Time(position.InvoiceDate),
		position.InvoiceNumber,
		position.TotalAmountCents,
		position.Currency,
		position.Tax,
		position.FiscalPeriodId,
		position.Description,
		position.AttachmentPath,
		position.Id).Run()

	b, err := json.Marshal(position)
	// fmt.Println(string(b))
	if err == nil && updateError == nil {
		io.WriteString(w, string(b))
	} else {
		fmt.Println("ERRRRRORRR %v, %v", err, updateError)
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "{}")
	}
}

func (app *App) FiscalPeriodCreatePositionHandler(w http.ResponseWriter, req *http.Request, vars map[string]string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	log.Println("POST /fiscalPeriods/%v/positions", vars["year"])

	var fiscalPeriods []FiscalPeriod
	app.Db.Query(`SELECT * FROM "fiscal_periods" WHERE year = $1 LIMIT 1`, vars["year"]).Rows(&fiscalPeriods)
	var fiscalPeriod FiscalPeriod = fiscalPeriods[0]

	dec := json.NewDecoder(req.Body)
	var position Position
	if err := dec.Decode(&position); err != nil && err != io.EOF {
		log.Println("decode error", err)
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprintf(`{ "errors": "%v" }`, err))
		return
	}

	position.FiscalPeriodId = fiscalPeriod.Id

	if !position.IsValid() {
		log.Println("INFO: unable to insert position due to validation errors: %+v", position.Errors)
		w.WriteHeader(http.StatusBadRequest)

		if b, err := json.Marshal(position); err == nil {
			io.WriteString(w, string(b))
		}
		return
	}

	insertError := app.Db.Query(`INSERT INTO "positions"
        (category, account_code, type, invoice_date, invoice_number, total_amount_cents, currency, tax, fiscal_period_id, description)
      VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $11) RETURNING *`,
		position.Category,
		position.AccountCode,
		position.PositionType,
		time.Time(position.InvoiceDate),
		position.InvoiceNumber,
		position.TotalAmountCents,
		position.Currency,
		position.Tax,
		position.FiscalPeriodId,
		position.Description).Rows(&position)

	if storeErr := position.StoreAttachment(fmt.Sprintf("uploads/%v", vars["year"])); storeErr != nil {
		log.Println("INFO: unable to store attachment: %v", storeErr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if position.AttachmentPath != "" {
		app.Db.Query(`UPDATE positions SET attachment_path = $1 WHERE id = $2`, position.AttachmentPath, position.Id).Run()
	}

	b, err := json.Marshal(position)

	// fmt.Println(string(b))
	if err == nil && insertError == nil {
		io.WriteString(w, string(b))
	} else {
		fmt.Println("INSERT ERRR %v, %v", err, insertError)
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "{}")
	}
}

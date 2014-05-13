package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/melvinmt/gt"
)

type ShortDate time.Time

func (date ShortDate) MarshalJSON() ([]byte, error) {
	if time.Time(date).Format("2006-01-02") == "0001-01-01" {
		return json.Marshal("")
	}
	return json.Marshal(time.Time(date).Format("2006-01-02"))
}

func (date *ShortDate) UnmarshalJSON(data []byte) (err error) {
	strDate := string(data)
	time, err := time.Parse("2006-01-02", strDate[1:len(strDate)-1])
	if err != nil {
		date = &ShortDate{}
		err = nil
	} else {
		*date = ShortDate(time)
	}
	return err
}

type Position struct {
	Id                    int                 `json:"id,omitempty"`
	AccountCodeFrom       string              `json:"accountCodeFrom"`
	AccountCodeTo         string              `json:"accountCodeTo"`
	PositionType          string              `json:"type"`
	InvoiceDate           ShortDate           `json:"invoiceDate"`
	BookingDate           ShortDate           `json:"bookingDate"`
	InvoiceNumber         string              `json:"invoiceNumber"`
	TotalAmountCents      int                 `json:"totalAmountCents"`
	TotalAmountCentsInEur int                 `json:"totalAmountCentsEur"`
	Currency              string              `json:"currency"`
	Tax                   int                 `json:"tax"`
	FiscalPeriodId        int                 `json:"fiscalPeriodId"`
	Description           string              `json:"description"`
	CreatedAt             time.Time           `json:"createdAt"`
	UpdatedAt             time.Time           `json:"updatedAt"`
	AttachmentPath        string              `json:"attachmentPath"`
	Errors                map[string][]string `json:"errors,omitempty"`
}

func (p *Position) IsValid(g *gt.Build) bool {
	p.Errors = make(map[string][]string)

	addError := func(attr string, msg string) {
		p.Errors[attr] = append(p.Errors[attr], g.T(fmt.Sprintf("validations.attribute.%s", msg)))
	}

	if p.PositionType != "income" && p.PositionType != "expense" {
		addError("type", "inclusion")
	}
	if p.Currency == "" {
		addError("currency", "missing")
	}
	if p.AccountCodeFrom == "" {
		addError("accountCodeFrom", "missing")
	}
	if p.AccountCodeTo == "" {
		addError("accountCodeTo", "missing")
	}
	if p.InvoiceDate == (ShortDate{}) {
		addError("invoiceDate", "missing")
	}
	if p.InvoiceNumber == "" {
		addError("invoiceNumber", "missing")
	}

	return len(p.Errors) == 0
}

type shortExchangeInfo struct {
	Currency string  `json:"currency"`
	Rate     float32 `json:"rate"`
}

type exchangeInfo struct {
	Date string `json:"date"`
	shortExchangeInfo
}

func setTotalAmountCentsInEur(p *Position) error {
	if p.Currency == "EUR" {
		p.TotalAmountCentsInEur = p.TotalAmountCents
	} else {
		url := fmt.Sprintf(`http://127.0.0.1:8081/%v/%v`, time.Time(p.InvoiceDate).Format("2006-01-02"), p.Currency)

		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		body, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			return err2
		}

		var v exchangeInfo
		json.Unmarshal(body, &v)
		fmt.Printf(`%#v`, v)
		p.TotalAmountCentsInEur = int(float32(p.TotalAmountCents) / v.Rate)
	}
	return nil
}

func (app *App) FiscalPeriodPositionIndexHandler(w http.ResponseWriter, req *http.Request, vars map[string]string) {
	var fiscalPeriod FiscalPeriod
	app.Db.Query(`SELECT * FROM "fiscal_periods" WHERE year = $1 LIMIT 1`, vars["year"]).Rows(&fiscalPeriod)

	var positions []Position
	if err := app.Db.Query(`SELECT *, type as position_type FROM positions WHERE fiscal_period_id = $1 ORDER BY invoice_date ASC`, fiscalPeriod.Id).Rows(&positions); err != nil {
		log.Println("database error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	b, err := json.Marshal(positions)

	if err == nil {
		io.WriteString(w, string(b))
	} else {
		fmt.Println("ERRRRRORRR %v, %v", err)
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "{}")
	}
}

func (app *App) FiscalPeriodDeletePositionHandler(w http.ResponseWriter, req *http.Request, vars map[string]string) {
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

	if !position.IsValid(app.Build) {
		log.Println("INFO: unable to update position due to validation errors: %v", position.Errors)
		w.WriteHeader(http.StatusBadRequest)

		if b, err := json.Marshal(position); err == nil {
			io.WriteString(w, string(b))
		}
	}

	if err := setTotalAmountCentsInEur(&position); err != nil {
		fmt.Println("currency lookup error %v", err)
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "{}")
		return
	}

	updateError := app.Db.Query(`UPDATE "positions" SET
        account_code_from = $1,
        account_code_to = $2,
        type = $3,
        invoice_date = $4,
        booking_date = $5,
        invoice_number = $6,
        total_amount_cents = $7,
        total_amount_cents_in_eur = $8,
        currency = $9,
        tax = $10,
        fiscal_period_id = $11,
        description = $12,
        attachment_path = $13,
        updated_at = now()::timestamptz
        WHERE ID = $14`,
		position.AccountCodeFrom,
		position.AccountCodeTo,
		position.PositionType,
		time.Time(position.InvoiceDate),
		time.Time(position.BookingDate),
		position.InvoiceNumber,
		position.TotalAmountCents,
		position.TotalAmountCentsInEur,
		position.Currency,
		position.Tax,
		position.FiscalPeriodId,
		position.Description,
		position.AttachmentPath,
		position.Id).Run()

	b, err := json.Marshal(position)

	if err == nil && updateError == nil {
		io.WriteString(w, string(b))
	} else {
		fmt.Printf(`Error updating position: %v, %v\n`, err, updateError)
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "{}")
	}
}

func (app *App) FiscalPeriodCreatePositionHandler(w http.ResponseWriter, req *http.Request, vars map[string]string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	app.SetLocale(req)

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

	if !position.IsValid(app.Build) {
		log.Println("INFO: unable to insert position due to validation errors: %+v", position.Errors)
		w.WriteHeader(http.StatusBadRequest)

		if b, err := json.Marshal(position); err == nil {
			io.WriteString(w, string(b))
		}
		return
	}

	if err := setTotalAmountCentsInEur(&position); err != nil {
		fmt.Println("currency lookup error %v", err)
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "{}")
		return
	}

	insertError := app.Db.Query(`INSERT INTO "positions"
        (account_code_from, account_code_to, type, invoice_date, booking_date, invoice_number, total_amount_cents, total_amount_cents_in_eur, currency, tax, fiscal_period_id, description, attachment_path)
      VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING *`,
		position.AccountCodeFrom,
		position.AccountCodeTo,
		position.PositionType,
		time.Time(position.InvoiceDate),
		time.Time(position.BookingDate),
		position.InvoiceNumber,
		position.TotalAmountCents,
		position.TotalAmountCentsInEur,
		position.Currency,
		position.Tax,
		position.FiscalPeriodId,
		position.Description,
		position.AttachmentPath).Rows(&position)

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

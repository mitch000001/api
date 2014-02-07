package main

import (
	 // "errors"
	_ "database/sql"
	"encoding/json"
	"fmt"
	"./models"
	"github.com/eaigner/jet"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/user"
	//  "log"
	//  "math"
	"syscall"
	"time"
)

var jetDb *jet.Db

func SetupDb() *jet.Db {
	var revDsn = os.Getenv("REV_DSN")
	if revDsn == "" {
		user, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		revDsn = "user=" + user.Username + " dbname=umsatz sslmode=disable"
	}

	newDb, err := jet.Open("postgres", revDsn)
	if err != nil {
		log.Fatal("failed to connect to postgres", err)
	}
	newDb.SetMaxIdleConns(100)

	return newDb
}

func init() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
	log.SetPrefix(fmt.Sprintf("pid:%d ", syscall.Getpid()))

	jetDb = SetupDb()
}

func FiscalPeriodDeletePositionHandler(w http.ResponseWriter, req *http.Request, vars map[string]string) {
	log.Println("DELETE /fiscalPeriods/%v/positions/%v", vars["year"], vars["id"])

	var fiscalPeriods []models.FiscalPeriod
	jetDb.Query(`SELECT * FROM "fiscal_periods" WHERE year = $1 LIMIT 1`, vars["year"]).Rows(&fiscalPeriods)
	var fiscalPeriod models.FiscalPeriod = fiscalPeriods[0]

	var positions []models.Position
	if err := jetDb.Query(`SELECT * FROM positions WHERE fiscal_period_id = $1 AND id = $2`, fiscalPeriod.Id, vars["id"]).Rows(&positions); err != nil {
		log.Println("database error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var position models.Position = positions[0]

	if err := jetDb.Query(`DELETE FROM positions
    WHERE fiscal_period_id = (SELECT id FROM fiscal_periods WHERE year = $1 LIMIT 1) AND positions.id = $2`, vars["year"], vars["id"]).Run(); err != nil {
		log.Println("database error", err)
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprintf(`{ "errors": "%v" }`, err))
		return
	}
	os.Remove( position.AttachmentPath )

	io.WriteString(w, "")
}

func FiscalPeriodUpdatePositionHandler(w http.ResponseWriter, req *http.Request, vars map[string]string) {
	log.Println("PUT /fiscalPeriods/%v/positions/%v", vars["year"], vars["id"])
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var fiscalPeriods []models.FiscalPeriod
	jetDb.Query(`SELECT * FROM "fiscal_periods" WHERE year = $1 LIMIT 1`, vars["year"]).Rows(&fiscalPeriods)
	var fiscalPeriod models.FiscalPeriod = fiscalPeriods[0]

	var positions []models.Position
	err := jetDb.Query(`SELECT * FROM positions WHERE fiscal_period_id = $1 AND id = $2`, fiscalPeriod.Id, vars["id"]).Rows(&positions)
	var position models.Position = positions[0]
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

	updateError := jetDb.Query(`UPDATE "positions" SET
        category = $1,
        account = $2,
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
		position.Account,
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

func FiscalPeriodCreatePositionHandler(w http.ResponseWriter, req *http.Request, vars map[string]string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	log.Println("POST /fiscalPeriods/%v/positions", vars["year"])

	var fiscalPeriods []models.FiscalPeriod
	jetDb.Query(`SELECT * FROM "fiscal_periods" WHERE year = $1 LIMIT 1`, vars["year"]).Rows(&fiscalPeriods)
	var fiscalPeriod models.FiscalPeriod = fiscalPeriods[0]

	dec := json.NewDecoder(req.Body)
	var position models.Position
	if err := dec.Decode(&position); err != nil && err != io.EOF {
		log.Println("decode error", err)
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprintf(`{ "errors": "%v" }`, err))
		return
	}

	position.FiscalPeriodId = fiscalPeriod.Id

	if !position.IsValid() {
		log.Println("INFO: unable to insert position due to validation errors: %v", position.Errors)
		w.WriteHeader(http.StatusBadRequest)

		if b, err := json.Marshal(position); err == nil {
			io.WriteString(w, string(b))
			// fmt.Fprint(w, string(b))
		}
		return
	}

	insertError := jetDb.Query(`INSERT INTO "positions"
        (category, account, type, invoice_date, invoice_number, total_amount_cents, currency, tax, fiscal_period_id, description)
      VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $11) RETURNING *`,
		position.Category,
		position.Account,
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
		jetDb.Query(`UPDATE positions SET attachment_path = $1 WHERE id = $2`, position.AttachmentPath, position.Id).Run()
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

func FiscalPeriodIndexHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("GET /fiscalPeriods")
	var fiscalPeriods []models.FiscalPeriod
	err := jetDb.Query(`SELECT * FROM "fiscal_periods" ORDER BY year ASC`).Rows(&fiscalPeriods)

	if err != nil {
		log.Fatal("unable to load fiscalPeriods", err)
	}

	for i, fiscalPeriod := range fiscalPeriods {
		var positions []models.Position
		err = jetDb.Query(`SELECT *, type as position_type FROM positions WHERE fiscal_period_id = $1`, fiscalPeriod.Id).Rows(&positions)
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

type RequestHandlerWithVars func(http.ResponseWriter, *http.Request, map[string]string)

func (requestHandler RequestHandlerWithVars) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	requestHandler(w, req, vars)
}

func main() {
	var port string = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	l, err := net.Listen("tcp", "0.0.0.0:"+port)
	if nil != err {
		log.Fatalln(err)
	}
	log.Println("listening on %v", l.Addr())

	r := mux.NewRouter()
	r.HandleFunc("/fiscalPeriods", FiscalPeriodIndexHandler).Methods("GET")
	r.Handle("/fiscalPeriods/{year}/positions", RequestHandlerWithVars(FiscalPeriodCreatePositionHandler)).Methods("POST")
	r.Handle("/fiscalPeriods/{year}/positions/{id}", RequestHandlerWithVars(FiscalPeriodDeletePositionHandler)).Methods("DELETE")
	r.Handle("/fiscalPeriods/{year}/positions/{id}", RequestHandlerWithVars(FiscalPeriodUpdatePositionHandler)).Methods("PUT")

	http.Handle("/", r)
	http.Serve(l, r)
}

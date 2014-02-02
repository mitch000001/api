package main

import (
	// 	"errors"
	// 	"encoding/base64"
	_ "database/sql"
	"encoding/json"
	"fmt"
	"github.com/eaigner/jet"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/user"
	// 	"io"
	// 	"log"
	// 	"math"
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

type FiscalPeriod struct {
	Id        int        `json:"-"`
	Year      int        `json:"year"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Positions []Position `json:"positions"`
}

type Position struct {
	Id 							int 				`json:"id"`
	Category  			string 			`json:"category"`
	Account   			string 			`json:"account"`
	Type 						string 			`json:"type"`
	InvoiceDate 		time.Time 	`json:"invoiceDate"`
	InvoiceNumber 	string 			`json:"invoiceNumber"`
	TotalAmount   	int 				`json:"totalAmount"`
	Currency 				string 			`json:"currency"`
	Tax							int    			`json:"tax"`
	FiscalPeriodId  int 				`json:"fiscalPeriodId"`
	Description     string 			`json:"description"`
	CreatedAt  			time.Time   `json:"createdAt"`
	UpdatedAt 			time.Time   `json:"updatedAt"`
}

func FiscalPeriodDeletePositionHandler(w http.ResponseWriter, req *http.Request, vars map[string]string) {
	jetDb.Query(`DELETE FROM "positions"
		INNER JOIN "fiscal_periods" ON "fiscal_periods".id = "positions.fiscal_period_id"
		WHERE "fiscal_periods".year = $1 AND "positions".id = $2 LIMIT 1`, vars["year"], vars["id"]).Run()
	io.WriteString(w, "")
}

func FiscalPeriodCreatePositionHandler(w http.ResponseWriter, req *http.Request, vars map[string]string) {
	var fiscalPeriods []FiscalPeriod
	jetDb.Query(`SELECT * FROM "fiscal_periods" WHERE year = $1 LIMIT 1`, vars["year"]).Rows(&fiscalPeriods)
	var fiscalPeriod FiscalPeriod = fiscalPeriods[0]

	dec := json.NewDecoder(req.Body)
	var position Position
	if err := dec.Decode(&position); err != nil && err != io.EOF {
		log.Fatal("decode error", err)
	}

	position.FiscalPeriodId = fiscalPeriod.Id

	err := jetDb.Query(`INSERT INTO "positions"
				(category, account, type, invoice_date, invoice_number, total_amount, currency, tax, fiscal_period_id, description)
			VALUES
				($1			 , $2			, $3	, $4					, $5						, $6					, $7			, $8 , $9							 , $10)`,
			position.Category,
			position.Account,
			position.Type,
			position.InvoiceDate,
			position.InvoiceNumber,
			position.TotalAmount,
			position.Currency,
			position.Tax,
			position.FiscalPeriodId,
			position.Description).Run()

	if err != nil {
		log.Fatal(err)
	}

	b, err := json.Marshal(fiscalPeriod)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err == nil {
		io.WriteString(w, string(b))
	} else {
		io.WriteString(w, "{}")
	}
}

func FiscalPeriodIndexHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("GET /fiscalPeriods")
	var fiscalPeriods []FiscalPeriod
	err := jetDb.Query(`SELECT * FROM "fiscal_periods" ORDER BY year ASC`).Rows(&fiscalPeriods)

	if err != nil {
		log.Fatal("unable to load fiscalPeriods", err)
	}

	for i, fiscalPeriod := range fiscalPeriods {
		var positions []Position
		err = jetDb.Query(`SELECT * FROM positions WHERE fiscal_period_id = $1`, fiscalPeriod.Id).Rows(&positions)
		fiscalPeriods[i].Positions = positions
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
	log.Printf("listening on %v", l.Addr())

	r := mux.NewRouter()
	r.HandleFunc("/fiscalPeriods", FiscalPeriodIndexHandler).Methods("GET")
	r.Handle("/fiscalPeriods/{year}/positions",  RequestHandlerWithVars(FiscalPeriodCreatePositionHandler)).Methods("POST")
	r.Handle("/fiscalPeriods/{year}/positions/{id}",  RequestHandlerWithVars(FiscalPeriodDeletePositionHandler)).Methods("DELETE")

	http.Handle("/", r)
	http.Serve(l, r)
}

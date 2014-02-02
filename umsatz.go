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
	Id        int       `json:"-"`
	Year      int       `json:"year"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func FiscalPeriodIndexHandler(w http.ResponseWriter, req *http.Request) {
	var fiscalPeriods []FiscalPeriod
	err := jetDb.Query(`SELECT * FROM "fiscal_periods" ORDER BY year ASC`).Rows(&fiscalPeriods)

	if err != nil {
		log.Fatal("unable to load fiscalPeriods", err)
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

	http.Handle("/", r)
	http.Serve(l, r)
}

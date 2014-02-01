package main

import (
	// 	"errors"
	// 	"encoding/base64"
	// 	"encoding/json"
	// 	"fmt"
	"database/sql"
	"github.com/eaigner/hood"
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
	"time"
)

type FiscalPeriod struct {
  Id        hood.Id   `json:"-"`
  Year      int 	    `json:"year"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
}

func FiscalPeriodIndexHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "[]")
}

var db *sql.DB
var hd *hood.Hood

func SetupHood() *hood.Hood {
	var revDsn = os.Getenv("REV_DSN")
	if revDsn == "" {
		user, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		revDsn = "user=" + user.Username + " dbname=revisioneer sslmode=disable"
	}

	var err error
	db, err = sql.Open("postgres", revDsn)
	if err != nil {
		log.Fatal("failed to connect to postgres", err)
	}
	db.SetMaxIdleConns(100)

	newHd := hood.New(db, hood.NewPostgres())
	newHd.Log = true
	return newHd
}

func init() {
	SetupHood()
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
	s := r.PathPrefix("/timeframes").Subrouter()
	s.HandleFunc("/fiscalPeriods", FiscalPeriodIndexHandler).
		Methods("GET")

	http.Handle("/", r)
	http.Serve(l, r)
}

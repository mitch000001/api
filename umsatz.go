package main

import (
	_ "database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/melvinmt/gt"
	. "github.com/umsatz/api/controllers"
)

var app *App

func I18nInit() *gt.Build {
	bytes, err := ioutil.ReadFile("locales/all.json")
	if err != nil {
		panic("error reading locales")
	}

	dec := json.NewDecoder(strings.NewReader(string(bytes)))
	var locales gt.Strings
	if err := dec.Decode(&locales); err != nil && err != io.EOF {
		panic("unable to parse json")
	}

	g := &gt.Build{
		Index:  locales,
		Origin: "en",
	}
	g.SetTarget("de")
	return g
}

func init() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
	log.SetPrefix(fmt.Sprintf("pid:%d ", syscall.Getpid()))

	app = &App{SetupDb(), I18nInit()}
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
	r.HandleFunc("/accounts", app.AccountIndexHandler).Methods("GET")
	r.HandleFunc("/accounts", app.CreateAccountHandler).Methods("POST")
	r.HandleFunc("/fiscalPeriods", app.FiscalPeriodIndexHandler).Methods("GET")
	r.Handle("/fiscalPeriods/{year}/positions", RequestHandlerWithVars(app.FiscalPeriodPositionIndexHandler)).Methods("GET")
	r.Handle("/fiscalPeriods/{year}/positions", RequestHandlerWithVars(app.FiscalPeriodCreatePositionHandler)).Methods("POST")
	r.Handle("/fiscalPeriods/{year}/positions/{id}", RequestHandlerWithVars(app.FiscalPeriodDeletePositionHandler)).Methods("DELETE")
	r.Handle("/fiscalPeriods/{year}/positions/{id}", RequestHandlerWithVars(app.FiscalPeriodUpdatePositionHandler)).Methods("PUT")

	http.Handle("/", r)
	http.Serve(l, r)
}

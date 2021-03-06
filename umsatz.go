package main

import (
	_ "database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/melvinmt/gt"
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

func logHandler(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v %v", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	}
}

func routingHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/vnd.umsatz+json; charset=utf-8")

	var routes []Link = []Link{
		NewLink(&req.Header, "index.backups", "/backups"),
		NewLink(&req.Header, "index.accounts", "/accounts"),
		NewLink(&req.Header, "index.fiscalPeriods", "/fiscalPeriods"),
	}

	bytes, err := json.Marshal(routes)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	io.WriteString(w, string(bytes))
}

func main() {
	var (
		httpAddress = flag.String("http.addr", ":8080", "HTTP listen address")
	)
	flag.Parse()

	r := mux.NewRouter()
	r.Handle("/accounts", logHandler(http.HandlerFunc(app.AccountIndexHandler))).Methods("GET")
	r.Handle("/accounts", logHandler(http.HandlerFunc(app.CreateAccountHandler))).Methods("POST")
	r.Handle("/accounts/{id}", logHandler(RequestHandlerWithVars(app.UpdateAccountHandler))).Methods("PUT")
	r.Handle("/fiscalPeriods", logHandler(http.HandlerFunc(app.FiscalPeriodIndexHandler))).Methods("GET")
	r.Handle("/fiscalPeriods/{year}/positions", logHandler(RequestHandlerWithVars(app.FiscalPeriodPositionIndexHandler))).Methods("GET")
	r.Handle("/fiscalPeriods/{year}/positions", logHandler(RequestHandlerWithVars(app.FiscalPeriodCreatePositionHandler))).Methods("POST")
	r.Handle("/fiscalPeriods/{year}/positions/{id}", logHandler(RequestHandlerWithVars(app.FiscalPeriodDeletePositionHandler))).Methods("DELETE")
	r.Handle("/fiscalPeriods/{year}/positions/{id}", logHandler(RequestHandlerWithVars(app.FiscalPeriodUpdatePositionHandler))).Methods("PUT")
	r.Handle("/", logHandler(http.HandlerFunc(routingHandler))).Methods("GET")

	http.Handle("/", r)
	log.Printf("listening on %s", *httpAddress)
	log.Fatal(http.ListenAndServe(*httpAddress, http.Handler(r)))
}

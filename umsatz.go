package main

import (
	. "./controllers"
	_ "database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
	"os"
	"syscall"
)

var app *App

func init() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
	log.SetPrefix(fmt.Sprintf("pid:%d ", syscall.Getpid()))

	app = &App{}
	app.SetupDb()
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
	r.HandleFunc("/fiscalPeriods", app.FiscalPeriodIndexHandler).Methods("GET")
	r.Handle("/fiscalPeriods/{year}/positions", RequestHandlerWithVars(app.FiscalPeriodCreatePositionHandler)).Methods("POST")
	r.Handle("/fiscalPeriods/{year}/positions/{id}", RequestHandlerWithVars(app.FiscalPeriodDeletePositionHandler)).Methods("DELETE")
	r.Handle("/fiscalPeriods/{year}/positions/{id}", RequestHandlerWithVars(app.FiscalPeriodUpdatePositionHandler)).Methods("PUT")

	http.Handle("/", r)
	http.Serve(l, r)
}

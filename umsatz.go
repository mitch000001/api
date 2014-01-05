package main

import (
	// "encoding/json"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net"
	"net/http"
	"os"
)

func TimeframesIndexHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "[]")
}

func TimeframesCreationHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "[]")
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
	s.HandleFunc("/", TimeframesIndexHandler).
		Methods("GET")
	s.HandleFunc("/", TimeframesCreationHandler).
		Methods("POST")

	http.Handle("/", r)
	http.Serve(l, r)
}

package main

import (
	// "encoding/json"
	_ "database/sql"
	"github.com/eaigner/hood"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

// import (
// 	"crypto/rand"
// 	"encoding/base64"
// 	"encoding/json"
// 	"errors"
// 	"fmt"

// 	"io"
// 	"log"
// 	"math"
// 	"net"
// 	"net/http"
// 	"os"
// 	"os/user"
// 	"strconv"
// 	"syscall"
// )

type FiscalPeriod struct {
  Id        hood.Id   `json:"-"`
  Year      int 	    `json:"year"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
}

func FiscalPeriodIndexHandler(w http.ResponseWriter, req *http.Request) {
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
	s.HandleFunc("/fiscalPeriods", FiscalPeriodIndexHandler).
		Methods("GET")

	http.Handle("/", r)
	http.Serve(l, r)
}

package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/umsatz/api/models"
)

func (app *App) AccountIndexHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	log.Println("GET /accounts")

	var accounts []models.Account
	if err := app.Db.Query(`SELECT * FROM accounts ORDER BY code ASC`).Rows(&accounts); err != nil {
		log.Println("database error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bytes, err := json.Marshal(accounts)
	if err != nil {
		log.Println("json marshal error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	io.WriteString(w, string(bytes))
}

func (app *App) UpdateAccountHandler(w http.ResponseWriter, req *http.Request, vars map[string]string) {
	log.Println("PUT /accounts/%v", vars["id"])
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var account models.Account
	err := app.Db.Query(`SELECT * FROM accounts WHERE id = $1`, vars["id"]).Rows(&account)
	if err != nil {
		log.Fatal("unknown account", err)
	}

	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&account); err != nil && err != io.EOF {
		log.Fatal("decode error", err)
	}

	if !account.IsValid() {
		log.Println("INFO: unable to update account due to validation errors: %v", account.Errors)
		w.WriteHeader(http.StatusBadRequest)

		if b, err := json.Marshal(account); err == nil {
			io.WriteString(w, string(b))
		}
	}

	updateError := app.Db.Query(`UPDATE "accounts" SET
        code = $1,
        label = $2,
        WHERE ID = $3`,
		account.Code,
		account.Label,
		account.Id).Run()

	b, err := json.Marshal(account)

	if err == nil {
		io.WriteString(w, string(b))
	} else {
		fmt.Println("ERRRRRORRR %v, %v", err, updateError)
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "{}")
	}
}

func (app *App) CreateAccountHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	log.Println("POST /accounts")

	dec := json.NewDecoder(req.Body)
	var account models.Account
	if err := dec.Decode(&account); err != nil && err != io.EOF {
		log.Println("decode error", err)
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprintf(`{ "errors": "%v" }`, err))
		return
	}

	if !account.IsValid() {
		log.Println("INFO: unable to insert account due to validation errors: %+v", account.Errors)
		w.WriteHeader(http.StatusBadRequest)

		if b, err := json.Marshal(account); err == nil {
			io.WriteString(w, string(b))
		}
		return
	}

	insertError := app.Db.Query(`INSERT INTO "accounts"
        (code, label)
      VALUES ($1, $2) RETURNING *`,
		account.Code,
		account.Label).Rows(&account)

	b, err := json.Marshal(account)

	// fmt.Println(string(b))
	if err == nil && insertError == nil {
		io.WriteString(w, string(b))
	} else {
		fmt.Println("INSERT ERRR %v, %v", err, insertError)
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "{}")
	}
}

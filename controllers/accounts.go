package controllers

import (
	"../models"
	"encoding/json"
	"io"
	"log"
	"net/http"
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

package main

import (
	"log"
	"net/http"
	"os"
	"os/user"
	"strings"

	_ "github.com/lib/pq"
	"github.com/melvinmt/gt"
	"github.com/splicers/jet"
)

type Link struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

type App struct {
	Db *jet.Db
	*gt.Build
}

func SetupDb() *jet.Db {
	database := os.Getenv("DATABASE")
	if database == "" {
		database = "umsatz"
	}

	revDsn := os.Getenv("REV_DSN")
	if revDsn == "" {
		user, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		revDsn = "user=" + user.Username + " dbname=" + database + " sslmode=disable"
	}

	newDb, err := jet.Open("postgres", revDsn)
	if err != nil {
		log.Fatal("failed to connect to postgres", err)
	}
	newDb.SetMaxIdleConns(100)

	return newDb
}

func (app *App) SetLocale(req *http.Request) {
	if len(req.Header["Accept-Language"]) > 0 {
		locale := strings.Split(req.Header["Accept-Language"][0], "-")[0]
		app.SetTarget(locale)
	}
}

func (app *App) ClearDb() {
	app.Db.Query("DELETE FROM accounts").Run()
	app.Db.Query("DELETE FROM positions").Run()
	app.Db.Query("DELETE FROM fiscal_periods").Run()
}

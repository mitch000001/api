package controllers

import (
	"github.com/eaigner/jet"
	_ "github.com/lib/pq"
	"log"
	"os"
	"os/user"
)

type App struct {
	Db *jet.Db
}

func (app *App) SetupDb() *jet.Db {
	if app.Db != nil {
		return app.Db
	}

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

	app.Db = newDb
	return app.Db
}

func (app *App) ClearDb() {
	app.Db.Query("DELETE FROM accounts").Run()
	app.Db.Query("DELETE FROM positions").Run()
	app.Db.Query("DELETE FROM fiscal_periods").Run()
}

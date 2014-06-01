package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/user"
	"strings"

	_ "github.com/lib/pq"
	"github.com/melvinmt/gt"
	"github.com/splicers/jet"
)

// Hypermedia link structure
type Link struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

func BaseUri(h *http.Header) string {
	baseUri := h.Get("X-Requested-Uri")
	if strings.HasSuffix(baseUri, "/") {
		baseUri = baseUri[:len(baseUri)-1]
	}
	return baseUri
}

func NewLink(h *http.Header, rel string, href string) Link {
	absoluteHref := fmt.Sprintf("%v%v", BaseUri(h), href)
	return Link{rel, absoluteHref}
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

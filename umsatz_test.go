package main

import (
	"./models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	// "io/ioutil"
	"fmt"
)

func init() {
	jetDb = SetupDb()
}

func ClearDb() {
	jetDb.Query("DELETE FROM positions").Run()
	jetDb.Query("DELETE FROM fiscal_periods").Run()
}
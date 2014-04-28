package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/melvinmt/gt"
)

type ShortDate time.Time

func (date ShortDate) MarshalJSON() ([]byte, error) {
	if time.Time(date).Format("2006-01-02") == "0001-01-01" {
		return json.Marshal("")
	}
	return json.Marshal(time.Time(date).Format("2006-01-02"))
}

func (date *ShortDate) UnmarshalJSON(data []byte) (err error) {
	strDate := string(data)
	time, err := time.Parse("2006-01-02", strDate[1:len(strDate)-1])
	if err != nil {
		date = &ShortDate{}
		err = nil
	} else {
		*date = ShortDate(time)
	}
	return err
}

type Position struct {
	Id               int                 `json:"id,omitempty"`
	AccountCodeFrom  string              `json:"accountCodeFrom"`
	AccountCodeTo    string              `json:"accountCodeTo"`
	PositionType     string              `json:"type"`
	InvoiceDate      ShortDate           `json:"invoiceDate"`
	BookingDate      ShortDate           `json:"bookingDate"`
	InvoiceNumber    string              `json:"invoiceNumber"`
	TotalAmountCents int                 `json:"totalAmountCents"`
	Currency         string              `json:"currency"`
	Tax              int                 `json:"tax"`
	FiscalPeriodId   int                 `json:"fiscalPeriodId"`
	Description      string              `json:"description"`
	CreatedAt        time.Time           `json:"createdAt"`
	UpdatedAt        time.Time           `json:"updatedAt"`
	AttachmentPath   string              `json:"attachmentPath"`
	Errors           map[string][]string `json:"errors,omitempty"`
}

func (p *Position) IsValid(g *gt.Build) bool {
	p.Errors = make(map[string][]string)

	addError := func(attr string, msg string) {
		p.Errors[attr] = append(p.Errors[attr], g.T(fmt.Sprintf("validations.attribute.%s", msg)))
	}

	if p.PositionType != "income" && p.PositionType != "expense" {
		addError("type", "inclusion")
	}
	if p.Currency == "" {
		addError("currency", "missing")
	}
	if p.AccountCodeFrom == "" {
		addError("accountCodeFrom", "missing")
	}
	if p.AccountCodeTo == "" {
		addError("accountCodeTo", "missing")
	}
	if p.InvoiceDate == (ShortDate{}) {
		addError("invoiceDate", "missing")
	}
	if p.InvoiceNumber == "" {
		addError("invoiceNumber", "missing")
	}

	return len(p.Errors) == 0
}

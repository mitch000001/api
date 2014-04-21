package models

import (
	"encoding/json"
	"time"
)

type ShortDate time.Time

func (date ShortDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(date).Format("2006-01-02"))
}

func (date *ShortDate) UnmarshalJSON(data []byte) (err error) {
	strDate := string(data)
	time, err := time.Parse("2006-01-02", strDate[1:len(strDate)-1])
	*date = ShortDate(time)
	return err
}

type Position struct {
	Id                   int       `json:"id,omitempty"`
	Category             string    `json:"category"`
	AccountCode          string    `json:"accountCode"`
	PositionType         string    `json:"type"`
	InvoiceDate          ShortDate `json:"invoiceDate"`
	InvoiceNumber        string    `json:"invoiceNumber"`
	TotalAmountCents     int       `json:"totalAmountCents"`
	Currency             string    `json:"currency"`
	Tax                  int       `json:"tax"`
	FiscalPeriodId       int       `json:"fiscalPeriodId"`
	Description          string    `json:"description"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
	EncodedFileExtension string    `json:"encodedFileExtension,omitempty"`
	EncodedAttachment    string    `json:"encodedAttachment,omitempty"`
	AttachmentPath       string    `json:"attachmentPath"`
	Errors               []string  `json:"errors,omitempty"`
}

func (p *Position) IsValid() bool {
	p.Errors = make([]string, 0)

	if p.PositionType != "income" && p.PositionType != "expense" {
		p.AddError("type", "must be either income or expense")
	}
	if p.Currency == "" {
		p.AddError("currency", "must be present")
	}
	if p.AccountCode == "" {
		p.AddError("accountCode", "must be present")
	}
	if p.InvoiceDate == (ShortDate{}) {
		p.AddError("invoiceDate", "must be present")
	}
	if p.InvoiceNumber == "" {
		p.AddError("invoiceNumber", "must be present")
	}

	return len(p.Errors) == 0
}

func (p *Position) AddError(attr string, errorMsg string) {
	p.Errors = append(p.Errors, attr+":"+errorMsg)
}
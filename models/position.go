package umsatz

import (
  "time"
  "encoding/json"
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
  Id             int       `json:"id"`
  Category       string    `json:"category"`
  Account        string    `json:"account"`
  PositionType   string    `json:"type"`
  InvoiceDate    ShortDate `json:"invoiceDate"`
  InvoiceNumber  string    `json:"invoiceNumber"`
  TotalAmount    int       `json:"totalAmount"`
  Currency       string    `json:"currency"`
  Tax            int       `json:"tax"`
  FiscalPeriodId int       `json:"fiscalPeriodId"`
  Description    string    `json:"description"`
  CreatedAt      time.Time `json:"createdAt"`
  UpdatedAt      time.Time `json:"updatedAt"`
  Errors       []string    `json:"errors,omitempty"`
}

func (p *Position) IsValid() (bool) {
  p.Errors = make([]string, 0)

  if p.PositionType != "income" && p.PositionType != "expense" {
    p.AddError("type", "must be either income or expense")
  }
  if p.Category == "" {
    p.AddError("category", "must be present")
  }
  if p.Currency == "" {
    p.AddError("currency", "must be present")
  }

  return len(p.Errors) == 0
}

func (p *Position) AddError(attr string, errorMsg string) () {
  p.Errors = append(p.Errors, attr + ":" + errorMsg)
}
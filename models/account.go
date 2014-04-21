package models

type Account struct {
  Code        string        `json:"code"`
  Label       string        `json:"label"`
  Errors      []string      `json:"errors,omitempty"`
}

func (a *Account) AddError(attr string, errorMsg string) {
  a.Errors = append(a.Errors, attr+":"+errorMsg)
}

func (a *Account) IsValid() bool {
  a.Errors = make([]string, 0)

  if a.Code == "" {
    a.AddError("code", "must be present")
  }
  if a.Label == "" {
    a.AddError("label", "must be present")
  }

  return len(a.Errors) == 0
}
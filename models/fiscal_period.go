package models

import (
  "time"
)

type FiscalPeriod struct {
  Id        int        `json:"-"`
  Year      int        `json:"year"`
  CreatedAt time.Time  `json:"created_at"`
  UpdatedAt time.Time  `json:"updated_at"`
  Positions []Position `json:"positions"`
}
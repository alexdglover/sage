package models

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	Date        string
	Description string
	Amount      decimal.Decimal
	Excluded    string
	AccountId   uint
	Account     Account
	CategoryId  uint
	Category    Category
}

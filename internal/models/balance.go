package models

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Balance struct {
	gorm.Model
	Date               string
	EffectiveStartDate string
	EffectiveEndDate   *string
	Balance            decimal.Decimal
	AccountId          uint
	Account            Account
}

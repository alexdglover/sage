package models

import "gorm.io/gorm"

type Balance struct {
	gorm.Model
	Date               string
	EffectiveStartDate string
	EffectiveEndDate   *string
	Amount             int64
	AccountID          uint
	Account            Account
}

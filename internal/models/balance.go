package models

import "gorm.io/gorm"

type Balance struct {
	gorm.Model
	Date               string
	EffectiveStartDate string
	EffectiveEndDate   *string
	Amount             int64
	AccountId          uint
	Account            Account
}

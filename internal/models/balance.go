package models

import "gorm.io/gorm"

type Balance struct {
	gorm.Model
	Date               string
	EffectiveStartDate string
	EffectiveEndDate   *string
	Balance            int64
	AccountId          uint
	Account            Account
}

package models

import "gorm.io/gorm"

type Balance struct {
	gorm.Model
	Date               string
	EffectiveStartDate string
	EffectiveEndDate   string
	Balance            float32
	AccountId          uint
	Account            Account
}

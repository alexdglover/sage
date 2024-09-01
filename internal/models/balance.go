package models

import "gorm.io/gorm"

type Balance struct {
	gorm.Model
	EffectiveDate string
	Amount        int64
	AccountID     uint
	Account       Account
}

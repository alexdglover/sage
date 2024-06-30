package models

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	Name             string
	AccountType      string
	AssetOrLiability string
}

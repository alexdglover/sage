package models

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	Name            string
	AccountCategory string
	AccountType     string
}

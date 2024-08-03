package models

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	Date        string
	Description string
	Amount      int64
	Excluded    string
	AccountId   uint
	Account     Account
	CategoryId  uint
	Category    Category
}

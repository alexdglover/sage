package models

import (
	"gorm.io/gorm"
)

type AccountType struct {
	gorm.Model
	Name            string
	LedgerType      string  // asset vs liability
	AccountCategory string  // checking, brokerage, credit card, loan, etc
	DefaultParser   *string // institution-specific parser for CSVs
}

type AccountTypeRepository struct {
	DB *gorm.DB
}

func (atr *AccountTypeRepository) GetAllAccountTypes() ([]AccountType, error) {
	var accountTypes []AccountType
	result := atr.DB.Find(&accountTypes)
	return accountTypes, result.Error
}

func (atr *AccountTypeRepository) GetAccountTypeByID(id uint) (AccountType, error) {
	var accountType AccountType
	result := atr.DB.Where("id = ?", id).First(&accountType)
	return accountType, result.Error
}

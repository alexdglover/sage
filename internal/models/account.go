package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Account struct {
	gorm.Model
	Name          string
	AccountTypeID uint
	AccountType   AccountType
}

type AccountRepository struct {
	DB *gorm.DB
}

func (ar *AccountRepository) GetAllAccounts() ([]Account, error) {
	var accounts []Account
	result := ar.DB.Find(&accounts)
	return accounts, result.Error
}

func (ar *AccountRepository) GetAccountByID(id uint) (Account, error) {
	var account Account
	result := ar.DB.Preload(clause.Associations).Where("id = ?", id).Find(&account)
	return account, result.Error
}

// Save is an UPSERT operation, returning the ID of the record and an optional error
func (ar *AccountRepository) Save(account Account) (id uint, err error) {
	result := ar.DB.Save(&account).Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}})
	return account.ID, result.Error
}

// Soft deletes an account and all associated transactions and balances
func (ar *AccountRepository) DeleteAccountByID(accountID uint) (err error) {
	ar.DB.Where("account_id = ?", accountID).Delete(&Balance{})
	ar.DB.Where("account_id = ?", accountID).Delete(&Transaction{})
	result := ar.DB.Delete(&Account{}, accountID)
	return result.Error
}

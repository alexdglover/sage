package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Account struct {
	gorm.Model
	Name            string
	AccountCategory string
	AccountType     string
	DefaultParser   *string
}

type AccountRepository struct{}

var accountRepository *AccountRepository

func GetAccountRepository() *AccountRepository {
	if accountRepository == nil {
		accountRepository = &AccountRepository{}
	}
	return accountRepository
}

func (ar *AccountRepository) GetAllAccounts() ([]Account, error) {
	var accounts []Account
	result := db.Find(&accounts)
	return accounts, result.Error
}

func (ar *AccountRepository) GetAccountByID(id uint) (Account, error) {
	var account Account
	result := db.Where("id = ?", id).First(&account)
	return account, result.Error
}

// Save is an UPSERT operation, returning the ID of the record and an optional error
func (*AccountRepository) Save(account Account) (id uint, err error) {
	result := db.Save(&account).Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}})
	return account.ID, result.Error
}

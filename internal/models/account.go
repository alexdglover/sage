package models

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	Name            string
	AccountCategory string
	AccountType     string
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

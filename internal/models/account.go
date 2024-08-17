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

func (ar *AccountRepository) GetByID(id uint) (Account, error) {
	var account Account
	result := db.Where("id = ?", id).First(&account)
	return account, result.Error
}

package models

import (
	"errors"
	"fmt"
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
	if account.ID != 0 {
		var existingAccount Account
		result := ar.DB.Where("id = ?", account.ID).First(&existingAccount)
		if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			fmt.Printf("Check Error: %v\n", result.Error)
			return 0, result.Error
		}
		if result.RowsAffected == 0 {
			return 0, fmt.Errorf("Account with ID %d not found", account.ID)
		}
	}
	result := ar.DB.Save(&account)
	// FOR DEBUGGING PURPOSES ONLY
	// result := ar.DB.Debug().Save(&account)
	if result.Error != nil {
		fmt.Printf("Save Error: %v\n", result.Error)
		return 0, result.Error
	}
	return account.ID, result.Error
}

// Soft deletes an account and all associated transactions and balances
func (ar *AccountRepository) DeleteAccountByID(accountID uint) (err error) {
	ar.DB.Where("account_id = ?", accountID).Delete(&Balance{})
	ar.DB.Where("account_id = ?", accountID).Delete(&Transaction{})
	result := ar.DB.Delete(&Account{}, accountID)
	return result.Error
}

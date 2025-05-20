package services

import (
	"github.com/alexdglover/sage/internal/models"
)

type AccountNameAndID struct {
	AccountName string
	AccountID   uint
}

type AccountManager struct {
	AccountRepository *models.AccountRepository
}

func (am *AccountManager) GetAccountNamesAndIDs() (result []AccountNameAndID, err error) {
	account, err := am.AccountRepository.GetAllAccounts()
	if err != nil {
		return nil, err
	}
	for _, acc := range account {
		result = append(result, AccountNameAndID{AccountName: acc.Name, AccountID: acc.ID})
	}
	return result, nil
}

// GetAccountTypeByID returns the account type for a given account ID
// If the account type is not found, it returns an error
// 
// Arguments:
// id - The ID of the account type to retrieve
//
// Returns:
// - The account type object
// - An error if the account type is not found or if there was an error retrieving it
func (am *AccountManager) GetAccountTypeByID(id uint) (models.AccountType, error) {
    var accountType models.AccountType
    if err := am.AccountRepository.DB.Where("id = ?", id).First(&accountType).Error; err != nil {
        return models.AccountType{}, err
    }
    return accountType, nil
}
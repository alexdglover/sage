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

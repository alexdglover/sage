package services

import (
	"github.com/alexdglover/sage/internal/models"
)

type AccountNameAndID struct {
	AccountName string
	AccountID   uint
}

func GetAccountNamesAndIDs() (result []AccountNameAndID, err error) {
	ar := models.GetAccountRepository()
	account, err := ar.GetAllAccounts()
	if err != nil {
		return nil, err
	}
	for _, acc := range account {
		result = append(result, AccountNameAndID{AccountName: acc.Name, AccountID: acc.ID})
	}
	return result, nil
}

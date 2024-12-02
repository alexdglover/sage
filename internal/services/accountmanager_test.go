package services

import (
	"fmt"

	"github.com/alexdglover/sage/internal/models"
	"github.com/alexdglover/sage/internal/utils"
)

type TestAccountRepository struct{}

func (tar TestAccountRepository) GetAllAccounts() ([]models.Account, error) {
	var accounts []models.Account
	accounts = append(accounts, models.Account{
		// ID:              666,
		Name:            "foo",
		AccountCategory: "Checking",
		AccountType:     "asset",
		DefaultParser:   utils.StrPointer("schwabChecking"),
	})
	return accounts, nil
}

func test_getAccountNamesAndIds() error {
	am := AccountManager{
		AccountRepository: TestAccountRepository{},
	}
	results, err := am.GetAccountNamesAndIDs()

	if err != nil {
		fmt.Println("test failed")
	}

	if len(results) != 1 {
		fmt.Println("test failed")
	}

	result := results[0]
	if result.AccountName != "foo" {
		fmt.Println("test failed")
	}
	if result.AccountID != 666 {
		fmt.Println("test failed")
	}

	return nil
}

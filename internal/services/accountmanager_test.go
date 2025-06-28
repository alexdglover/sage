package services

import (
	"errors"
	"testing"

	"github.com/alexdglover/sage/internal/models"
	"gorm.io/gorm"
)

func (m *MockAccountRepository) GetAllAccounts() ([]models.Account, error) {
	return m.Accounts, m.Err
}

type MockAccountTypeRepository struct {
	AccountType    models.AccountType
	AccountTypeErr error
}

func (m *MockAccountTypeRepository) GetAccountTypeByID(id uint) (models.AccountType, error) {
	if m.AccountTypeErr != nil {
		return models.AccountType{}, m.AccountTypeErr
	}
	return m.AccountType, nil
}

func TestGetAccountNamesAndIDs(t *testing.T) {
	accounts := []models.Account{
		{Model: gorm.Model{ID: 1}, Name: "Checking"},
		{Model: gorm.Model{ID: 2}, Name: "Savings"},
	}
	repo := &MockAccountRepository{Accounts: accounts}
	am := &AccountManager{AccountRepository: repo}
	result, err := am.GetAccountNamesAndIDs()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 accounts, got %d", len(result))
	}
	if result[0].AccountName != "Checking" || result[1].AccountName != "Savings" {
		t.Errorf("unexpected account names: %+v", result)
	}
}

func TestGetAccountNamesAndIDs_Error(t *testing.T) {
	repo := &MockAccountRepository{Err: errors.New("fail")}
	am := &AccountManager{AccountRepository: repo}
	_, err := am.GetAccountNamesAndIDs()
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetAccountTypeByID(t *testing.T) {
	typeRepo := &MockAccountTypeRepository{AccountType: models.AccountType{Model: gorm.Model{ID: 3}, Name: "Brokerage"}}
	am := &AccountManager{AccountTypeRepository: typeRepo}
	at, err := am.AccountTypeRepository.GetAccountTypeByID(3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if at.Name != "Brokerage" {
		t.Errorf("expected Brokerage, got %s", at.Name)
	}
}

func TestGetAccountTypeByID_Error(t *testing.T) {
	typeRepo := &MockAccountTypeRepository{AccountTypeErr: errors.New("not found")}
	am := &AccountManager{AccountTypeRepository: typeRepo}
	_, err := am.AccountTypeRepository.GetAccountTypeByID(99)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

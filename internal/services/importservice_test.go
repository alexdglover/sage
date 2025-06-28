package services

import (
	"errors"
	"testing"

	"github.com/alexdglover/sage/internal/models"
)

func TestImportStatement_Success(t *testing.T) {
	parserName := "mock"
	account := models.Account{Name: "Test Account", AccountTypeID: 1, AccountType: models.AccountType{DefaultParser: &parserName}}
	parsersByInstitution[parserName] = &MockParser{
		Txns:     []models.Transaction{{Amount: 100, Date: "2024-01-01", Description: "Test txn"}},
		Balances: []models.Balance{{Amount: 1000}},
	}
	is := &ImportService{
		AccountRepository:          &MockAccountRepository{Account: account},
		BalanceRepository:          &MockBalanceRepository{},
		ImportSubmissionRepository: &MockImportSubmissionRepository{},
		TransactionRepository:      &MockTransactionRepository{TxnsByHash: map[string][]models.Transaction{}},
		Categorizer:                &MockCategorizer{Category: models.Category{Name: "Test Category"}},
	}
	res, err := is.ImportStatement("file.csv", "statement", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == nil {
		t.Error("expected result, got nil")
	}
}

func TestImportStatement_AccountNotFound(t *testing.T) {
	is := &ImportService{
		AccountRepository:          &MockAccountRepository{Err: errors.New("not found")},
		ImportSubmissionRepository: &MockImportSubmissionRepository{},
	}
	_, err := is.ImportStatement("file.csv", "statement", 1)
	if err == nil || err.Error() != "Could not find an account with ID 0" {
		t.Errorf("expected AccountNotFoundError, got %v", err)
	}
}

func TestImportStatement_NoParser(t *testing.T) {
	account := models.Account{Name: "Test Account", AccountTypeID: 1, AccountType: models.AccountType{DefaultParser: nil}}
	is := &ImportService{
		AccountRepository:          &MockAccountRepository{Account: account},
		ImportSubmissionRepository: &MockImportSubmissionRepository{},
	}
	_, err := is.ImportStatement("file.csv", "statement", 1)
	if err == nil || err.Error() != "No parser was found for the provided account" {
		t.Errorf("expected NoParserError, got %v", err)
	}
}

func TestImportStatement_ParserError(t *testing.T) {
	parserName := "mock"
	account := models.Account{Name: "Test Account", AccountTypeID: 1, AccountType: models.AccountType{DefaultParser: &parserName}}
	parsersByInstitution[parserName] = &MockParser{ParseErr: errors.New("parse fail")}
	is := &ImportService{
		AccountRepository:          &MockAccountRepository{Account: account},
		ImportSubmissionRepository: &MockImportSubmissionRepository{},
	}
	_, err := is.ImportStatement("file.csv", "statement", 1)
	if err == nil || err.Error() != "parse fail" {
		t.Errorf("expected parse fail, got %v", err)
	}
}

func TestImportStatement_DuplicateTransaction(t *testing.T) {
	parserName := "mock"
	hash := "duphash"
	account := models.Account{Name: "Test Account", AccountTypeID: 1, AccountType: models.AccountType{DefaultParser: &parserName}}
	parsersByInstitution[parserName] = &MockParser{
		Txns: []models.Transaction{{Amount: 100, Date: "2024-01-01", Description: "Test txn"}},
	}
	is := &ImportService{
		AccountRepository:          &MockAccountRepository{Account: account},
		ImportSubmissionRepository: &MockImportSubmissionRepository{},
		TransactionRepository:      &MockTransactionRepository{TxnsByHash: map[string][]models.Transaction{hash: {{}}}},
		Categorizer:                &MockCategorizer{Category: models.Category{Name: "Test Category"}},
	}
	_, err := is.ImportStatement("file.csv", "statement", 1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

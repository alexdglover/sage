package services

import "github.com/alexdglover/sage/internal/models"

type MockAccountRepository struct {
	Account  models.Account
	Accounts []models.Account
	Err      error
}

func (m *MockAccountRepository) GetAccountByID(id uint) (models.Account, error) {
	return m.Account, m.Err
}

type MockBalanceRepository struct {
	Saved []models.Balance
}

func (m *MockBalanceRepository) Save(balance models.Balance) (uint, error) {
	m.Saved = append(m.Saved, balance)
	return 1, nil
}

type MockCategoryRepository struct {
	Category models.Category
	Err      error
}

type MockImportSubmissionRepository struct {
	Saved   []models.ImportSubmission
	SaveErr error
}

func (m *MockImportSubmissionRepository) Save(sub models.ImportSubmission) (uint, error) {
	m.Saved = append(m.Saved, sub)
	if m.SaveErr != nil {
		return 0, m.SaveErr
	}
	return 1, nil
}

type MockTransactionRepository struct {
	Err        error
	SaveErr    error
	Sum        int
	Totals     []models.TotalByMonth
	TxnsByHash map[string][]models.Transaction
}

func (m *MockTransactionRepository) GetTransactionsByHash(hash string, submissionID uint) ([]models.Transaction, error) {
	return m.TxnsByHash[hash], nil
}

func (m *MockTransactionRepository) Save(txn models.Transaction) (uint, error) {
	if m.SaveErr != nil {
		return 0, m.SaveErr
	}
	return 1, nil
}

type MockCategorizer struct {
	Category models.Category
	CatErr   error
	BuildErr error
}

func (m *MockCategorizer) BuildModel() error { return m.BuildErr }
func (m *MockCategorizer) CategorizeTransaction(txn *models.Transaction) (models.Category, error) {
	return m.Category, m.CatErr
}

type MockParser struct {
	Txns     []models.Transaction
	Balances []models.Balance
	ParseErr error
}

func (m *MockParser) Parse(statement string) ([]models.Transaction, []models.Balance, error) {
	return m.Txns, m.Balances, m.ParseErr
}

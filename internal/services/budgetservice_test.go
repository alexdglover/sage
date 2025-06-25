package services

import (
	"errors"
	"testing"
	"time"

	"github.com/alexdglover/sage/internal/models"
	"gorm.io/gorm"
)

type MockBudgetRepository struct {
	Budget models.Budget
	Err    error
}

func (m *MockBudgetRepository) GetBudgetByID(id uint) (models.Budget, error) {
	return m.Budget, m.Err
}

func (m *MockBudgetRepository) GetAllBudgets() ([]models.Budget, error) {
	return []models.Budget{m.Budget}, m.Err
}

// Only implement methods needed for the test

type MockTransactionRepository struct {
	Totals []models.TotalByMonth
	Err    error
	Sum    int
}

func (m *MockTransactionRepository) GetSumOfTransactionsByCategoryAndMonth(categoryID uint, startDate time.Time, endDate time.Time) ([]models.TotalByMonth, error) {
	return m.Totals, m.Err
}

func (m *MockTransactionRepository) GetSumOfTransactionsByCategoryID(categoryID uint, startDate time.Time, endDate time.Time) (int, error) {
	return m.Sum, m.Err
}

func TestGetMeanAndStandardDeviation(t *testing.T) {
	budget := models.Budget{
		Model:      gorm.Model{ID: 1},
		Amount:     1000,
		CategoryID: 2,
	}

	now := time.Now()
	// Use the last three months for test data
	mockTotals := []models.TotalByMonth{
		{Amount: 100, Month: now.AddDate(0, -2, 0)},
		{Amount: 200, Month: now.AddDate(0, -1, 0)},
		{Amount: 300, Month: now},
	}
	mockBudgetRepo := &MockBudgetRepository{Budget: budget}
	mockTxnRepo := &MockTransactionRepository{Totals: mockTotals}

	service := &BudgetService{
		BudgetRepository:      mockBudgetRepo,
		TransactionRepository: mockTxnRepo,
	}

	avg, stddev, err := service.GetMeanAndStandardDeviation(1, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if avg != 200 {
		t.Errorf("expected average 200, got %d", avg)
	}
	if stddev == 0 {
		t.Errorf("expected non-zero stddev, got %d", stddev)
	}
}

func TestGetMeanAndStandardDeviation_Error(t *testing.T) {
	mockBudgetRepo := &MockBudgetRepository{Err: errors.New("fail")}
	mockTxnRepo := &MockTransactionRepository{}
	service := &BudgetService{
		BudgetRepository:      mockBudgetRepo,
		TransactionRepository: mockTxnRepo,
	}
	_, _, err := service.GetMeanAndStandardDeviation(1, 3)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

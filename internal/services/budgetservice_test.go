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

func (m *MockTransactionRepository) GetSumOfTransactionsByCategoryAndMonth(categoryID uint, startDate time.Time, endDate time.Time) ([]models.TotalByMonth, error) {
	return m.Totals, m.Err
}

func (m *MockTransactionRepository) GetSumOfTransactionsByCategoryID(categoryID uint, startDate time.Time, endDate time.Time) (int, error) {
	return m.Sum, m.Err
}

func (m *MockCategoryRepository) GetCategoryByID(id uint) (models.Category, error) {
	return m.Category, m.Err
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

func TestGetAllBudgetsAndCurrentSpend(t *testing.T) {
	budget := models.Budget{
		Model:      gorm.Model{ID: 1},
		Amount:     1000,
		CategoryID: 2,
		Category:   models.Category{Model: gorm.Model{ID: 2}, Name: "Food"},
	}
	mockBudgetRepo := &MockBudgetRepository{Budget: budget}
	mockCategoryRepo := &MockCategoryRepository{Category: budget.Category}
	mockTxnRepo := &MockTransactionRepository{Sum: 500}

	service := &BudgetService{
		BudgetRepository:      mockBudgetRepo,
		CategoryRepository:    mockCategoryRepo,
		TransactionRepository: mockTxnRepo,
	}

	budgetsAndSpend, err := service.GetAllBudgetsAndCurrentSpend()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(budgetsAndSpend) != 1 {
		t.Errorf("expected 1 budget, got %d", len(budgetsAndSpend))
	}
	if budgetsAndSpend[0].Spend != 500 {
		t.Errorf("expected spend 500, got %d", budgetsAndSpend[0].Spend)
	}
	if budgetsAndSpend[0].PercentUsed != 50 {
		t.Errorf("expected percent used 50, got %d", budgetsAndSpend[0].PercentUsed)
	}
}

func TestGetAllBudgetsAndCurrentSpend_ErrorOnGetBudgets(t *testing.T) {
	mockBudgetRepo := &MockBudgetRepository{Err: errors.New("fail")}
	mockCategoryRepo := &MockCategoryRepository{}
	mockTxnRepo := &MockTransactionRepository{}
	service := &BudgetService{
		BudgetRepository:      mockBudgetRepo,
		CategoryRepository:    mockCategoryRepo,
		TransactionRepository: mockTxnRepo,
	}
	_, err := service.GetAllBudgetsAndCurrentSpend()
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetAllBudgetsAndCurrentSpend_ErrorOnGetCategory(t *testing.T) {
	budget := models.Budget{
		Model:      gorm.Model{ID: 1},
		Amount:     1000,
		CategoryID: 2,
	}
	mockBudgetRepo := &MockBudgetRepository{Budget: budget}
	mockCategoryRepo := &MockCategoryRepository{Err: errors.New("fail cat")}
	mockTxnRepo := &MockTransactionRepository{}
	service := &BudgetService{
		BudgetRepository:      mockBudgetRepo,
		CategoryRepository:    mockCategoryRepo,
		TransactionRepository: mockTxnRepo,
	}
	_, err := service.GetAllBudgetsAndCurrentSpend()
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetAllBudgetsAndCurrentSpend_ErrorOnGetSum(t *testing.T) {
	budget := models.Budget{
		Model:      gorm.Model{ID: 1},
		Amount:     1000,
		CategoryID: 2,
		Category:   models.Category{Model: gorm.Model{ID: 2}, Name: "Food"},
	}
	mockBudgetRepo := &MockBudgetRepository{Budget: budget}
	mockCategoryRepo := &MockCategoryRepository{Category: budget.Category}
	mockTxnRepo := &MockTransactionRepository{Err: errors.New("fail sum")}
	service := &BudgetService{
		BudgetRepository:      mockBudgetRepo,
		CategoryRepository:    mockCategoryRepo,
		TransactionRepository: mockTxnRepo,
	}
	_, err := service.GetAllBudgetsAndCurrentSpend()
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetBudgetAndMonthlySpend(t *testing.T) {
	budget := models.Budget{
		Model:      gorm.Model{ID: 1},
		Amount:     1000,
		CategoryID: 2,
		Category:   models.Category{Model: gorm.Model{ID: 2}, Name: "Food"},
	}
	mockBudgetRepo := &MockBudgetRepository{Budget: budget}
	mockTxnRepo := &MockTransactionRepository{Sum: 300}
	service := &BudgetService{
		BudgetRepository:      mockBudgetRepo,
		TransactionRepository: mockTxnRepo,
	}
	budgetsAndSpend, err := service.GetBudgetAndMonthlySpend(1, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(budgetsAndSpend) != 2 {
		t.Errorf("expected 2 months, got %d", len(budgetsAndSpend))
	}
	for _, b := range budgetsAndSpend {
		if b.Spend != 300 {
			t.Errorf("expected spend 300, got %d", b.Spend)
		}
		if b.PercentUsed != 30 {
			t.Errorf("expected percent used 30, got %d", b.PercentUsed)
		}
	}
}

func TestGetBudgetAndMonthlySpend_ErrorOnGetBudget(t *testing.T) {
	mockBudgetRepo := &MockBudgetRepository{Err: errors.New("fail budget")}
	mockTxnRepo := &MockTransactionRepository{}
	service := &BudgetService{
		BudgetRepository:      mockBudgetRepo,
		TransactionRepository: mockTxnRepo,
	}
	_, err := service.GetBudgetAndMonthlySpend(1, 2)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetBudgetAndMonthlySpend_ErrorOnGetSum(t *testing.T) {
	budget := models.Budget{
		Model:      gorm.Model{ID: 1},
		Amount:     1000,
		CategoryID: 2,
		Category:   models.Category{Model: gorm.Model{ID: 2}, Name: "Food"},
	}
	mockBudgetRepo := &MockBudgetRepository{Budget: budget}
	mockTxnRepo := &MockTransactionRepository{Err: errors.New("fail sum")}
	service := &BudgetService{
		BudgetRepository:      mockBudgetRepo,
		TransactionRepository: mockTxnRepo,
	}
	_, err := service.GetBudgetAndMonthlySpend(1, 2)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

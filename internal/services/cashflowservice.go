package services

import (
	"context"
	"time"

	"github.com/alexdglover/sage/internal/models"
)

type CashFlowData struct {
	TotalIncome   int
	TotalExpenses int
	Expenses      []ExpenseByCategory
}

type ExpenseByCategory struct {
	Name   string
	Amount int
}

type CashFlowService struct {
	transactionRepo *models.TransactionRepository
	categoryRepo    *models.CategoryRepository
}

func NewCashFlowService(tr *models.TransactionRepository) *CashFlowService {
	return &CashFlowService{
		transactionRepo: tr,
	}
}

func (s *CashFlowService) GetCashFlowData(ctx context.Context, startDate time.Time, endDate time.Time) (cashFlowData *CashFlowData, err error) {
	cashFlowData = &CashFlowData{} // Initialize the struct

	incomeCategory, err := s.categoryRepo.GetCategoryByName("Income")
	if err != nil {
		return nil, err
	}
	cashFlowData.TotalIncome, err = s.transactionRepo.GetSumOfTransactionsByCategoryID(incomeCategory.ID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	expensesByCategory, err := s.transactionRepo.GetSumOfTransactionsByCategory(startDate, endDate)
	if err != nil {
		return nil, err
	}
	totalExpenses := 0
	for _, expense := range expensesByCategory {
		totalExpenses += expense.Amount
		cashFlowData.Expenses = append(cashFlowData.Expenses, ExpenseByCategory{
			Name:   expense.Category,
			Amount: expense.Amount,
		})
	}
	cashFlowData.TotalExpenses = totalExpenses

	return cashFlowData, nil
}

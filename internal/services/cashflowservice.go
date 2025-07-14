package services

import (
	"context"
	"time"

	"github.com/alexdglover/sage/internal/models"
)

type CashFlowData struct {
	StartDate time.Time
	EndDate   time.Time
	Income    []CategoryFlow
	Expenses  []CategoryFlow
}

type CategoryFlow struct {
	Category string
	Amount   int
}

type CashFlowService struct {
	transactionRepo *models.TransactionRepository
}

func NewCashFlowService(tr *models.TransactionRepository) *CashFlowService {
	return &CashFlowService{
		transactionRepo: tr,
	}
}

func (s *CashFlowService) GetCashFlowData(ctx context.Context, startDate time.Time, endDate time.Time) (*CashFlowData, error) {
	// Get income transactions
	incomeData, err := s.transactionRepo.GetSumOfTransactionsByCategory(startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Split into income and expenses
	cashFlowData := &CashFlowData{
		StartDate: startDate,
		EndDate:   endDate,
		Income:    make([]CategoryFlow, 0),
		Expenses:  make([]CategoryFlow, 0),
	}

	for _, total := range incomeData {
		flow := CategoryFlow{
			Category: total.Category,
			Amount:   total.Amount,
		}

		if total.Category == "Income" {
			cashFlowData.Income = append(cashFlowData.Income, flow)
		} else if total.Category != "Transfers" {
			// Convert expense amounts to positive for visualization
			flow.Amount = -flow.Amount
			cashFlowData.Expenses = append(cashFlowData.Expenses, flow)
		}
	}

	return cashFlowData, nil
}

package services

import (
	"fmt"
	"time"

	"github.com/alexdglover/sage/internal/models"
)

type BudgetService struct {
	BudgetRepository      *models.BudgetRepository
	CategoryRepository    *models.CategoryRepository
	TransactionRepository *models.TransactionRepository
}

type BudgetAndSpend struct {
	ID           uint
	CategoryName string
	Amount       int
	Spend        int
	PercentUsed  int
}

func (bs *BudgetService) GetAllBudgetsAndCurrentSpend() (budgetsAndSpend []BudgetAndSpend, err error) {
	budgets, err := bs.BudgetRepository.GetAllBudgets()
	if err != nil {
		fmt.Println("Unable to get budgets:", err)
		return budgetsAndSpend, err
	}

	now := time.Now()
	firstOfMonth := now.AddDate(0, 0, 1-now.Day())
	lastOfMonth := now.AddDate(0, 1, 0-now.Day())

	for _, budget := range budgets {
		category, err := bs.CategoryRepository.GetCategoryByID(budget.CategoryID)
		if err != nil {
			fmt.Println("Unable to get category by ID:", err)
			return budgetsAndSpend, err
		}
		sum, err := bs.TransactionRepository.GetSumOfTransactionsByCategoryID(category.ID, firstOfMonth, lastOfMonth)
		if err != nil {
			fmt.Printf("Unable to get transactions with category ID %v: %v", category.ID, err)
			return budgetsAndSpend, err
		}
		budgetsAndSpend = append(budgetsAndSpend, BudgetAndSpend{
			ID:           budget.ID,
			CategoryName: budget.Category.Name,
			Amount:       budget.Amount,
			Spend:        sum,
			PercentUsed:  int(float64(sum) / float64(budget.Amount) * 100),
		})
	}
	return budgetsAndSpend, nil
}

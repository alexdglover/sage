package services

import (
	"fmt"
	"time"

	"github.com/alexdglover/sage/internal/models"
	"gonum.org/v1/gonum/stat"
)

type BudgetRepositoryInterface interface {
	GetBudgetByID(id uint) (models.Budget, error)
	GetAllBudgets() ([]models.Budget, error)
}

type TransactionRepositoryInterface interface {
	GetSumOfTransactionsByCategoryAndMonth(categoryID uint, startDate time.Time, endDate time.Time) ([]models.TotalByMonth, error)
	GetSumOfTransactionsByCategoryID(categoryID uint, startDate time.Time, endDate time.Time) (int, error)
}

type BudgetService struct {
	BudgetRepository      BudgetRepositoryInterface
	CategoryRepository    *models.CategoryRepository // unchanged for now
	TransactionRepository TransactionRepositoryInterface
}

type BudgetAndSpend struct {
	ID           uint
	CategoryName string
	Amount       int
	Spend        int
	PercentUsed  int
	Month        time.Time
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
			fmt.Printf("Unable to get category by ID %v: %v\n", budget.CategoryID, err)
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

func (bs *BudgetService) GetBudgetAndMonthlySpend(budgetID uint, numOfMonths int) (budgetsAndSpend []BudgetAndSpend, err error) {
	budget, err := bs.BudgetRepository.GetBudgetByID(budgetID)
	if err != nil {
		fmt.Println("Unable to get budgets:", err)
		return budgetsAndSpend, err
	}

	now := time.Now()
	firstOfMonth := now.AddDate(0, -int(numOfMonths-1), 1-now.Day())
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	for i := int(0); i < numOfMonths; i++ {
		sum, err := bs.TransactionRepository.GetSumOfTransactionsByCategoryID(budget.CategoryID, firstOfMonth, lastOfMonth)
		if err != nil {
			fmt.Printf("Unable to get transactions with category ID %v: %v", budget.CategoryID, err)
			return budgetsAndSpend, err
		}
		budgetsAndSpend = append(budgetsAndSpend, BudgetAndSpend{
			ID:           budget.ID,
			CategoryName: budget.Category.Name,
			Amount:       budget.Amount,
			Spend:        sum,
			PercentUsed:  int(float64(sum) / float64(budget.Amount) * 100),
			Month:        firstOfMonth,
		})
		firstOfMonth = firstOfMonth.AddDate(0, 1, 0)
		lastOfMonth = lastOfMonth.AddDate(0, 1, 0)
	}
	return budgetsAndSpend, nil
}

func (bs *BudgetService) GetMeanAndStandardDeviation(budgetID uint, numOfMonths int) (averageSpend int, standardDeviation int, err error) {
	budget, err := bs.BudgetRepository.GetBudgetByID(budgetID)
	if err != nil {
		fmt.Println("Unable to get budgets:", err)
		return averageSpend, standardDeviation, err
	}

	now := time.Now()
	endDate := now.AddDate(0, 1, 0-now.Day())
	startDate := now.AddDate(0, -int(numOfMonths), 1-now.Day())

	spendByMonth, err := bs.TransactionRepository.GetSumOfTransactionsByCategoryAndMonth(budget.CategoryID, startDate, endDate)
	var amounts []float64
	for i := range spendByMonth {
		amounts = append(amounts, float64(spendByMonth[i].Amount))
	}
	averageSpendFloat, standardDeviationFloat := stat.MeanStdDev(amounts, nil)
	averageSpend = int(averageSpendFloat)
	standardDeviation = int(standardDeviationFloat)
	return averageSpend, standardDeviation, err
}

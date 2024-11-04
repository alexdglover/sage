package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Budget struct {
	gorm.Model
	Amount     int64
	CategoryID uint
	Category   Category
}

type BudgetRepository struct{}

var budgetRepository *BudgetRepository

func GetBudgetRepository() *BudgetRepository {
	if budgetRepository == nil {
		budgetRepository = &BudgetRepository{}
	}
	return budgetRepository
}

func (*BudgetRepository) GetAllBudgets() ([]Budget, error) {
	var budgets []Budget
	result := db.Preload(clause.Associations).Find(&budgets)
	return budgets, result.Error
}

func (*BudgetRepository) GetBudgetByID(id uint) (Budget, error) {
	var budget Budget
	result := db.Preload(clause.Associations).Where("id = ?", id).First(&budget)
	return budget, result.Error
}

// Save is an UPSERT operation, returning the ID of the record and an optional error
func (*BudgetRepository) Save(budget Budget) (id uint, err error) {
	result := db.Save(&budget).Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}})

	return budget.ID, result.Error
}

package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Category struct {
	gorm.Model
	Name string `gorm:"uniqueIndex"`
}

type CategoryAndBudgetStatus struct {
	Category
	HasBudget bool
}

type CategoryRepository struct {
	DB *gorm.DB
}

func (cr *CategoryRepository) GetAllCategories() ([]Category, error) {
	var categories []Category
	result := db.Order("name asc").Find(&categories)
	return categories, result.Error
}

func (cr *CategoryRepository) GetAllCategoriesAndBudgetStatus() (categories []CategoryAndBudgetStatus, err error) {
	cr.DB.Raw(`SELECT
		c.ID,
		c.Name,
		CASE
			WHEN b.amount IS NOT NULL THEN true
			ELSE false
		END AS has_budget
		FROM categories c
		LEFT JOIN budgets b ON c.ID = b.category_id and b.deleted_at IS NULL
		WHERE c.deleted_at IS NULL;`).Scan(&categories)
	return categories, nil
}

func (cr *CategoryRepository) GetCategoryByID(id uint) (Category, error) {
	var category Category
	result := db.Where("id = ?", id).First(&category)
	return category, result.Error
}

func (cr *CategoryRepository) GetCategoryByName(name string) (Category, error) {
	var category Category
	result := db.Where("name = ?", name).First(&category)
	return category, result.Error
}

// Save is an UPSERT operation, returning the ID of the record and an optional error
func (cr *CategoryRepository) Save(category Category) (id uint, err error) {
	result := cr.DB.Save(&category).Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}})
	return category.ID, result.Error
}

// Soft deletes a category and sets all associated transactions to Unknown
func (cr *CategoryRepository) DeleteCategoryByID(categoryID uint) (err error) {
	// delete any associated budgets
	cr.DB.Where("category_id = ?", categoryID).Delete(&Budget{})
	// bulk update all transactions to "Unknown" category
	cr.DB.Model(&Transaction{}).Where("category_id = ?", categoryID).Update("category_id", 1)
	result := cr.DB.Delete(&Category{}, categoryID)
	return result.Error
}

package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name string `gorm:"uniqueIndex"`
}

type CategoryRepository struct{}

var categoryRepository *CategoryRepository

func GetCategoryRepository() *CategoryRepository {
	if categoryRepository == nil {
		categoryRepository = &CategoryRepository{}
	}
	return categoryRepository
}

func (*CategoryRepository) GetAllCategories() ([]Category, error) {
	var categories []Category
	result := db.Find(&categories)
	return categories, result.Error
}

func (*CategoryRepository) GetCategoryByID(id uint) (Category, error) {
	var category Category
	result := db.Where("id = ?", id).First(&category)
	return category, result.Error
}

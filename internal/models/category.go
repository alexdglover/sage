package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name        string
	DisplayName string
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

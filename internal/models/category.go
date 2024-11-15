package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name string `gorm:"uniqueIndex"`
}

type CategoryRepository struct {
	DB *gorm.DB
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

func (*CategoryRepository) GetCategoryByName(name string) (Category, error) {
	var category Category
	result := db.Where("name = ?", name).First(&category)
	return category, result.Error
}

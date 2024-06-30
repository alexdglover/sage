package models

import "gorm.io/gorm"

type Budget struct {
	gorm.Model
	Name       string
	Amount     float64
	CategoryId uint
	Category   Category
}

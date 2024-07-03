package models

import "gorm.io/gorm"

type Budget struct {
	gorm.Model
	Name       string
	Amount     float32
	CategoryId uint
	Category   Category
}

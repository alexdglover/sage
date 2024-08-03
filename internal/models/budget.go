package models

import "gorm.io/gorm"

type Budget struct {
	gorm.Model
	Name       string
	Amount     int64
	CategoryId uint
	Category   Category
}

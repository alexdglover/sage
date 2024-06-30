package internal

import (
	"context"

	models "github.com/alexdglover/sage/internal/models"
	"gorm.io/gorm"
)

func BootstrapTables(ctx context.Context, db *gorm.DB) {
	// Migrate the schema
	db.AutoMigrate(&models.Account{})
	db.AutoMigrate(&models.Balance{})
	db.AutoMigrate(&models.Budget{})
	db.AutoMigrate(&models.Category{})
	db.AutoMigrate(&models.Transaction{})

	// Insert seed data
	for _, name := range []string{"Home", "Income", "Auto", "Food", "Dining"} {
		db.Create(&models.Category{Name: name})
	}

}

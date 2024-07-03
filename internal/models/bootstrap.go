package models

import (
	"context"
)

func BootstrapDatabase(ctx context.Context) {
	// Instantiate the sqlite client singleton
	createDbClient()

	// Migrate the schema
	db.AutoMigrate(&Account{})
	db.AutoMigrate(&Balance{})
	db.AutoMigrate(&Budget{})
	db.AutoMigrate(&Category{})
	db.AutoMigrate(&Transaction{})

	// Insert seed data
	for _, name := range []string{"Home", "Income", "Auto", "Food", "Dining"} {
		db.Create(&Category{Name: name})
	}

}

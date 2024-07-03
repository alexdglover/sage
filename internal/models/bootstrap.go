package models

import (
	"context"
	"os"
)

func BootstrapDatabase(ctx context.Context) {
	// Instantiate the sqlite client singleton
	createDbClient()

	// Conditionally drop all tables to start from scratch
	if os.Getenv("DROP_TABLES") != "" {
		db.Migrator().DropTable(&Account{})
		db.Migrator().DropTable(&Balance{})
		db.Migrator().DropTable(&Budget{})
		db.Migrator().DropTable(&Category{})
		db.Migrator().DropTable(&Transaction{})
	}

	// Migrate the schema
	db.AutoMigrate(&Account{})
	db.AutoMigrate(&Balance{})
	db.AutoMigrate(&Budget{})
	db.AutoMigrate(&Category{})
	db.AutoMigrate(&Transaction{})

	// Insert common seed data
	for _, name := range []string{"Home", "Income", "Auto", "Food", "Dining"} {
		db.Create(&Category{Name: name})
	}

	// Conditionally insert sample date for testing purposes
	if os.Getenv("ADD_SAMPLE_DATA") != "" {
		// Create one asset account and one liability account
		db.Create(&Account{Name: "schwaby", AccountCategory: "savings", AccountType: "asset"})
		db.Create(&Account{Name: "schmisa", AccountCategory: "credit card", AccountType: "liability"})

		// Create balances
		db.Create(&Balance{Date: "2024-03-13", EffectiveStartDate: "2024-03-13", Balance: 420.32, AccountId: 1})
		db.Create(&Balance{Date: "2024-03-17", EffectiveStartDate: "2024-03-17", Balance: 103.87, AccountId: 2})
	}

}

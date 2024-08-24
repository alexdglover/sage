package models

import (
	"context"
	"os"

	"github.com/alexdglover/sage/internal/utils"
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
		db.Migrator().DropTable(&ImportSubmission{})
	}

	// Migrate the schema
	db.AutoMigrate(&Account{})
	db.AutoMigrate(&Balance{})
	db.AutoMigrate(&Budget{})
	db.AutoMigrate(&Category{})
	db.AutoMigrate(&Transaction{})
	db.AutoMigrate(&ImportSubmission{})

	// Insert common seed data
	for _, name := range []string{"Home", "Income", "Auto", "Food", "Dining"} {
		db.Create(&Category{Name: name})
	}

	// Conditionally insert sample date for testing purposes
	if os.Getenv("ADD_SAMPLE_DATA") != "" {
		// Create one normal asset account, one normal liability account, and one infrequently updated account
		// of each type
		db.Create(&Account{Name: "Schwab", AccountCategory: "checking", AccountType: "asset", DefaultParser: utils.StrPointer("schwab")})
		db.Create(&Account{Name: "Fidelity Visa", AccountCategory: "credit card", AccountType: "liability", DefaultParser: utils.StrPointer("fidelity visa")})
		db.Create(&Account{Name: "My House", AccountCategory: "real estate", AccountType: "asset"})
		db.Create(&Account{Name: "Mortgage", AccountCategory: "loan", AccountType: "liability"})

		// Create open-ended balances for infrequently updated accounts
		// db.Create(&Balance{Date: "2024-01-17", EffectiveStartDate: "2024-01-17", Amount: 2500, AccountId: 3})
		// db.Create(&Balance{Date: "2024-01-17", EffectiveStartDate: "2024-01-17", Amount: 1250, AccountId: 4})

		// Create monthly balances for normal accounts
		db.Create(&Balance{Date: "2024-02-13", EffectiveStartDate: "2024-02-01", EffectiveEndDate: utils.StrPointer("2024-02-29"), Amount: 21013, AccountId: 1})
		db.Create(&Balance{Date: "2024-03-13", EffectiveStartDate: "2024-03-01", EffectiveEndDate: utils.StrPointer("2024-03-31"), Amount: 41062, AccountId: 1})
		db.Create(&Balance{Date: "2024-04-13", EffectiveStartDate: "2024-04-01", EffectiveEndDate: utils.StrPointer("2024-04-30"), Amount: 42032, AccountId: 1})
		db.Create(&Balance{Date: "2024-05-13", EffectiveStartDate: "2024-05-01", EffectiveEndDate: utils.StrPointer("2024-05-31"), Amount: 49032, AccountId: 1})
		db.Create(&Balance{Date: "2024-06-13", EffectiveStartDate: "2024-06-01", EffectiveEndDate: utils.StrPointer("2024-06-30"), Amount: 64097, AccountId: 1})
		db.Create(&Balance{Date: "2024-07-13", EffectiveStartDate: "2024-07-01", Amount: 63201, AccountId: 1})

		db.Create(&Balance{Date: "2024-02-17", EffectiveStartDate: "2024-02-01", EffectiveEndDate: utils.StrPointer("2024-02-29"), Amount: 10111, AccountId: 2})
		db.Create(&Balance{Date: "2024-03-17", EffectiveStartDate: "2024-03-01", EffectiveEndDate: utils.StrPointer("2024-03-31"), Amount: 17387, AccountId: 2})
		db.Create(&Balance{Date: "2024-04-17", EffectiveStartDate: "2024-04-01", EffectiveEndDate: utils.StrPointer("2024-04-30"), Amount: 10387, AccountId: 2})
		db.Create(&Balance{Date: "2024-05-17", EffectiveStartDate: "2024-05-01", EffectiveEndDate: utils.StrPointer("2024-05-31"), Amount: 13312, AccountId: 2})
		db.Create(&Balance{Date: "2024-06-17", EffectiveStartDate: "2024-06-01", EffectiveEndDate: utils.StrPointer("2024-06-30"), Amount: 14044, AccountId: 2})
		db.Create(&Balance{Date: "2024-07-17", EffectiveStartDate: "2024-07-01", Amount: 13255, AccountId: 2})

	}

}

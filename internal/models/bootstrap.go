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
		db.Create(&Account{Name: "Fidelity Visa", AccountCategory: "creditCard", AccountType: "liability", DefaultParser: utils.StrPointer("fidelity")})
		db.Create(&Account{Name: "My House", AccountCategory: "realEstate", AccountType: "asset"})
		db.Create(&Account{Name: "Mortgage", AccountCategory: "loan", AccountType: "liability"})

		// Create open-ended balances for infrequently updated accounts
		// db.Create(&Balance{EffectiveDate: "2024-01-17", Amount: 2500, AccountID: 3})
		// db.Create(&Balance{EffectiveDate: "2024-01-17", Amount: 1250, AccountID: 4})

		// Create monthly balances for normal accounts
		db.Create(&Balance{EffectiveDate: "2024-02-01", Amount: 21013, AccountID: 1})
		db.Create(&Balance{EffectiveDate: "2024-03-01", Amount: 41062, AccountID: 1})
		db.Create(&Balance{EffectiveDate: "2024-04-01", Amount: 42032, AccountID: 1})
		db.Create(&Balance{EffectiveDate: "2024-05-01", Amount: 49032, AccountID: 1})
		db.Create(&Balance{EffectiveDate: "2024-06-01", Amount: 64097, AccountID: 1})
		db.Create(&Balance{EffectiveDate: "2024-07-01", Amount: 63201, AccountID: 1})

		db.Create(&Balance{EffectiveDate: "2024-02-01", Amount: 10111, AccountID: 2})
		db.Create(&Balance{EffectiveDate: "2024-03-01", Amount: 17387, AccountID: 2})
		db.Create(&Balance{EffectiveDate: "2024-04-01", Amount: 10387, AccountID: 2})
		db.Create(&Balance{EffectiveDate: "2024-05-01", Amount: 13312, AccountID: 2})
		db.Create(&Balance{EffectiveDate: "2024-06-01", Amount: 14044, AccountID: 2})
		db.Create(&Balance{EffectiveDate: "2024-07-01", Amount: 13255, AccountID: 2})

	}

}

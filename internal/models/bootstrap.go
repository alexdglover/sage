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
		// Create one normal asset account, one normal liability account, and one infrequently updated account
		// of each type
		db.Create(&Account{Name: "schwaby", AccountCategory: "savings", AccountType: "asset"})
		db.Create(&Account{Name: "schmisa", AccountCategory: "credit card", AccountType: "liability"})
		db.Create(&Account{Name: "My House", AccountCategory: "real estate", AccountType: "asset"})
		db.Create(&Account{Name: "Mortgage", AccountCategory: "loan", AccountType: "liability"})

		// Create open-ended balances for infrequently updated accounts
		// db.Create(&Balance{Date: "2024-01-17", EffectiveStartDate: "2024-01-17", Balance: 2500, AccountId: 3})
		// db.Create(&Balance{Date: "2024-01-17", EffectiveStartDate: "2024-01-17", Balance: 1250, AccountId: 4})

		// Create monthly balances for normal accounts
		db.Create(&Balance{Date: "2024-02-13", EffectiveStartDate: "2024-02-01", EffectiveEndDate: utils.StrPointer("2024-02-29"), Balance: 21013, AccountId: 1})
		db.Create(&Balance{Date: "2024-03-13", EffectiveStartDate: "2024-03-01", EffectiveEndDate: utils.StrPointer("2024-03-31"), Balance: 41062, AccountId: 1})
		db.Create(&Balance{Date: "2024-04-13", EffectiveStartDate: "2024-04-01", EffectiveEndDate: utils.StrPointer("2024-04-30"), Balance: 42032, AccountId: 1})
		db.Create(&Balance{Date: "2024-05-13", EffectiveStartDate: "2024-05-01", EffectiveEndDate: utils.StrPointer("2024-05-31"), Balance: 49032, AccountId: 1})
		db.Create(&Balance{Date: "2024-06-13", EffectiveStartDate: "2024-06-01", EffectiveEndDate: utils.StrPointer("2024-06-30"), Balance: 64097, AccountId: 1})
		db.Create(&Balance{Date: "2024-07-13", EffectiveStartDate: "2024-07-01", Balance: 63201, AccountId: 1})

		db.Create(&Balance{Date: "2024-02-17", EffectiveStartDate: "2024-02-01", EffectiveEndDate: utils.StrPointer("2024-02-29"), Balance: 10111, AccountId: 2})
		db.Create(&Balance{Date: "2024-03-17", EffectiveStartDate: "2024-03-01", EffectiveEndDate: utils.StrPointer("2024-03-31"), Balance: 17387, AccountId: 2})
		db.Create(&Balance{Date: "2024-04-17", EffectiveStartDate: "2024-04-01", EffectiveEndDate: utils.StrPointer("2024-04-30"), Balance: 10387, AccountId: 2})
		db.Create(&Balance{Date: "2024-05-17", EffectiveStartDate: "2024-05-01", EffectiveEndDate: utils.StrPointer("2024-05-31"), Balance: 13312, AccountId: 2})
		db.Create(&Balance{Date: "2024-06-17", EffectiveStartDate: "2024-06-01", EffectiveEndDate: utils.StrPointer("2024-06-30"), Balance: 14044, AccountId: 2})
		db.Create(&Balance{Date: "2024-07-17", EffectiveStartDate: "2024-07-01", Balance: 13255, AccountId: 2})

	}

}

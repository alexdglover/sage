package models

import (
	"context"
	"os"

	"github.com/alexdglover/sage/internal/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Bootstrapper struct {
	db *gorm.DB
}

func NewBootstrapper(db *gorm.DB) *Bootstrapper {
	return &Bootstrapper{db: db}
}

func (b *Bootstrapper) BootstrapDatabase(ctx context.Context) {
	// Instantiate the sqlite client singleton
	createDbClient()

	// Conditionally drop all tables to start from scratch
	if os.Getenv("DROP_TABLES") != "" {
		b.db.Migrator().DropTable(&Account{})
		b.db.Migrator().DropTable(&Balance{})
		b.db.Migrator().DropTable(&Budget{})
		b.db.Migrator().DropTable(&Category{})
		b.db.Migrator().DropTable(&Transaction{})
		b.db.Migrator().DropTable(&ImportSubmission{})
	}

	// Migrate the schema
	b.db.AutoMigrate(&Account{})
	b.db.AutoMigrate(&Balance{})
	b.db.AutoMigrate(&Budget{})
	b.db.AutoMigrate(&Category{})
	b.db.AutoMigrate(&Transaction{})
	b.db.AutoMigrate(&ImportSubmission{})

	// Insert common seed data
	for _, name := range []string{"Unknown", "Home", "Income", "Auto", "Food", "Dining"} {
		// The Category table has a unique index on the Name column, so we can use the DoNothing option
		// to safely attempt to insert a record that may already exist
		b.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&Category{Name: name})
	}

	// Conditionally insert sample date for testing purposes
	if os.Getenv("ADD_SAMPLE_DATA") != "" {
		// Create one normal asset account, one normal liability account, and one infrequently updated account
		// of each type
		b.db.Create(&Account{Name: "Schwab", AccountCategory: "checking", AccountType: "asset", DefaultParser: utils.StrPointer("schwabChecking")})
		b.db.Create(&Account{Name: "Fidelity Visa", AccountCategory: "creditCard", AccountType: "liability", DefaultParser: utils.StrPointer("fidelityCreditCard")})
		b.db.Create(&Account{Name: "My House", AccountCategory: "realEstate", AccountType: "asset"})
		b.db.Create(&Account{Name: "Mortgage", AccountCategory: "loan", AccountType: "liability"})

		// Create open-ended balances for infrequently updated accounts
		// b.db.Create(&Balance{EffectiveDate: "2024-01-17", Amount: 2500, AccountID: 3})
		// b.db.Create(&Balance{EffectiveDate: "2024-01-17", Amount: 1250, AccountID: 4})

		// Create monthly balances for normal accounts
		b.db.Create(&Balance{EffectiveDate: "2024-02-01", Amount: 21013, AccountID: 1})
		b.db.Create(&Balance{EffectiveDate: "2024-03-01", Amount: 41062, AccountID: 1})
		b.db.Create(&Balance{EffectiveDate: "2024-04-01", Amount: 42032, AccountID: 1})
		b.db.Create(&Balance{EffectiveDate: "2024-05-01", Amount: 49032, AccountID: 1})
		b.db.Create(&Balance{EffectiveDate: "2024-06-01", Amount: 64097, AccountID: 1})
		b.db.Create(&Balance{EffectiveDate: "2024-07-01", Amount: 63201, AccountID: 1})

		b.db.Create(&Balance{EffectiveDate: "2024-02-01", Amount: 10111, AccountID: 2})
		b.db.Create(&Balance{EffectiveDate: "2024-03-01", Amount: 17387, AccountID: 2})
		b.db.Create(&Balance{EffectiveDate: "2024-04-01", Amount: 10387, AccountID: 2})
		b.db.Create(&Balance{EffectiveDate: "2024-05-01", Amount: 13312, AccountID: 2})
		b.db.Create(&Balance{EffectiveDate: "2024-06-01", Amount: 14044, AccountID: 2})
		b.db.Create(&Balance{EffectiveDate: "2024-07-01", Amount: 13255, AccountID: 2})

	}

}

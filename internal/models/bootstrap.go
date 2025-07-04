package models

import (
	"context"
	"os"

	"github.com/alexdglover/sage/internal/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const Asset = "asset"
const Liability = "liability"

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
		err := b.db.Migrator().DropTable(&Account{})
		if err != nil {
			panic("Error dropping Account table: " + err.Error())
		}
		err = b.db.Migrator().DropTable(&AccountType{})
		if err != nil {
			panic("Error dropping AccountType table: " + err.Error())
		}
		err = b.db.Migrator().DropTable(&Balance{})
		if err != nil {
			panic("Error dropping Balance table: " + err.Error())
		}
		err = b.db.Migrator().DropTable(&Budget{})
		if err != nil {
			panic("Error dropping Budget table: " + err.Error())
		}
		err = b.db.Migrator().DropTable(&Category{})
		if err != nil {
			panic("Error dropping Category table: " + err.Error())
		}
		err = b.db.Migrator().DropTable(&ImportSubmission{})
		if err != nil {
			panic("Error dropping ImportSubmission table: " + err.Error())
		}
		err = b.db.Migrator().DropTable(&Settings{})
		if err != nil {
			panic("Error dropping Settings table: " + err.Error())
		}
		err = b.db.Migrator().DropTable(&Transaction{})
		if err != nil {
			panic("Error dropping Transaction table: " + err.Error())
		}

	}

	// Migrate the schema
	err := b.db.AutoMigrate(&Account{})
	if err != nil {
		panic("Error dropping migrationg Account table: " + err.Error())
	}
	err = b.db.AutoMigrate(&AccountType{})
	if err != nil {
		panic("Error dropping migrationg Account table: " + err.Error())
	}
	err = b.db.AutoMigrate(&Balance{})
	if err != nil {
		panic("Error dropping migrationg Account table: " + err.Error())
	}
	err = b.db.AutoMigrate(&Budget{})
	if err != nil {
		panic("Error dropping migrationg Account table: " + err.Error())
	}
	err = b.db.AutoMigrate(&Category{})
	if err != nil {
		panic("Error dropping migrationg Account table: " + err.Error())
	}
	err = b.db.AutoMigrate(&ImportSubmission{})
	if err != nil {
		panic("Error dropping migrationg Account table: " + err.Error())
	}
	err = b.db.AutoMigrate(&Settings{})
	if err != nil {
		panic("Error dropping migrationg Account table: " + err.Error())
	}
	err = b.db.AutoMigrate(&Transaction{})
	if err != nil {
		panic("Error dropping migrationg Account table: " + err.Error())
	}

	// Seed data for common categories, if they don't exist already
	for _, name := range []string{"Unknown", "Transfers", "Home", "Income", "Auto", "Food", "Dining"} {
		// The Category table has a unique index on the Name column, so we can use the DoNothing option
		// to safely attempt to insert a record that may already exist
		b.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&Category{Name: name})
	}

	accountTypesData := map[string]map[string]string{
		"Bank of America Credit Card": {"ledgerType": Liability, "accountCategory": "creditCard", "defaultParser": "bankOfAmericaCreditCard"},
		"Capital One Credit Card":     {"ledgerType": Liability, "accountCategory": "creditCard", "defaultParser": "capitalOneCreditCard"},
		"Capital One Savings":         {"ledgerType": Asset, "accountCategory": "savings", "defaultParser": "capitalOneSavings"},
		"Chase Checking":              {"ledgerType": Asset, "accountCategory": "checking", "defaultParser": "chaseChecking"},
		"Chase Credit Card":           {"ledgerType": Liability, "accountCategory": "creditCard", "defaultParser": "chaseCreditCard"},
		"Schwab Brokerage":            {"ledgerType": Asset, "accountCategory": "brokerage", "defaultParser": "schwabBrokerage"},
		"Schwab Checking":             {"ledgerType": Asset, "accountCategory": "checking", "defaultParser": "schwabChecking"},
		"Fidelity Credit Card":        {"ledgerType": Liability, "accountCategory": "creditCard", "defaultParser": "fidelityCreditCard"},
		"Fidelity Brokerage":          {"ledgerType": Asset, "accountCategory": "brokerage", "defaultParser": "fidelityBrokerage"},
		"Real Estate":                 {"ledgerType": Asset, "accountCategory": "realEstate"},
		"Mortgage":                    {"ledgerType": Liability, "accountCategory": "loan"},
		"Target Credit Card":          {"ledgerType": Liability, "accountCategory": "creditCard", "defaultParser": "targetCreditCard"},
		"UWCU Mortgage":               {"ledgerType": Liability, "accountCategory": "loan", "defaultParser": "uwcuMortgage"},

		"Misc Asset":     {"ledgerType": Asset, "accountCategory": Asset},
		"Misc Liability": {"ledgerType": Liability, "accountCategory": Asset},
	}

	// Seed data for supported account types, if they don't exist already
	for name, accountTypeDetails := range accountTypesData {
		// The account_types table has a unique index on the Name column, so we can use the DoNothing option
		// to safely attempt to insert a record that may already exist
		b.db.Clauses(clause.OnConflict{DoNothing: true}).Create(
			&AccountType{
				Name:            name,
				AccountCategory: accountTypeDetails["accountCategory"],
				DefaultParser:   utils.StrPointer(accountTypeDetails["defaultParser"]),
				LedgerType:      accountTypeDetails["ledgerType"],
			})
	}

	// Conditionally insert sample date for testing purposes
	if os.Getenv("ADD_SAMPLE_DATA") != "" {
		// Create one normal asset account, one normal liability account, and one infrequently updated account
		// of each type
		b.db.Create(&Account{Name: "Schwab", AccountTypeID: 1})
		b.db.Create(&Account{Name: "Fidelity Visa", AccountTypeID: 2})
		b.db.Create(&Account{Name: "My House", AccountTypeID: 3})
		b.db.Create(&Account{Name: "Mortgage", AccountTypeID: 4})

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

	// Seed default settings data
	// The Settings table has a unique index on the Name column, so we can use the DoNothing option
	// to safely attempt to insert a record that may already exist
	b.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&Settings{Model: gorm.Model{ID: 1}, LaunchBrowserOnStartup: true})

}

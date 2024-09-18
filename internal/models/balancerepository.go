package models

import (
	"context"
	"time"

	"github.com/alexdglover/sage/internal/utils"
	"gorm.io/gorm/clause"
)

type BalanceRepository struct{}

var balanceRepository *BalanceRepository

type BalancesWithDate struct {
	Date     time.Time
	Balances []Balance
}

func GetBalanceRepository() *BalanceRepository {
	if balanceRepository == nil {
		balanceRepository = &BalanceRepository{}
	}
	return balanceRepository
}

func (*BalanceRepository) GetAllBalances(ctx context.Context) ([]Balance, error) {
	var balances []Balance
	result := db.Find(&balances)
	return balances, result.Error
}

func (*BalanceRepository) GetBalancesOfAllAssetsByMonth(ctx context.Context, startYearMonth time.Time, endYearMonth time.Time) []BalancesWithDate {
	return balanceRepository.GetBalancesByMonth(ctx, "asset", startYearMonth, endYearMonth)
}

func (*BalanceRepository) GetBalancesOfAllLiabilitiesByMonth(ctx context.Context, startYearMonth time.Time, endYearMonth time.Time) []BalancesWithDate {
	return balanceRepository.GetBalancesByMonth(ctx, "liability", startYearMonth, endYearMonth)
}

func (*BalanceRepository) GetBalancesByMonth(ctx context.Context, accountType string, startYearMonth time.Time, endYearMonth time.Time) []BalancesWithDate {
	// TODO: implement a better way of limiting input options (like an enum)
	if accountType != "asset" && accountType != "liability" {
		panic("only `asset` or `liability` are valid accountType options")
	}

	var result []BalancesWithDate
	accountIDs := db.Select("id").Where("account_type=?", accountType).Table("accounts")

	//create a slice of months in Go instead of relying on SQL
	months := []time.Time{}
	for month := startYearMonth; month.Before(endYearMonth); month = month.AddDate(0, 1, 0) {
		// Set the date to the first of the month
		month = month.AddDate(0, 0, 1-month.Day())
		months = append(months, month)
	}

	for _, month := range months {
		var balances []Balance

		// Convert dates to YYYY-MM-DD so date comparisons work consistently with strings in SQLite
		lastDayOfMonth := utils.TimeToISO8601DateString(month.AddDate(0, 1, -1))

		db.Raw("select *, max(effective_date) from balances "+
			"where account_id in (?) "+
			"and effective_date <= (?) "+
			"group by account_id "+
			"order by effective_date", accountIDs, lastDayOfMonth).Scan(&balances)

		result = append(result, BalancesWithDate{
			Date:     month,
			Balances: balances,
		})
	}

	return result
}

func (*BalanceRepository) GetLatestBalanceForAccount(ctx context.Context, accountID uint) Balance {
	var balance Balance
	db.Where("account_id = ?", accountID).Order("effective_date desc").Limit(1).Find(&balance)
	return balance
}

func (*BalanceRepository) GetBalancesForAccount(ctx context.Context, accountID uint) []Balance {
	var balances []Balance
	db.Preload(clause.Associations).Where("account_id = ?", accountID).Order("effective_date desc").Find(&balances)
	return balances
}

func (*BalanceRepository) GetBalanceByID(ctx context.Context, balanceID uint) Balance {
	var balance Balance
	db.Preload(clause.Associations).Where("id = ?", balanceID).Find(&balance)
	return balance
}

// Save is an UPSERT operation, returning the ID of the record and an optional error
func (*BalanceRepository) Save(balance Balance) (id uint, err error) {
	result := db.Save(&balance).Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}})
	return balance.ID, result.Error
}

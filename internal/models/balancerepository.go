package models

import (
	"context"
	"fmt"
	"time"

	"github.com/alexdglover/sage/internal/utils"
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

func (br BalanceRepository) GetAllBalances(ctx context.Context) ([]Balance, error) {
	var balances []Balance
	result := db.Find(&balances)
	return balances, result.Error
}

func (br BalanceRepository) GetBalancesOfAllAssetsByMonth(ctx context.Context, startYearMonth time.Time, endYearMonth time.Time) []BalancesWithDate {
	return balanceRepository.GetBalancesByMonth(ctx, "asset", startYearMonth, endYearMonth)
}

func (br BalanceRepository) GetBalancesOfAllLiabilitiesByMonth(ctx context.Context, startYearMonth time.Time, endYearMonth time.Time) []BalancesWithDate {
	return balanceRepository.GetBalancesByMonth(ctx, "liability", startYearMonth, endYearMonth)
}

func (br BalanceRepository) GetBalancesByMonth(ctx context.Context, accountType string, startYearMonth time.Time, endYearMonth time.Time) []BalancesWithDate {
	// TODO: implement a better way of limiting input options (like an enum)
	if accountType != "asset" && accountType != "liability" {
		panic("only `asset` or `liability` are valid accountType options")
	}

	fmt.Println("startYearMonth is", startYearMonth)
	fmt.Println("endYearMonth is", endYearMonth)

	var result []BalancesWithDate
	assetAccountIds := db.Select("id").Where("account_type=?", accountType).Table("accounts")

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
		effectiveStartDate := utils.TimeToISO8601DateString(month)
		effectiveEndDate := utils.TimeToISO8601DateString(month.AddDate(0, 1, -1))

		db.Where("account_id IN (?)", assetAccountIds).
			Where("effective_start_date <= ?", effectiveStartDate).
			Where("effective_end_date IS NULL OR effective_end_date >= ?", effectiveEndDate).Find(&balances)

		result = append(result, BalancesWithDate{
			Date:     month,
			Balances: balances,
		})
	}

	return result
}

package models

import (
	"context"
	"fmt"
	"time"
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

func (br BalanceRepository) GetBalancesOfAllAssets(ctx context.Context, startYearMonth string, endYearMonth string) []Balance {
	var result []Balance
	assetAccountIds := db.Select("id").Where("account_type=?", "asset").Table("accounts")
	db.Select("*").
		Where("account_id IN (?)", assetAccountIds).
		Where("effective_start_date >= ?", startYearMonth).
		Where(db.Where("effective_end_date < ?", endYearMonth).Or("effective_end_date IS NULl")).Find(&result)

	return result
}

func (br BalanceRepository) GetBalancesOfAllAssetsByMonth(ctx context.Context, startYearMonth time.Time, endYearMonth time.Time) []BalancesWithDate {
	var result []BalancesWithDate
	assetAccountIds := db.Select("id").Where("account_type=?", "asset").Table("accounts")

	//create a slice of months in Go instead of relying on SQL
	months := []time.Time{}
	for month := startYearMonth; month.Before(endYearMonth.AddDate(0, 1, 0)); month = month.AddDate(0, 1, 0) {
		months = append(months, month)
	}

	fmt.Println("months are ", months)

	for _, month := range months {
		var balances []Balance
		nextMonth := month.AddDate(0, 1, 0)
		db.Where("account_id IN (?)", assetAccountIds).
			Where("effective_start_date <= ?", nextMonth).
			Where("effective_end_date IS NULL OR effective_end_date >= ?", month).Find(&balances)
		result = append(result, BalancesWithDate{
			Date:     month,
			Balances: balances,
		})
	}

	return result
}

func (br BalanceRepository) GetBalancesOfAllLiabilities(ctx context.Context, startYearMonth time.Time, endYearMonth time.Time) []Balance {
	var result []Balance
	liabilityAccountIds := db.Select("id").Where("account_type=?", "liability").Table("accounts")
	db.Select("*").
		Where("account_id IN (?)", liabilityAccountIds).
		Where("effective_start_date >= ?", startYearMonth).
		Where(db.Where("effective_end_date < ?", endYearMonth).Or("effective_end_date IS NULl")).Find(&result)

	return result
}

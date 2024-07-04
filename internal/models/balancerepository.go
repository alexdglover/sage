package models

import "context"

type BalanceRepository struct{}

var balanceRepository *BalanceRepository

type BalanceAmount struct {
	Date    string
	Balance int
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

func (br BalanceRepository) GetBalancesOfAllLiabilities(ctx context.Context, startYearMonth string, endYearMonth string) []Balance {
	var result []Balance
	liabilityAccountIds := db.Select("id").Where("account_type=?", "liability").Table("accounts")
	db.Select("*").
		Where("account_id IN (?)", liabilityAccountIds).
		Where("effective_start_date >= ?", startYearMonth).
		Where(db.Where("effective_end_date < ?", endYearMonth).Or("effective_end_date IS NULl")).Find(&result)

	return result
}

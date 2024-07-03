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
	// TODO: Replace this with something more idiomatic to GORM if possible
	var result []Balance
	db.Raw(`SELECT *
		FROM balances
		WHERE account_id in (
			SELECT id
			FROM accounts
			WHERE account_type="asset"
		)
		AND effective_start_date >= ?
		AND (effective_end_date < ? or effective_end_date is null)`,
		startYearMonth,
		endYearMonth,
	).Scan(&result)
	return result
}

func (br BalanceRepository) GetBalancesOfAllLiabilities(ctx context.Context, startYearMonth string, endYearMonth string) []Balance {
	// TODO: Replace this with something more idiomatic to GORM if possible
	var result []Balance
	db.Raw(`SELECT *
		FROM balances
		WHERE account_id in (
			SELECT id
			FROM accounts
			WHERE account_type="liability"
		)
		AND effective_start_date >= ?
		AND (effective_end_date < ? or effective_end_date is null)`,
		startYearMonth,
		endYearMonth,
	).Scan(&result)
	return result
}

// func (br BalanceRepository) getBalanceAllAssets(ctx context.Context, startDate string, endDate string) BalanceAmount {
// 	// var balanceAmount float32
// 	// var result BalanceAmount
// 	db.Raw(`SELECT date, SUM(balance) as balance
// 	FROM balances
// 	WHERE
// 		account_id in (SELECT id from accounts where AssetOrLiability='Asset')
// 		AND effective_start_date >= ?
// 		AND (effective_end_date < ? or effective_end_date is null)`,
// 		startDate, endDate).Scan(&result)
// 	return result
// }

// func (br BalanceRepository) getBalanceAllLiabilities(ctx context.Context, startDate string, endDate string) BalanceAmount {
// 	// var balanceAmount float32
// 	var result BalanceAmount
// 	db.Raw(`SELECT date, SUM(balance) as balance
// 	FROM balances
// 	WHERE
// 		account_id in (SELECT id from accounts where AssetOrLiability='Liability')
// 		AND effective_start_date >= ?
// 		AND (effective_end_date < ? or effective_end_date is null)`,
// 		startDate, endDate).Scan(&result)
// 	return result
// }

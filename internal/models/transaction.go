package models

import (
	"context"
	"time"

	"github.com/alexdglover/sage/internal/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Transaction struct {
	gorm.Model
	Date               string
	Description        string
	Amount             int64
	Excluded           bool // Will be stored as 0 or 1 in SQLite
	Hash               string
	UseForTraining     bool
	AccountID          uint
	Account            Account
	CategoryID         uint
	Category           Category
	ImportSubmissionID *uint
	ImportSubmission   *ImportSubmission
}

type TransactionsByDate struct {
	Date         time.Time
	Transactions []Transaction
}

type TTMRollingAverageByDate struct {
	Date              time.Time
	TTMRollingAverage int64
}

type TransactionRepository struct {
	DB *gorm.DB
}

func (tr *TransactionRepository) GetAllTransactions() ([]Transaction, error) {
	// TODO: Need to implement pagination
	var txns []Transaction
	result := tr.DB.Preload(clause.Associations).Order("date desc").Find(&txns)
	return txns, result.Error
}

func (tr *TransactionRepository) GetTransactionsByHash(hash string, submissionID uint) ([]Transaction, error) {
	// Implement GORM query to look up transactions by hash
	var transactions []Transaction
	result := tr.DB.Where("import_submission_id != ?", submissionID).Where("hash = ?", hash).Find(&transactions)
	return transactions, result.Error
}

func (tr *TransactionRepository) Create(txn *Transaction) error {
	result := tr.DB.Create(txn)
	return result.Error
}

// Save is an UPSERT operation, returning the ID of the record and an optional error
func (tr *TransactionRepository) Save(txn Transaction) (id uint, err error) {
	result := tr.DB.Save(&txn).Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}})

	return txn.ID, result.Error
}

func (tr *TransactionRepository) GetTransactionsByImportSubmission(id uint) ([]Transaction, error) {
	var transactions []Transaction
	result := tr.DB.Preload(clause.Associations).Where("import_submission_id = ?", id).Find(&transactions)
	return transactions, result.Error
}

func (tr *TransactionRepository) GetTransactionByID(id uint) (Transaction, error) {
	var transaction Transaction
	result := tr.DB.Preload(clause.Associations).Where("id = ?", id).Find(&transaction)
	return transaction, result.Error
}

func (tr *TransactionRepository) GetTransactionsForTraining() ([]Transaction, error) {
	var transactions []Transaction
	result := tr.DB.Preload(clause.Associations).Where("use_for_training = ?", 1).Find(&transactions)
	return transactions, result.Error
}

func (tr *TransactionRepository) GetAllIncomeTransactions(ctx context.Context, startYearMonth time.Time, endYearMonth time.Time) (txnsByDate []TransactionsByDate, err error) {
	return tr.GetAllTransactionsByNetIncomeType(ctx, "income", startYearMonth, endYearMonth)
}

func (tr *TransactionRepository) GetAllExpenseTransactions(ctx context.Context, startYearMonth time.Time, endYearMonth time.Time) (txnsByDate []TransactionsByDate, err error) {
	return tr.GetAllTransactionsByNetIncomeType(ctx, "expense", startYearMonth, endYearMonth)
}

func (tr *TransactionRepository) GetAllTransactionsByNetIncomeType(ctx context.Context, incomeOrExpense string, startYearMonth time.Time, endYearMonth time.Time) (txnsByDate []TransactionsByDate, err error) {
	//create a slice of months in Go instead of relying on SQL
	months := []time.Time{}
	for month := startYearMonth; month.Before(endYearMonth); month = month.AddDate(0, 1, 0) {
		// Set the date to the first of the month
		month = month.AddDate(0, 0, 1-month.Day())
		months = append(months, month)
	}
	for _, month := range months {
		var transactions []Transaction
		// Convert dates to YYYY-MM-DD so date comparisons work consistently with strings in SQLite
		lastDayOfMonth := utils.TimeToISO8601DateString(month.AddDate(0, 1, -1))
		if incomeOrExpense == "income" {
			queryResult := tr.DB.Raw(`SELECT t.*
			FROM transactions AS t
			JOIN categories AS c
			ON c.id=t.category_id
			WHERE c.name="Income"
			AND date >= (?)
			AND date <= (?)`, month, lastDayOfMonth).Scan(&transactions)

			if queryResult.Error != nil {
				return []TransactionsByDate{}, queryResult.Error
			}
		} else if incomeOrExpense == "expense" {
			queryResult := tr.DB.Raw(`SELECT t.*
			FROM transactions AS t
			JOIN categories AS c
			ON c.id=t.category_id
			WHERE c.name not in ("Income", "Transfers")
			AND date >= (?)
			AND date <= (?)`, month, lastDayOfMonth).Scan(&transactions)

			if queryResult.Error != nil {
				return []TransactionsByDate{}, queryResult.Error
			}
		}

		txnsByDate = append(txnsByDate, TransactionsByDate{
			Date:         month,
			Transactions: transactions,
		})
	}

	return txnsByDate, nil
}

func (tr *TransactionRepository) GetTTMRollingAverage(ctx context.Context, yearMonth time.Time) (average int64, err error) {
	var netIncomeAmounts []int64

	twelveMonthsEarlier := yearMonth.AddDate(0, -12, 0)

	for month := yearMonth; month.After(twelveMonthsEarlier); month = month.AddDate(0, -1, 0) {
		var netIncome int64
		// Get dates for beginning and end of the month
		firstDayOfTheMonth := month.AddDate(0, 0, 1-month.Day())
		lastDayOfTheMonth := firstDayOfTheMonth.AddDate(0, 1, -1)
		firstDayOfTheMonthISO := utils.TimeToISO8601DateString(firstDayOfTheMonth)
		lastDayOfTheMonthISO := utils.TimeToISO8601DateString(lastDayOfTheMonth)

		netIncomeQuery := tr.DB.Raw(`WITH income AS (
				SELECT sum(t.amount) as amount,
				STRFTIME('%Y-%m', t.date) as yearmonth
				FROM transactions AS t
				JOIN categories AS c
				ON c.id=t.category_id
				WHERE c.name="Income"
				AND date >= (?)
				AND date <= (?)
				GROUP BY yearmonth
			),
			expenses AS (
				SELECT sum(t.amount) as amount,
				STRFTIME('%Y-%m', t.date) as yearmonth
				FROM transactions AS t
				JOIN categories AS c
				ON c.id=t.category_id
				WHERE c.name not in ("Income", "Transfers")
				AND date >= (?)
				AND date <= (?)
				GROUP BY yearmonth
			)
			SELECT COALESCE(income.amount, 0) - COALESCE(expenses.amount, 0) as net_income
			FROM income FULL OUTER JOIN expenses 
			ON income.yearmonth = expenses.yearmonth
		`, firstDayOfTheMonthISO, lastDayOfTheMonthISO, firstDayOfTheMonthISO, lastDayOfTheMonthISO).Scan(&netIncome)

		if netIncomeQuery.Error != nil {
			return 0, netIncomeQuery.Error
		}

		// a user may not have transactions going back 12 months, so we need to skip months where there is no data
		// to avoid adding artificial zeros to the average
		if netIncome != 0 {
			netIncomeAmounts = append(netIncomeAmounts, netIncome)
		}
	}

	// Calculate the average of the net income amounts
	total := int64(0)
	for _, amount := range netIncomeAmounts {
		total += amount
	}
	if len(netIncomeAmounts) == 0 {
		return 0, nil
	}
	average = total / int64(len(netIncomeAmounts))
	return average, nil
}

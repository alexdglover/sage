package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Transaction struct {
	gorm.Model
	Date        string
	Description string
	Amount      int64
	// Excluded           int // 0 = false, 1 = true. SQLite doesn't have a boolean type
	Excluded           bool // Will be stored as 0 or 1 in SQLite
	Hash               string
	AccountID          uint
	Account            Account
	CategoryID         uint
	Category           Category
	ImportSubmissionID *uint
	ImportSubmission   *ImportSubmission
}

type TransactionRepository struct{}

var transactionRepository *TransactionRepository

func GetTransactionRepository() *TransactionRepository {
	if transactionRepository == nil {
		transactionRepository = &TransactionRepository{}
	}
	return transactionRepository
}

func (*TransactionRepository) GetAllTransactions() ([]Transaction, error) {
	// TODO: Need to implement pagination
	var txns []Transaction
	result := db.Preload(clause.Associations).Order("date desc").Find(&txns)
	return txns, result.Error
}

func (*TransactionRepository) GetTransactionsByHash(hash string, submissionID uint) ([]Transaction, error) {
	// Implement GORM query to look up transactions by hash
	var transactions []Transaction
	result := db.Where("import_submission_id != ?", submissionID).Where("hash = ?", hash).Find(&transactions)
	return transactions, result.Error
}

func (*TransactionRepository) Create(txn *Transaction) error {
	result := db.Create(txn)
	return result.Error
}

// Save is an UPSERT operation, returning the ID of the record and an optional error
func (*TransactionRepository) Save(txn Transaction) (id uint, err error) {
	result := db.Save(&txn).Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}})

	return txn.ID, result.Error
}

func (*TransactionRepository) GetTransactionsByImportSubmission(id uint) ([]Transaction, error) {
	var transactions []Transaction
	result := db.Preload(clause.Associations).Where("import_submission_id = ?", id).Find(&transactions)
	return transactions, result.Error
}

func (*TransactionRepository) GetTransactionByID(id uint) (Transaction, error) {
	var transaction Transaction
	result := db.Preload(clause.Associations).Where("id = ?", id).Find(&transaction)
	return transaction, result.Error
}

package models

import (
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
	result := tr.DB.Preload(clause.Associations).Where("use_for_training = ", 1).Find(&transactions)
	return transactions, result.Error
}

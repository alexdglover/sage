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
	Excluded           string
	Hash               string
	AccountId          uint
	Account            Account
	CategoryId         uint
	Category           Category
	ImportSubmissionId *uint
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

// TODO: Implement this function. This expects a sha256 hash that has been hex encoded to string
func (tr *TransactionRepository) GetTransactionsByHash(hash string, submission ImportSubmission) ([]Transaction, error) {
	// Implement GORM query to look up transactions by hash
	return []Transaction{}, nil
}

func (tr *TransactionRepository) Create(txn *Transaction) error {
	result := db.Create(txn)
	return result.Error
}

// Save is an UPSERT operation, returning the ID of the record and an optional error
func (*TransactionRepository) Save(txn Transaction) (id uint, err error) {
	result := db.Save(&txn).Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}})

	return txn.ID, result.Error
}

func (tr *TransactionRepository) GetTransactionsByImportSubmission(id uint) ([]Transaction, error) {
	var transactions []Transaction
	result := db.Where("import_submission_id = ?", id).Find(&transactions)
	return transactions, result.Error
}

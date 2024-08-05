package models

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	Date        string
	Description string
	Amount      int64
	Excluded    string
	Hash        string
	AccountId   uint
	Account     Account
	CategoryId  uint
	Category    Category
}

type TransactionRepository struct{}

var transactionRepository *TransactionRepository

func GetTransactionRepository() *TransactionRepository {
	if transactionRepository == nil {
		transactionRepository = &TransactionRepository{}
	}
	return transactionRepository
}

// This expects a sha256 hash that has been hex encoded to string
func (tr *TransactionRepository) FindTransactionsByHash(hash string) ([]Transaction, error) {
	// Implement GORM query to look up transactions by hash
	return []Transaction{}, nil
}

func (tr *TransactionRepository) Upsert(txn *Transaction) *gorm.DB {
	var result *gorm.DB
	if txn.ID != 0 {
		result = db.Save(txn)
	} else {
		result = db.Create(txn)
	}
	return result
}

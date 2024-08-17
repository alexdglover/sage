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

// TODO: Implement this function. This expects a sha256 hash that has been hex encoded to string
func (tr *TransactionRepository) FindTransactionsByHash(hash string, submission ImportSubmission) ([]Transaction, error) {
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

// func (tr *TransactionRepository) Upsert(txn *Transaction) *gorm.DB {
// 	var result *gorm.DB
// 	if txn.ID != 0 {
// 		result = db.Save(txn)
// 	} else {
// 		result = db.Create(txn)
// 	}
// 	return result
// }

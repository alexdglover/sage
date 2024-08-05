package services

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/alexdglover/sage/internal/models"
)

type ImportSubmissionResult struct {
	TransactionsImported int
	TransactionsSkipped  int
	BalancesImported     int
	BalancesSkipped      int
}

type Parser interface {
	Parse(string) ([]models.Transaction, []models.Balance, error)
}

type SchwabCSVParser struct{}

func (s SchwabCSVParser) Parse(statement string) (transactions []models.Transaction, balances []models.Balance, err error) {
	// parse the string into a CSV
	// for each row in the CSV, parse the columns and add it to transactions
	return []models.Transaction{}, []models.Balance{}, nil
}

var accountParser map[string]Parser = map[string]Parser{
	"schwab": SchwabCSVParser{},
}

func ImportStatement(statement string, accountID uint, parserID string) (result *ImportSubmissionResult, err error) {

	result = &ImportSubmissionResult{
		TransactionsImported: 0,
		TransactionsSkipped:  0,
		BalancesImported:     0,
		BalancesSkipped:      0,
	}
	var transactions []models.Transaction
	// balances := []models.Balance{}

	// parse the statement using the appropriate parser, getting a slice of transactions and balances
	parser := accountParser[parserID]
	transactions, _, err = parser.Parse(statement)
	if err != nil {
		return nil, err
	}

	hasher := sha256.New()
	tr := models.GetTransactionRepository()

	// for each transaction in the slice
	for _, transaction := range transactions {
		// Create a hash of all the relevant fields - date, amount, description
		builder := strings.Builder{}
		builder.WriteString(string(transaction.AccountId))
		builder.WriteString(string(transaction.Amount))
		builder.WriteString(transaction.Date)
		builder.WriteString(transaction.Description)

		hasher.Write([]byte(builder.String()))
		hash := hasher.Sum(nil)
		hashHex := hex.EncodeToString(hash)

		// check if hash already exists on existing transaction
		txns, err := tr.FindTransactionsByHash(hashHex)
		if err != nil {
			return nil, err
		}

		if len(txns) > 0 {
			fmt.Println("found existing transaction with same data, not adding it to database")
			result.TransactionsSkipped = result.TransactionsSkipped + 1
		} else {
			transaction.Hash = hashHex
			dbResult := tr.Upsert(&transaction)
			if dbResult.Error != nil {
				return nil, dbResult.Error
			}
			result.TransactionsImported = result.TransactionsImported + 1
		}
	}

	return result, nil
}

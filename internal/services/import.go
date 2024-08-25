package services

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/alexdglover/sage/internal/models"
)

type NoParserError struct{}

func (*NoParserError) Error() string {
	return "No parser was found for the provided account"
}

type AccountNotFoundError struct {
	AccountID uint
}

func (a *AccountNotFoundError) Error() string {
	return fmt.Sprintf("Could not find an account with ID %v", a.AccountID)
}

func ImportStatement(filename string, statement string, accountID uint) (result *models.ImportSubmission, err error) {
	isr := models.GetImportSubmissionRepository()

	submission := models.ImportSubmission{
		FileName:             filename,
		SubmissionDateTime:   time.Now().String(),
		Status:               models.Submitted,
		TransactionsImported: 0,
		TransactionsSkipped:  0,
		BalancesImported:     0,
		BalancesSkipped:      0,
		AccountID:            accountID,
	}
	id, err := isr.Save(submission)
	if err != nil {
		return nil, err
	}
	submission.ID = id

	var transactions []models.Transaction
	// var balances []models.Balance

	// parse the statement using the appropriate parser, getting a slice of transactions and balances
	ar := models.GetAccountRepository()
	account, err := ar.GetAccountByID(accountID)
	if err != nil {
		submission.Status = models.Failed
		isr.Save(submission)
		return nil, &AccountNotFoundError{}
	}
	if account.DefaultParser == nil {
		submission.Status = models.Failed
		isr.Save(submission)
		return nil, &NoParserError{}
	}
	parser := parsersByInstitution[*account.DefaultParser]
	transactions, _, err = parser.Parse(statement)
	if err != nil {
		submission.Status = models.Failed
		isr.Save(submission)
		return nil, err
	}

	hasher := sha256.New()
	tr := models.GetTransactionRepository()

	for idx, transaction := range transactions {
		if idx == 0 {
			submission.Status = models.Processing
			isr.Save(submission)
		}
		// Create a hash of all the relevant fields - date, amount, description
		builder := strings.Builder{}
		builder.WriteString(fmt.Sprint(transaction.AccountID))
		builder.WriteString(fmt.Sprint(transaction.Amount))
		builder.WriteString(transaction.Date)
		builder.WriteString(transaction.Description)
		hasher.Write([]byte(builder.String()))
		hash := hasher.Sum(nil)
		hashHex := hex.EncodeToString(hash)

		// use hash to check if this is a duplicate transaction, but ignore
		// duplicates from the statement currently being imported since it is possible
		// to have a transcation with same date, amount, description, and account
		txns, err := tr.GetTransactionsByHash(hashHex, submission)
		if err != nil {
			submission.Status = models.Failed
			isr.Save(submission)
			return nil, err
		}

		if len(txns) > 0 {
			fmt.Println("found existing transaction with same data, not adding it to database")
			submission.TransactionsSkipped = submission.TransactionsSkipped + 1
			continue
		}

		// Set the fields not directly sourced from the statement
		transaction.Hash = hashHex
		transaction.AccountID = accountID
		transaction.ImportSubmissionID = &submission.ID
		_, dbError := tr.Save(transaction)
		if dbError != nil {
			submission.Status = models.Failed
			isr.Save(submission)
			return nil, dbError
		}
		submission.TransactionsImported = submission.TransactionsImported + 1
	}

	submission.Status = models.Completed
	isr.Save(submission)

	result = &submission
	return result, nil
}

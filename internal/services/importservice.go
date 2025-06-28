package services

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/alexdglover/sage/internal/models"
)

type BalanceRepositoryInterface interface {
	Save(balance models.Balance) (uint, error)
}

type ImportSubmissionRepositoryInterface interface {
	Save(sub models.ImportSubmission) (uint, error)
}

type CategorizerInterface interface {
	BuildModel() error
	CategorizeTransaction(txn *models.Transaction) (models.Category, error)
}

// AccountRepositoryInterface specifically for ImportService
type ImportAccountRepositoryInterface interface {
	GetAccountByID(id uint) (models.Account, error)
}

// TransactionRepositoryInterface specifically for ImportService
type ImportTransactionRepositoryInterface interface {
	GetTransactionsByHash(hash string, submissionID uint) ([]models.Transaction, error)
	Save(txn models.Transaction) (uint, error)
}

type ImportService struct {
	AccountRepository          ImportAccountRepositoryInterface
	BalanceRepository          BalanceRepositoryInterface
	Categorizer                CategorizerInterface
	ImportSubmissionRepository ImportSubmissionRepositoryInterface
	TransactionRepository      ImportTransactionRepositoryInterface
}

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

func (is *ImportService) ImportStatement(filename string, statement string, accountID uint) (result *models.ImportSubmission, err error) {

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
	id, err := is.ImportSubmissionRepository.Save(submission)
	if err != nil {
		return nil, err
	}
	submission.ID = id

	var transactions []models.Transaction
	var balances []models.Balance

	// parse the statement using the appropriate parser, getting a slice of transactions and balances
	account, err := is.AccountRepository.GetAccountByID(accountID)
	if err != nil {
		submission.Status = models.Failed
		is.ImportSubmissionRepository.Save(submission)
		return nil, &AccountNotFoundError{}
	}
	if account.AccountType.DefaultParser == nil {
		submission.Status = models.Failed
		is.ImportSubmissionRepository.Save(submission)
		return nil, &NoParserError{}
	}
	parser := parsersByInstitution[*account.AccountType.DefaultParser]
	transactions, balances, err = parser.Parse(statement)
	if err != nil {
		submission.Status = models.Failed
		is.ImportSubmissionRepository.Save(submission)
		return nil, err
	}

	hasher := sha256.New()

	for idx, transaction := range transactions {
		if idx == 0 {
			submission.Status = models.Processing
			is.ImportSubmissionRepository.Save(submission)
		}
		transaction.AccountID = account.ID

		// Create a hash of all the relevant fields - date, amount, description
		builder := strings.Builder{}
		builder.WriteString(fmt.Sprint(transaction.AccountID))
		builder.WriteString(" ")
		builder.WriteString(fmt.Sprint(transaction.Amount))
		builder.WriteString(" ")
		builder.WriteString(transaction.Date)
		builder.WriteString(" ")
		builder.WriteString(transaction.Description)
		identifyingContent := builder.String()

		// Reset the hasher before generating a new hash
		hasher.Reset()
		hasher.Write([]byte(identifyingContent))
		hash := hasher.Sum(nil)
		hashHex := hex.EncodeToString(hash)

		// use hash to check if this is a duplicate transaction, but ignore
		// duplicates from the statement currently being imported since it is possible
		// to have a transcation with same date, amount, description, and account
		txns, err := is.TransactionRepository.GetTransactionsByHash(hashHex, submission.ID)
		if err != nil {
			submission.Status = models.Failed
			is.ImportSubmissionRepository.Save(submission)
			return nil, err
		}

		if len(txns) > 0 {
			submission.TransactionsSkipped = submission.TransactionsSkipped + 1
			continue
		}

		// Set the fields not directly sourced from the statement
		transaction.Hash = hashHex
		transaction.ImportSubmissionID = &submission.ID

		// TODO: Add a check for the category and set it to the default category if it is not set

		// Rebuild the model first. Consider making this optional
		is.Categorizer.BuildModel()
		category, err := is.Categorizer.CategorizeTransaction(&transaction)
		if err != nil {
			fmt.Println("error while categorizing transaction")
			fmt.Println(err)
			return nil, err
		}
		transaction.CategoryID = category.ID

		_, dbError := is.TransactionRepository.Save(transaction)
		if dbError != nil {
			submission.Status = models.Failed
			is.ImportSubmissionRepository.Save(submission)
			return nil, dbError
		}
		submission.TransactionsImported = submission.TransactionsImported + 1
	}

	for _, balance := range balances {
		balance.AccountID = accountID
		balance.ImportSubmissionID = &submission.ID
		is.BalanceRepository.Save(balance)
		submission.BalancesImported = submission.BalancesImported + 1
	}

	submission.Status = models.Completed
	is.ImportSubmissionRepository.Save(submission)

	result = &submission
	return result, nil
}

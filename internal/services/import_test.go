package services

import (
	"fmt"

	"github.com/alexdglover/sage/internal/models"
)

type TestAccountRepositoryForImport struct{}

func (tar TestAccountRepositoryForImport) GetAccountByID(id uint) (models.Account, error) {
	return models.Account{
		// ID:              id,
		Name:            "foo",
		AccountCategory: "checking",
		AccountType:     "asset",
	}, nil
}

type TestBalanceRepository struct{}

func (tbr TestBalanceRepository) Save(balance models.Balance) (id uint, err error) {
	return 1, nil
}

type TestImportSubmissionRepository struct{}

func (tisr TestImportSubmissionRepository) Save(submission models.ImportSubmission) (id uint, err error) {
	return 1, nil
}

type TestTransactionRepository struct{}

func (ttr TestTransactionRepository) GetTransactionsByHash(hash string, submissionID uint) ([]models.Transaction, error) {
	// For cases where we want to find a conflicting transaction
	if hash == "666" {
		return []models.Transaction{models.Transaction{
			// ID:                 1,
			Date:               "2024-01-02",
			Description:        "Coffee Shop",
			Amount:             1332,
			Excluded:           false,
			Hash:               "666",
			AccountID:          1,
			CategoryID:         1,
			ImportSubmissionID: nil,
		}}, nil
	} else {
		return []models.Transaction{}, nil
	}
}

func (ttr TestTransactionRepository) Save(txn models.Transaction) (id uint, err error) {
	return 1, nil
}

func ImportHappyPath() error {
	is := ImportService{
		AccountRepository:          TestAccountRepositoryForImport{},
		BalanceRepository:          TestBalanceRepository{},
		ImportSubmissionRepository: TestImportSubmissionRepository{},
		TransactionRepository:      TestTransactionRepository{},
	}

	result, err := is.ImportStatement("./some/file.path", "statement contents", 1)
	if err != nil {
		fmt.Println("test failed", err)
	}
	if result.Status != "whatever" {
		fmt.Println("test failed")
	}

	return nil
}

/*
	Tests to write
	* NoParserError
	* AccountNotFoundError
*/

// Ensure that duplicate transactions _within_ a single import are retained
// TODO: Write the actual test
func ImportWithDuplicateTransactions_test() error {
	is := ImportService{
		AccountRepository:          TestAccountRepositoryForImport{},
		BalanceRepository:          TestBalanceRepository{},
		ImportSubmissionRepository: TestImportSubmissionRepository{},
		TransactionRepository:      TestTransactionRepository{},
	}

	result, err := is.ImportStatement("./some/file.path", "statement contents", 1)
	if err != nil {
		fmt.Println("test failed", err)
	}
	if result.Status != "whatever" {
		fmt.Println("test failed")
	}

	return nil
}

// Ensure that duplicate transactions from a _different_ import are _not_ retained
// TODO: Write the actual test
func ImportWithConflictingTransactions_test() error {
	is := ImportService{
		AccountRepository:          TestAccountRepositoryForImport{},
		BalanceRepository:          TestBalanceRepository{},
		ImportSubmissionRepository: TestImportSubmissionRepository{},
		TransactionRepository:      TestTransactionRepository{},
	}

	result, err := is.ImportStatement("./some/file.path", "statement contents", 1)
	if err != nil {
		fmt.Println("test failed", err)
	}
	if result.Status != "whatever" {
		fmt.Println("test failed")
	}
	if result.TransactionsSkipped != 1 {
		fmt.Printf("test failed, expected 1 skipped transactions, got %v", result.TransactionsSkipped)
	}

	return nil
}

// func (is *ImportService) ImportStatement(filename string, statement string, accountID uint) (result *models.ImportSubmission, err error) {

// 	submission := models.ImportSubmission{
// 		FileName:             filename,
// 		SubmissionDateTime:   time.Now().String(),
// 		Status:               models.Submitted,
// 		TransactionsImported: 0,
// 		TransactionsSkipped:  0,
// 		BalancesImported:     0,
// 		BalancesSkipped:      0,
// 		AccountID:            accountID,
// 	}
// 	id, err := is.ImportSubmissionRepository.Save(submission)
// 	if err != nil {
// 		return nil, err
// 	}
// 	submission.ID = id

// 	var transactions []models.Transaction
// 	var balances []models.Balance

// 	// parse the statement using the appropriate parser, getting a slice of transactions and balances
// 	account, err := is.AccountRepository.GetAccountByID(accountID)
// 	if err != nil {
// 		submission.Status = models.Failed
// 		is.ImportSubmissionRepository.Save(submission)
// 		return nil, &AccountNotFoundError{}
// 	}
// 	if account.DefaultParser == nil {
// 		submission.Status = models.Failed
// 		is.ImportSubmissionRepository.Save(submission)
// 		return nil, &NoParserError{}
// 	}
// 	parser := parsersByInstitution[*account.DefaultParser]
// 	transactions, balances, err = parser.Parse(statement)
// 	if err != nil {
// 		submission.Status = models.Failed
// 		is.ImportSubmissionRepository.Save(submission)
// 		return nil, err
// 	}

// 	hasher := sha256.New()

// 	for idx, transaction := range transactions {
// 		if idx == 0 {
// 			submission.Status = models.Processing
// 			is.ImportSubmissionRepository.Save(submission)
// 		}
// 		// Create a hash of all the relevant fields - date, amount, description
// 		builder := strings.Builder{}
// 		builder.WriteString(fmt.Sprint(transaction.AccountID))
// 		builder.WriteString(fmt.Sprint(transaction.Amount))
// 		builder.WriteString(transaction.Date)
// 		builder.WriteString(transaction.Description)
// 		hasher.Write([]byte(builder.String()))
// 		hash := hasher.Sum(nil)
// 		hashHex := hex.EncodeToString(hash)

// 		// use hash to check if this is a duplicate transaction, but ignore
// 		// duplicates from the statement currently being imported since it is possible
// 		// to have a transcation with same date, amount, description, and account
// 		txns, err := is.TransactionRepository.GetTransactionsByHash(hashHex, submission.ID)
// 		if err != nil {
// 			submission.Status = models.Failed
// 			is.ImportSubmissionRepository.Save(submission)
// 			return nil, err
// 		}

// 		if len(txns) > 0 {
// 			submission.TransactionsSkipped = submission.TransactionsSkipped + 1
// 			continue
// 		}

// 		// Set the fields not directly sourced from the statement
// 		transaction.Hash = hashHex
// 		transaction.AccountID = accountID
// 		transaction.ImportSubmissionID = &submission.ID

// 		// TODO: Add a check for the category and set it to the default category if it is not set
// 		transaction.CategoryID = 1 // Default to Unknown

// 		_, dbError := is.TransactionRepository.Save(transaction)
// 		if dbError != nil {
// 			submission.Status = models.Failed
// 			is.ImportSubmissionRepository.Save(submission)
// 			return nil, dbError
// 		}
// 		submission.TransactionsImported = submission.TransactionsImported + 1
// 	}

// 	for _, balance := range balances {
// 		balance.AccountID = accountID
// 		balance.ImportSubmissionID = &submission.ID
// 		is.BalanceRepository.Save(balance)
// 		submission.BalancesImported = submission.BalancesImported + 1
// 	}

// 	submission.Status = models.Completed
// 	is.ImportSubmissionRepository.Save(submission)

// 	result = &submission
// 	return result, nil
// }

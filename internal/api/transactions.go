package api

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/alexdglover/sage/internal/models"
	"github.com/alexdglover/sage/internal/utils"
)

type TransactionController struct {
	AccountRepository     *models.AccountRepository
	CategoryRepository    *models.CategoryRepository
	TransactionRepository *models.TransactionRepository
}

//go:embed transactions.html
var transactionsPageTmpl string

//go:embed transactionForm.html
var transactionFormTmpl string

type TransactionDTO struct {
	ID                 uint
	Date               string
	Description        string
	Amount             string
	Excluded           bool
	AccountName        string
	CategoryName       string
	ImportSubmissionID string
}

type TransactionsPageDto struct {
	Transactions              []TransactionDTO
	TransactionUpdated        bool
	TransactionUpdatedMessage string
	Accounts                  []models.Account
	SelectedAccountID         uint
	Categories                []models.Category
	SelectedCategoryID        uint
	Description               string
	StartDate                 string
	EndDate                   string
}

type TransactionFormDTO struct {
	// If we're editing an existing transaction, Editing will be true
	// If we're creating a new transaction, Editing will be false
	Editing            bool
	TransactionID      uint
	Date               string
	Description        string
	Amount             string
	Excluded           bool
	AccountName        string
	CategoryName       string
	ImportSubmissionID string
	Accounts           []models.Account
	Categories         []models.Category
}

func (tc *TransactionController) generateTransactionsView(w http.ResponseWriter, req *http.Request) {
	tc.generateTransactionsViewContent(w, req, "")
}

func (tc *TransactionController) generateTransactionsViewContent(w http.ResponseWriter, req *http.Request, transactionUpdateMessage string) {
	transactionsDTO := []TransactionDTO{}
	dto := TransactionsPageDto{}

	// Parse the query parameters
	var accountID, categoryID uint
	var descriptionQueryParameter string
	var startDate, endDate *time.Time
	var err error
	if req != nil {
		query := req.URL.Query()
		accountIDQueryParameter := query.Get("accountID")
		categoryIDQueryParameter := query.Get("categoryID")
		descriptionQueryParameter = query.Get("description")
		startDateQueryParameter := query.Get("startDate")
		endDateQueryParameter := query.Get("endDate")

		if accountIDQueryParameter != "" {
			accountID, err = utils.StringToUint(accountIDQueryParameter)
			if err != nil {
				http.Error(w, "Unable to parse account ID", http.StatusInternalServerError)
				return
			}
			dto.SelectedAccountID = accountID
		}

		if categoryIDQueryParameter != "" {
			categoryID, err = utils.StringToUint(categoryIDQueryParameter)
			if err != nil {
				http.Error(w, "Unable to parse category ID", http.StatusInternalServerError)
				return
			}
			dto.SelectedCategoryID = categoryID
		}

		dto.Description = descriptionQueryParameter

		if startDateQueryParameter != "" {
			startDateValue := utils.ISO8601DateStringToTime(startDateQueryParameter)
			startDate = &startDateValue
			dto.StartDate = startDateQueryParameter
		}
		if endDateQueryParameter != "" {
			endDateValue := utils.ISO8601DateStringToTime(endDateQueryParameter)
			endDate = &endDateValue
			dto.EndDate = endDateQueryParameter
		}
	}

	// If startDate wasn't explicitly set, we will set it to 3 months ago to limit the number of transactions returned
	if startDate == nil {
		startDateValue := time.Now().AddDate(0, -3, 0)
		startDate = &startDateValue
		dto.StartDate = startDate.Format("2006-01-02")
	}

	// Get all Transactions with filters
	transactions, err := tc.TransactionRepository.GetAllTransactions(accountID, categoryID, descriptionQueryParameter, startDate, endDate)
	if err != nil {
		http.Error(w, "Unable to get transactions", http.StatusInternalServerError)
		return
	}

	// Build Transactions DTO
	for _, txn := range transactions {
		transactionsDTO = append(transactionsDTO, TransactionDTO{
			ID:                 txn.ID,
			Date:               txn.Date,
			Description:        txn.Description,
			Amount:             utils.CentsToDollarStringHumanized(txn.Amount),
			Excluded:           txn.Excluded,
			AccountName:        txn.Account.Name,
			CategoryName:       txn.Category.Name,
			ImportSubmissionID: utils.UintPointerToString(txn.ImportSubmissionID),
		})
	}
	dto.Transactions = transactionsDTO

	if transactionUpdateMessage != "" {
		dto.TransactionUpdated = true
		dto.TransactionUpdatedMessage = transactionUpdateMessage
	}

	// Get all accounts
	accounts, err := tc.AccountRepository.GetAllAccounts()
	if err != nil {
		http.Error(w, "Unable to get accounts", http.StatusInternalServerError)
		return
	}
	dto.Accounts = accounts
	// Get all categories
	categories, err := tc.CategoryRepository.GetAllCategories()
	if err != nil {
		http.Error(w, "Unable to get categories", http.StatusInternalServerError)
		return
	}
	dto.Categories = categories

	tmpl := template.Must(template.New("TransactionsPage").Funcs(template.FuncMap{
		"mod": func(i, j int) int { return i % j },
	}).Parse(transactionsPageTmpl))

	err = utils.RenderTemplateAsHTML(w, tmpl, dto)
	if err != nil {
		panic(err)
	}
}

func (tc *TransactionController) generateTransactionForm(w http.ResponseWriter, req *http.Request) {
	var dto TransactionFormDTO

	txnIDQueryParameter := req.URL.Query().Get("id")

	if txnIDQueryParameter != "" {
		txnID, err := utils.StringToUint(txnIDQueryParameter)
		if err != nil {
			http.Error(w, "Unable to parse transaction ID", http.StatusInternalServerError)
			return
		}
		txn, err := tc.TransactionRepository.GetTransactionByID(txnID)
		if err != nil {
			http.Error(w, "Unable to get Transaction", http.StatusInternalServerError)
			return
		}

		dto = TransactionFormDTO{
			Editing:            true,
			TransactionID:      txn.ID,
			Date:               txn.Date,
			Description:        txn.Description,
			Amount:             utils.CentsToDollarStringHumanized(txn.Amount),
			Excluded:           txn.Excluded,
			AccountName:        txn.Account.Name,
			CategoryName:       txn.Category.Name,
			ImportSubmissionID: utils.UintPointerToString(txn.ImportSubmissionID),
		}
	}

	accounts, err := tc.AccountRepository.GetAllAccounts()
	if err != nil {
		http.Error(w, "Unable to get accounts", http.StatusInternalServerError)
	}
	dto.Accounts = accounts

	categories, err := tc.CategoryRepository.GetAllCategories()
	if err != nil {
		http.Error(w, "Unable to get categories", http.StatusInternalServerError)
	}
	dto.Categories = categories

	tmpl, err := template.New("transactionForm").Parse(transactionFormTmpl)
	if err != nil {
		panic(err)
	}

	err = utils.RenderTemplateAsHTML(w, tmpl, dto)
	if err != nil {
		panic(err)
	}
}

func (tc *TransactionController) upsertTransaction(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	date := req.FormValue("date")
	description := req.FormValue("description")
	amount := req.FormValue("amount")
	excludedSelector := req.FormValue("excluded")
	var excluded bool
	if excludedSelector == "on" {
		excluded = true
	} else if excludedSelector == "off" || excludedSelector == "" {
		excluded = false
	} else {
		http.Error(w, "Unable to parse excluded", http.StatusInternalServerError)
		return
	}
	accountID, err := utils.StringToUint(req.FormValue("account"))
	if err != nil {
		http.Error(w, "Unable to parse accountID", http.StatusInternalServerError)
		return
	}
	categoryID, err := utils.StringToUint(req.FormValue("category"))
	if err != nil {
		http.Error(w, "Unable to parse categoryID", http.StatusInternalServerError)
		return
	}

	// If there is a transactionID, we are editing an existing transaction
	// so we should pull the existing record from the DB
	var transaction models.Transaction
	transactionID := req.FormValue("transactionID")
	if transactionID != "" {
		id, err := utils.StringToUint(transactionID)
		if err != nil {
			http.Error(w, "Unable to parse transactionID", http.StatusBadRequest)
			return
		}
		transaction, err = tc.TransactionRepository.GetTransactionByID(id)
		if err != nil {
			errorMessage := "Unable to get transaction with ID " + transactionID + " from the database"
			http.Error(w, errorMessage, http.StatusBadRequest)
			return
		}
	} else {
		transaction = models.Transaction{}
	}

	category, err := tc.CategoryRepository.GetCategoryByID(categoryID)
	if err != nil {
		http.Error(w, "Unable to get category by categoryID", http.StatusInternalServerError)
		return
	}

	transaction.Date = date
	transaction.Description = description
	transaction.Amount = utils.DollarStringToCents(amount)
	transaction.Excluded = excluded
	transaction.AccountID = accountID
	transaction.CategoryID = categoryID
	transaction.Category = category
	transaction.UseForTraining = true

	_, err = tc.TransactionRepository.Save(transaction)
	if err != nil {
		http.Error(w, "Unable to save transaction", http.StatusBadRequest)
		return
	}

	tc.generateTransactionsViewContent(w, nil, "Transaction saved successfully")
}

func (tc *TransactionController) deleteTransaction(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	transactionID, err := utils.StringToUint(req.FormValue("transactionID"))
	if err != nil {
		http.Error(w, "Unable to parse transactionID", http.StatusInternalServerError)
		return
	}

	err = tc.TransactionRepository.DeleteTransactionByID(context.TODO(), transactionID)

	if err != nil {
		http.Error(w, "Unable to delete transaction", http.StatusInternalServerError)
		return
	}

	tc.generateTransactionsViewContent(w, nil, fmt.Sprintf("Transaction %v deleted successfully", transactionID))
}

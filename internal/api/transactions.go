package api

import (
	_ "embed"
	"net/http"
	"net/url"
	"text/template"

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

type TransactionsPageDTO struct {
	Transactions           []TransactionDTO
	TransactionSaved       bool
	CreatedTransactionName string
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
	// Get all Transactions
	transactions, err := tc.TransactionRepository.GetAllTransactions()
	if err != nil {
		http.Error(w, "Unable to get transactions", http.StatusInternalServerError)
		return
	}

	// Build Transactions DTO
	transactionsDTO := make([]TransactionDTO, len(transactions))
	for i, txn := range transactions {
		transactionsDTO[i] = TransactionDTO{
			ID:                 txn.ID,
			Date:               txn.Date,
			Description:        txn.Description,
			Amount:             utils.CentsToDollarStringHumanized(txn.Amount),
			Excluded:           txn.Excluded,
			AccountName:        txn.Account.Name,
			CategoryName:       txn.Category.Name,
			ImportSubmissionID: utils.UintPointerToString(txn.ImportSubmissionID),
		}
	}
	TransactionsPageDTO := TransactionsPageDTO{
		Transactions: transactionsDTO,
	}
	if req.URL.Query().Get("TransactionSaved") != "" {
		TransactionsPageDTO.TransactionSaved = true
		TransactionsPageDTO.CreatedTransactionName = req.URL.Query().Get("TransactionSaved")
	}

	tmpl := template.Must(template.New("TransactionsPage").Funcs(template.FuncMap{
		"mod": func(i, j int) int { return i % j },
	}).Parse(transactionsPageTmpl))

	err = tmpl.Execute(w, TransactionsPageDTO)
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

	err = tmpl.Execute(w, dto)
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

	queryValues := url.Values{}
	queryValues.Add("transactionSaved", transactionID)
	// TODO: Consider moving the TransactionView to a function that accepts an extra argument
	// instead of invoking the endpoint with a custom request
	transactionViewReq := http.Request{
		Method: "GET",
		URL: &url.URL{
			RawQuery: queryValues.Encode(),
		},
	}
	transactionViewReq.URL.RawQuery = queryValues.Encode()

	tc.generateTransactionsView(w, &transactionViewReq)
}

package api

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"net/url"
	"text/template"

	"github.com/alexdglover/sage/internal/models"
	"github.com/alexdglover/sage/internal/utils"
)

//go:embed accounts.html
var accountsPageTmpl string

//go:embed accountForm.html
var accountFormTmpl string

type AccountDTO struct {
	ID                 uint
	Name               string
	AccountCategory    string
	AccountType        string
	DefaultParser      *string
	Balance            string
	BalanceLastUpdated string
}

type AccountsPageDTO struct {
	Accounts           []AccountDTO
	AccountSaved       bool
	CreatedAccountName string
}

type AccountFormDTO struct {
	// If we're editing an existing account, Editing will be true
	// If we're creating a new account, Editing will be false
	Editing         bool
	AccountID       string
	AccountName     string
	AccountCategory string
	AccountType     string
	DefaultParser   string
}

func accountsHandler(w http.ResponseWriter, req *http.Request) {
	// Get all accounts
	ar := models.GetAccountRepository()
	accounts, err := ar.GetAllAccounts()
	if err != nil {
		http.Error(w, "Unable to get accounts", http.StatusInternalServerError)
		return
	}

	br := models.GetBalanceRepository()

	// Build accounts DTO
	accountsDTO := make([]AccountDTO, len(accounts))
	for i, account := range accounts {
		balance := br.GetLatestBalanceForAccount(context.TODO(), account.ID)
		accountsDTO[i] = AccountDTO{
			ID:                 account.ID,
			Name:               account.Name,
			AccountCategory:    account.AccountCategory,
			AccountType:        account.AccountType,
			DefaultParser:      account.DefaultParser,
			Balance:            utils.CentsToDollarString(balance.Amount),
			BalanceLastUpdated: balance.Date,
		}
	}
	accountsPageDTO := AccountsPageDTO{
		Accounts: accountsDTO,
	}
	if req.URL.Query().Get("accountSaved") != "" {
		accountsPageDTO.AccountSaved = true
		accountsPageDTO.CreatedAccountName = req.URL.Query().Get("accountSaved")
	}

	tmpl := template.Must(template.New("accountsPage").Funcs(template.FuncMap{
		"mod": func(i, j int) int { return i % j },
	}).Parse(accountsPageTmpl))

	err = tmpl.Execute(w, accountsPageDTO)
	if err != nil {
		panic(err)
	}
}

func accountFormHandler(w http.ResponseWriter, req *http.Request) {
	var dto AccountFormDTO

	accountIDQueryParameter := req.URL.Query().Get("accountID")
	if accountIDQueryParameter != "" {
		ar := models.GetAccountRepository()
		accountID, err := utils.StringToUint(accountIDQueryParameter)
		if err != nil {
			http.Error(w, "Unable to parse account ID", http.StatusInternalServerError)
			return
		}
		account, err := ar.GetAccountByID(accountID)
		if err != nil {
			http.Error(w, "Unable to get account", http.StatusInternalServerError)
			return
		}

		dto = AccountFormDTO{
			Editing:         true,
			AccountID:       fmt.Sprint(account.ID),
			AccountName:     account.Name,
			AccountCategory: account.AccountCategory,
			AccountType:     account.AccountType,
		}
		// If the account has a default parser, set it. Otherwise let it default to empty string
		if account.DefaultParser != nil {
			dto.DefaultParser = *account.DefaultParser
		}
	}

	tmpl, err := template.New("accountForm").Parse(accountFormTmpl)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(w, dto)
	if err != nil {
		panic(err)
	}
}

func accountController(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	accountID := req.FormValue("accountID")
	accountName := req.FormValue("accountName")
	accountCategory := req.FormValue("accountCategory")
	accountType := req.FormValue("accountType")
	defaultParser := req.FormValue("defaultParser")

	ar := models.GetAccountRepository()
	var account models.Account

	if accountID != "" {
		id, err := utils.StringToUint(accountID)
		if err != nil {
			http.Error(w, "Unable to parse account ID", http.StatusBadRequest)
			return
		}
		account, err = ar.GetAccountByID(id)
		if err != nil {
			http.Error(w, "Unable to get account", http.StatusBadRequest)
			return
		}
	} else {
		account = models.Account{}
	}

	account.Name = accountName
	account.AccountCategory = accountCategory
	account.AccountType = accountType
	account.DefaultParser = &defaultParser

	_, err := ar.Save(account)
	if err != nil {
		http.Error(w, "Unable to save account", http.StatusBadRequest)
		return
	}

	queryValues := url.Values{}
	queryValues.Add("accountSaved", accountName)
	// TODO: Consider moving the accountView to a function that accepts an extra argument
	// instead of invoking the endpoint with a custom request
	accountViewReq := http.Request{
		Method: "GET",
		URL: &url.URL{
			RawQuery: queryValues.Encode(),
		},
	}
	accountViewReq.URL.RawQuery = queryValues.Encode()

	accountsHandler(w, &accountViewReq)
}

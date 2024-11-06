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
	humanize "github.com/dustin/go-humanize"
)

type AccountController struct {
	AccountRepository *models.AccountRepository
	BalanceRepository *models.BalanceRepository
}

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
	// If we're updating an existing account in the form, Updating will be true
	// If we're creating a new account, Updating will be false
	Updating        bool
	AccountID       string
	AccountName     string
	AccountCategory string
	AccountType     string
	DefaultParser   string
}

func (ac *AccountController) generateAccountsView(w http.ResponseWriter, req *http.Request) {
	// Get all accounts
	accounts, err := ac.AccountRepository.GetAllAccounts()
	if err != nil {
		http.Error(w, "Unable to get accounts", http.StatusInternalServerError)
		return
	}

	// Build accounts DTO
	accountsDTO := make([]AccountDTO, len(accounts))
	for i, account := range accounts {
		balance := ac.BalanceRepository.GetLatestBalanceForAccount(context.TODO(), account.ID)

		// there may not be a balance associated with an account
		// in those cases, we want to display "Never" as the last updated date
		var balanceLastUpdated string
		// A real Balance will never have ID 0, but an unpopulated model.Balance
		// struct will use the default value for the ID field, which is 0
		if balance.ID == 0 {
			balanceLastUpdated = "Never"
		} else {
			balanceLastUpdated = humanize.Time(balance.UpdatedAt)
		}

		accountsDTO[i] = AccountDTO{
			ID:                 account.ID,
			Name:               account.Name,
			AccountCategory:    account.AccountCategory,
			AccountType:        account.AccountType,
			DefaultParser:      account.DefaultParser,
			Balance:            utils.CentsToDollarString(balance.Amount),
			BalanceLastUpdated: balanceLastUpdated,
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

func (ac *AccountController) generateAccountForm(w http.ResponseWriter, req *http.Request) {
	var dto AccountFormDTO

	accountIDQueryParameter := req.URL.Query().Get("accountID")
	if accountIDQueryParameter != "" {
		accountID, err := utils.StringToUint(accountIDQueryParameter)
		if err != nil {
			http.Error(w, "Unable to parse account ID", http.StatusInternalServerError)
			return
		}
		account, err := ac.AccountRepository.GetAccountByID(accountID)
		if err != nil {
			http.Error(w, "Unable to get account", http.StatusInternalServerError)
			return
		}

		dto = AccountFormDTO{
			Updating:        true,
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

func (ac *AccountController) upsertAccount(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	accountID := req.FormValue("accountID")
	accountName := req.FormValue("accountName")
	accountCategory := req.FormValue("accountCategory")
	accountType := req.FormValue("accountType")
	defaultParser := req.FormValue("defaultParser")

	var account models.Account

	if accountID != "" {
		id, err := utils.StringToUint(accountID)
		if err != nil {
			http.Error(w, "Unable to parse account ID", http.StatusBadRequest)
			return
		}
		account, err = ac.AccountRepository.GetAccountByID(id)
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

	_, err := ac.AccountRepository.Save(account)
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

	ac.generateAccountsView(w, &accountViewReq)
}

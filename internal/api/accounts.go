package api

import (
	"context"
	_ "embed"
	"net/http"
	"text/template"

	"github.com/alexdglover/sage/internal/models"
	"github.com/alexdglover/sage/internal/utils"
)

//go:embed accounts.html
var accountsPageTmpl string

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
	Accounts []AccountDTO
}

// TODO: Consider moving this into a service class that returns just the data needed
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

	tmpl := template.Must(template.New("accountsPage").Funcs(template.FuncMap{
		"mod": func(i, j int) int { return i % j },
	}).Parse(accountsPageTmpl))

	err = tmpl.Execute(w, accountsPageDTO)
	if err != nil {
		panic(err)
	}
}

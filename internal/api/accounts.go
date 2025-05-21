package api

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"text/template"

	"github.com/alexdglover/sage/internal/models"
	"github.com/alexdglover/sage/internal/services"
	"github.com/alexdglover/sage/internal/utils"
	humanize "github.com/dustin/go-humanize"
)

type AccountController struct {
	AccountManager      	*services.AccountManager
	AccountRepository   	*models.AccountRepository
	AccountTypeRepository 	*models.AccountTypeRepository
	BalanceRepository     	*models.BalanceRepository
	TransactionRepository 	*models.TransactionRepository
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
	TxnLastUpdated     string
}

type AccountsPageDTO struct {
	Accounts              []AccountDTO
	AccountUpdated        bool
	AccountUpdatedMessage string
}

type AccountFormDTO struct {
	// If we're updating an existing account in the form, Updating will be true
	// If we're creating a new account, Updating will be false
	Updating        bool
	AccountID       string
	AccountName     string
	AccountTypeName string
	AccountTypes    []models.AccountType // the DTO probably shouldn't be using the models
	DefaultParser   string
}

func (ac *AccountController) generateAccountsView(w http.ResponseWriter, req *http.Request) {
	ac.generateAccountsViewContent(w, "")
}

func (ac *AccountController) generateAccountsViewContent(w http.ResponseWriter, AccountUpdatedMessage string) {
	// Get all accounts
	accounts, err := ac.AccountRepository.GetAllAccounts()
	if err != nil {
		http.Error(w, "Unable to get accounts", http.StatusInternalServerError)
		return
	}

	// Build accounts DTO
	accountsDTO := make([]AccountDTO, len(accounts))
	for i, account := range accounts {
		latestBalance, err := ac.BalanceRepository.GetLatestBalanceForAccount(context.TODO(), account.ID)
		// there may not be a balance associated with an account
		// in those cases, we want to display "Never" as the last updated date
		var balanceLastUpdated string
		if err != nil {
			if err.Error() == "record not found" {
				balanceLastUpdated = "Never"
			} else {
				errorMessage := fmt.Sprintf("Error getting latest balance for account ID %d - %v", account.ID, err)
				http.Error(w, errorMessage, http.StatusInternalServerError)
				return
			}
		} else {
			balanceLastUpdated = humanize.Time(latestBalance.UpdatedAt)
		}

		latestTransaction, err := ac.TransactionRepository.GetLatestTransactionForAccount(account.ID)
		var txnLastUpdated string
		if err != nil {
			if err.Error() == "record not found" {
				txnLastUpdated = "Never"
			} else {
				errorMessage := fmt.Sprintf("Error getting latest transaction for account ID %d - %v", account.ID, err)
				http.Error(w, errorMessage, http.StatusInternalServerError)
				return
			}
		} else {
			txnLastUpdated = humanize.Time(latestTransaction.UpdatedAt)
		}

		accountsDTO[i] = AccountDTO{
			ID:                 account.ID,
			Name:               account.Name,
			AccountCategory:    account.AccountType.AccountCategory,
			AccountType:        account.AccountType.LedgerType,
			DefaultParser:      account.AccountType.DefaultParser,
			Balance:            utils.CentsToDollarStringHumanized(latestBalance.Amount),
			BalanceLastUpdated: balanceLastUpdated,
			TxnLastUpdated:     txnLastUpdated,
		}
	}
	accountsPageDTO := AccountsPageDTO{
		Accounts: accountsDTO,
	}
	if AccountUpdatedMessage != "" {
		accountsPageDTO.AccountUpdated = true
		accountsPageDTO.AccountUpdatedMessage = AccountUpdatedMessage
	}

	tmpl := template.Must(template.New("accountsPage").Funcs(template.FuncMap{
		"mod": func(i, j int) int { return i % j },
	}).Parse(accountsPageTmpl))

	err = utils.RenderTemplateAsHTML(w, tmpl, accountsPageDTO)
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
			AccountTypeName: account.AccountType.Name,
		}
	}

	accountTypes, err := ac.AccountTypeRepository.GetAllAccountTypes()
	if err != nil {
		http.Error(w, "Unable to get account types", http.StatusInternalServerError)
		return
	}
	dto.AccountTypes = accountTypes

	tmpl, err := template.New("accountForm").Parse(accountFormTmpl)
	if err != nil {
		panic(err)
	}

	err = utils.RenderTemplateAsHTML(w, tmpl, dto)
	if err != nil {
		panic(err)
	}
}

func (ac *AccountController) upsertAccount(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		http.Error(w, "Unable to Parse Form ", http.StatusBadRequest)
		return
	}

	accountID := req.FormValue("accountID")
	accountName := req.FormValue("accountName")
	accountTypeIDFormValue := req.FormValue("accountTypeID")
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
	accountTypeID, err := utils.StringToUint(accountTypeIDFormValue)
	if err != nil {
		http.Error(w, "Unable to find parse an account type ID", http.StatusBadRequest)
		return
	}
	account.AccountTypeID = accountTypeID

	// Fetch AccountType to align with AccountTypeID using AccountManager
	accountType, err := ac.AccountManager.GetAccountTypeByID(accountTypeID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid account type ID %d: %v", accountTypeID, err), http.StatusBadRequest)
		return
	}
	account.AccountType = accountType

	_, err = ac.AccountRepository.Save(account)
	if err != nil {
		http.Error(w, "Unable to save account", http.StatusBadRequest)
		return
	}

	ac.generateAccountsViewContent(w, fmt.Sprintf("'%s' account saved", account.Name))
}

func (ac *AccountController) deleteAccount(w http.ResponseWriter, req *http.Request) {
	accountIDInput := req.FormValue("accountID")

	accountID, err := utils.StringToUint(accountIDInput)
	if err != nil {
		http.Error(w, "Unable to parse an account ID from input", http.StatusBadRequest)
		return
	}
	account, err := ac.AccountRepository.GetAccountByID(accountID)
	if err != nil {
		http.Error(w, "Unable to get account", http.StatusBadRequest)
		return
	}

	err = ac.AccountRepository.DeleteAccountByID(accountID)
	if err != nil {
		http.Error(w, "Unable to delete account", http.StatusBadRequest)
		return
	}

	ac.generateAccountsViewContent(w, fmt.Sprintf("'%s' account deleted", account.Name))
}

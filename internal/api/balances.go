package api

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"text/template"

	"github.com/alexdglover/sage/internal/models"
	"github.com/alexdglover/sage/internal/utils"
)

//go:embed balances.html
var balancesPageTmpl string

//go:embed balanceForm.html
var balanceFormTmpl string

type BalanceDTO struct {
	ID            uint
	UpdatedAt     string
	Date          string
	EffectiveDate string
	Amount        string
	AccountID     uint
	AccountName   string
}

type BalancesPageDTO struct {
	AccountID           uint
	Balances            []BalanceDTO
	BalanceSaved        bool
	BalanceSavedMessage string
}

type BalanceFormDTO struct {
	// If we're editing an existing account, Editing will be true
	// If we're creating a new account, Editing will be false
	Editing bool
	// we set the AccountID for the to populate the balanceForm with the accountID
	// in HTML forms, so the account ID can be passed in the form submission even
	// if there is no Balance object set
	AccountID    uint
	BalanceDTO   BalanceDTO
	ErrorMessage string
}

func generateBalancesView(w http.ResponseWriter, req *http.Request) {
	// Get all balances for a given account
	br := models.GetBalanceRepository()
	accountIDQueryParameter := req.URL.Query().Get("accountID")
	accountID, err := utils.StringToUint(accountIDQueryParameter)
	if err != nil {
		http.Error(w, "Unable to parse account ID", http.StatusInternalServerError)
		return
	}
	balances := br.GetBalancesForAccount(context.TODO(), accountID)

	// Create balance DTO for each balance
	balancesDTO := make([]BalanceDTO, len(balances))
	for i, balance := range balances {
		balancesDTO[i] = BalanceDTO{
			ID:            balance.ID,
			UpdatedAt:     balance.UpdatedAt.String(),
			EffectiveDate: balance.EffectiveDate,
			Amount:        utils.CentsToDollarString(balance.Amount),
			AccountID:     balance.AccountID,
			AccountName:   balance.Account.Name,
		}
	}
	balancesPageDTO := BalancesPageDTO{
		AccountID: accountID,
		Balances:  balancesDTO,
	}
	if req.URL.Query().Get("balanceSaved") != "" {
		balancesPageDTO.BalanceSaved = true
		balancesPageDTO.BalanceSavedMessage = req.URL.Query().Get("balanceSaved")
	}

	tmpl, err := template.New("balancesPage").Parse(balancesPageTmpl)
	if err != nil {
		http.Error(w, "Unable to parse balancesPage template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, balancesPageDTO)
	if err != nil {
		http.Error(w, "Unable to render balancesPage template", http.StatusInternalServerError)
		return
	}
}

func generateBalanceForm(w http.ResponseWriter, req *http.Request) {
	var balanceID uint
	var accountID uint
	var err error
	balanceIDQueryParameter := req.URL.Query().Get("balanceID")
	if balanceIDQueryParameter != "" {
		balanceID, err = utils.StringToUint(balanceIDQueryParameter)
		if err != nil {
			http.Error(w, "Unable to parse balance ID", http.StatusInternalServerError)
			return
		}
	}

	accountIDQueryParameter := req.URL.Query().Get("accountID")
	if accountIDQueryParameter != "" {
		accountID, err = utils.StringToUint(accountIDQueryParameter)
		if err != nil {
			http.Error(w, "Unable to parse account ID", http.StatusInternalServerError)
			return
		}
	}
	errorMessage := req.URL.Query().Get("errorMessage")
	balanceFormContent(w, balanceID, accountID, errorMessage)
}

func balanceFormContent(w http.ResponseWriter, balanceID uint, accountID uint, errorMessage string) {
	var dto BalanceFormDTO
	if balanceID != 0 {
		br := models.GetBalanceRepository()
		balance := br.GetBalanceByID(context.TODO(), balanceID)

		dto = BalanceFormDTO{
			Editing:   true,
			AccountID: accountID,
			BalanceDTO: BalanceDTO{
				ID:            balance.ID,
				UpdatedAt:     balance.UpdatedAt.String(),
				EffectiveDate: balance.EffectiveDate,
				Amount:        utils.CentsToDollarString(balance.Amount),
				AccountID:     balance.AccountID,
				AccountName:   balance.Account.Name,
			},
		}
	} else {
		dto.AccountID = accountID
	}

	if errorMessage != "" {
		dto.ErrorMessage = errorMessage
	}

	tmpl, err := template.New("balanceForm").Parse(balanceFormTmpl)
	if err != nil {
		http.Error(w, "Unable to parse balanceForm template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, dto)
	if err != nil {
		http.Error(w, "Unable to render balanceForm template", http.StatusInternalServerError)
		return
	}
}

func upsertBalance(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	balanceIDFormValue := req.FormValue("balanceID")
	balanceID, err := utils.StringToUint(balanceIDFormValue)
	if err != nil {
		http.Error(w, "Unable to parse balance ID", http.StatusBadRequest)
		return
	}
	accountIDFormValue := req.FormValue("accountID")
	accountID, err := utils.StringToUint(accountIDFormValue)
	if err != nil {
		http.Error(w, "Unable to parse account ID", http.StatusBadRequest)
		return
	}

	amount := req.FormValue("amount")
	// TODO: move this all into a utility function
	amount = strings.Replace(amount, ",", "", -1)
	amount = strings.Replace(amount, "$", "", -1)
	amount = strings.Replace(amount, " ", "", -1)
	if !utils.AmountValid(amount) {
		balanceFormContent(w, balanceID, accountID, fmt.Sprintf("%s is not a valid amount format", amount))
		return
	}

	effectiveDate := req.FormValue("effectiveDate")
	if !utils.DateValid(effectiveDate) {
		balanceFormContent(w, balanceID, accountID, fmt.Sprintf("%s is not a valid date format - please use YYYY-MM-DD", effectiveDate))
		return
	}

	br := models.GetBalanceRepository()
	var balance models.Balance

	if balanceID != 0 {
		balance.ID = balanceID
	}

	balance.Amount = utils.DollarStringToCents(amount)
	balance.EffectiveDate = effectiveDate
	balance.AccountID = accountID

	_, err = br.Save(balance)
	if err != nil {
		http.Error(w, "Unable to save balance", http.StatusBadRequest)
		return
	}

	// Redirect to the balances page with the balanceSaved query parameter set to true
	ar := models.GetAccountRepository()
	var balanceSavedMessage string
	account, err := ar.GetAccountByID(accountID)
	// If we can't get the account, we'll just show a generic message. Nothing actually broke
	if err != nil {
		balanceSavedMessage = "Balanced saved"
	} else {
		balanceSavedMessage = "Balanced saved for " + account.Name
	}

	queryValues := url.Values{}
	queryValues.Add("balanceSaved", balanceSavedMessage)
	queryValues.Add("accountID", accountIDFormValue)
	// TODO: Consider moving the balanceView to a function that accepts an extra argument
	// instead of invoking the endpoint with a custom request
	balanceViewReq := http.Request{
		Method: "GET",
		URL: &url.URL{
			RawQuery: queryValues.Encode(),
		},
	}
	balanceViewReq.URL.RawQuery = queryValues.Encode()

	generateBalancesView(w, &balanceViewReq)
}

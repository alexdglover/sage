package api

import (
	"context"
	_ "embed"
	"net/http"
	"text/template"

	"github.com/alexdglover/sage/internal/models"
	"github.com/alexdglover/sage/internal/utils"
)

//go:embed balances.html
var balancesPageTmpl string

//go:embed balanceForm.html
var balanceFormTmpl string

type BalanceDTO struct {
	ID                 uint
	UpdatedAt          string
	Date               string
	EffectiveStartDate string
	EffectiveEndDate   string
	Amount             string
	AccountID          uint
	AccountName        string
}

type BalancesPageDTO struct {
	AccountID    uint
	Balances     []BalanceDTO
	BalanceSaved bool
}

type BalanceFormDTO struct {
	// If we're editing an existing account, Editing will be true
	// If we're creating a new account, Editing will be false
	Editing    bool
	BalanceDTO BalanceDTO
}

func balancesHandler(w http.ResponseWriter, req *http.Request) {
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
		var eed string
		if balance.EffectiveEndDate != nil {
			eed = *balance.EffectiveEndDate
		}
		balancesDTO[i] = BalanceDTO{
			ID:                 balance.ID,
			UpdatedAt:          balance.UpdatedAt.String(),
			Date:               balance.Date,
			EffectiveStartDate: balance.EffectiveStartDate,
			EffectiveEndDate:   eed,
			Amount:             utils.CentsToDollarString(balance.Amount),
			AccountID:          balance.AccountID,
			AccountName:        balance.Account.Name,
		}
	}
	balancesPageDTO := BalancesPageDTO{
		AccountID: accountID,
		Balances:  balancesDTO,
	}
	if req.URL.Query().Get("balanceSaved") != "" {
		balancesPageDTO.BalanceSaved = true
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

func balanceFormHandler(w http.ResponseWriter, req *http.Request) {
	var dto BalanceFormDTO

	balanceIDQueryParameter := req.URL.Query().Get("balanceID")
	accountIDQueryParameter := req.URL.Query().Get("accountID")
	if balanceIDQueryParameter != "" {
		br := models.GetBalanceRepository()
		balanceID, err := utils.StringToUint(balanceIDQueryParameter)
		if err != nil {
			http.Error(w, "Unable to parse balance ID", http.StatusInternalServerError)
			return
		}
		balance := br.GetBalanceByID(context.TODO(), balanceID)

		var eed string
		if balance.EffectiveEndDate != nil {
			eed = *balance.EffectiveEndDate
		}
		dto = BalanceFormDTO{
			Editing: true,
			BalanceDTO: BalanceDTO{
				ID:                 balance.ID,
				UpdatedAt:          balance.UpdatedAt.String(),
				Date:               balance.Date,
				EffectiveStartDate: balance.EffectiveStartDate,
				EffectiveEndDate:   eed,
				Amount:             utils.CentsToDollarString(balance.Amount),
				AccountID:          balance.AccountID,
				AccountName:        balance.Account.Name,
			},
		}
	} else {
		accountID, err := utils.StringToUint(accountIDQueryParameter)
		if err != nil {
			http.Error(w, "Unable to parse account ID", http.StatusInternalServerError)
			return
		}
		dto.BalanceDTO.AccountID = accountID
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

// func accountController(w http.ResponseWriter, req *http.Request) {
// 	req.ParseForm()

// 	accountName := req.FormValue("accountName")
// 	accountCategory := req.FormValue("accountCategory")
// 	accountType := req.FormValue("accountType")
// 	defaultParser := req.FormValue("defaultParser")

// 	account := models.Account{
// 		Name:            accountName,
// 		AccountCategory: accountCategory,
// 		AccountType:     accountType,
// 		DefaultParser:   &defaultParser,
// 	}

// 	accountID := req.FormValue("accountID")
// 	if accountID != "" {
// 		id, err := utils.StringToUint(accountID)
// 		if err != nil {
// 			http.Error(w, "Unable to parse account ID", http.StatusBadRequest)
// 			return
// 		}
// 		account.ID = id
// 	}

// 	ar := models.GetAccountRepository()

// 	_, err := ar.Save(account)
// 	if err != nil {
// 		http.Error(w, "Unable to save account", http.StatusBadRequest)
// 		return
// 	}

// 	queryValues := url.Values{}
// 	queryValues.Add("accountSaved", accountName)
// 	// TODO: Consider moving the accountView to a function that accepts an extra argument
// 	// instead of invoking the endpoint with a custom request
// 	accountViewReq := http.Request{
// 		Method: "GET",
// 		URL: &url.URL{
// 			RawQuery: queryValues.Encode(),
// 		},
// 	}
// 	accountViewReq.URL.RawQuery = queryValues.Encode()

// 	accountsHandler(w, &accountViewReq)
// }

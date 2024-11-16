package api

import (
	"context"
	_ "embed"
	"log"
	"net/http"
)

//go:embed main.html.tmpl
var mainPageTmpl string

type ApiServer struct {
	AccountController     *AccountController
	BalanceController     *BalanceController
	BudgetController      *BudgetController
	ImportController      *ImportController
	MainController        *MainController
	NetWorthController    *NetWorthController
	TransactionController *TransactionController
}

func (as *ApiServer) StartApiServer(ctx context.Context) {
	http.HandleFunc("/", as.MainController.mainPageHandler)
	http.HandleFunc("/dashboard", dashboardHandler)
	http.HandleFunc("GET /net-worth", as.NetWorthController.netWorthHandler)
	http.HandleFunc("GET /import-form", as.ImportController.importStatementFormHandler)
	http.HandleFunc("POST /import-submission", as.ImportController.importSubmissionHandler)

	http.HandleFunc("GET /accounts", as.AccountController.generateAccountsView)
	http.HandleFunc("POST /accounts", as.AccountController.upsertAccount)
	http.HandleFunc("GET /accountForm", as.AccountController.generateAccountForm)

	http.HandleFunc("GET /balances", as.BalanceController.generateBalancesView)
	http.HandleFunc("POST /balances", as.BalanceController.upsertBalance)
	http.HandleFunc("GET /balanceForm", as.BalanceController.generateBalanceForm)

	http.HandleFunc("GET /budgets", as.BudgetController.generateBudgetsView)
	http.HandleFunc("POST /budgets", as.BudgetController.upsertBudget)
	http.HandleFunc("GET /budgetForm", as.BudgetController.generateBudgetForm)

	http.HandleFunc("GET /transactions", as.TransactionController.generateTransactionsView)
	http.HandleFunc("POST /transactions", as.TransactionController.upsertTransaction)
	http.HandleFunc("GET /transactionForm", as.TransactionController.generateTransactionForm)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

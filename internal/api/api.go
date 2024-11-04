package api

import (
	"context"
	_ "embed"
	"log"
	"net/http"
)

//go:embed main.html.tmpl
var mainPageTmpl string

func StartApiServer(ctx context.Context) {
	http.HandleFunc("/", mainPageHandler)
	http.HandleFunc("/dashboard", dashboardHandler)
	http.HandleFunc("GET /net-worth", netWorthHandler)
	http.HandleFunc("GET /import-form", importStatementFormHandler)
	http.HandleFunc("POST /import-submission", importSubmissionHandler)

	http.HandleFunc("GET /accounts", generateAccountsView)
	http.HandleFunc("POST /accounts", upsertAccount)
	http.HandleFunc("GET /accountForm", generateAccountForm)

	http.HandleFunc("GET /balances", generateBalancesView)
	http.HandleFunc("POST /balances", upsertBalance)
	http.HandleFunc("GET /balanceForm", generateBalanceForm)

	http.HandleFunc("GET /budgets", generateBudgetsView)
	http.HandleFunc("POST /budgets", upsertBudget)
	http.HandleFunc("GET /budgetForm", generateBudgetForm)

	http.HandleFunc("GET /transactions", generateTransactionsView)
	http.HandleFunc("POST /transactions", upsertTransaction)
	http.HandleFunc("GET /transactionForm", generateTransactionForm)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

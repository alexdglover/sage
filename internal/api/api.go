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
	http.HandleFunc("GET /accounts", accountsHandler)
	http.HandleFunc("POST /accounts", accountController)
	http.HandleFunc("GET /accountForm", accountFormHandler)
	http.HandleFunc("GET /transactions", transactionsHandler)
	http.HandleFunc("POST /transactions", transactionController)
	http.HandleFunc("GET /transactionForm", transactionFormHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

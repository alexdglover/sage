package api

import (
	"context"
	"embed"
	_ "embed"
	"log"
	"net/http"

	"github.com/alexdglover/sage/internal/utils/logger"
)

type ApiServer struct {
	AccountController     *AccountController
	BalanceController     *BalanceController
	BudgetController      *BudgetController
	CategoryController    *CategoryController
	ImportController      *ImportController
	NetIncomeController   *NetIncomeController
	NetWorthController    *NetWorthController
	SpendingController    *SpendingController
	TransactionController *TransactionController
	SettingsController    *SettingsController
}

//go:embed assets
var assets embed.FS

//go:embed pageComponents.html
var pageComponents string

func (as *ApiServer) StartApiServer(ctx context.Context) {

	http.Handle("/assets/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "max-age=86400, public")
		http.FileServer(http.FS(assets)).ServeHTTP(w, r)
	}))

	http.HandleFunc("/", as.AccountController.generateAccountsView)
	http.HandleFunc("/dashboard", dashboardHandler)
	http.HandleFunc("GET /net-worth", as.NetWorthController.netWorthHandler)
	http.HandleFunc("GET /net-income", as.NetIncomeController.netIncomeHandler)
	http.HandleFunc("GET /import-form", as.ImportController.importStatementFormHandler)
	http.HandleFunc("POST /import-submission", as.ImportController.importSubmissionHandler)

	http.HandleFunc("GET /spending-by-category", as.SpendingController.spendingByCategoryHandler)

	http.HandleFunc("GET /accounts", as.AccountController.generateAccountsView)
	http.HandleFunc("POST /accounts", as.AccountController.upsertAccount)
	http.HandleFunc("DELETE /accounts", as.AccountController.deleteAccount)
	http.HandleFunc("GET /accountForm", as.AccountController.generateAccountForm)

	http.HandleFunc("GET /balances", as.BalanceController.generateBalancesView)
	http.HandleFunc("POST /balances", as.BalanceController.upsertBalance)
	http.HandleFunc("GET /balanceForm", as.BalanceController.generateBalanceForm)

	http.HandleFunc("GET /budgets", as.BudgetController.generateBudgetsView)
	http.HandleFunc("POST /budgets", as.BudgetController.upsertBudget)
	http.HandleFunc("DELETE /budgets", as.BudgetController.deleteBudget)
	http.HandleFunc("GET /budgetForm", as.BudgetController.generateBudgetForm)
	http.HandleFunc("GET /budgetDetails", as.BudgetController.generateBudgetDetailsView)

	http.HandleFunc("GET /categories", as.CategoryController.generateCategoriesView)
	http.HandleFunc("POST /categories", as.CategoryController.upsertCategory)
	http.HandleFunc("DELETE /categories", as.CategoryController.deleteCategory)
	http.HandleFunc("GET /categoryForm", as.CategoryController.generateCategoryForm)

	http.HandleFunc("GET /transactions", as.TransactionController.generateTransactionsView)
	http.HandleFunc("POST /transactions", as.TransactionController.upsertTransaction)
	http.HandleFunc("DELETE /transactions", as.TransactionController.deleteTransaction)
	http.HandleFunc("GET /transactionForm", as.TransactionController.generateTransactionForm)

	http.HandleFunc("GET /settings", as.SettingsController.generateSettingsView)
	http.HandleFunc("POST /settings", as.SettingsController.upsertSettings)

	logger := logger.Get()
	logger.Info("Starting Server on http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

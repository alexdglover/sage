package api

import (
	_ "embed"
	"fmt"
	"net/http"
	"text/template"

	"github.com/alexdglover/sage/internal/models"
	"github.com/alexdglover/sage/internal/utils"
)

//go:embed budgets.html
var budgetsPageTmpl string

//go:embed budgetForm.html
var budgetsFormTmpl string

type BudgetDTO struct {
	ID           uint
	CategoryName string
	Amount       string
}

type BudgetsPageDTO struct {
	Budgets     []BudgetDTO
	BudgetSaved bool
}

type BudgetFormDTO struct {
	// If we're updating an existing budget in the form, Updating will be true
	// If we're creating a new budget, Updating will be false
	Updating     bool
	BudgetID     string
	CategoryName string
	Amount       string
	Categories   []models.Category
}

func generateBudgetsView(w http.ResponseWriter, req *http.Request) {
	sendViewResponse(w, false)
}

func generateBudgetForm(w http.ResponseWriter, req *http.Request) {
	var dto BudgetFormDTO

	cr := models.GetCategoryRepository()
	categories, err := cr.GetAllCategories()
	if err != nil {
		http.Error(w, "Unable to get categories", http.StatusInternalServerError)
	}
	dto.Categories = categories

	budgetIDQueryParameter := req.URL.Query().Get("budgetID")
	if budgetIDQueryParameter != "" {
		br := models.GetBudgetRepository()
		budgetID, err := utils.StringToUint(budgetIDQueryParameter)
		if err != nil {
			http.Error(w, "Unable to parse budget ID", http.StatusInternalServerError)
			return
		}
		budget, err := br.GetBudgetByID(budgetID)
		if err != nil {
			http.Error(w, "Unable to get budget", http.StatusInternalServerError)
			return
		}

		dto.Updating = true
		dto.BudgetID = fmt.Sprint(budget.ID)
		dto.CategoryName = budget.Category.Name
		dto.Amount = utils.CentsToDollarString(budget.Amount)
	} else {
		dto.Updating = false
	}

	tmpl, err := template.New("budgetForm").Parse(budgetsFormTmpl)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(w, dto)
	if err != nil {
		panic(err)
	}
}

func upsertBudget(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	budgetID := req.FormValue("budgetID")
	budgetCategory := req.FormValue("budgetCategory")
	amount := req.FormValue("amount")

	br := models.GetBudgetRepository()
	var budget models.Budget

	if budgetID != "" {
		id, err := utils.StringToUint(budgetID)
		if err != nil {
			http.Error(w, "Unable to parse budget ID", http.StatusBadRequest)
			return
		}
		budget, err = br.GetBudgetByID(id)
		if err != nil {
			http.Error(w, "Unable to get budget", http.StatusBadRequest)
			return
		}
	} else {
		budget = models.Budget{}
	}

	budgetCategoryID, err := utils.StringToUint(budgetCategory)
	if err != nil {
		http.Error(w, "Unable to parse category ID", http.StatusBadRequest)
		return
	}
	budget.CategoryID = budgetCategoryID
	cr := models.GetCategoryRepository()
	category, err := cr.GetCategoryByID(budget.CategoryID)
	if err != nil {
		panic("failed to get category from CategoryID")
	}
	budget.Category = category

	budget.Amount = utils.DollarStringToCents(amount)

	_, err = br.Save(budget)
	if err != nil {
		http.Error(w, "Unable to save budget", http.StatusBadRequest)
		return
	}

	sendViewResponse(w, true)
}

// Generic function to send the view response
func sendViewResponse(w http.ResponseWriter, update bool) {
	// Get all budgets
	br := models.GetBudgetRepository()
	budgets, err := br.GetAllBudgets()
	if err != nil {
		http.Error(w, "Unable to get budgets", http.StatusInternalServerError)
		return
	}

	// Build budgets DTO
	budgetsDTO := make([]BudgetDTO, len(budgets))
	for i, budget := range budgets {
		budgetsDTO[i] = BudgetDTO{
			ID:           budget.ID,
			CategoryName: budget.Category.Name,
			Amount:       utils.CentsToDollarString(budget.Amount),
		}
	}
	budgetsPageDTO := BudgetsPageDTO{
		Budgets: budgetsDTO,
	}
	if update {
		budgetsPageDTO.BudgetSaved = true
	}

	tmpl := template.Must(template.New("budgetsPage").Funcs(template.FuncMap{
		"mod": func(i, j int) int { return i % j },
	}).Parse(budgetsPageTmpl))

	err = tmpl.Execute(w, budgetsPageDTO)
	if err != nil {
		panic(err)
	}
}

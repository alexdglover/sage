package api

import (
	_ "embed"
	"fmt"
	"net/http"
	"text/template"

	"github.com/alexdglover/sage/internal/models"
	"github.com/alexdglover/sage/internal/services"
	"github.com/alexdglover/sage/internal/utils"
)

type BudgetController struct {
	BudgetRepository   *models.BudgetRepository
	BudgetService      *services.BudgetService
	CategoryRepository *models.CategoryRepository
}

//go:embed budgets.html
var budgetsPageTmpl string

//go:embed budgetForm.html
var budgetsFormTmpl string

type BudgetDTO struct {
	ID           uint
	CategoryName string
	Amount       string
	Spend        string
	PercentUsed  int
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

func (bc *BudgetController) generateBudgetsView(w http.ResponseWriter, req *http.Request) {
	bc.sendViewResponse(w, false)
}

func (bc *BudgetController) generateBudgetForm(w http.ResponseWriter, req *http.Request) {
	var dto BudgetFormDTO

	categories, err := bc.CategoryRepository.GetAllCategories()
	if err != nil {
		http.Error(w, "Unable to get categories", http.StatusInternalServerError)
	}
	dto.Categories = categories

	budgetIDQueryParameter := req.URL.Query().Get("budgetID")
	if budgetIDQueryParameter != "" {
		budgetID, err := utils.StringToUint(budgetIDQueryParameter)
		if err != nil {
			http.Error(w, "Unable to parse budget ID", http.StatusInternalServerError)
			return
		}
		budget, err := bc.BudgetRepository.GetBudgetByID(budgetID)
		if err != nil {
			http.Error(w, "Unable to get budget", http.StatusInternalServerError)
			return
		}

		dto.Updating = true
		dto.BudgetID = fmt.Sprint(budget.ID)
		dto.CategoryName = budget.Category.Name
		dto.Amount = utils.CentsToDollarStringHumanized(budget.Amount)
	} else {
		dto.Updating = false
	}

	// If a user creates a budget from the categories page, we should pre-populate the category dropdown even
	// though a budget doesn't exist yet
	categoryNameQueryParameter := req.URL.Query().Get("categoryName")
	if categoryNameQueryParameter != "" {
		dto.CategoryName = categoryNameQueryParameter
	}

	tmpl, err := template.New("budgetForm").Parse(budgetsFormTmpl)
	if err != nil {
		panic(err)
	}

	err = utils.RenderTemplateAsHTML(w, tmpl, dto)
	if err != nil {
		panic(err)
	}
}

func (bc *BudgetController) upsertBudget(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	budgetID := req.FormValue("budgetID")
	budgetCategory := req.FormValue("budgetCategory")
	amount := req.FormValue("amount")

	var budget models.Budget

	if budgetID != "" {
		id, err := utils.StringToUint(budgetID)
		if err != nil {
			http.Error(w, "Unable to parse budget ID", http.StatusBadRequest)
			return
		}
		budget, err = bc.BudgetRepository.GetBudgetByID(id)
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
	category, err := bc.CategoryRepository.GetCategoryByID(budget.CategoryID)
	if err != nil {
		panic("failed to get category from CategoryID")
	}
	budget.Category = category

	budget.Amount = utils.DollarStringToCents(amount)

	_, err = bc.BudgetRepository.Save(budget)
	if err != nil {
		http.Error(w, "Unable to save budget", http.StatusBadRequest)
		return
	}

	bc.sendViewResponse(w, true)
}

// Generic function to send the view response
func (bc *BudgetController) sendViewResponse(w http.ResponseWriter, update bool) {
	// Get all budget and associated spend
	budgetsAndSpend, err := bc.BudgetService.GetAllBudgetsAndCurrentSpend()

	// Build budgets DTO
	budgetsDTO := make([]BudgetDTO, len(budgetsAndSpend))
	for i, budget := range budgetsAndSpend {
		budgetsDTO[i] = BudgetDTO{
			ID:           budget.ID,
			CategoryName: budget.CategoryName,
			Amount:       utils.CentsToDollarStringHumanized(budget.Amount),
			Spend:        utils.CentsToDollarStringHumanized(budget.Spend),
			PercentUsed:  budget.PercentUsed,
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

	err = utils.RenderTemplateAsHTML(w, tmpl, budgetsPageDTO)
	if err != nil {
		panic(err)
	}
}
